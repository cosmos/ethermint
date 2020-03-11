package evm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authutils "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	emint "github.com/cosmos/ethermint/types"
	"github.com/cosmos/ethermint/x/evm/types"

	tmtypes "github.com/tendermint/tendermint/types"
)

// NewHandler returns a handler for Ethermint type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.EthereumTxMsg:
			return handleEthTxMsg(ctx, k, msg)
		case *types.EmintMsg:
			return handleEmintMsg(ctx, k, *msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized ethermint message type: %T", msg)
		}
	}
}

// handleEthTxMsg handles an Ethereum specific tx
func handleEthTxMsg(ctx sdk.Context, k Keeper, msg types.EthereumTxMsg) (*sdk.Result, error) {
	// TODO: move to client
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// parse the chainID from a string to a base-10 integer
	intChainID, ok := new(big.Int).SetString(ctx.ChainID(), 10)
	if !ok {
		return nil, sdkerrors.Wrap(emint.ErrInvalidChainID, ctx.ChainID())
	}

	// Verify signature and retrieve sender address
	sender, err := msg.VerifySig(intChainID)
	if err != nil {
		return nil, err
	}

	// Encode transaction by default Tx encoder
	txEncoder := authutils.GetTxEncoder(types.ModuleCdc)
	txBytes, err := txEncoder(msg)
	if err != nil {
		return nil, err
	}
	txHash := tmtypes.Tx(txBytes).Hash()
	ethHash := common.BytesToHash(txHash)

	st := types.StateTransition{
		Sender:       sender,
		AccountNonce: msg.Data.AccountNonce,
		Price:        msg.Data.Price,
		GasLimit:     msg.Data.GasLimit,
		Recipient:    msg.Data.Recipient,
		Amount:       msg.Data.Amount,
		Payload:      msg.Data.Payload,
		Csdb:         k.CommitStateDB.WithContext(ctx),
		ChainID:      intChainID,
		THash:        &ethHash,
		Simulate:     ctx.IsCheckTx(),
	}
	// Prepare db for logs
	k.CommitStateDB.Prepare(ethHash, common.Hash{}, k.TxCount.Get())
	k.TxCount.Increment()

	// TODO: move to keeper
	bloom, res, err := st.TransitionCSDB(ctx)
	if err != nil {
		return nil, err
	}

	// update block bloom filter
	k.Bloom.Or(k.Bloom, bloom)

	return res, nil
}

func handleEmintMsg(ctx sdk.Context, k Keeper, msg types.EmintMsg) (*sdk.Result, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// parse the chainID from a string to a base-10 integer
	intChainID, ok := new(big.Int).SetString(ctx.ChainID(), 10)
	if !ok {
		return nil, sdkerrors.Wrap(emint.ErrInvalidChainID, ctx.ChainID())
	}

	st := types.StateTransition{
		Sender:       common.BytesToAddress(msg.From.Bytes()),
		AccountNonce: msg.AccountNonce,
		Price:        msg.Price.BigInt(),
		GasLimit:     msg.GasLimit,
		Amount:       msg.Amount.BigInt(),
		Payload:      msg.Payload,
		Csdb:         k.CommitStateDB.WithContext(ctx),
		ChainID:      intChainID,
		Simulate:     ctx.IsCheckTx(),
	}

	if msg.Recipient != nil {
		to := common.BytesToAddress(msg.Recipient.Bytes())
		st.Recipient = &to
	}

	// Prepare db for logs
	k.CommitStateDB.Prepare(common.Hash{}, common.Hash{}, k.TxCount.Get()) // Cannot provide tx hash
	k.TxCount.Increment()

	_, res, err := st.TransitionCSDB(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}
