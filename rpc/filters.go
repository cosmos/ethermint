package rpc

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

/*
	- Filter functions derived from go-ethereum
	Used to set the criteria passed in from RPC params
*/

// Backend implements the functionality needed to filter changes.
// Implemented by PublicEthAPI.
type Backend interface {
	BlockNumber() (hexutil.Uint64, error)
	GetBlockByNumber(blockNum BlockNumber, fullTx bool) (map[string]interface{}, error)
	PendingTransactions() ([]*Transaction, error)
	GetTxLogs(txHash common.Hash) ([]*ethtypes.Log, error)
}

// Filter can be used to retrieve and filter logs.
type Filter struct {
	backend            Backend
	fromBlock, toBlock *big.Int         // start and end block numbers
	addresses          []common.Address // contract addresses to watch
	topics             [][]common.Hash  // log topics to watch for
	block              common.Hash      // Block hash if filtering a single block
}

// NewFilter returns a new Filter
func NewFilter(backend Backend, fromBlock, toBlock *big.Int, addresses []common.Address, topics [][]common.Hash) *Filter {
	return &Filter{
		backend:   backend,
		fromBlock: fromBlock,
		toBlock:   toBlock,
		addresses: addresses,
		topics:    topics,
	}
}

// NewBlockFilter creates a new filter which directly inspects the contents of
// a block to figure out whether it is interesting or not.
func NewBlockFilter(backend Backend, block common.Hash, addresses []common.Address, topics [][]common.Hash) *Filter {
	filter := NewFilter(backend, nil, nil, addresses, topics)
	filter.block = block
	return filter
}

// NewRangeFilter creates a new filter with a start and an end range.
func NewRangeFilter(backend Backend, addresses []common.Address, topics [][]common.Hash, fromBlock, toBlock *big.Int) *Filter {
	filter := NewFilter(backend, fromBlock, toBlock, addresses, topics)
	filter.fromBlock = fromBlock
	filter.toBlock = toBlock
	return filter
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
