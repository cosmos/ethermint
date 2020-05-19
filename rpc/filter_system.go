package rpc

import (
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"

	abci "github.com/tendermint/tendermint/abci/types"
)

// EventSystem creates subscriptions, processes events and broadcasts them to the
// subscription which match the subscription criteria.
type EventSystem interface {
	SubscribeLogs(crit ethereum.FilterQuery, logs chan []*types.Log) (*filters.Subscription, error)
	SubscribeNewHeads(headers chan abci.Header) *filters.Subscription
	SubscribePendingTxs(hashes chan []common.Hash) *filters.Subscription
}

// TODO: create concrete type
