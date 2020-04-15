package ante

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/types"

	emint "github.com/cosmos/ethermint/types"
	evmtypes "github.com/cosmos/ethermint/x/evm/types"

	ethcore "github.com/ethereum/go-ethereum/core"
)

type EthSetupContextDecorator struct{}

// NewEthSetupContextDecorator creates a new EthSetupContextDecorator
func NewEthSetupContextDecorator() EthSetupContextDecorator {
	return EthSetupContextDecorator{}
}

// AnteHandle verifies that enough fees have been provided by the
// Ethereum transaction that meet the minimum threshold set by the block
// proposer.
//
// NOTE: This should only be ran during a CheckTx mode.
func (escd EthSetupContextDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	// This is done to ignore costs in Ante handler checks
	ctx = ctx.WithBlockGasMeter(sdk.NewInfiniteGasMeter())

	return next(ctx, tx, simulate)
}

// EthMempoolFeeDecorator validates that sufficient fees have been provided that
// meet a minimum threshold defined by the proposer (for mempool purposes during CheckTx).
type EthMempoolFeeDecorator struct{}

// NewEthMempoolFeeDecorator creates a new EthMempoolFeeDecorator
func NewEthMempoolFeeDecorator() EthMempoolFeeDecorator {
	return EthMempoolFeeDecorator{}
}

// AnteHandle verifies that enough fees have been provided by the
// Ethereum transaction that meet the minimum threshold set by the block
// proposer.
//
// NOTE: This should only be ran during a CheckTx mode.
func (emfd EthMempoolFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if !ctx.IsCheckTx() {
		return next(ctx, tx, simulate)
	}

	ethTxMsg, ok := tx.(evmtypes.MsgEthereumTx)
	if !ok {
		return ctx, sdk.ErrInternal(fmt.Sprintf("invalid tx type %T", tx))
	}

	// fee = GP * GL
	fee := sdk.NewInt64DecCoin(emint.DenomDefault, ethTxMsg.Fee().Int64())

	minGasPrices := ctx.MinGasPrices()

	allGTE := true
	for _, v := range minGasPrices {
		if !fee.IsGTE(v) {
			allGTE = false
		}
	}

	// it is assumed that the minimum fees will only include the single valid denom
	if !ctx.MinGasPrices().IsZero() && !allGTE {
		// reject the transaction that does not meet the minimum fee
		return ctx, sdk.ErrInsufficientFee(
			fmt.Sprintf(
				"insufficient fee, got: %q required: %q", fee, ctx.MinGasPrices(),
			),
		)
	}

	return next(ctx, tx, simulate)
}

// EthIntrinsicGasDecorator validates enough intrinsic gas for the transaction.
type EthIntrinsicGasDecorator struct{}

// NewEthIntrinsicGasDecorator creates a new EthIntrinsicGasDecorator
func NewEthIntrinsicGasDecorator() EthIntrinsicGasDecorator {
	return EthIntrinsicGasDecorator{}
}

// AnteHandle validates that the Ethereum tx message has enough to
// cover intrinsic gas. Intrinsic gas for a transaction is the amount of gas
// that the transaction uses before the transaction is executed. The gas is a
// constant value of 21000 plus any cost inccured by additional bytes of data
// supplied with the transaction.
//
// NOTE: This should only be ran during a CheckTx mode.
func (eigd EthIntrinsicGasDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if !ctx.IsCheckTx() {
		return next(ctx, tx, simulate)
	}

	ethTxMsg, ok := tx.(evmtypes.MsgEthereumTx)
	if !ok {
		return ctx, sdk.ErrInternal(fmt.Sprintf("invalid tx type %T", tx))
	}

	gas, err := ethcore.IntrinsicGas(ethTxMsg.Data.Payload, ethTxMsg.To() == nil, true)
	if err != nil {
		return ctx, sdk.ErrInternal(fmt.Sprintf("failed to compute intrinsic gas cost: %s", err))
	}

	if ethTxMsg.Data.GasLimit < gas {
		return ctx, sdk.ErrInternal(
			fmt.Sprintf("intrinsic gas too low: %d < %d", ethTxMsg.Data.GasLimit, gas),
		)
	}

	return next(ctx, tx, simulate)
}

// EthSigVerificationDecorator validates an ethereum signature
type EthSigVerificationDecorator struct{}

// NewEthSigVerificationDecorator creates a new EthSigVerificationDecorator
func NewEthSigVerificationDecorator() EthSigVerificationDecorator {
	return EthSigVerificationDecorator{}
}

// AnteHandle validates the signature and returns sender address
func (esvd EthSigVerificationDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	ethTxMsg, ok := tx.(evmtypes.MsgEthereumTx)
	if !ok {
		return ctx, sdk.ErrInternal(fmt.Sprintf("invalid tx type %T", tx))
	}

	// parse the chainID from a string to a base-10 integer
	chainID, ok := new(big.Int).SetString(ctx.ChainID(), 10)
	if !ok {
		return ctx, emint.ErrInvalidChainID(fmt.Sprintf("invalid chainID: %s", ctx.ChainID()))
	}

	// validate sender/signature
	// NOTE: signer is retrieved from the transaction on the next AnteDecorator
	_, err = ethTxMsg.VerifySig(chainID)
	if err != nil {
		return ctx, sdk.ErrUnauthorized(fmt.Sprintf("signature verification failed: %s", err))
	}

	return next(ctx, ethTxMsg, simulate)
}

// AccountVerificationDecorator validates an account balance checks
type AccountVerificationDecorator struct {
	ak auth.AccountKeeper
}

// NewAccountVerificationDecorator creates a new AccountVerificationDecorator
func NewAccountVerificationDecorator(ak auth.AccountKeeper) AccountVerificationDecorator {
	return AccountVerificationDecorator{
		ak: ak,
	}
}

// AnteHandle validates the signature and returns sender address
func (avd AccountVerificationDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if !ctx.IsCheckTx() {
		return next(ctx, tx, simulate)
	}

	ethTxMsg, ok := tx.(evmtypes.MsgEthereumTx)
	if !ok {
		return ctx, sdk.ErrInternal(fmt.Sprintf("invalid tx type %T", tx))
	}

	// sender address should be in the tx cache
	address := ethTxMsg.From()
	if address == nil {
		panic("sender address is nil")
	}

	acc := avd.ak.GetAccount(ctx, address)
	if acc == nil {
		return ctx, sdk.ErrInternal(fmt.Sprintf("account %s is nil", address))
	}

	// on InitChain make sure account number == 0
	if ctx.BlockHeight() == 0 && acc.GetAccountNumber() != 0 {
		return ctx, sdk.ErrInternal(
			fmt.Sprintf(
				"invalid account number for height zero (got %d)", acc.GetAccountNumber(),
			),
		)
	}

	// validate sender has enough funds
	balance := acc.GetCoins().AmountOf(emint.DenomDefault)
	if balance.BigInt().Cmp(ethTxMsg.Cost()) < 0 {
		return ctx, sdk.ErrInsufficientFunds(
			fmt.Sprintf("insufficient funds: %s < %s", balance, ethTxMsg.Cost()),
		)
	}

	return next(ctx, tx, simulate)
}

// NonceVerificationDecorator that the nonce matches
type NonceVerificationDecorator struct {
	ak auth.AccountKeeper
}

// NewNonceVerificationDecorator creates a new NonceVerificationDecorator
func NewNonceVerificationDecorator(ak auth.AccountKeeper) NonceVerificationDecorator {
	return NonceVerificationDecorator{
		ak: ak,
	}
}

// AnteHandle validates that the transaction nonce is valid (equivalent to the sender account’s
// current nonce).
func (nvd NonceVerificationDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	ethTxMsg, ok := tx.(evmtypes.MsgEthereumTx)
	if !ok {
		return ctx, sdk.ErrInternal(fmt.Sprintf("invalid tx type %T", tx))
	}

	// sender address should be in the tx cache
	address := ethTxMsg.From()
	if address == nil {
		panic("sender address is nil")
	}

	acc := nvd.ak.GetAccount(ctx, address)
	if acc == nil {
		return ctx, sdk.ErrInternal(fmt.Sprintf("account %s is nil", address))
	}

	seq := acc.GetSequence()
	if ethTxMsg.Data.AccountNonce != seq {
		return ctx, sdk.ErrInvalidSequence(
			fmt.Sprintf("invalid nonce; got %d, expected %d", ethTxMsg.Data.AccountNonce, seq),
		)
	}

	return next(ctx, tx, simulate)
}

// EthGasConsumeDecorator
type EthGasConsumeDecorator struct {
	ak auth.AccountKeeper
	sk types.SupplyKeeper
}

// NewEthGasConsumeDecorator creates a new EthGasConsumeDecorator
func NewEthGasConsumeDecorator(ak auth.AccountKeeper, sk types.SupplyKeeper) EthGasConsumeDecorator {
	return EthGasConsumeDecorator{
		ak: ak,
		sk: sk,
	}
}

func (egcd EthGasConsumeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	ethTxMsg, ok := tx.(evmtypes.MsgEthereumTx)
	if !ok {
		return ctx, sdk.ErrInternal(fmt.Sprintf("invalid tx type %T", tx))
	}

	// sender address should be in the tx cache
	address := ethTxMsg.From()
	if address == nil {
		panic("sender address is nil")
	}

	// Fetch sender account from signature
	senderAcc, err := auth.GetSignerAcc(ctx, egcd.ak, address)
	if err != nil {
		return ctx, err
	}

	if senderAcc == nil {
		return ctx, sdk.ErrInternal(fmt.Sprintf("sender account %s is nil", address))
	}

	// Charge sender for gas up to limit
	if ethTxMsg.Data.GasLimit != 0 {
		// Cost calculates the fees paid to validators based on gas limit and price
		cost := new(big.Int).Mul(ethTxMsg.Data.Price, new(big.Int).SetUint64(ethTxMsg.Data.GasLimit))

		feeAmt := sdk.NewCoins(
			sdk.NewCoin(emint.DenomDefault, sdk.NewIntFromBigInt(cost)),
		)

		err = auth.DeductFees(egcd.sk, ctx, senderAcc, feeAmt)
		if err != nil {
			return ctx, err
		}
	}

	// Set gas meter after ante handler to ignore gaskv costs
	newCtx = auth.SetGasMeter(simulate, ctx, ethTxMsg.Data.GasLimit)

	gas, err := ethcore.IntrinsicGas(ethTxMsg.Data.Payload, ethTxMsg.To() == nil, true)
	if err != nil {
		return newCtx, err
	}

	newCtx.GasMeter().ConsumeGas(gas, "eth intrinsic gas")

	return next(newCtx, tx, simulate)
}

// IncrementSenderSequenceDecorator handles incrementing the sequence of the sender.
type IncrementSenderSequenceDecorator struct {
	ak auth.AccountKeeper
}

// NewIncrementSenderSequenceDecorator creates a new IncrementSenderSequenceDecorator.
func NewIncrementSenderSequenceDecorator(ak auth.AccountKeeper) IncrementSenderSequenceDecorator {
	return IncrementSenderSequenceDecorator{
		ak: ak,
	}
}

func (issd IncrementSenderSequenceDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	ethTxMsg, ok := tx.(evmtypes.MsgEthereumTx)
	if !ok {
		return ctx, sdk.ErrInternal(fmt.Sprintf("invalid tx type %T", tx))
	}

	// sender address should be in the tx cache
	address := ethTxMsg.From()
	if address == nil {
		panic("sender address is nil")
	}

	// Fetch sender account from signature
	senderAcc, err := auth.GetSignerAcc(ctx, issd.ak, address)
	if err != nil {
		return ctx, err
	}

	if senderAcc == nil {
		return ctx, sdk.ErrInternal(fmt.Sprintf("sender account %s is nil", address))
	}

	// Increment sequence of sender
	if err := senderAcc.SetSequence(senderAcc.GetSequence() + 1); err != nil {
		panic(err)
	}
	issd.ak.SetAccount(ctx, senderAcc)

	return next(ctx, tx, simulate)
}
