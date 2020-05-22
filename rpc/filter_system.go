package rpc

import (
	"context"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rpc"
)

// EventSystem creates subscriptions, processes events and broadcasts them to the
// subscription which match the subscription criteria.
type EventSystem interface {
	WithContext(ctx context.Context) EventSystem
	SubscribeLogs(subscriberID rpc.ID) (eventCh <-chan coretypes.ResultEvent, err error)
	UnsubscribeLogs(subscriberID rpc.ID) (err error)
	SubscribeNewHeads(subscriberID rpc.ID) (eventCh <-chan coretypes.ResultEvent, err error)
	UnsubscribeHeads(subscriberID rpc.ID) (err error)
	SubscribePendingTxs(hashes chan []common.Hash) *filters.Subscription
}

var _ EventSystem = &TendermintEvents{}

// TendermintEvents implements the EventSystem using Tendermint's RPC client.
type TendermintEvents struct {
	ctx    context.Context
	client rpcclient.Client
}

// WithContext sets the a given context to the
func (te *TendermintEvents) WithContext(ctx context.Context) EventSystem {
	te.ctx = ctx
	return te
}

func (te TendermintEvents) SubscribeLogs(subscriberID rpc.ID) (eventCh <-chan coretypes.ResultEvent, err error) {
	// var from, to rpc.BlockNumber
	// if crit.FromBlock == nil {
	// 	from = rpc.LatestBlockNumber
	// } else {
	// 	from = rpc.BlockNumber(crit.FromBlock.Int64())
	// }
	// if crit.ToBlock == nil {
	// 	to = rpc.LatestBlockNumber
	// } else {
	// 	to = rpc.BlockNumber(crit.ToBlock.Int64())
	// }

	// // TODO: filter logs

	// // only interested in pending logs
	// if from == rpc.PendingBlockNumber && to == rpc.PendingBlockNumber {
	// 	return es.subscribePendingLogs(crit, logs), nil
	// }
	// // only interested in new mined logs
	// if from == rpc.LatestBlockNumber && to == rpc.LatestBlockNumber {
	// 	return te.subscribeLogs(subscriberID)
	// }
	// // only interested in mined logs within a specific block range
	// if from >= 0 && to >= 0 && to >= from {
	// 	return te.subscribeLogs(subscriberID)
	// }
	// // interested in mined logs from a specific block number, new logs and pending logs
	// if from >= rpc.LatestBlockNumber && to == rpc.PendingBlockNumber {
	// 	return tc.subscribeMinedPendingLogs(subscriberID), nil
	// }
	// // interested in logs from a specific block number to new mined blocks
	// if from >= 0 && to == rpc.LatestBlockNumber {
	// 	return te.subscribeLogs(subscriberID)
	// }

	return te.subscribeLogs(subscriberID)

	// return nil, fmt.Errorf("invalid from and to block combination: from > to")
}

func (te TendermintEvents) SubscribeNewHeads(subscriberID rpc.ID) (eventCh <-chan coretypes.ResultEvent, err error) {
	return te.client.Subscribe(
		te.ctx, string(subscriberID),
		tmtypes.QueryForEvent(tmtypes.EventNewBlockHeader).String(),
	)
}

func (te TendermintEvents) UnsubscribeHeads(subscriberID rpc.ID) (err error) {
	return te.client.Unsubscribe(
		te.ctx, string(subscriberID),
		tmtypes.QueryForEvent(tmtypes.EventNewBlockHeader).String(),
	)
}

func (te TendermintEvents) SubscribePendingTxs(hashes chan []common.Hash) *filters.Subscription {
	return &filters.Subscription{}
}

func (te TendermintEvents) subscribeLogs(subscriberID rpc.ID) (eventCh <-chan coretypes.ResultEvent, err error) {
	return te.client.Subscribe(
		te.ctx, string(subscriberID),
		tmtypes.QueryForEvent(tmtypes.EventTx).String(),
	)
}

func (te TendermintEvents) UnsubscribeLogs(subscriberID rpc.ID) (err error) {
	return te.client.Unsubscribe(
		te.ctx, string(subscriberID),
		tmtypes.QueryForEvent(tmtypes.EventTx).String(),
	)
}
