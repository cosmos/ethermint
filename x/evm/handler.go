package evm

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	emint "github.com/cosmos/ethermint/types"
	"github.com/cosmos/ethermint/x/evm/types"

	tmtypes "github.com/tendermint/tendermint/types"
)

// NewHandler returns a handler for Ethermint type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgEthereumTx:
			return handleMsgEthereumTx(ctx, k, msg)
		case types.MsgEthermint:
			return handleMsgEthermint(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

// handleMsgEthereumTx handles an Ethereum specific tx
func handleMsgEthereumTx(ctx sdk.Context, k Keeper, msg types.MsgEthereumTx) (*sdk.Result, error) {
	fmt.Println("handleMsgEthereumTx")

	// parse the chainID from a string to a base-10 integer
	intChainID, ok := new(big.Int).SetString(ctx.ChainID(), 10)
	if !ok {
		fmt.Println("invalid chain ID")
		return nil, sdkerrors.Wrap(emint.ErrInvalidChainID, ctx.ChainID())
	}

	// Verify signature and retrieve sender address
	sender, err := msg.VerifySig(intChainID)
	if err != nil {
		fmt.Println("invalid sig")
		return nil, err
	}

	fmt.Println("from balance before", k.GetBalance(ctx, sender))
	fmt.Println("to balance before", k.GetBalance(ctx, *msg.Data.Recipient))

	txHash := tmtypes.Tx(ctx.TxBytes()).Hash()
	ethHash := common.BytesToHash(txHash)

	st := types.StateTransition{
		AccountNonce: msg.Data.AccountNonce,
		Price:        msg.Data.Price,
		GasLimit:     msg.Data.GasLimit,
		Recipient:    msg.Data.Recipient,
		Amount:       msg.Data.Amount,
		Payload:      msg.Data.Payload,
		Csdb:         k.CommitStateDB.WithContext(ctx),
		ChainID:      intChainID,
		TxHash:       &ethHash,
		Sender:       sender,
		Simulate:     ctx.IsCheckTx(),
	}

	// Prepare db for logs
	// TODO: block hash
	k.CommitStateDB.Prepare(ethHash, common.Hash{}, k.TxCount)
	k.TxCount++

	// TODO: move to keeper
	executionResult, err := st.TransitionDb(ctx)
	if err != nil {
		return nil, err
	}

	// update block bloom filter
	k.Bloom.Or(k.Bloom, executionResult.Bloom)

	// update transaction logs in KVStore
	err = k.SetLogs(ctx, common.BytesToHash(txHash), executionResult.Logs)
	if err != nil {
		panic(err)
	}

	// log successful execution
	k.Logger(ctx).Info(executionResult.Result.Log)

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

	fmt.Println(executionResult.Result)

	fmt.Println("from balance after", k.GetBalance(ctx, sender))
	fmt.Println("to balance after", k.GetBalance(ctx, *msg.Data.Recipient))

	// set the events to the result
	executionResult.Result.Events = ctx.EventManager().Events().ToABCIEvents()
	return executionResult.Result, nil
}

// handleMsgEthermint handles an sdk.StdTx for an Ethereum state transition
func handleMsgEthermint(ctx sdk.Context, k Keeper, msg types.MsgEthermint) (*sdk.Result, error) {
	fmt.Println("handleMsgEthermint")

	// parse the chainID from a string to a base-10 integer
	intChainID, ok := new(big.Int).SetString(ctx.ChainID(), 10)
	if !ok {
		return nil, sdkerrors.Wrap(emint.ErrInvalidChainID, ctx.ChainID())
	}

	txHash := tmtypes.Tx(ctx.TxBytes()).Hash()
	ethHash := common.BytesToHash(txHash)

	st := types.StateTransition{
		AccountNonce: msg.AccountNonce,
		Price:        msg.Price.BigInt(),
		GasLimit:     msg.GasLimit,
		Amount:       msg.Amount.BigInt(),
		Payload:      msg.Payload,
		Csdb:         k.CommitStateDB.WithContext(ctx),
		ChainID:      intChainID,
		TxHash:       &ethHash,
		Sender:       common.BytesToAddress(msg.From.Bytes()),
		Simulate:     ctx.IsCheckTx(),
	}

	if msg.Recipient != nil {
		to := common.BytesToAddress(msg.Recipient.Bytes())
		st.Recipient = &to
	}

	fmt.Println("from balance before", k.GetBalance(ctx, st.Sender))
	fmt.Println("to balance before", k.GetBalance(ctx, *st.Recipient))

	// Prepare db for logs
	k.CommitStateDB.Prepare(ethHash, common.Hash{}, k.TxCount)
	k.TxCount++

	executionResult, err := st.TransitionDb(ctx)
	if err != nil {
		return nil, err
	}

	// update block bloom filter
	k.Bloom.Or(k.Bloom, executionResult.Bloom)

	// update transaction logs in KVStore
	err = k.SetLogs(ctx, common.BytesToHash(txHash), executionResult.Logs)
	if err != nil {
		panic(err)
	}

	// log successful execution
	k.Logger(ctx).Info(executionResult.Result.Log)

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

	fmt.Println(executionResult.Result)
	fmt.Println("from balance after", k.GetBalance(ctx, st.Sender))
	fmt.Println("to balance after", k.GetBalance(ctx, *st.Recipient))

	// set the events to the result
	executionResult.Result.Events = ctx.EventManager().Events().ToABCIEvents()
	return executionResult.Result, nil
}
