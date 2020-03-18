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
		case types.MsgEthereumTx:
			return HandleEthTxMsg(ctx, k, msg)
		case types.MsgEthermint:
			return HandleMsgEthermint(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized ethermint message type: %T", msg)
		}
	}
}

// HandleEthTxMsg handles an Ethereum specific tx
func HandleEthTxMsg(ctx sdk.Context, k Keeper, msg types.MsgEthereumTx) (*sdk.Result, error) {
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
	k.CommitStateDB.Prepare(ethHash, common.Hash{}, k.TxCount)
	k.TxCount++

	// TODO: move to keeper
	returnData, err := st.TransitionCSDB(ctx)
	if err != nil {
		return nil, err
	}

	// update block bloom filter
	k.Bloom.Or(k.Bloom, returnData.Bloom)

	// update transaction logs in KVStore
	err = k.SetTransactionLogs(ctx, returnData.Logs, txHash)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEthereumTx,
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Data.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
		),
	})

	if msg.Data.Recipient != nil {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeEthereumTx,
				sdk.NewAttribute(types.AttributeKeyRecipient, msg.Data.Recipient.String()),
			),
		)
	}

	// set the events to the result
	returnData.Result.Events = ctx.EventManager().Events()
	return returnData.Result, nil
}

func HandleMsgEthermint(ctx sdk.Context, k Keeper, msg types.MsgEthermint) (*sdk.Result, error) {
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
	k.CommitStateDB.Prepare(common.Hash{}, common.Hash{}, k.TxCount) // Cannot provide tx hash
	k.TxCount++

	returnData, err := st.TransitionCSDB(ctx)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEthermint,
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})

	if msg.Recipient != nil {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeEthermint,
				sdk.NewAttribute(types.AttributeKeyRecipient, msg.Recipient.String()),
			),
		)
	}

	// set the events to the result
	returnData.Result.Events = ctx.EventManager().Events()
	return returnData.Result, nil
}
