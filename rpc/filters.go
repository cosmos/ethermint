package rpc

import (
	"context"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"


	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/bloombits"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

/*
	- Filter functions derived from go-ethereum
	Used to set the criteria passed in from RPC params
*/

type Backend interface {
	ChainDb() ethdb.Database
	EventMux() *event.TypeMux
	HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*ethtypes.Header, error)
	HeaderByHash(ctx context.Context, blockHash common.Hash) (*ethtypes.Header, error)
	GetReceipts(ctx context.Context, blockHash common.Hash) (*ethtypes.Receipt, error)
	GetLogs(ctx context.Context, blockHash common.Hash) ([][]*ethtypes.Log, error)
	SubscribeNewTxsEvent(chan<- core.NewTxsEvent) event.Subscription
	SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription
	SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription
	SubscribeLogsEvent(ch chan<- []*ethtypes.Log) event.Subscription
	BloomStatus() (uint64, uint64)
	ServiceFilter(ctx context.Context, session *bloombits.MatcherSession)
}

// Filter can be used to retrieve and filter logs.
type Filter struct {
	backend Backend
	addresses []common.Address
	topics    [][]common.Hash

	block      common.Hash // Block hash if filtering a single block
	begin, end int64       // Range interval if filtering multiple blocks

	matcher *bloombits.Matcher
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
//func NewRangeFilter(backend Backend, begin, end int64, addresses []common.Address, topics [][]common.Hash) *Filter {
//	// Flatten the address and topic filter clauses into a single bloombits filter
//	// system. Since the bloombits are not positional, nil topics are permitted,
//	// which get flattened into a nil byte slice.
//	var filters [][][]byte
//	if len(addresses) > 0 {
//		filter := make([][]byte, len(addresses))
//		for i, address := range addresses {
//			filter[i] = address.Bytes()
//		}
//		filters = append(filters, filter)
//	}
//	for _, topicList := range topics {
//		filter := make([][]byte, len(topicList))
//		for i, topic := range topicList {
//			filter[i] = topic.Bytes()
//		}
//		filters = append(filters, filter)
//	}
//	fmt.Println("BLOOM STATUS:::")
//	size, _ := BloomStatus()
//	fmt.Println("size ::: ", size)
//
//	// Create a generic filter and convert it into a range filter
//	filter := newFilter(backend, addresses, topics)
//
//	filter.matcher = bloombits.NewMatcher(size, filters)
//	fmt.Println("filter matcher", filter.matcher)
//	filter.begin = begin
//	filter.end = end
//
//	return filter
//}

//func BloomStatus() (uint64, uint64) { return 4096, 0 }

// newFilter creates a generic filter that can either filter based on a block hash,
// or based on range queries. The search criteria needs to be explicitly set.
func newFilter(backend Backend, addresses []common.Address, topics [][]common.Hash) *Filter {
	return &Filter{
		backend: backend,
		addresses: addresses,
		topics:    topics,
	}
}

//func bn(cliCtx context2.CLIContext) (int64, error) {
//	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/blockNumber", types.ModuleName), nil)
//	if err != nil {
//		return int64(0), err
//	}
//
//	var out types.QueryResBlockNumber
//	cliCtx.Codec.MustUnmarshalJSON(res, &out)
//	return out.Number, nil
//}
//
//// Logs searches the blockchain for matching log entries, returning all from the
//// first block that contains matches, updating the start of the filter accordingly.
//func (f *Filter) Logs(ctx context2.CLIContext) ([]*ethtypes.Log, error) {
//
//	fmt.Println("LOGS:::::::::::::::::::::::::")
//	// If we're doing singleton block filtering, execute and return
//	if f.block != (common.Hash{}) {
//		//header, err := f.backend.HeaderByHash(ctx, f.block)
//		//if err != nil {
//		//	return nil, err
//		//}
//		//if header == nil {
//		//	return nil, errors.New("unknown block")
//		//}
//		//return f.blockLogs(ctx, header)
//	}
//	// Figure out the limits of the filter range
//	bn, _ := bn(ctx)
//	header, _ := ctx.Client.Block(&bn)
//	//header, _ := f.backend.HeaderByNumber(ctx, rpc.LatestBlockNumber)
//	if header == nil {
//		return nil, nil
//	}
//	head := header.Block.Height
//	if f.begin == -1 {
//		f.begin = int64(head)
//	}
//	end := uint64(f.end)
//	if f.end == -1 {
//		end = uint64(head)
//	}
//	// Gather all indexed logs, and finish with non indexed ones
//	var (
//		logs []*ethtypes.Log
//		err  error
//	)
//	size, sections := BloomStatus()
//	if indexed := sections * size; indexed > uint64(f.begin) {
//		if indexed > end {
//			fmt.Println("INDEXED:::::", indexed)
//			logs, err = f.indexedLogs(context.Background(), end)
//		} else {
//			fmt.Println("ELSE:::::")
//			logs, err = f.indexedLogs(context.Background(), indexed-1)
//		}
//		if err != nil {
//			return logs, err
//		}
//	}
//	fmt.Println("outside::::")
//	rest, err := f.unindexedLogs(ctx, end)
//	logs = append(logs, rest...)
//	return logs, err
//}
//
//// indexedLogs returns the logs matching the filter criteria based on the bloom
//// bits indexed available locally or via the network.
//func (f *Filter) indexedLogs(ctx context.Context, end uint64) ([]*ethtypes.Log, error) {
//	// Create a matcher session and request servicing from the backend
//	matches := make(chan uint64, 64)
//
//	session, err := f.matcher.Start(ctx, uint64(f.begin), end, matches)
//	if err != nil {
//		return nil, err
//	}
//	defer session.Close()
//
//	f.backend.ServiceFilter(ctx, session)
//
//	// Iterate over the matches until exhausted or context closed
//	var logs []*ethtypes.Log
//
//	for {
//		select {
//		case number, ok := <-matches:
//			// Abort if all matches have been fulfilled
//			if !ok {
//				err := session.Error()
//				if err == nil {
//					f.begin = int64(end) + 1
//				}
//				return logs, err
//			}
//			f.begin = int64(number) + 1
//
//			// Retrieve the suggested block and pull any truly matching logs
//			header, err := f.backend.HeaderByNumber(ctx, rpc.BlockNumber(number))
//			if header == nil || err != nil {
//				return logs, err
//			}
//
//			found, err := f.checkMatches(ctx, header)
//			if err != nil {
//				return logs, err
//			}
//			logs = append(logs, found...)
//
//		case <-ctx.Done():
//			return logs, ctx.Err()
//		}
//	}
//}
//
//// indexedLogs returns the logs matching the filter criteria based on raw block
//// iteration and bloom matching.
//func (f *Filter) unindexedLogs(ctx context2.CLIContext, end uint64) ([]*ethtypes.Log, error) {
//	var logs []*ethtypes.Log
//
//	for ; f.begin <= int64(end); f.begin++ {
//		num := rpc.BlockNumber(f.begin).Int64()
//		bl, err := ctx.Client.Block(&num)
//
//		if err != nil {
//			fmt.Println("err", err)
//		}
//		header := &bl.BlockMeta.Header
//		//header, err := f.backend.HeaderByNumber(ctx, rpc.BlockNumber(f.begin))
//		if header == nil || err != nil {
//			return logs, err
//		}
//		bloom, err := getBloom(ctx, bl)
//		found, err := f.blockLogs(context.Background(), header, bloom)
//		if err != nil {
//			return logs, err
//		}
//		logs = append(logs, found...)
//	}
//	return logs, nil
//}
//
//func getBloom(cliCtx context2.CLIContext, block *core_types.ResultBlock) (ethtypes.Bloom, error) {
//	res, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s/%s", types.ModuleName, evm.QueryLogsBloom, strconv.FormatInt(block.Block.Height, 10)))
//	if err != nil {
//		return ethtypes.Bloom{}, err
//	}
//
//	var out types.QueryBloomFilter
//	cliCtx.Codec.MustUnmarshalJSON(res, &out)
//	return out.Bloom, nil
//}
//
//// blockLogs returns the logs matching the filter criteria within a single block.
//func (f *Filter) blockLogs(ctx context.Context, header *types2.Header, bloom ethtypes.Bloom) (logs []*ethtypes.Log, err error) {
//	if bloomFilter(bloom, f.addresses, f.topics) {
//		found, err := f.checkMatches(ctx, header)
//		if err != nil {
//			return logs, err
//		}
//		logs = append(logs, found...)
//	}
//	return logs, nil
//}
//
//// checkMatches checks if the receipts belonging to the given header contain any log events that
//// match the filter criteria. This function is called when the bloom filter signals a potential match.
//func (f *Filter) checkMatches(ctx context.Context, header *types2.Header) (logs []*ethtypes.Log, err error) {
//	// Get the logs of the block
//	fmt.Println("here")
//	logsList, err := f.backend.GetLogs(ctx, common.BytesToHash(header.Hash()))
//	if err != nil {
//		return nil, err
//	}
//	var unfiltered []*ethtypes.Log
//	for _, logs := range logsList {
//		unfiltered = append(unfiltered, logs...)
//	}
//	logs = filterLogs(unfiltered, nil, nil, f.addresses, f.topics)
//	if len(logs) > 0 {
//		// We have matching logs, check if we need to resolve full logs via the light client
//		if logs[0].TxHash == (common.Hash{}) {
//			receipts, err := f.backend.GetReceipts(ctx, common.BytesToHash(header.Hash()))
//			if err != nil {
//				return nil, err
//			}
//			fmt.Println("receipts:::: ", receipts)
//			unfiltered = unfiltered[:0]
//			//for _, receipt := range receipts {
//			//	unfiltered = append(unfiltered, receipt.Logs...)
//			//}
//			logs = filterLogs(unfiltered, nil, nil, f.addresses, f.topics)
//		}
//		return logs, nil
//	}
//	return nil, nil
//}
//
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

func bloomFilter(bloom ethtypes.Bloom, addresses []common.Address, topics [][]common.Hash) bool {
	if len(addresses) > 0 {
		var included bool
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
