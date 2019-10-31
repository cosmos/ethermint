package rpc

import (
	"context"
	"fmt"
	context2 "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/ethereum/go-ethereum/rpc"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/bloombits"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

/*
	- Filter functions derived from go-ethereum
	Used to set the criteria passed in from RPC params
*/

// Backend contains function signatures for querying blockchain data
type Backend interface {
	GetBlockNumber(cliCtx context2.CLIContext) (int64, error)
	GetBlockNumberFromHash(cliCtx context2.CLIContext, bhash common.Hash) (int64, error)
	GetBloom(cliCtx context2.CLIContext, block *core_types.ResultBlock) (ethtypes.Bloom, error)
	GetReceipts(cliCtx context2.CLIContext, block *core_types.ResultBlock) []Receipt
	GetTransactionReceipt(cliCtx context2.CLIContext, hash common.Hash) (*Receipt, error)
	GetBlockLogs(ctx context2.CLIContext, logs []*ethtypes.Log, bhash common.Hash, blockNumber int64) ([]*ethtypes.Log, error)
	GetLogs(cliCtx context2.CLIContext) (results []*ethtypes.Log)
	BloomStatus() (uint64, uint64)
	ServiceFilter(session *bloombits.MatcherSession)
}

// Filter can be used to retrieve and filter logs.
type Filter struct {
	backend   Backend
	addresses []common.Address
	topics    [][]common.Hash
	block      common.Hash // Block hash if filtering a single block
	begin, end int64       // Range interval if filtering multiple blocks

	matcher *bloombits.Matcher
}

// newFilter creates a generic filter that can either filter based on a block hash,
// or based on range queries. The search criteria needs to be explicitly set.
func newFilter(backend Backend, addresses []common.Address, topics [][]common.Hash) *Filter {
	return &Filter{
		backend:   backend,
		addresses: addresses,
		topics:    topics,
	}
}

// NewBlockFilter creates a new filter which directly inspects the contents of
// a block to figure out whether it is interesting or not.
func NewBlockFilter(backend Backend, block common.Hash, addresses []common.Address, topics [][]common.Hash) *Filter {
	// Create a generic filter and convert it into a block filter
	filter := newFilter(backend, addresses, topics)
	filter.block = block
	return filter
}

// NewRangeFilter creates a new filter which uses a bloom filter on blocks to
// figure out whether a particular block is interesting or not.
func NewRangeFilter(backend Backend, begin, end int64, addresses []common.Address, topics [][]common.Hash) *Filter {
	// Flatten the address and topic filter clauses into a single bloombits filter
	// system. Since the bloombits are not positional, nil topics are permitted,
	// which get flattened into a nil byte slice.
	var filters [][][]byte
	if len(addresses) > 0 {
		filter := make([][]byte, len(addresses))
		for i, address := range addresses {
			filter[i] = address.Bytes()
		}
		filters = append(filters, filter)
	}
	for _, topicList := range topics {
		filter := make([][]byte, len(topicList))
		for i, topic := range topicList {
			filter[i] = topic.Bytes()
		}
		filters = append(filters, filter)
	}
	// BloomStatus is currently stubbed to return 4096, 0
	// ----------------------------------------------------------------------------
	// Returns size and sections variables which are required for bloom matching
	// see https://github.com/ethereum/go-ethereum/blob/master/core/bloombits/matcher.go AND
	// https://github.com/ethereum/go-ethereum/blob/master/eth/api_backend.go
	// ----------------------------------------------------------------------------
	size, _ := backend.BloomStatus()
	// Create a generic filter and convert it into a range filter
	filter := newFilter(backend, addresses, topics)

	filter.matcher = bloombits.NewMatcher(size, filters)
	filter.begin = begin
	filter.end = end

	return filter
}

//// Logs searches the blockchain for matching log entries, returning all from the
//// first block that contains matches, updating the start of the filter accordingly.
func (f *Filter) Logs(ctx context2.CLIContext) ([]*ethtypes.Log, error) {
	// If we're doing singleton block filtering, execute and return
	if f.block != (common.Hash{}) {
		bn, err := f.backend.GetBlockNumberFromHash(ctx, f.block)
		if err != nil {
			 return nil, err
		}
		bl, err := ctx.Client.Block(&bn)
		if err != nil {
			return nil, err
		}

		// GetBloom from block height
		bloom, err := f.backend.GetBloom(ctx, bl)
		if err != nil {
			return nil, err
		}
		return f.blockLogs(ctx, bl, bloom, f.block)
	}
	// Figure out the limits of the filter range
	bn, _ := f.backend.GetBlockNumber(ctx)
	bl, _ := ctx.Client.Block(&bn)

	head := bl.Block.Height
	if f.begin == -1 {
		f.begin = int64(head) -1
	}
	end := uint64(f.end)
	if f.end == -1 {
		end = uint64(head) -1
	}
	// Gather all indexed logs, and finish with non indexed ones
	var (
		logs []*ethtypes.Log
		err  error
	)

	// Not really meaningful right now since indexed will never be greater than f.begin
	size, sections := f.backend.BloomStatus()

	// ----------------------------------------------------------------------------
	// indexedLogs is currently never hit since sections * size will always equal 0 as it is stubbed;
	// change manually or pull out indexedLogs func and call explicitly
	// when user inputs address/topic it should hit indexedLogs function
	// @todo context
	// ----------------------------------------------------------------------------
	if indexed := sections * size; indexed > uint64(f.begin) {
		if indexed > end {
			logs, err = f.indexedLogs(ctx, context.Background(), end)
		} else {
			logs, err = f.indexedLogs(ctx, context.Background(), indexed-1)
		}
		if err != nil {
			return logs, err
		}
	}

	rest, err := f.unindexedLogs(ctx, end)
	logs = append(logs, rest...)
	return logs, err
}

//// indexedLogs returns the logs matching the filter criteria based on the bloom
//// bits indexed available locally or via the network.
func (f *Filter) indexedLogs(cliCtx context2.CLIContext, ctx context.Context, end uint64) ([]*ethtypes.Log, error) {
	// Create a matcher session and request servicing from the backend
	matches := make(chan uint64, 64)
	// ----------------------------------------------------------------------------
	// f.matcher.Start starts the matching process and returns a stream of bloom matches in
	// a given range of blocks
	// see see https://github.com/ethereum/go-ethereum/blob/master/core/bloombits/matcher.go
	// ----------------------------------------------------------------------------
	// f.matcher.Start does not do much right now since ServiceFilter is not implemented
	session, err := f.matcher.Start(ctx, uint64(f.begin), end, matches)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	// @TODO implement ServiceFilter
	f.backend.ServiceFilter(session)
	// Iterate over the matches until exhausted or context closed
	var logs []*ethtypes.Log

	// ----------------------------------------------------------------------------
	// never receive matches currently which is required to continue to filter logs and
	// return logs given filter params
	// ----------------------------------------------------------------------------
	for {
		select {
		case number, ok := <-matches:
			// Abort if all matches have been fulfilled
			if !ok {
				err := session.Error()
				if err == nil {
					f.begin = int64(end) + 1
				}
				return logs, err
			}
			f.begin = int64(number) + 1

			// Retrieve the suggested block and pull any truly matching logs
			num := int64(number)
			bl, err := cliCtx.Client.Block(&num)
			if err != nil {
				return nil, err
			}

			//bloom, err := f.backend.GetBloom(cliCtx, bl)
			//if err != nil {
			//	return logs, err
			//}

			found, err := f.checkMatches(cliCtx, bl, f.block)
			if err != nil {
				return logs, err
			}
			logs = append(logs, found...)

		case <-ctx.Done():
			return logs, ctx.Err()
		}
	}

}

// unindexedLogs...
func (f *Filter) unindexedLogs(ctx context2.CLIContext, end uint64) ([]*ethtypes.Log, error) {
	fmt.Println("unindexed logs ---> ")
	var logs []*ethtypes.Log

	for ; f.begin <= int64(end); f.begin++ {
		num := rpc.BlockNumber(f.begin).Int64()
		bl, err := ctx.Client.Block(&num)
		if err != nil {
			return nil, err
		}

		bloom, err := f.backend.GetBloom(ctx, bl)
		if err != nil {
			return logs, err
		}

		found, err := f.blockLogs(ctx, bl, bloom, f.block)
		if err != nil {
			return logs, err
		}
		logs = append(logs, found...)
	}

	return logs, nil
}

//// blockLogs returns the logs matching the filter criteria within a single block.
func (f *Filter) blockLogs(ctx context2.CLIContext, block *core_types.ResultBlock, bloom ethtypes.Bloom, bhash common.Hash) (logs []*ethtypes.Log, err error) {
	if bloomFilter(bloom, f.addresses, f.topics) {
		found, err := f.checkMatches(ctx, block, bhash)
		if err != nil {
			return logs, err
		}
		logs = append(logs, found...)
	}
	return logs, nil
}

//// checkMatches checks if the receipts belonging to the given header contain any log events that
//// match the filter criteria. This function is called when the bloom filter signals a potential match.
func (f *Filter) checkMatches(ctx context2.CLIContext, block *core_types.ResultBlock, bhash common.Hash) (logs []*ethtypes.Log, err error) {
	// Get the logs of the block
	lg := f.backend.GetLogs(ctx)
	logsList, err := f.backend.GetBlockLogs(ctx, lg, bhash, block.Block.Height)
	if err != nil {
		return nil, err
	}
	logs = filterLogs(logsList, nil, nil, f.addresses, f.topics)

	if len(logs) > 0 {
		// @TODO not sure if this is necessary given its for light client
		// We have matching logs, check if we need to resolve full logs via the light client
		if logs[0].TxHash == (common.Hash{}) {
			receipts := f.backend.GetReceipts(ctx, block)
			logsList = logsList[:0]
			for _, receipt := range receipts {
				logsList = append(logsList, receipt.Logs...)
			}
			logs = filterLogs(logsList, nil, nil, f.addresses, f.topics)
		}
		return logs, nil
	}
	return nil, nil
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
	for _, lg := range logs {
		if fromBlock != nil && fromBlock.Int64() >= 0 && fromBlock.Uint64() > lg.BlockNumber {
			continue
		}
		if toBlock != nil && toBlock.Int64() >= 0 && toBlock.Uint64() < lg.BlockNumber {
			continue
		}
		if len(addresses) > 0 && !includes(addresses, lg.Address) {
			continue
		}
		// If the to filtered topics is greater than the amount of topics in logs, skip.
		if len(topics) > len(lg.Topics) {
			continue Logs
		}
		for i, sub := range topics {
			match := len(sub) == 0 // empty rule set == wildcard
			for _, topic := range sub {
				if lg.Topics[i] == topic {
					match = true
					break
				}
			}
			if !match {
				continue Logs
			}
		}
		ret = append(ret, lg)
	}
	return ret
}

func bloomFilter(bloom ethtypes.Bloom, addresses []common.Address, topics [][]common.Hash) bool {
	if len(addresses) > 0 {
		var included bool
		for _, addr := range addresses {
			// ----------------------------------------------------------------------------
			// BloomLookup never returns true even though the addr
			// and bloom should provide a positive check
			// ----------------------------------------------------------------------------
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
		included := len(sub) == 0 // empty rule set == wildcard
		for _, topic := range sub {
			if ethtypes.BloomLookup(bloom, topic) {
				included = true
				break
			}
		}
		if !included {
			return false
		}
	}
	return true
}
