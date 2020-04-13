package rpc

import (
	"errors"
	"fmt"
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

const blockFilter = "block"
const pendingTxFilter = "pending"
const logFilter = "log"

// Filter can be used to retrieve and filter logs, blocks, or pending transactions.
type Filter struct {
	backend            Backend
	fromBlock, toBlock *big.Int         // start and end block numbers
	addresses          []common.Address // contract addresses to watch
	topics             [][]common.Hash  // log topics to watch for
	blockHash          *common.Hash     // Block hash if filtering a single block

	typ     string
	hashes  []common.Hash   // filtered block or transaction hashes
	logs    []*ethtypes.Log //nolint // filtered logs
	stopped bool            // set to true once filter in uninstalled

	err error
}

// NewFilter returns a new Filter
func NewFilter(backend Backend, criteria *filters.FilterCriteria) *Filter {
	filter := &Filter{
		backend:   backend,
		fromBlock: criteria.FromBlock,
		toBlock:   criteria.ToBlock,
		addresses: criteria.Addresses,
		topics:    criteria.Topics,
		typ:       logFilter,
		stopped:   false,
	}

	return filter
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
	filter := NewFilter(backend, &filters.FilterCriteria{})
	filter.typ = blockFilter

	go func() {
		err := filter.pollForBlocks()
		if err != nil {
			filter.err = err
		}
	}()

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

		if num == prev {
			continue
		}

		block, err := f.backend.GetBlockByNumber(BlockNumber(num), false)
		if err != nil {
			return err
		}

		hashBytes, ok := block["hash"].(hexutil.Bytes)
		if !ok {
			return errors.New("could not convert block hash to hexutil.Bytes")
		}

		hash := common.BytesToHash([]byte(hashBytes))
		f.hashes = append(f.hashes, hash)

		prev = num

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

func (f *Filter) getFilterChanges() (interface{}, error) {
	switch f.typ {
	case blockFilter:
		if f.err != nil {
			return nil, f.err
		}

		blocks := make([]common.Hash, len(f.hashes))
		copy(blocks, f.hashes)
		f.hashes = []common.Hash{}

		return blocks, nil
	case pendingTxFilter:
		// TODO
	case logFilter:
		fmt.Println("getFilterLogs logFilter")
		return f.getFilterLogs()
	}

	return nil, errors.New("unsupported filter")
}

func (f *Filter) getFilterLogs() ([]*ethtypes.Log, error) {
	ret := []*ethtypes.Log{}

	// filter specific block only
	if f.blockHash != nil {
		block, err := f.backend.GetBlockByHash(*f.blockHash, true)
		if err != nil {
			return nil, err
		}

		// if the logsBloom == 0, there are no logs in that block
		if bloom, ok := block["logsBloom"].(ethtypes.Bloom); ok && big.NewInt(0).SetBytes(bloom[:]).Cmp(big.NewInt(0)) == 0 {
			return ret, nil
		} else if ok {
			return f.checkMatches(block)
		} else {
			return nil, errors.New("invalid logsBloom returned")
		}
	}

	// filter range of blocks
	// TODO: check if toBlock or fromBlock is "latest", "pending", or "earliest"
	num, err := f.backend.BlockNumber()
	if err != nil {
		return nil, err
	}

	// if f.fromBlock is set to 0, set it to the latest block number
	if f.fromBlock == nil || f.fromBlock.Cmp(big.NewInt(0)) == 0 {
		f.fromBlock = big.NewInt(int64(num))
	}

	// if f.toBlock is set to 0, set it to the latest block number
	if f.toBlock == nil || f.toBlock.Cmp(big.NewInt(0)) == 0 {
		f.toBlock = big.NewInt(int64(num))
	}

	fmt.Printf("fromBlock=%d\n", f.fromBlock)
	fmt.Printf("toBlock=%d\n", f.toBlock)
	fmt.Printf("topics=%v\n", f.topics)
	fmt.Printf("addresses=%v\n", f.addresses)

	from := f.fromBlock.Int64()
	to := f.toBlock.Int64()

	for i := from; i <= to; i++ {
		fmt.Printf("i=%d\n", i)

		block, err := f.backend.GetBlockByNumber(NewBlockNumber(big.NewInt(i)), true)
		if err != nil {
			f.err = err
			fmt.Printf("cannot get block %d err=%s\n", block["number"], err)
			continue
		}

		fmt.Printf("got block=%v\n", block)

		// TODO: block logsBloom is often set in the wrong block
		// if the logsBloom == 0, there are no logs in that block
		//if bloom, ok := block["logsBloom"].(ethtypes.Bloom); !ok {
		if txs, ok := block["transactions"].([]common.Hash); !ok {
			fmt.Printf("could not cast transactions\n")
			continue
			//} else if big.NewInt(0).SetBytes(bloom[:]).Cmp(big.NewInt(0)) != 0 {
		} else if len(txs) != 0 {
			logs, err := f.checkMatches(block)
			if err != nil {
				f.err = err // return or just keep for later?
				fmt.Printf("checkMatches err=%s\n", err)
				continue
			}

			fmt.Printf("block %d logs=%v\n", block["number"], logs)
			ret = append(ret, logs...)
		} else {
			continue
		}
	}

	return ret, nil
}

func (f *Filter) checkMatches(block map[string]interface{}) ([]*ethtypes.Log, error) {
	transactions, ok := block["transactions"].([]common.Hash)
	if !ok {
		return nil, errors.New("invalid block transactions")
	}

	unfiltered := []*ethtypes.Log{}

	fmt.Println("searching block txs")

	for _, tx := range transactions {
		fmt.Printf("txhash=%x\n", tx)

		logs, err := f.backend.GetTxLogs(common.BytesToHash(tx[:]))
		if err != nil {
			return nil, err
		}

		fmt.Println("logs len=", len(logs))

		unfiltered = append(unfiltered, logs...)
	}

	//return unfiltered, nil
	return filterLogs(unfiltered, f.fromBlock, f.toBlock, f.addresses, f.topics), nil
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
