package ante

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/cosmos/ethermint/crypto"
	evmtypes "github.com/cosmos/ethermint/x/evm/types"

	tmcrypto "github.com/tendermint/tendermint/crypto"
)

const (
	// TODO: Use this cost per byte through parameter or overriding NewConsumeGasForTxSizeDecorator
	// which currently defaults at 10, if intended
	// memoCostPerByte     sdk.Gas = 3
	secp256k1VerifyCost uint64 = 21000
)

// NewAnteHandler returns an ante handler responsible for attempting to route an
// Ethereum or SDK transaction to an internal ante handler for performing
// transaction-level processing (e.g. fee payment, signature verification) before
// being passed onto it's respective handler.
func NewAnteHandler(ak auth.AccountKeeper, sk types.SupplyKeeper) sdk.AnteHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx, sim bool,
	) (newCtx sdk.Context, err error) {
		var anteHandler sdk.AnteHandler
		switch tx.(type) {
		case auth.StdTx:
			anteHandler = sdk.ChainAnteDecorators(
				ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
				ante.NewMempoolFeeDecorator(),
				ante.NewValidateBasicDecorator(),
				ante.NewValidateMemoDecorator(ak),
				ante.NewConsumeGasForTxSizeDecorator(ak),
				ante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
				ante.NewValidateSigCountDecorator(ak),
				ante.NewDeductFeeDecorator(ak, sk),
				ante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
				ante.NewSigVerificationDecorator(ak),
				ante.NewIncrementSequenceDecorator(ak), // innermost AnteDecorator
			)

		case evmtypes.MsgEthereumTx:
			anteHandler = sdk.ChainAnteDecorators(
				NewEthSetupContextDecorator(), // outermost AnteDecorator. EthSetUpContext must be called first
				NewEthMempoolFeeDecorator(),
				NewEthSigVerificationDecorator(),
				NewAccountVerificationDecorator(ak),
				NewNonceVerificationDecorator(ak),
				NewEthGasConsumeDecorator(ak, sk),
				NewIncrementSenderSequenceDecorator(ak), // innermost AnteDecorator.
			)
		default:
			return ctx, sdk.ErrInternal(fmt.Sprintf("transaction type invalid: %T", tx))
		}

		return anteHandler(ctx, tx, sim)
	}
}

// sigGasConsumer overrides the DefaultSigVerificationGasConsumer from the x/auth
// module on the SDK. It doesn't allow ed25519 nor multisig thresholds.
func sigGasConsumer(
	meter sdk.GasMeter, sig []byte, pubkey tmcrypto.PubKey, params types.Params,
) error {
	switch pubkey.(type) {
	case crypto.PubKeySecp256k1:
		meter.ConsumeGas(secp256k1VerifyCost, "ante verify: secp256k1")
		return nil
	case tmcrypto.PubKey:
		meter.ConsumeGas(secp256k1VerifyCost, "ante verify: tendermint secp256k1")
		return nil
	default:
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidPubKey, "unrecognized public key type: %T", pubkey)
	}
}
