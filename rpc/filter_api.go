package rpc

import (
	"github.com/cosmos/cosmos-sdk/client/context"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rpc"
)

// PublicFilterAPI is the eth_ prefixed set of APIs in the Web3 JSON-RPC spec.
type PublicFilterAPI struct {
	cliCtx  context.CLIContext
	backend Backend
}

// NewPublicEthAPI creates an instance of the public ETH Web3 API.
func NewPublicFilterAPI(cliCtx context.CLIContext) *PublicFilterAPI {
	return &PublicFilterAPI{
		cliCtx:  cliCtx,
		backend: &EmintAPIBackend{},
	}
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
		filter = NewBlockFilter(e.backend, *criteria.BlockHash, criteria.Addresses, criteria.Topics)
	} else {
		// Convert the RPC block numbers into internal representations
		begin := rpc.LatestBlockNumber.Int64()
		if criteria.FromBlock != nil {
			begin = criteria.FromBlock.Int64()
		}
		end := rpc.LatestBlockNumber.Int64()
		if criteria.ToBlock != nil {
			end = criteria.ToBlock.Int64()
		}
		filter = NewRangeFilter(e.backend, begin, end, criteria.Addresses, criteria.Topics)
	}
	logs, err := filter.Logs(e.cliCtx)
	if err != nil {
		return nil, err
	}

	return returnLogs(logs), nil
}

// returnLogs is a helper that will return an empty log array in case the given logs array is nil,
// otherwise the given logs array is returned.
func returnLogs(logs []*ethtypes.Log) []*ethtypes.Log {
	if logs == nil {
		return []*ethtypes.Log{}
	}
	return logs
}
