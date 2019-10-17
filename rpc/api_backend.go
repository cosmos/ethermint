package rpc

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/bloombits"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/net/context"
)

type EmintAPIBackend struct {
	bloomRequests chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer  *core.ChainIndexer
}

func (b *EmintAPIBackend) ChainDb() ethdb.Database {
	return nil
}

func (b *EmintAPIBackend) EventMux() *event.TypeMux {
	return nil
}

func (b *EmintAPIBackend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*ethtypes.Header, error) {
	return nil, nil
}

func (b *EmintAPIBackend) HeaderByHash(ctx context.Context, blockHash common.Hash) (*ethtypes.Header, error) {
	return nil, nil
}

func (b *EmintAPIBackend) GetReceipts(ctx context.Context, blockHash common.Hash) (*ethtypes.Receipt, error) {
	return nil, nil
}

func (b *EmintAPIBackend) GetLogs(ctx context.Context, blockHash common.Hash) ([][]*ethtypes.Log, error) {
	return nil, nil
}

func (b *EmintAPIBackend) SubscribeNewTxsEvent(chan<- core.NewTxsEvent) event.Subscription {
	return nil
}

func (b *EmintAPIBackend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return nil
}

func (b *EmintAPIBackend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return nil
}

func (b *EmintAPIBackend) SubscribeLogsEvent(ch chan<- []*ethtypes.Log) event.Subscription {
	return nil
}

func (b *EmintAPIBackend) BloomStatus() (uint64, uint64) {
	fmt.Println("bloomstatus")
	sections, _, _ := b.bloomIndexer.Sections()
	fmt.Println("indexer__+_+)+")
	return params.BloomBitsBlocks, sections
}

func (b *EmintAPIBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	return
}