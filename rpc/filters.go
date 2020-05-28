package rpc

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/bloombits"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rpc"
)

/*
	- Filter functions derived from go-ethereum
	Used to set the criteria passed in from RPC params
*/

type Filter struct {
	backend Backend

	typ      filters.Type  // filter type
	deadline *time.Timer   // filter is inactiv when deadline triggers
	hashes   []common.Hash // filtered block or transaction hashes
	criteria filters.FilterCriteria

	matcher *bloombits.Matcher

	subscription bool // associated subscription in event system
}

// NewFilter returns a new Filter
func NewFilter(backend Backend, filterType filters.Type, criteria filters.FilterCriteria) *Filter {
	return &Filter{
		backend:  backend,
		typ:      filterType,
		deadline: time.NewTimer(deadline),
		criteria: criteria,
	}
}

func (f *Filter) Unsubscribe() {
	if !f.subscription {
		return
	}

	// switch f.typ {
	// case:
	// }

}

// Logs searches the blockchain for matching log entries, returning all from the
// first block that contains matches, updating the start of the filter accordingly.
func (f *Filter) Logs() ([]*ethtypes.Log, error) {
	logs := []*ethtypes.Log{}
	var err error
	// If we're doing singleton block filtering, execute and return
	if f.criteria.BlockHash != nil || f.criteria.BlockHash != (&common.Hash{}) {
		header, err := f.backend.HeaderByHash(*f.criteria.BlockHash)
		if err != nil {
			return nil, err
		}
		if header == nil {
			return nil, fmt.Errorf("unknown block header %s", f.criteria.BlockHash.String())
		}
		return f.blockLogs(header)
	}

	// Figure out the limits of the filter range
	header, err := f.backend.HeaderByNumber(rpc.LatestBlockNumber)
	if err != nil {
		return nil, err
	}

	if header == nil {
		return nil, nil
	}
	head := header.Number.Uint64()

	if f.criteria.FromBlock.Int64() == -1 {
		f.criteria.FromBlock = big.NewInt(int64(head))
	}
	if f.criteria.ToBlock.Int64() == -1 {
		f.criteria.ToBlock = big.NewInt(int64(head))
	}

	for i := f.criteria.FromBlock.Int64(); i <= f.criteria.ToBlock.Int64(); i++ {
		block, err := f.backend.GetBlockByNumber(rpc.BlockNumber(i), true)
		if err != nil {
			return logs, err
		}

		txs, ok := block["transactions"].([]common.Hash)
		if !ok || len(txs) == 0 {
			continue
		}

		logsMatched, err := f.checkMatches(txs)
		if err != nil {
			return logs, err
		}

		logs = append(logs, logsMatched...)
	}

	return logs, nil
}

// blockLogs returns the logs matching the filter criteria within a single block.
func (f *Filter) blockLogs(header *ethtypes.Header) ([]*ethtypes.Log, error) {
	if !bloomFilter(header.Bloom, f.criteria.Addresses, f.criteria.Topics) {
		return []*ethtypes.Log{}, nil
	}

	return f.checkMatches(header)
}

func (f *Filter) checkMatches(transactions []common.Hash) ([]*ethtypes.Log, error) {
	unfiltered := []*ethtypes.Log{}
	for _, tx := range transactions {
		logs, err := f.backend.GetTransactionLogs(tx)
		if err != nil {
			// ignore error if transaction didn't set any logs (eg: when tx type is not
			// MsgEthereumTx or MsgEthermint)
			continue
		}

		unfiltered = append(unfiltered, logs...)
	}

	return filterLogs(unfiltered, f.criteria.FromBlock, f.criteria.ToBlock, f.criteria.Addresses, f.criteria.Topics), nil
}

// filterLogs creates a slice of logs matching the given criteria.
// [] -> anything
// [A] -> A in first position of log topics, anything after
// [null, B] -> anything in first position, B in second position
// [A, B] -> A in first position and B in second position
// [[A, B], [A, B]] -> A or B in first position, A or B in second position
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

func includes(addresses []common.Address, a common.Address) bool {
	for _, addr := range addresses {
		if addr == a {
			return true
		}
	}

	return false
}

func bloomFilter(bloom ethtypes.Bloom, addresses []common.Address, topics [][]common.Hash) bool {
	var included bool
	if len(addresses) > 0 {
		for _, addr := range addresses {
			if ethtypes.BloomLookup(bloom, addr) {
				included = true
				break
			}
		}
		if !included {
			return false
		}
	}

	for _, sub := range topics {
		included = len(sub) == 0 // empty rule set == wildcard
		for _, topic := range sub {
			if ethtypes.BloomLookup(bloom, topic) {
				included = true
				break
			}
		}
	}
	return included
}
