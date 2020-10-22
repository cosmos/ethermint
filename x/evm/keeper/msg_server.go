package keeper

import (
	"context"
	"math/big"

	"github.com/armon/go-metrics"
	tmtypes "github.com/tendermint/tendermint/types"

	ethcmn "github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ethermint "github.com/cosmos/ethermint/types"
	"github.com/cosmos/ethermint/x/evm/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the evm MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) EthereumTx(goCtx context.Context, msg *types.MsgEthereumTx) (*types.MsgEthereumTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// parse the chainID from a string to a base-10 integer
	chainIDEpoch, err := ethermint.ParseChainID(ctx.ChainID())
	if err != nil {
		return nil, err
	}

	// Verify signature and retrieve sender address
	sender, err := msg.VerifySig(chainIDEpoch)
	if err != nil {
		return nil, err
	}

	txHash := tmtypes.Tx(ctx.TxBytes()).Hash()
	ethHash := ethcmn.BytesToHash(txHash)

	var recipient *ethcmn.Address

	labels := []metrics.Label{telemetry.NewLabel("operation", "create")}

	if msg.Data.Recipient != "" {
		addr := ethcmn.HexToAddress(msg.Data.Recipient)
		recipient = &addr
		labels = []metrics.Label{telemetry.NewLabel("operation", "call")}
	}

	st := types.StateTransition{
		AccountNonce: msg.Data.AccountNonce,
		Price:        new(big.Int).SetBytes(msg.Data.Price),
		GasLimit:     msg.Data.GasLimit,
		Recipient:    recipient,
		Amount:       new(big.Int).SetBytes(msg.Data.Amount),
		Payload:      msg.Data.Payload,
		Csdb:         k.CommitStateDB.WithContext(ctx),
		ChainID:      chainIDEpoch,
		TxHash:       &ethHash,
		Sender:       sender,
		Simulate:     ctx.IsCheckTx(),
	}

	// since the txCount is used by the stateDB, and a simulated tx is run only on the node it's submitted to,
	// then this will cause the txCount/stateDB of the node that ran the simulated tx to be different than the
	// other nodes, causing a consensus error
	if !st.Simulate {
		// Prepare db for logs
		k.Prepare(ctx, ethHash, ethcmn.Hash{}, k.TxCount)
		k.TxCount++
	}

	config, found := k.GetChainConfig(ctx)
	if !found {
		return nil, types.ErrChainConfigNotFound
	}

	executionResult, err := st.TransitionDb(ctx, config)
	if err != nil {
		return nil, err
	}

	if !st.Simulate {
		// update block bloom filter
		k.Bloom.Or(k.Bloom, executionResult.Bloom)

		// update transaction logs in KVStore
		err = k.SetLogs(ctx, ethcmn.BytesToHash(txHash), executionResult.Logs)
		if err != nil {
			panic(err)
		}
	}

	// add metrics for the transaction
	defer func() {
		if st.Amount.IsInt64() {
			telemetry.SetGaugeWithLabels(
				[]string{"tx", "msg", "ethereum"},
				float32(st.Amount.Int64()),
				labels,
			)
		}
	}()

	// emit events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEthereumTx,
			sdk.NewAttribute(sdk.AttributeKeyAmount, st.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
		),
	})

	if msg.Data.Recipient != "" {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeEthereumTx,
				sdk.NewAttribute(types.AttributeKeyRecipient, msg.Data.Recipient),
			),
		)
	}

	return executionResult.Response, nil
}
