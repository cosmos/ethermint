package provider

import (
	"context"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	eth "github.com/cosmos/ethermint/ethereum/util"
)

type EVMProvider interface {
	bind.ContractCaller
	bind.ContractFilterer

	Balance(ctx context.Context, account common.Address) (*eth.Wei, error)
	BalanceAt(ctx context.Context, account common.Address, blockNum uint64) (*eth.Wei, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error)
	EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	TransactionByHash(ctx context.Context, hash common.Hash) (*TxInfo, error)
	TransactionReceiptByHash(ctx context.Context, hash common.Hash) (*TxReceipt, error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	BlockHeaderByNumber(ctx context.Context, blockNum *big.Int) (*BlockHeader, error)

	ChainID() uint64
	GasLimit() uint64

	Nodes() []string
	Close() error
}
