package rpc

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/ethermint/x/evm/types"

	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rpc"
)

// PublicFilterAPI is the eth_ prefixed set of APIs in the Web3 JSON-RPC spec.
type PublicFilterAPI struct {
	cliCtx  context.CLIContext
	backend Backend
	filters map[uint64]*Filter // ID to filter
	nextID  uint64
}

// NewPublicEthAPI creates an instance of the public ETH Web3 API.
func NewPublicFilterAPI(cliCtx context.CLIContext, backend Backend) *PublicFilterAPI {
	return &PublicFilterAPI{
		cliCtx:  cliCtx,
		backend: backend,
		filters: make(map[uint64]*Filter),
		nextID:  0,
	}
}

// NewFilter instantiates a new filter.
func (e *PublicFilterAPI) NewFilter(fromBlock, toBlock *big.Int, addresses []common.Address, topics [][]common.Hash) uint64 {
	e.filters[e.nextID] = NewFilter(e.backend, fromBlock, toBlock, addresses, topics)
	e.nextID++
	return e.nextID
}

// NewBlockFilter instantiates a new block filter.
func (e *PublicFilterAPI) NewBlockFilter() uint64 {
	e.filters[e.nextID] = NewBlockFilter(e.backend)
	e.nextID++
	return e.nextID
}

// NewPendingTransactionFilter instantiates a new pending transaction filter.
func (e *PublicFilterAPI) NewPendingTransactionFilter() uint64 {
	e.filters[e.nextID] = NewPendingTransactionFilter(e.backend)
	e.nextID++
	return e.nextID
}

// UninstallFilter uninstalls a filter with the given ID.
func (e *PublicFilterAPI) UninstallFilter(id uint64) bool {
	// TODO
	delete(e.filters, id)
	return false
}

// GetFilterChanges returns an array of changes since the last poll.
// If the filter is a log filter, it returns an array of Logs.
// If the filter is a block filter, it returns an array of block hashes.
// If the filter is a pending transaction filter, it returns an array of transaction hashes.
func (e *PublicFilterAPI) GetFilterChanges(id uint64) interface{} {
	return nil
}

// GetFilterLogs returns an array of all logs matching filter with given id.
func (e *PublicFilterAPI) GetFilterLogs(id uint64) []*ethtypes.Log {
	return nil
}

// GetLogs returns logs matching the given argument that are stored within the state.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getlogs
func (e *PublicFilterAPI) GetLogs(criteria filters.FilterCriteria) ([]*ethtypes.Log, error) {
	var filter *Filter
	if criteria.BlockHash != nil {
		/*
			Still need to add blockhash in prepare function for log entry
		*/
		filter = NewFilterWithBlockHash(e.backend, criteria.FromBlock, criteria.ToBlock, *criteria.BlockHash, criteria.Addresses, criteria.Topics)
		results := e.getLogs()
		logs := filterLogs(results, nil, nil, filter.addresses, filter.topics)
		return logs, nil
	}
	// Convert the RPC block numbers into internal representations
	begin := rpc.LatestBlockNumber.Int64()
	if criteria.FromBlock != nil {
		begin = criteria.FromBlock.Int64()
	}
	from := big.NewInt(begin)
	end := rpc.LatestBlockNumber.Int64()
	if criteria.ToBlock != nil {
		end = criteria.ToBlock.Int64()
	}
	to := big.NewInt(end)
	results := e.getLogs()
	logs := filterLogs(results, from, to, criteria.Addresses, criteria.Topics)

	return returnLogs(logs), nil
}

func (e *PublicFilterAPI) getLogs() (results []*ethtypes.Log) {
	l, _, err := e.cliCtx.QueryWithData(fmt.Sprintf("custom/%s/logs", types.ModuleName), nil)
	if err != nil {
		fmt.Printf("error from querier %e ", err)
	}

	if err := json.Unmarshal(l, &results); err != nil {
		panic(err)
	}
	return results
}

// returnLogs is a helper that will return an empty log array in case the given logs array is nil,
// otherwise the given logs array is returned.
func returnLogs(logs []*ethtypes.Log) []*ethtypes.Log {
	if logs == nil {
		return []*ethtypes.Log{}
	}
	return logs
}
