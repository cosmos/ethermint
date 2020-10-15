package provider

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/cosmos/ethermint/metrics"
	"github.com/pkg/errors"
	log "github.com/xlab/suplog"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TxInfo struct {
	Hash        common.Hash     `json:"hash"`
	BlockNumber *hexutil.Big    `json:"blockNumber"`
	Nonce       hexutil.Uint64  `json:"nonce"`
	From        common.Address  `json:"from"`
	To          *common.Address `json:"to"`
	Value       *hexutil.Big    `json:"value"`
}

func (p *ethProvider) TransactionByHash(ctx context.Context, hash common.Hash) (*TxInfo, error) {
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

	var raw json.RawMessage
	var info TxInfo

	if err := rpc.CallContext(ctx, &raw, "eth_getTransactionByHash", hash); err != nil {
		log.WithError(err).
			WithField("fn", "TransactionByHash").
			WithField("txHash", hash.Hex()).
			Warningln("failed to retrieve Tx info")
		metrics.ReportFuncError(p.svcTags)
		return nil, err
	} else if raw == nil {
		metrics.ReportFuncError(p.svcTags)
		return nil, errNotFound
	}

	if err := json.Unmarshal(raw, &info); err != nil {
		err = errors.Wrap(err, "failed to unmarshal Tx info")
		metrics.ReportFuncError(p.svcTags)
		return nil, err
	}

	return &info, nil
}

type TxReceipt struct {
	BlockNumber       hexutil.Uint64  `json:"blockNumber"`
	Status            hexutil.Uint64  `json:"status"`
	CumulativeGasUsed hexutil.Uint64  `json:"cumulativeGasUsed"`
	Logs              []*types.Log    `json:"logs"`
	Hash              common.Hash     `json:"transactionHash"`
	ContractAddress   *common.Address `json:"contractAddress"`
	GasUsed           hexutil.Uint64  `json:"gasUsed"`
}

func (p *ethProvider) TransactionReceiptByHash(ctx context.Context, hash common.Hash) (*TxReceipt, error) {
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

	var raw json.RawMessage
	var receipt TxReceipt

	if err := rpc.CallContext(ctx, &raw, "eth_getTransactionReceipt", hash); err != nil {
		log.WithError(err).
			WithField("fn", "TransactionReceiptByHash").
			WithField("txHash", hash).
			Warningln("failed to retrieve Tx receipt")
		metrics.ReportFuncError(p.svcTags)
		return nil, err
	} else if raw == nil {
		metrics.ReportFuncError(p.svcTags)
		return nil, errNotFound
	}

	if err := json.Unmarshal(raw, &receipt); err != nil {
		err = errors.Wrap(err, "failed to unmarshal Tx receipt")
		metrics.ReportFuncError(p.svcTags)
		return nil, err
	}

	return &receipt, nil
}

type BlockHeader struct {
	ParentHash common.Hash    `json:"parentHash"`
	Coinbase   common.Address `json:"miner"`
	Root       common.Hash    `json:"stateRoot"`
	Difficulty *hexutil.Big   `json:"difficulty"`
	Number     *hexutil.Big   `json:"number"`
	GasLimit   hexutil.Uint64 `json:"gasLimit"`
	Time       hexutil.Uint64 `json:"timestamp"`
}

func (b *BlockHeader) GetNumber() *big.Int         { return b.Number.ToInt() }
func (b *BlockHeader) GetGasLimit() uint64         { return uint64(b.GasLimit) }
func (b *BlockHeader) GetDifficulty() *big.Int     { return b.Difficulty.ToInt() }
func (b *BlockHeader) GetTime() uint64             { return uint64(b.Time) }
func (b *BlockHeader) GetNumberU64() uint64        { return b.Number.ToInt().Uint64() }
func (b *BlockHeader) GetCoinbase() common.Address { return b.Coinbase }
func (b *BlockHeader) GetRoot() common.Hash        { return b.Root }
func (b *BlockHeader) GetParentHash() common.Hash  { return b.ParentHash }

func (p *ethProvider) BlockHeaderByNumber(ctx context.Context, blockNum *big.Int) (*BlockHeader, error) {
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

	blockNumArg := toBlockNumArg(blockNum)

	var raw json.RawMessage
	var blockHeader BlockHeader

	if err := rpc.CallContext(ctx, &raw, "eth_getBlockByNumber", blockNumArg, false); err != nil {
		log.WithError(err).
			WithField("fn", "BlockHeaderByNumber").
			WithField("blockNum", blockNumArg).
			Warningln("failed to retrieve Block header")
		metrics.ReportFuncError(p.svcTags)
		return nil, err
	} else if raw == nil {
		metrics.ReportFuncError(p.svcTags)
		return nil, errNotFound
	}

	if err := json.Unmarshal(raw, &blockHeader); err != nil {
		err = errors.Wrap(err, "failed to unmarshal Block header")
		metrics.ReportFuncError(p.svcTags)
		return nil, err
	}

	return &blockHeader, nil
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}

	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}

	return hexutil.EncodeBig(number)
}

func (p *ethProvider) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	metrics.ReportFuncCall(p.svcTags)
	doneFn := metrics.ReportFuncTiming(p.svcTags)
	defer doneFn()

	ctx, cancelFn := contextWithCloseChan(ctx, p.closeC)
	defer cancelFn()

	rpc, _, ok := p.rpcClient(ctx)
	if !ok {
		metrics.ReportFuncError(p.svcTags)
		return errNodeUnavailable
	}

	cli := ethclient.NewClient(rpc)
	if err := cli.SendTransaction(ctx, tx); err != nil {
		metrics.ReportFuncError(p.svcTags)
		return err
	}

	return nil
}
