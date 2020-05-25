package rpc

import (
	"context"
	"time"

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
	GetTimeout() time.Duration
	SubscribeLogs(subscriberID rpc.ID) (eventCh <-chan coretypes.ResultEvent, err error)
	UnsubscribeLogs(subscriberID rpc.ID) (err error)
	SubscribeNewHeads(subscriberID rpc.ID) (eventCh <-chan coretypes.ResultEvent, err error)
	UnsubscribeHeads(subscriberID rpc.ID) (err error)
	// SubscribePendingTxs(hashes chan []common.Hash) *filters.Subscription
}

var _ EventSystem = &TendermintEvents{}

// TendermintEvents implements the EventSystem using Tendermint's RPC client.
type TendermintEvents struct {
	ctx     context.Context
	client  rpcclient.Client
	timeout time.Duration
}

// WithContext sets the a given context to the
func (te *TendermintEvents) WithContext(ctx context.Context) EventSystem {
	te.ctx = ctx
	return te
}

// GetTimeout returns the default timeout in seconds
func (te TendermintEvents) GetTimeout() time.Duration {
	return te.timeout
}

// SubscribeLogs subscribes to new incoming MsgEthereumTx or MsgEthermint transactions
// TODO:
// - subscribe based on Msg Type
// - subscribe to logs based on filter criteria
func (te TendermintEvents) SubscribeLogs(subscriberID rpc.ID) (eventCh <-chan coretypes.ResultEvent, err error) {
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
