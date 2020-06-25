package rpc

import (
	"context"
	"fmt"
	"log"
	"time"

	rpcclient "github.com/tendermint/tendermint/rpc/client"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rpc"

	sdk "github.com/cosmos/cosmos-sdk/types"

	evmtypes "github.com/cosmos/ethermint/x/evm/types"
)

var (
	txEvents     = fmt.Sprintf("%s = '%s'", tmtypes.EventTypeKey, tmtypes.EventTx)
	evmEvents    = fmt.Sprintf("%s = '%s' AND %s.%s = '%s'", tmtypes.EventTypeKey, tmtypes.EventTx, sdk.EventTypeMessage, sdk.AttributeKeyModule, evmtypes.ModuleName)
	headerEvents = fmt.Sprintf("%s = '%s'", tmtypes.EventTypeKey, tmtypes.EventNewBlockHeader)
)

// EventSystem creates subscriptions, processes events and broadcasts them to the
// subscription which match the subscription criteria using the Tendermint's RPC client.
type EventSystem struct {
	ctx    context.Context
	client rpcclient.Client

	// light client mode
	lightMode bool

	index filterIndex

	// Subscriptions
	txsSub  *Subscription // Subscription for new transaction event
	logsSub *Subscription // Subscription for new log event
	// rmLogsSub      *Subscription // Subscription for removed log event
	pendingLogsSub *Subscription // Subscription for pending log event
	chainSub       *Subscription // Subscription for new chain event

	// Channels
	install       chan *Subscription           // install filter for event notification
	uninstall     chan *Subscription           // remove filter for event notification
	txsCh         <-chan coretypes.ResultEvent // Channel to receive new pending transactions event
	logsCh        <-chan coretypes.ResultEvent // Channel to receive new log event
	pendingLogsCh <-chan coretypes.ResultEvent // Channel to receive new log event
	rmLogsCh      <-chan coretypes.ResultEvent // Channel to receive removed log event
	chainCh       <-chan coretypes.ResultEvent // Channel to receive new chain event
}

// NewEventSystem creates a new manager that listens for event on the given mux,
// parses and filters them. It uses the all map to retrieve filter changes. The
// work loop holds its own index that is used to forward events to filters.
//
// The returned manager has a loop that needs to be stopped with the Stop function
// or by stopping the given mux.
func NewEventSystem(client rpcclient.Client) *EventSystem {
	index := make(filterIndex)
	for i := filters.UnknownSubscription; i < filters.LastIndexSubscription; i++ {
		index[i] = make(map[rpc.ID]*Subscription)
	}

	es := &EventSystem{
		ctx:       context.Background(),
		client:    client,
		lightMode: false,
		index:     index,
	}

	go es.eventLoop()
	return es
}

// WithContext sets the a given context to the
func (es *EventSystem) WithContext(ctx context.Context) {
	es.ctx = ctx
}

func (es *EventSystem) subscribe(sub *Subscription) (*Subscription, context.CancelFunc, error) {
	var (
		err      error
		cancelFn context.CancelFunc
	)

	es.ctx, cancelFn = context.WithTimeout(context.Background(), deadline)

	switch sub.typ {
	case filters.PendingTransactionsSubscription:
		es.txsCh, err = es.client.Subscribe(es.ctx, string(sub.id), sub.event, 1000)
		log.Println("subscribed to pending txs")
	case filters.PendingLogsSubscription, filters.MinedAndPendingLogsSubscription:
		es.pendingLogsCh, err = es.client.Subscribe(es.ctx, string(sub.id), sub.event, 1000)
		log.Println("subscribed to pending logs")
	case filters.LogsSubscription:
		es.logsCh, err = es.client.Subscribe(es.ctx, string(sub.id), sub.event, 1000)
		log.Println("subscribed to logs")
	case filters.BlocksSubscription:
		es.chainCh, err = es.client.Subscribe(es.ctx, string(sub.id), sub.event, 1000)
		log.Println("subscribed to headers")
	default:
		err = fmt.Errorf("invalid filter subscription type %d", sub.typ)
	}

	if err != nil {
		return nil, cancelFn, err
	}

	go func() {
	subscribeLoop:
		for {
			select {
			case es.install <- sub:
				break subscribeLoop
			case <-sub.installed:
				break subscribeLoop
			}
		}
	}()

	return sub, cancelFn, nil
}

// SubscribeLogs creates a subscription that will write all logs matching the
// given criteria to the given logs channel. Default value for the from and to
// block is "latest". If the fromBlock > toBlock an error is returned.
func (es *EventSystem) SubscribeLogs(crit filters.FilterCriteria) (*Subscription, context.CancelFunc, error) {
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
	// only interested in pending logs
	case from == rpc.PendingBlockNumber && to == rpc.PendingBlockNumber:
		return es.subscribePendingLogs(crit)

	// only interested in new mined logs, mined logs within a specific block range, or
	// logs from a specific block number to new mined blocks
	case (from == rpc.LatestBlockNumber && to == rpc.LatestBlockNumber),
		(from >= 0 && to >= 0 && to >= from):
		return es.subscribeLogs(crit)

	// interested in mined logs from a specific block number, new logs and pending logs
	case from >= rpc.LatestBlockNumber && to == rpc.PendingBlockNumber:
		return es.subscribeMinedPendingLogs(crit)

	default:
		return nil, nil, fmt.Errorf("invalid from and to block combination: from > to (%d > %d)", from, to)
	}
}

// subscribeMinedPendingLogs creates a subscription that returned mined and
// pending logs that match the given criteria.
func (es *EventSystem) subscribeMinedPendingLogs(crit filters.FilterCriteria) (*Subscription, context.CancelFunc, error) {
	sub := &Subscription{
		id:        rpc.NewID(),
		typ:       filters.MinedAndPendingLogsSubscription,
		event:     evmEvents,
		logsCrit:  crit,
		created:   time.Now(),
		logs:      make(chan []*ethtypes.Log, 1000),
		installed: make(chan struct{}, 1),
	}
	return es.subscribe(sub)
}

// subscribeLogs creates a subscription that will write all logs matching the
// given criteria to the given logs channel.
func (es *EventSystem) subscribeLogs(crit filters.FilterCriteria) (*Subscription, context.CancelFunc, error) {
	sub := &Subscription{
		id:        rpc.NewID(),
		typ:       filters.LogsSubscription,
		event:     evmEvents,
		logsCrit:  crit,
		created:   time.Now(),
		logs:      make(chan []*ethtypes.Log, 1000),
		installed: make(chan struct{}, 1),
	}
	return es.subscribe(sub)
}

// subscribePendingLogs creates a subscription that writes transaction hashes for
// transactions that enter the transaction pool.
func (es *EventSystem) subscribePendingLogs(crit filters.FilterCriteria) (*Subscription, context.CancelFunc, error) {
	sub := &Subscription{
		id:        rpc.NewID(),
		typ:       filters.PendingLogsSubscription,
		event:     evmEvents,
		logsCrit:  crit,
		created:   time.Now(),
		logs:      make(chan []*ethtypes.Log, 1000),
		installed: make(chan struct{}, 1),
	}
	return es.subscribe(sub)
}

// SubscribeNewHeads subscribes to new block headers events.
func (es EventSystem) SubscribeNewHeads() (*Subscription, context.CancelFunc, error) {
	sub := &Subscription{
		id:        rpc.NewID(),
		typ:       filters.BlocksSubscription,
		event:     tmtypes.EventNewBlockHeader,
		created:   time.Now(),
		headers:   make(chan *ethtypes.Header, 1000),
		installed: make(chan struct{}, 1),
	}
	return es.subscribe(sub)
}

// SubscribePendingTxs subscribes to new pending transactions events from the mempool.
func (es EventSystem) SubscribePendingTxs() (*Subscription, context.CancelFunc, error) {
	sub := &Subscription{
		id:        rpc.NewID(),
		typ:       filters.PendingTransactionsSubscription,
		event:     txEvents,
		created:   time.Now(),
		hashes:    make(chan []common.Hash, 1000),
		installed: make(chan struct{}, 1),
	}
	return es.subscribe(sub)
}

type filterIndex map[filters.Type]map[rpc.ID]*Subscription

func (es *EventSystem) handleLogs(ev coretypes.ResultEvent) {
	data, _ := ev.Data.(tmtypes.EventDataTx)
	resultData, err := evmtypes.DecodeResultData(data.TxResult.Result.Data)
	if err != nil {
		return
	}

	if len(resultData.Logs) == 0 {
		return
	}
	for _, f := range es.index[filters.LogsSubscription] {
		matchedLogs := filterLogs(resultData.Logs, f.logsCrit.FromBlock, f.logsCrit.ToBlock, f.logsCrit.Addresses, f.logsCrit.Topics)
		if len(matchedLogs) > 0 {
			f.logs <- matchedLogs
		}
	}
}

func (es *EventSystem) handleTxsEvent(ev coretypes.ResultEvent) {
	data, _ := ev.Data.(tmtypes.EventDataTx)
	for _, f := range es.index[filters.PendingTransactionsSubscription] {
		f.hashes <- []common.Hash{common.BytesToHash(data.Tx.Hash())}
	}
}

func (es *EventSystem) handleChainEvent(ev coretypes.ResultEvent) {
	data, _ := ev.Data.(tmtypes.EventDataNewBlockHeader)
	for _, f := range es.index[filters.BlocksSubscription] {
		f.headers <- EthHeaderFromTendermint(data.Header)
	}
	// TODO: light client
}

// eventLoop (un)installs filters and processes mux events.
func (es *EventSystem) eventLoop() {
	var (
		err                                                    error
		cancelPendingTxsSubs, cancelLogsSubs, cancelHeaderSubs context.CancelFunc
	)

	// Subscribe events
	es.txsSub, cancelPendingTxsSubs, err = es.SubscribePendingTxs()
	if err != nil {
		panic(fmt.Errorf("failed to subscribe pending txs: %w", err))
	}

	defer cancelPendingTxsSubs()

	es.logsSub, cancelLogsSubs, err = es.SubscribeLogs(filters.FilterCriteria{})
	if err != nil {
		panic(fmt.Errorf("failed to subscribe logs: %w", err))
	}

	defer cancelLogsSubs()

	es.chainSub, cancelHeaderSubs, err = es.SubscribeNewHeads()
	if err != nil {
		panic(fmt.Errorf("failed to subscribe headers: %w", err))
	}

	defer cancelHeaderSubs()

	// Ensure all subscriptions get cleaned up
	defer func() {
		_ = es.txsSub.Unsubscribe(es)
		_ = es.logsSub.Unsubscribe(es)
		// _ = es.rmLogsSub.Unsubscribe(es)
		// _ = es.pendingLogsSub.Unsubscribe(es)
		_ = es.chainSub.Unsubscribe(es)
	}()

	go func() {
		for {
			select {
			case txEvent := <-es.txsCh:
				log.Println("received tx event", txEvent)
				go es.handleTxsEvent(txEvent)
			case headerEv := <-es.chainCh:
				log.Println("received header event", headerEv)
				go es.handleChainEvent(headerEv)
			case logsEv := <-es.logsCh:
				log.Println("received logs event", logsEv)
				go es.handleLogs(logsEv)

			case f := <-es.install:
				if f.typ == filters.MinedAndPendingLogsSubscription {
					// the type are logs and pending logs subscriptions
					es.index[filters.LogsSubscription][f.id] = f
					es.index[filters.PendingLogsSubscription][f.id] = f
				} else {
					es.index[f.typ][f.id] = f
				}
				close(f.installed)
				log.Println("filter installed", f.id)

			case f := <-es.uninstall:
				if f.typ == filters.MinedAndPendingLogsSubscription {
					// the type are logs and pending logs subscriptions
					delete(es.index[filters.LogsSubscription], f.id)
					delete(es.index[filters.PendingLogsSubscription], f.id)
				} else {
					delete(es.index[f.typ], f.id)
				}
				log.Println("filter installed", f.id)
			}
		}
	}()
}

// Subscription defines a wrapper for the private subscription
type Subscription struct {
	id        rpc.ID
	typ       filters.Type
	event     string
	created   time.Time
	logsCrit  filters.FilterCriteria
	logs      chan []*ethtypes.Log
	hashes    chan []common.Hash
	headers   chan *ethtypes.Header
	installed chan struct{} // closed when the filter is installed
}

// ID returns the underlying subscription RPC identifier.
func (s Subscription) ID() rpc.ID {
	return s.id
}

// Unsubscribe the current subscription from Tendermint Websocket.
func (s *Subscription) Unsubscribe(es *EventSystem) error {
	if err := es.client.Unsubscribe(es.ctx, string(s.ID()), s.event); err != nil {
		return err
	}

	go func() {
	uninstallLoop:
		for {
			// write uninstall request and consume logs/hashes. This prevents
			// the eventLoop broadcast method to deadlock when writing to the
			// filter event channel while the subscription loop is waiting for
			// this method to return (and thus not reading these events).
			select {
			case es.uninstall <- s:
				break uninstallLoop
			case <-s.logs:
			case <-s.hashes:
			case <-s.headers:
			}
		}
	}()

	return nil
}
