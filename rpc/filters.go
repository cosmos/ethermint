package rpc

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
)

/*
	- Filter functions derived from go-ethereum
	Used to set the criteria passed in from RPC params
*/

type filterType = string

var blockFilter filterType = "block"
var pendingTxFilter filterType = "pending"
var logFilter filterType = "log"

// Filter can be used to retrieve and filter logs, blocks, or pending transactions.
type Filter struct {
	backend            Backend
	fromBlock, toBlock *big.Int         // start and end block numbers
	addresses          []common.Address // contract addresses to watch
	topics             [][]common.Hash  // log topics to watch for
	blockHash          *common.Hash     // Block hash if filtering a single block

	typ     filterType
	hashes  []common.Hash   // filtered block or transaction hashes
	logs    []*ethtypes.Log // filtered logs
	stopped bool            // set to true once filter in uninstalled
}

// NewFilter returns a new Filter
func NewFilter(backend Backend, criteria *filters.FilterCriteria) *Filter {
	return &Filter{
		backend:   backend,
		fromBlock: criteria.FromBlock,
		toBlock:   criteria.ToBlock,
		addresses: criteria.Addresses,
		topics:    criteria.Topics,
		typ:       logFilter,
		stopped:   false,
	}
}

// NewFilterWithBlockHash returns a new Filter with a blockHash.
func NewFilterWithBlockHash(backend Backend, criteria *filters.FilterCriteria) *Filter {
	return &Filter{
		backend:   backend,
		fromBlock: criteria.FromBlock,
		toBlock:   criteria.ToBlock,
		addresses: criteria.Addresses,
		topics:    criteria.Topics,
		blockHash: criteria.BlockHash,
		typ:       logFilter,
	}
}

// NewBlockFilter creates a new filter that notifies when a block arrives.
func NewBlockFilter(backend Backend) *Filter {
	filter := NewFilter(backend, nil)
	filter.typ = blockFilter

	// TODO: check error
	go filter.pollForBlocks()

	return filter
}

func (f *Filter) pollForBlocks() error {
	prev := hexutil.Uint64(0)

	for {
		if f.stopped {
			return nil
		}

		num, err := f.backend.BlockNumber()
		if err != nil {
			return err
		}

		if num != prev {
			block, err := f.backend.GetBlockByNumber(BlockNumber(num), false)
			if err != nil {
				return err
			}

			hash := common.BytesToHash(block["hash"].([]byte))
			f.hashes = append(f.hashes, hash)

			prev = num
		}

		// TODO: should we add a delay?
	}
}

// NewPendingTransactionFilter creates a new filter that notifies when a pending transaction arrives.
func NewPendingTransactionFilter(backend Backend) *Filter {
	// TODO: finish
	filter := NewFilter(backend, nil)
	filter.typ = pendingTxFilter
	return filter
}

func (f *Filter) uninstallFilter() {
	f.stopped = true
}

func (f *Filter) getFilterChanges() interface{} {
	switch f.typ {
	case blockFilter:
		blocks := f.hashes
		f.hashes = []common.Hash{}
		return blocks
	case pendingTxFilter:
		// TODO
	case logFilter:
		// TODO
	}

	return nil
}

func (f *Filter) getFilterLogs() []*ethtypes.Log {
	// TODO
	return nil
}

func includes(addresses []common.Address, a common.Address) bool {
	for _, addr := range addresses {
		if addr == a {
			return true
		}
	}

	return false
}

// filterLogs creates a slice of logs matching the given criteria.
func filterLogs(logs []*ethtypes.Log, fromBlock, toBlock *big.Int, addresses []common.Address, topics [][]common.Hash) []*ethtypes.Log {
	var ret []*ethtypes.Log
Logs:
	for _, log := range logs {
		if fromBlock != nil && fromBlock.Int64() >= 0 && fromBlock.Uint64() > log.BlockNumber {
			continue
		}
		if toBlock != nil && toBlock.Int64() >= 0 && toBlock.Uint64() < log.BlockNumber {
			continue
		}
		if len(addresses) > 0 && !includes(addresses, log.Address) {
			continue
		}
		// If the to filtered topics is greater than the amount of topics in logs, skip.
		if len(topics) > len(log.Topics) {
			continue Logs
		}
		for i, sub := range topics {
			match := len(sub) == 0 // empty rule set == wildcard
			for _, topic := range sub {
				if log.Topics[i] == topic {
					match = true
					break
				}
			}
			if !match {
				continue Logs
			}
		}
		ret = append(ret, log)
	}
	return ret
}
