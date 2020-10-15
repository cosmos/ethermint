package provider

import (
	"context"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	eth "github.com/cosmos/ethermint/ethereum/util"
	"github.com/cosmos/ethermint/metrics"
)

func (p *ethProvider) Balance(ctx context.Context, account common.Address) (*eth.Wei, error) {
	metrics.ReportFuncCall(p.svcTags)
	doneFn := metrics.ReportFuncTiming(p.svcTags)
	defer doneFn()

	balance, err := p.balanceAt(ctx, account, 0)
	if err != nil {
		metrics.ReportFuncError(p.svcTags)
		return nil, err
	}

	return balance, nil
}

func (p *ethProvider) BalanceAt(ctx context.Context, account common.Address, blockNum uint64) (*eth.Wei, error) {
	metrics.ReportFuncCall(p.svcTags)
	doneFn := metrics.ReportFuncTiming(p.svcTags)
	defer doneFn()

	balance, err := p.balanceAt(ctx, account, blockNum)
	if err != nil {
		metrics.ReportFuncError(p.svcTags)
		return nil, err
	}

	return balance, nil
}

func (p *ethProvider) balanceAt(ctx context.Context, account common.Address, blockNum uint64) (*eth.Wei, error) {
	ctx, cancelFn := contextWithCloseChan(ctx, p.closeC)
	defer cancelFn()

	rpc, _, ok := p.rpcClient(ctx)
	if !ok {
		return nil, errNodeUnavailable
	}

	var blockBigNum *big.Int
	if blockNum > 0 {
		blockBigNum = big.NewInt(0)
		blockBigNum.SetUint64(blockNum)
	}

	cli := ethclient.NewClient(rpc)

	bigint, err := cli.BalanceAt(ctx, account, blockBigNum)
	if err != nil {
		return nil, err
	}

	wei := eth.BigWei(bigint)

	return wei, nil
}

func (p *ethProvider) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	metrics.ReportFuncCall(p.svcTags)
	doneFn := metrics.ReportFuncTiming(p.svcTags)
	defer doneFn()

	ctx, cancelFn := contextWithCloseChan(ctx, p.closeC)
	defer cancelFn()

	rpc, _, ok := p.rpcClient(ctx)
	if !ok {
		metrics.ReportFuncError(p.svcTags)
		return 0, errNodeUnavailable
	}

	cli := ethclient.NewClient(rpc)

	nonce, err := cli.PendingNonceAt(ctx, account)
	if err != nil {
		metrics.ReportFuncError(p.svcTags)
		return 0, err
	}

	return nonce, nil
}

func (p *ethProvider) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	metrics.ReportFuncCall(p.svcTags)
	doneFn := metrics.ReportFuncTiming(p.svcTags)
	defer doneFn()

	ctx, cancelFn := contextWithCloseChan(ctx, p.closeC)
	defer cancelFn()

	rpc, _, ok := p.rpcClient(ctx)
	if !ok {
		metrics.ReportFuncError(p.svcTags)
		return nil, errNodeUnavailable
	}

	cli := ethclient.NewClient(rpc)
	code, err := cli.PendingCodeAt(ctx, account)
	if err != nil {
		metrics.ReportFuncError(p.svcTags)
		return nil, err
	}

	return code, nil
}

func (p *ethProvider) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	metrics.ReportFuncCall(p.svcTags)
	doneFn := metrics.ReportFuncTiming(p.svcTags)
	defer doneFn()

	ctx, cancelFn := contextWithCloseChan(ctx, p.closeC)
	defer cancelFn()

	rpc, _, ok := p.rpcClient(ctx)
	if !ok {
		metrics.ReportFuncError(p.svcTags)
		return 0, errNodeUnavailable
	}

	cli := ethclient.NewClient(rpc)

	gas, err := cli.EstimateGas(ctx, msg)
	if err != nil {
		metrics.ReportFuncError(p.svcTags)
		return 0, err
	}

	return gas, nil
}

func (p *ethProvider) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	metrics.ReportFuncCall(p.svcTags)
	doneFn := metrics.ReportFuncTiming(p.svcTags)
	defer doneFn()

	ctx, cancelFn := contextWithCloseChan(ctx, p.closeC)
	defer cancelFn()

	rpc, _, ok := p.rpcClient(ctx)
	if !ok {
		metrics.ReportFuncError(p.svcTags)
		return nil, errNodeUnavailable
	}

	cli := ethclient.NewClient(rpc)

	price, err := cli.SuggestGasPrice(ctx)
	if err != nil {
		metrics.ReportFuncError(p.svcTags)
		return nil, err
	}

	return price, nil
}
