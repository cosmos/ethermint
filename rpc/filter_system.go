package rpc

import (
	"context"
	"fmt"
	"sync"
	"time"

	evmtypes "github.com/cosmos/ethermint/x/evm/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rpc"
)

type subscription struct {
	id       rpc.ID
	typ      filters.Type
	created  time.Time
	logsCrit filters.FilterCriteria
	logs     chan []*ethtypes.Log
	// hashes    chan []common.Hash
	headers   chan *ethtypes.Header
	installed chan struct{} // closed when the filter is installed
}

func (s subscription) event() string {
	switch s.typ {
	case filters.LogsSubscription, filters.PendingTransactionsSubscription, filters.MinedAndPendingLogsSubscription, filters.PendingLogsSubscription:
		return tmtypes.EventTx
	case filters.BlocksSubscription:
		return tmtypes.EventNewBlockHeader
	default:
		return ""
	}
}

// EventSystem creates subscriptions, processes events and broadcasts them to the
// subscription which match the subscription criteria using the Tendermint's RPC client.
type EventSystem struct {
	ctx    context.Context
	client rpcclient.Client

	// Subscriptions
	// txsSub         *Subscription // Subscription for new transaction event
	logsSub *Subscription // Subscription for new log event
	// rmLogsSub      *Subscription // Subscription for removed log event
	// pendingLogsSub *Subscription // Subscription for pending log event
	chainSub *Subscription // Subscription for new chain event

	// Channels
	install       chan *subscription           // install filter for event notification
	uninstall     chan *subscription           // remove filter for event notification
	eventsChannel <-chan coretypes.ResultEvent // channel to receive tendermint event results
}

// NewEventSystem creates a new manager that listens for event on the given mux,
// parses and filters them. It uses the all map to retrieve filter changes. The
// work loop holds its own index that is used to forward events to filters.
//
// The returned manager has a loop that needs to be stopped with the Stop function
// or by stopping the given mux.
func NewEventSystem(client rpcclient.Client) *EventSystem {
	es := &EventSystem{
		ctx:    context.Background(),
		client: client,
	}

	// Subscribe events
	// es.txsSub = es.SubscribeNewTxsEvent(m.txsCh)
	es.logsSub, _ = es.SubscribeLogs(filters.FilterCriteria{})
	es.chainSub, _ = es.SubscribeNewHeads()

	go es.eventLoop()
	return es
}

// WithContext sets the a given context to the
func (es *EventSystem) WithContext(ctx context.Context) {
	es.ctx = ctx
}

func (es EventSystem) subscribe(sub *subscription) (*Subscription, error) {
	es.install <- sub
	<-sub.installed

	var err error
	subscription := &Subscription{
		subscription: sub,
	}
	subscription.eventChannel, err = es.client.Subscribe(es.ctx, string(sub.id), tmtypes.QueryForEvent(sub.event()).String())
	return subscription, err
}

// SubscribeLogs creates a subscription that will write all logs matching the
// given criteria to the given logs channel. Default value for the from and to
// block is "latest". If the fromBlock > toBlock an error is returned.
func (es *EventSystem) SubscribeLogs(crit filters.FilterCriteria) (*Subscription, error) {
	var from, to rpc.BlockNumber
	if crit.FromBlock == nil {
		from = rpc.LatestBlockNumber
	} else {
		from = rpc.BlockNumber(crit.FromBlock.Int64())
	}
	if crit.ToBlock == nil {
		to = rpc.LatestBlockNumber
	} else {
		to = rpc.BlockNumber(crit.ToBlock.Int64())
	}

	switch {
	// // only interested in pending logs
	// case from == rpc.PendingBlockNumber && to == rpc.PendingBlockNumber:
	// 	return es.subscribePendingLogs(crit, logs)

	// only interested in new mined logs, mined logs within a specific block range, or
	// logs from a specific block number to new mined blocks
	case (from == rpc.LatestBlockNumber && to == rpc.LatestBlockNumber),
		(from >= 0 && to >= 0 && to >= from):
		return es.subscribeLogs(crit)

	// // interested in mined logs from a specific block number, new logs and pending logs
	// case from >= rpc.LatestBlockNumber && to == rpc.PendingBlockNumber:
	// 	return es.subscribeMinedPendingLogs(crit, logs)

	default:
		return nil, fmt.Errorf("invalid from and to block combination: from > to (%d > %d)", from, to)
	}
}

// subscribeMinedPendingLogs creates a subscription that returned mined and
// pending logs that match the given criteria.
func (es *EventSystem) subscribeMinedPendingLogs(crit filters.FilterCriteria) (*Subscription, error) {
	sub := &subscription{
		id:        rpc.NewID(),
		typ:       filters.MinedAndPendingLogsSubscription,
		logsCrit:  crit,
		created:   time.Now(),
		logs:      make(chan []*ethtypes.Log),
		installed: make(chan struct{}),
	}
	return es.subscribe(sub)
}

// subscribeLogs creates a subscription that will write all logs matching the
// given criteria to the given logs channel.
func (es *EventSystem) subscribeLogs(crit filters.FilterCriteria) (*Subscription, error) {
	sub := &subscription{
		id:        rpc.NewID(),
		typ:       filters.LogsSubscription,
		logsCrit:  crit,
		created:   time.Now(),
		logs:      make(chan []*ethtypes.Log),
		installed: make(chan struct{}),
	}
	return es.subscribe(sub)
}

// subscribePendingLogs creates a subscription that writes transaction hashes for
// transactions that enter the transaction pool.
func (es *EventSystem) subscribePendingLogs(crit filters.FilterCriteria) (*Subscription, error) {
	sub := &subscription{
		id:        rpc.NewID(),
		typ:       filters.PendingLogsSubscription,
		logsCrit:  crit,
		created:   time.Now(),
		logs:      make(chan []*ethtypes.Log),
		installed: make(chan struct{}),
	}
	return es.subscribe(sub)
}

// SubscribeNewHeads subscribes to new block headers events.
func (es EventSystem) SubscribeNewHeads() (*Subscription, error) {
	sub := &subscription{
		id:        rpc.NewID(),
		typ:       filters.BlocksSubscription,
		created:   time.Now(),
		headers:   make(chan *ethtypes.Header),
		installed: make(chan struct{}),
	}
	return es.subscribe(sub)
}

// SubscribePendingTxs subscribes to new pending transactions events from the mempool.
func (es EventSystem) SubscribePendingTxs(hashes chan []common.Hash) (*Subscription, error) {
	sub := &subscription{
		id:        rpc.NewID(),
		typ:       filters.PendingTransactionsSubscription,
		created:   time.Now(),
		installed: make(chan struct{}),
	}
	return es.subscribe(sub)
}

type filterIndex map[filters.Type]map[rpc.ID]*subscription

func (es *EventSystem) handleLogs(filterIdx filterIndex, ev coretypes.ResultEvent) {
	// filter only events from EVM module txs
	_, isMsgEthermint := ev.Events[evmtypes.TypeMsgEthermint]
	_, isMsgEthereumTx := ev.Events[evmtypes.TypeMsgEthereumTx]

	if !(isMsgEthermint || isMsgEthereumTx) {
		// ignore transaction
		return
	}

	data, _ := ev.Data.(tmtypes.EventDataTx)
	resultData, err := evmtypes.DecodeResultData(data.TxResult.Result.Data)
	if err != nil {
		return
	}

	if len(resultData.Logs) == 0 {
		return
	}
	for _, f := range filterIdx[filters.LogsSubscription] {
		matchedLogs := filterLogs(resultData.Logs, f.logsCrit.FromBlock, f.logsCrit.ToBlock, f.logsCrit.Addresses, f.logsCrit.Topics)
		if len(matchedLogs) > 0 {
			f.logs <- matchedLogs
		}
	}
}

func (es *EventSystem) handleChainEvent(filterIdx filterIndex, ev coretypes.ResultEvent) {
	data, _ := ev.Data.(tmtypes.EventDataNewBlockHeader)
	for _, f := range filterIdx[filters.BlocksSubscription] {
		f.headers <- EthHeaderFromTendermint(data.Header)
	}
	// TODO: light client
}

// eventLoop (un)installs filters and processes mux events.
func (es *EventSystem) eventLoop() {
	// Ensure all subscriptions get cleaned up
	defer func() {
		// es.txsSub.Unsubscribe(es)
		es.logsSub.Unsubscribe(es)
		// es.rmLogsSub.Unsubscribe(es)
		// es.pendingLogsSub.Unsubscribe(es)
		es.chainSub.Unsubscribe(es)
	}()

	index := make(filterIndex)
	for i := filters.UnknownSubscription; i < filters.LastIndexSubscription; i++ {
		index[i] = make(map[rpc.ID]*subscription)
	}

	for {
		select {
		case ev := <-es.eventsChannel:
			switch ev.Data.(type) {
			case tmtypes.EventDataTx:
				es.handleLogs(index, ev)
			case tmtypes.EventDataNewBlockHeader:
				es.handleChainEvent(index, ev)
			}

		case f := <-es.install:
			if f.typ == filters.MinedAndPendingLogsSubscription {
				// the type are logs and pending logs subscriptions
				index[filters.LogsSubscription][f.id] = f
				index[filters.PendingLogsSubscription][f.id] = f
			} else {
				index[f.typ][f.id] = f
			}
			close(f.installed)

		case f := <-es.uninstall:
			if f.typ == filters.MinedAndPendingLogsSubscription {
				// the type are logs and pending logs subscriptions
				delete(index[filters.LogsSubscription], f.id)
				delete(index[filters.PendingLogsSubscription], f.id)
			} else {
				delete(index[f.typ], f.id)
			}
		}
	}
}

// Subscription defines a wrapper for the private subscription
type Subscription struct {
	subscription *subscription
	eventChannel <-chan coretypes.ResultEvent
	unsubOnce    sync.Once
}

// ID returns the underlying subscription RPC identifier.
func (s Subscription) ID() rpc.ID {
	return s.subscription.id
}

// Unsubscribe the current subscription from Tendermint Websocket.
func (s Subscription) Unsubscribe(es *EventSystem) (err error) {
	s.unsubOnce.Do(func() {
	uninstallLoop:
		for {
			// write uninstall request and consume logs/hashes. This prevents
			// the eventLoop broadcast method to deadlock when writing to the
			// filter event channel while the subscription loop is waiting for
			// this method to return (and thus not reading these events).
			select {
			case es.uninstall <- s.subscription:
				break uninstallLoop
			case <-s.subscription.logs:
			// case <-s.subscription.hashes:
			case <-s.subscription.headers:
			}
		}

		err = es.client.Unsubscribe(
			es.ctx, string(s.ID()),
			tmtypes.QueryForEvent(s.subscription.event()).String(),
		)
	})

	return
}
