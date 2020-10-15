package provider

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/cosmos/ethermint/metrics"
	"github.com/serialx/hashring"
	log "github.com/xlab/suplog"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	errNodeUnavailable = errors.New("no EVM node available")
	errNotFound        = errors.New("the requested object not found")
)

var _ bind.ContractBackend = &ethProvider{}

func NewEthProvider(nodes []string, gasLimit uint64) EVMProvider {
	provider := &ethProvider{
		// default fallback chain id
		chainID: 888,

		session:  NewSessionID(),
		ring:     hashring.New(nodes),
		ringMux:  new(sync.RWMutex),
		fails:    make(map[string]int),
		closeC:   make(chan struct{}, 1),
		wg:       new(sync.WaitGroup),
		gasLimit: gasLimit,

		svcTags: metrics.Tags{
			"module": "eth_provider",
		},
	}

	if len(nodes) > 0 {
		if cli, _, ok := provider.rpcClient(context.Background()); ok {
			var id hexutil.Uint64
			if err := cli.Call(&id, "eth_chainId"); err != nil {
				log.WithError(err).
					WithField("fn", "NewEthProvider").
					Warningln("failed to get chainID from RPC endpoint")
			} else if uint64(id) != provider.chainID {
				log.WithField("chain_id", id).
					WithField("fn", "NewEthProvider").
					Info("setting new chainID")
				provider.chainID = uint64(id)
			}
		}
	}

	return provider
}

type ethProvider struct {
	session  string
	chainID  uint64
	ring     *hashring.HashRing
	ringMux  *sync.RWMutex
	fails    map[string]int
	closeC   chan struct{}
	wg       *sync.WaitGroup
	gasLimit uint64

	svcTags metrics.Tags
}

func (p *ethProvider) Nodes() []string {
	p.ringMux.RLock()
	nodes, _ := p.ring.GetNodes("", p.ring.Size())
	p.ringMux.RUnlock()
	return nodes
}

func (p *ethProvider) Close() error {
	close(p.closeC)
	p.wg.Wait()
	return nil
}

func (p *ethProvider) ChainID() uint64 {
	return p.chainID
}

func (p *ethProvider) GasLimit() uint64 {
	return p.gasLimit
}

func (p *ethProvider) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	metrics.ReportFuncCall(p.svcTags)
	doneFn := metrics.ReportFuncTiming(p.svcTags)
	defer doneFn()

	ctx, cancelFn := contextWithCloseChan(ctx, p.closeC)
	defer cancelFn()

	rpcClient, _, ok := p.rpcClient(ctx)
	if !ok {
		metrics.ReportFuncError(p.svcTags)
		return nil, errNodeUnavailable
	}

	cli := ethclient.NewClient(rpcClient)
	return cli.CodeAt(ctx, contract, blockNumber)
}

func (p *ethProvider) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	metrics.ReportFuncCall(p.svcTags)
	doneFn := metrics.ReportFuncTiming(p.svcTags)
	defer doneFn()

	ctx, cancelFn := contextWithCloseChan(ctx, p.closeC)
	defer cancelFn()

	rpcClient, _, ok := p.rpcClient(ctx)
	if !ok {
		metrics.ReportFuncError(p.svcTags)
		return nil, errNodeUnavailable
	}

	cli := ethclient.NewClient(rpcClient)
	return cli.CallContract(ctx, call, blockNumber)
}

func (p *ethProvider) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	metrics.ReportFuncCall(p.svcTags)
	doneFn := metrics.ReportFuncTiming(p.svcTags)
	defer doneFn()

	ctx, cancelFn := contextWithCloseChan(ctx, p.closeC)
	defer cancelFn()

	rpcClient, _, ok := p.rpcClient(ctx)
	if !ok {
		metrics.ReportFuncError(p.svcTags)
		return nil, errNodeUnavailable
	}

	cli := ethclient.NewClient(rpcClient)
	return cli.FilterLogs(ctx, query)
}

func (p *ethProvider) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	metrics.ReportFuncCall(p.svcTags)

	ctx, cancelFn := contextWithCloseChan(ctx, p.closeC)
	defer cancelFn()

	rpcClient, _, ok := p.rpcClient(ctx)
	if !ok {
		metrics.ReportFuncError(p.svcTags)
		return nil, errNodeUnavailable
	}

	cli := ethclient.NewClient(rpcClient)
	return cli.SubscribeFilterLogs(ctx, query, ch)
}

func (p *ethProvider) rpcClient(ctx context.Context) (cli *rpc.Client, addr string, ok bool) {
	for {
		p.ringMux.RLock()
		addr, ok = p.ring.GetNode(p.session)
		p.ringMux.RUnlock()
		if !ok {
			err := errors.New("no available RPC nodes in pool, all dead x_X")
			log.WithError(err).
				WithField("fn", "rpcClient").
				Warningln("failed to pick node")

			return nil, "", false
		}
		newCli, err := rpc.Dial(addr)
		if err == nil {
			cli = newCli
			break
		}

		log.WithError(err).
			WithField("addr", addr).
			WithField("fn", "rpcClient").
			Warningln("failed to connect to RPC node")

		p.failNode(addr)
		time.Sleep(3 * time.Second)
	}

	return cli, addr, ok
}

func (p *ethProvider) failNode(addr string) {
	p.ringMux.Lock()
	defer p.ringMux.Unlock()
	if p.fails[addr] < 0 {
		// node been removed
		return
	}
	p.fails[addr]++
	if p.fails[addr] < 3 {
		return
	}
	p.fails[addr] = -1
	p.ring = p.ring.RemoveNode(addr)

	log.WithField("addr", addr).
		WithField("fn", "failNode").
		Warningln("RPC node has been removed from pool, will attempt revival in 5 minutes")

	go func() {
		// schedule a revival
		time.Sleep(5 * time.Minute)
		p.reviveNode(addr)
	}()
}

func (p *ethProvider) reviveNode(addr string) {
	p.ringMux.Lock()
	defer p.ringMux.Unlock()
	if p.fails[addr] >= 0 {
		// node been restored
		return
	}

	log.WithField("addr", addr).
		WithField("fn", "reviveNode").
		Warningln("RPC node has been added back into pool")

	p.ring = p.ring.AddNode(addr)
	p.fails[addr] = 0
}
