package rpc

import (
	"context"
	"fmt"
	"time"

	evmtypes "github.com/cosmos/ethermint/x/evm/types"
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
	SubscribeLogs(crit filters.FilterCriteria) (*Subscription, error)
	SubscribeNewHeads() (*Subscription, error)
	// SubscribePendingTxs(hashes chan []common.Hash) *filters.Subscription
}

type subscription struct {
	id        rpc.ID
	typ       filters.Type
	created   time.Time
	logsCrit  filters.FilterCriteria
	installed chan struct{} // closed when the filter is installed
	err       chan error    // closed when the filter is uninstalled
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

var _ EventSystem = &TendermintEvents{}

// TendermintEvents implements the EventSystem using Tendermint's RPC client.
type TendermintEvents struct {
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

// NewTendermintEvents creates a new manager that listens for event on the given mux,
// parses and filters them. It uses the all map to retrieve filter changes. The
// work loop holds its own index that is used to forward events to filters.
//
// The returned manager has a loop that needs to be stopped with the Stop function
// or by stopping the given mux.
func NewTendermintEvents(client rpcclient.Client) *TendermintEvents {
	te := &TendermintEvents{
		ctx:    context.Background(),
		client: client,
	}

	// Subscribe events
	// te.txsSub = te.SubscribeNewTxsEvent(m.txsCh)
	te.logsSub, _ = te.SubscribeLogs(filters.FilterCriteria{})
	te.chainSub, _ = te.SubscribeNewHeads()

	go te.eventLoop()
	return te
}

// WithContext sets the a given context to the
func (te *TendermintEvents) WithContext(ctx context.Context) EventSystem {
	te.ctx = ctx
	return te
}

func (te TendermintEvents) subscribe(sub *subscription) (*Subscription, error) {
	var err error
	subscription := &Subscription{
		subscription: sub,
	}
	subscription.eventChannel, err = te.client.Subscribe(te.ctx, string(sub.id), tmtypes.QueryForEvent(sub.event()).String())
	return subscription, err
}

// SubscribeLogs creates a subscription that will write all logs matching the
// given criteria to the given logs channel. Default value for the from and to
// block is "latest". If the fromBlock > toBlock an error is returned.
func (te *TendermintEvents) SubscribeLogs(crit filters.FilterCriteria) (*Subscription, error) {
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
	// 	return te.subscribePendingLogs(crit, logs)

	// only interested in new mined logs, mined logs within a specific block range, or
	// logs from a specific block number to new mined blocks
	case (from == rpc.LatestBlockNumber && to == rpc.LatestBlockNumber),
		(from >= 0 && to >= 0 && to >= from):
		return te.subscribeLogs(crit)

	// // interested in mined logs from a specific block number, new logs and pending logs
	// case from >= rpc.LatestBlockNumber && to == rpc.PendingBlockNumber:
	// 	return te.subscribeMinedPendingLogs(crit, logs)

	default:
		return nil, fmt.Errorf("invalid from and to block combination: from > to (%d > %d)", from, to)
	}
}

// subscribeMinedPendingLogs creates a subscription that returned mined and
// pending logs that match the given criteria.
func (te *TendermintEvents) subscribeMinedPendingLogs(crit filters.FilterCriteria) (*Subscription, error) {
	sub := &subscription{
		id:        rpc.NewID(),
		typ:       filters.MinedAndPendingLogsSubscription,
		logsCrit:  crit,
		created:   time.Now(),
		installed: make(chan struct{}),
		err:       make(chan error),
	}
	return te.subscribe(sub)
}

// subscribeLogs creates a subscription that will write all logs matching the
// given criteria to the given logs channel.
func (te *TendermintEvents) subscribeLogs(crit filters.FilterCriteria) (*Subscription, error) {
	sub := &subscription{
		id:        rpc.NewID(),
		typ:       filters.LogsSubscription,
		logsCrit:  crit,
		created:   time.Now(),
		installed: make(chan struct{}),
		err:       make(chan error),
	}
	return te.subscribe(sub)
}

// subscribePendingLogs creates a subscription that writes transaction hashes for
// transactions that enter the transaction pool.
func (te *TendermintEvents) subscribePendingLogs(crit filters.FilterCriteria) (*Subscription, error) {
	sub := &subscription{
		id:        rpc.NewID(),
		typ:       filters.PendingLogsSubscription,
		logsCrit:  crit,
		created:   time.Now(),
		installed: make(chan struct{}),
		err:       make(chan error),
	}
	return te.subscribe(sub)
}

func (te TendermintEvents) SubscribeNewHeads() (*Subscription, error) {
	sub := &subscription{
		id:        rpc.NewID(),
		typ:       filters.BlocksSubscription,
		created:   time.Now(),
		installed: make(chan struct{}),
		err:       make(chan error),
	}
	return te.subscribe(sub)
}

func (te TendermintEvents) SubscribePendingTxs(hashes chan []common.Hash) *filters.Subscription {
	return &filters.Subscription{}
}

type filterIndex map[filters.Type]map[rpc.ID]*subscription

func (te *TendermintEvents) handleLogs(filterIdx filterIndex, ev coretypes.ResultEvent) {
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
			f.logs = matchedLogs
		}
	}
}

func (te *TendermintEvents) handleChainEvent(filterIdx filterIndex, ev coretypes.ResultEvent) {
	data, _ := ev.Data.(tmtypes.EventDataNewBlockHeader)
	for _, f := range filterIdx[filters.BlocksSubscription] {
		f.headers = data.Header
	}
	// TODO: light client
}

// eventLoop (un)installs filters and processes mux events.
func (te *TendermintEvents) eventLoop() {
	// Ensure all subscriptions get cleaned up
	defer func() {
		// te.txsSub.Unsubscribe(te.ctx, te.client)
		te.logsSub.Unsubscribe(te.ctx, te.client)
		// te.rmLogsSub.Unsubscribe(te.ctx, te.client)
		// te.pendingLogsSub.Unsubscribe(te.ctx, te.client)
		te.chainSub.Unsubscribe(te.ctx, te.client)
	}()

	index := make(filterIndex)
	for i := filters.UnknownSubscription; i < filters.LastIndexSubscription; i++ {
		index[i] = make(map[rpc.ID]*subscription)
	}

	for {
		select {
		case ev := <-te.eventsChannel:
			switch ev.Data.(type) {
			case tmtypes.EventDataTx:
				te.handleLogs(index, ev)
			case tmtypes.EventDataNewBlockHeader:
				te.handleChainEvent(index, ev)
			}

		case f := <-te.install:
			if f.typ == filters.MinedAndPendingLogsSubscription {
				// the type are logs and pending logs subscriptions
				index[filters.LogsSubscription][f.id] = f
				index[filters.PendingLogsSubscription][f.id] = f
			} else {
				index[f.typ][f.id] = f
			}
			close(f.installed)

		case f := <-te.uninstall:
			if f.typ == filters.MinedAndPendingLogsSubscription {
				// the type are logs and pending logs subscriptions
				delete(index[filters.LogsSubscription], f.id)
				delete(index[filters.PendingLogsSubscription], f.id)
			} else {
				delete(index[f.typ], f.id)
			}
			close(f.err)
		}
	}
}

type Subscription struct {
	subscription *subscription
	eventChannel <-chan coretypes.ResultEvent
}

func (s Subscription) ID() rpc.ID {
	return s.subscription.id
}

func (s Subscription) Unsubscribe(ctx context.Context, client rpcclient.Client) error {
	return client.Unsubscribe(
		ctx, string(s.ID()),
		tmtypes.QueryForEvent(s.subscription.event()).String(),
	)
}
