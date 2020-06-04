package rpc

import (
	"context"
	"fmt"
	"sync"
	"time"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"

	clientcontext "github.com/cosmos/cosmos-sdk/client/context"

	evmtypes "github.com/cosmos/ethermint/x/evm/types"
)

// consider a filter inactive if it has not been polled for within deadline
var deadline = 5 * time.Minute

// filter is a helper struct that holds meta information over the filter type
// and associated subscription in the event system.
type filter struct {
	typ      filters.Type
	deadline *time.Timer // filter is inactiv when deadline triggers
	hashes   []common.Hash
	crit     filters.FilterCriteria
	logs     []*types.Log
	s        *Subscription // associated subscription in event system
}

// PublicFilterAPI offers support to create and manage filters. This will allow external clients to retrieve various
// information related to the Ethereum protocol such as blocks, transactions and logs.
type PublicFilterAPI struct {
	cliCtx    clientcontext.CLIContext
	backend   Backend
	mux       *event.TypeMux
	quit      chan struct{}
	events    *EventSystem
	filtersMu sync.Mutex
	filters   map[rpc.ID]*filter
}

// NewPublicFilterAPI returns a new PublicFilterAPI instance.
func NewPublicFilterAPI(cliCtx clientcontext.CLIContext, backend Backend) *PublicFilterAPI {
	api := &PublicFilterAPI{
		cliCtx:  cliCtx,
		backend: backend,
		filters: make(map[rpc.ID]*filter),
		events:  NewEventSystem(cliCtx.Client),
	}

	go api.timeoutLoop()

	return api
}

// timeoutLoop runs every 5 minutes and deletes filters that have not been recently used.
// Tt is started when the api is created.
func (api *PublicFilterAPI) timeoutLoop() {
	ticker := time.NewTicker(deadline)
	defer ticker.Stop()
	for {
		<-ticker.C
		api.filtersMu.Lock()
		for id, f := range api.filters {
			select {
			case <-f.deadline.C:
				f.s.Unsubscribe(api.events)
				delete(api.filters, id)
			default:
				continue
			}
		}
		api.filtersMu.Unlock()
	}
}

// // NewPendingTransactionFilter creates a filter that fetches pending transaction hashes
// // as transactions enter the pending state.
// //
// // It is part of the filter package because this filter can be used through the
// // `eth_getFilterChanges` polling method that is also used for log filters.
// //
// // https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_newpendingtransactionfilter
// func (api *PublicFilterAPI) NewPendingTransactionFilter() rpc.ID {
// 	var (
// 		pendingTxs   = make(chan []common.Hash)
// 		pendingTxSub = api.events.SubscribePendingTxs(pendingTxs)
// 	)

// 	ctx, cancelFn := context.WithTimeout(context.Background(), api.events.GetTimeout())
// 	defer cancelFn()

// 	api.events.WithContext(ctx)

// 	api.filtersMu.Lock()
// 	api.filters[pendingTxSub.ID] = NewFilter(api.backend, &filters.FilterCriteria{}, filters.PendingTransactionsSubscription)
// 	api.filtersMu.Unlock()

// 	go func() {
// 		for {
// 			select {
// 			case ph := <-pendingTxs:
// 				api.filtersMu.Lock()
// 				if f, found := api.filters[pendingTxSub.ID]; found {
// 					f.hashes = append(f.hashes, ph...)
// 				}
// 				api.filtersMu.Unlock()
// 			case <-pendingTxSub.Err():
// 				api.filtersMu.Lock()
// 				delete(api.filters, pendingTxSub.ID)
// 				api.filtersMu.Unlock()
// 				return
// 			}
// 		}
// 	}()

// 	return pendingTxSub.ID
// }

// // NewPendingTransactions creates a subscription that is triggered each time a transaction
// // enters the transaction pool and was signed from one of the transactions this nodes manages.
// func (api *PublicFilterAPI) NewPendingTransactions(ctx context.Context) (*rpc.Subscription, error) {
// 	notifier, supported := rpc.NotifierFromContext(ctx)
// 	if !supported {
// 		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
// 	}

// 	ctx, cancelFn := context.WithTimeout(context.Background(), api.events.GetTimeout())
// 	defer cancelFn()

// 	api.events.WithContext(ctx)
// 	rpcSub := notifier.CreateSubscription()

// 	go func() {
// 		txHashes := make(chan []common.Hash, 128)
// 		pendingTxSub := api.events.SubscribePendingTxs(txHashes)

// 		for {
// 			select {
// 			case hashes := <-txHashes:
// 				// To keep the original behaviour, send a single tx hash in one notification.
// 				// TODO(rjl493456442) Send a batch of tx hashes in one notification
// 				for _, h := range hashes {
// 					notifier.Notify(rpcSub.ID, h)
// 				}
// 			case <-rpcSub.Err():
// 				pendingTxSub.Unsubscribe()
// 				return
// 			case <-notifier.Closed():
// 				pendingTxSub.Unsubscribe()
// 				return
// 			}
// 		}
// 	}()

// 	return rpcSub, nil
// }

// NewBlockFilter creates a filter that fetches blocks that are imported into the chain.
// It is part of the filter package since polling goes with eth_getFilterChanges.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_newblockfilter
func (api *PublicFilterAPI) NewBlockFilter() rpc.ID {
	headerSub, err := api.events.SubscribeNewHeads()
	if err != nil {
		// return an empty id
		return rpc.ID("")
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), deadline)
	defer cancelFn()

	api.events.WithContext(ctx)

	api.filtersMu.Lock()
	api.filters[headerSub.ID()] = &filter{typ: filters.BlocksSubscription, deadline: time.NewTimer(deadline), hashes: make([]common.Hash, 0), s: headerSub}
	api.filtersMu.Unlock()

	go func() {
		for {
			select {
			case event := <-headerSub.eventChannel:
				evHeader, ok := event.Data.(tmtypes.EventDataNewBlockHeader)
				if !ok {
					return
				}
				api.filtersMu.Lock()
				if f, found := api.filters[headerSub.ID()]; found {
					f.hashes = append(f.hashes, common.BytesToHash(evHeader.Header.Hash()))
				}
				api.filtersMu.Unlock()
			}
		}
	}()

	return headerSub.ID()
}

// NewHeads send a notification each time a new (header) block is appended to the chain.
func (api *PublicFilterAPI) NewHeads(ctx context.Context) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}

	ctx, cancelFn := context.WithTimeout(ctx, deadline)
	defer cancelFn()

	api.events.WithContext(ctx)
	rpcSub := notifier.CreateSubscription()

	headersSub, err := api.events.SubscribeNewHeads()
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case event := <-headersSub.eventChannel:
				evHeader, ok := event.Data.(tmtypes.EventDataNewBlockHeader)
				if !ok {
					err = fmt.Errorf("invalid event data %T, expected %s", event.Data, tmtypes.EventNewBlockHeader)
					return
				}
				notifier.Notify(rpcSub.ID, evHeader.Header)
			case <-rpcSub.Err():
				err = headersSub.Unsubscribe(api.events)
				return
			case <-notifier.Closed():
				err = headersSub.Unsubscribe(api.events)
				return
			}
		}
	}()

	return rpcSub, err
}

// Logs creates a subscription that fires for all new log that match the given filter criteria.
func (api *PublicFilterAPI) Logs(ctx context.Context, crit filters.FilterCriteria) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), deadline)
	defer cancelFn()

	api.events.WithContext(ctx)

	var (
		rpcSub = notifier.CreateSubscription()
	)

	logsSub, err := api.events.SubscribeLogs(crit)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case event := <-logsSub.eventChannel:
				// filter only events from EVM module txs
				_, isMsgEthermint := event.Events[evmtypes.TypeMsgEthermint]
				_, isMsgEthereumTx := event.Events[evmtypes.TypeMsgEthereumTx]

				if !(isMsgEthermint || isMsgEthereumTx) {
					// ignore transaction
					return
				}

				// get the
				dataTx, ok := event.Data.(tmtypes.EventDataTx)
				if !ok {
					err = fmt.Errorf("invalid event data %T, expected %s", event.Data, tmtypes.EventTx)
					return
				}

				resultData, err := evmtypes.DecodeResultData(dataTx.TxResult.Result.Data)
				if err != nil {
					return
				}

				logs := filterLogs(resultData.Logs, crit.FromBlock, crit.ToBlock, crit.Addresses, crit.Topics)

				for _, log := range logs {
					notifier.Notify(rpcSub.ID, &log)
				}
			case <-rpcSub.Err(): // client send an unsubscribe request
				err = logsSub.Unsubscribe(api.events)
				return
			case <-notifier.Closed(): // connection dropped
				err = logsSub.Unsubscribe(api.events)
				return
			}
		}
	}()

	return rpcSub, err
}

// NewFilter creates a new filter and returns the filter id. It can be
// used to retrieve logs when the state changes. This method cannot be
// used to fetch logs that are already stored in the state.
//
// Default criteria for the from and to block are "latest".
// Using "latest" as block number will return logs for mined blocks.
// Using "pending" as block number returns logs for not yet mined (pending) blocks.
// In case logs are removed (chain reorg) previously returned logs are returned
// again but with the removed property set to true.
//
// In case "fromBlock" > "toBlock" an error is returned.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_newfilter
func (api *PublicFilterAPI) NewFilter(criteria filters.FilterCriteria) (rpc.ID, error) {
	logsSub, err := api.events.SubscribeLogs(criteria)
	if err != nil {
		return rpc.ID(""), err
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), deadline)
	defer cancelFn()

	api.events.WithContext(ctx)

	api.filtersMu.Lock()
	api.filters[logsSub.ID()] = &filter{typ: filters.LogsSubscription, crit: criteria, deadline: time.NewTimer(deadline), logs: make([]*types.Log, 0), s: logsSub}
	api.filtersMu.Unlock()

	go func() {
		for {
			select {
			case event := <-logsSub.eventChannel:
				// filter only events from EVM module txs
				_, isMsgEthermint := event.Events[evmtypes.TypeMsgEthermint]
				_, isMsgEthereumTx := event.Events[evmtypes.TypeMsgEthereumTx]

				if !(isMsgEthermint || isMsgEthereumTx) {
					// ignore transaction
					return
				}

				dataTx, ok := event.Data.(tmtypes.EventDataTx)
				if !ok {
					err = fmt.Errorf("invalid event data %T, expected %s", event.Data, tmtypes.EventTx)
					return
				}

				resultData, err := evmtypes.DecodeResultData(dataTx.TxResult.Result.Data)
				if err != nil {
					return
				}

				logs := filterLogs(resultData.Logs, criteria.FromBlock, criteria.ToBlock, criteria.Addresses, criteria.Topics)

				api.filtersMu.Lock()
				if f, found := api.filters[logsSub.ID()]; found {
					f.logs = append(f.logs, logs...)
				}
				api.filtersMu.Unlock()
			}
		}
	}()

	return logsSub.ID(), nil
}

// GetLogs returns logs matching the given argument that are stored within the state.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getlogs
func (api *PublicFilterAPI) GetLogs(ctx context.Context, crit filters.FilterCriteria) ([]*ethtypes.Log, error) {
	var filter *Filter
	if crit.BlockHash != nil {
		// Block filter requested, construct a single-shot filter
		filter = NewBlockFilter(api.backend, crit)
	} else {
		// Convert the RPC block numbers into internal representations
		begin := rpc.LatestBlockNumber.Int64()
		if crit.FromBlock != nil {
			begin = crit.FromBlock.Int64()
		}
		end := rpc.LatestBlockNumber.Int64()
		if crit.ToBlock != nil {
			end = crit.ToBlock.Int64()
		}
		// Construct the range filter
		filter = NewRangeFilter(api.backend, begin, end, crit.Addresses, crit.Topics)
	}
	// Run the filter and return all the logs
	logs, err := filter.Logs(ctx)
	if err != nil {
		return nil, err
	}
	return returnLogs(logs), err
}

// UninstallFilter removes the filter with the given filter id.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_uninstallfilter
func (api *PublicFilterAPI) UninstallFilter(id rpc.ID) bool {
	api.filtersMu.Lock()
	f, found := api.filters[id]
	if found {
		delete(api.filters, id)
	}
	api.filtersMu.Unlock()
	if !found {
		return false
	}

	f.s.Unsubscribe(api.events)
	return true
}

// GetFilterLogs returns the logs for the filter with the given id.
// If the filter could not be found an empty array of logs is returned.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getfilterlogs
func (api *PublicFilterAPI) GetFilterLogs(ctx context.Context, id rpc.ID) ([]*ethtypes.Log, error) {
	api.filtersMu.Lock()
	f, found := api.filters[id]
	api.filtersMu.Unlock()

	if !found {
		return returnLogs(nil), fmt.Errorf("filter %s not found", id)
	}

	if f.typ != filters.LogsSubscription {
		return returnLogs(nil), fmt.Errorf("filter %s doesn't have a LogsSubscription type: got %d", id, f.typ)
	}

	var filter *Filter
	if f.crit.BlockHash != nil {
		// Block filter requested, construct a single-shot filter
		filter = NewBlockFilter(api.backend, f.crit)
	} else {
		// Convert the RPC block numbers into internal representations
		begin := rpc.LatestBlockNumber.Int64()
		if f.crit.FromBlock != nil {
			begin = f.crit.FromBlock.Int64()
		}
		end := rpc.LatestBlockNumber.Int64()
		if f.crit.ToBlock != nil {
			end = f.crit.ToBlock.Int64()
		}
		// Construct the range filter
		filter = NewRangeFilter(api.backend, begin, end, f.crit.Addresses, f.crit.Topics)
	}
	// Run the filter and return all the logs
	logs, err := filter.Logs(ctx)
	if err != nil {
		return nil, err
	}
	return returnLogs(logs), nil

	// return api.filters[id].Logs(ctx)
}

// GetFilterChanges returns the logs for the filter with the given id since
// last time it was called. This can be used for polling.
//
// For pending transaction and block filters the result is []common.Hash.
// (pending)Log filters return []Log.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getfilterchanges
func (api *PublicFilterAPI) GetFilterChanges(id rpc.ID) (interface{}, error) {
	api.filtersMu.Lock()
	defer api.filtersMu.Unlock()

	f, found := api.filters[id]
	if !found {
		return nil, fmt.Errorf("filter %s not found", id)
	}

	if !f.deadline.Stop() {
		// timer expired but filter is not yet removed in timeout loop
		// receive timer value and reset timer
		<-f.deadline.C
	}
	f.deadline.Reset(deadline)

	switch f.typ {
	case filters.PendingTransactionsSubscription, filters.BlocksSubscription:
		hashes := f.hashes
		f.hashes = nil
		return returnHashes(hashes), nil
	case filters.LogsSubscription, filters.MinedAndPendingLogsSubscription:
		logs := f.logs
		f.logs = nil
		return returnLogs(logs), nil
	default:
		return nil, fmt.Errorf("invalid filter %s type %d", id, f.typ)
	}
}

// returnHashes is a helper that will return an empty hash array case the given hash array is nil,
// otherwise the given hashes array is returned.
func returnHashes(hashes []common.Hash) []common.Hash {
	if hashes == nil {
		return []common.Hash{}
	}
	return hashes
}

// returnLogs is a helper that will return an empty log array in case the given logs array is nil,
// otherwise the given logs array is returned.
func returnLogs(logs []*types.Log) []*types.Log {
	if logs == nil {
		return []*types.Log{}
	}
	return logs
}
