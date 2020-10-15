package rpc

import (
	"context"
	"os"

	"github.com/tendermint/tendermint/libs/log"

	evmtypes "github.com/cosmos/ethermint/x/evm/types"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// Backend implements the functionality needed to filter changes.
// Implemented by EthermintBackend.
type Backend interface {
	// Used by block filter; also used for polling
	BlockNumber() (hexutil.Uint64, error)
	HeaderByNumber(blockNum BlockNumber) (*ethtypes.Header, error)
	HeaderByHash(blockHash common.Hash) (*ethtypes.Header, error)
	GetBlockByNumber(blockNum BlockNumber, fullTx bool) (map[string]interface{}, error)
	GetBlockByHash(hash common.Hash, fullTx bool) (map[string]interface{}, error)
	getEthBlockByNumber(height int64, fullTx bool) (map[string]interface{}, error)
	getGasLimit() (int64, error)
	// returns the logs of a given block
	GetLogs(blockHash common.Hash) ([][]*ethtypes.Log, error)

	// Used by pending transaction filter
	PendingTransactions() ([]*Transaction, error)

	// Used by log filter
	GetTransactionLogs(txHash common.Hash) ([]*ethtypes.Log, error)
	BloomStatus() (uint64, uint64)
}

var _ Backend = (*EthermintBackend)(nil)

// EthermintBackend implements the Backend interface
type EthermintBackend struct {
	ctx         context.Context
	clientCtx   client.Context
	queryClient *QueryClient // gRPC query client
	logger      log.Logger
}

// NewEthermintBackend creates a new EthermintBackend instance
func NewEthermintBackend(clientCtx client.Context) *EthermintBackend {
	return &EthermintBackend{
		ctx:         context.Background(),
		clientCtx:   clientCtx,
		queryClient: NewQueryClient(clientCtx),
		logger:      log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "json-rpc"),
	}
}

// BlockNumber returns the current block number.
func (e *EthermintBackend) BlockNumber() (hexutil.Uint64, error) {
	// NOTE: using 0 as min and max height returns the blockchain info up to the latest block.
	info, err := e.clientCtx.Client.BlockchainInfo(e.ctx, 0, 0)
	if err != nil {
		return hexutil.Uint64(0), err
	}

	return hexutil.Uint64(info.LastHeight), nil
}

// GetBlockByNumber returns the block identified by number.
func (e *EthermintBackend) GetBlockByNumber(blockNum BlockNumber, fullTx bool) (map[string]interface{}, error) {
	blockRes, err := e.clientCtx.Client.Block(e.ctx, blockNum.TmHeight())
	if err != nil {
		return nil, err
	}

	return e.getEthBlockByNumber(value, fullTx)
}

// GetBlockByHash returns the block identified by hash.
func (e *EthermintBackend) GetBlockByHash(hash common.Hash, fullTx bool) (map[string]interface{}, error) {
	resBlock, err := e.clientCtx.Client.BlockByHash(e.ctx, hash.Bytes())
	if err != nil {
		return nil, err
	}

	// TODO: gas, txs and bloom

	return formatBlock(resBlock.Block.Header, resBlock.Block.Size(), gasLimit, gasUsed, transactions, out.Bloom), nil
}

// HeaderByNumber returns the block header identified by height.
func (e *EthermintBackend) HeaderByNumber(blockNum BlockNumber) (*ethtypes.Header, error) {
	resBlock, err := e.clientCtx.Client.Block(e.ctx, blockNum.TmHeight())
	if err != nil {
		return nil, err
	}

	req := &evmtypes.QueryBlockBloomRequest{}

	res, err := e.queryClient.BlockBloom(ContextWithHeight(blockNum.Int64()), req)
	if err != nil {
		return nil, err
	}

	ethHeader := EthHeaderFromTendermint(resBlock.Block.Header)
	ethHeader.Bloom = ethtypes.BytesToBloom(res.Bloom)
	return ethHeader, nil
}

// HeaderByHash returns the block header identified by hash.
func (e *EthermintBackend) HeaderByHash(blockHash common.Hash) (*ethtypes.Header, error) {
	resBlock, err := e.clientCtx.Client.BlockByHash(e.ctx, blockHash.Bytes())
	if err != nil {
		return nil, err
	}

	req := &evmtypes.QueryBlockBloomRequest{}

	res, err := e.queryClient.BlockBloom(ContextWithHeight(resBlock.Block.Height), req)
	if err != nil {
		return nil, err
	}

	ethHeader := EthHeaderFromTendermint(resBlock.Block.Header)
	ethHeader.Bloom = ethtypes.BytesToBloom(res.Bloom)
	return ethHeader, nil
}

// GetTransactionLogs returns the logs given a transaction hash.
// It returns an error if there's an encoding error.
// If no logs are found for the tx hash, the error is nil.
func (e *EthermintBackend) GetTransactionLogs(txHash common.Hash) ([]*ethtypes.Log, error) {

	req := &evmtypes.QueryTxLogsRequest{
		Hash: txHash.String(),
	}

	res, err := e.queryClient.TxLogs(e.ctx, req)
	if err != nil {
		return nil, err
	}
	// TODO: logs to Ethereum
	return res.Logs, nil
}

// PendingTransactions returns the transactions that are in the transaction pool
// and have a from address that is one of the accounts this node manages.
func (e *EthermintBackend) PendingTransactions() ([]*Transaction, error) {
	limit := 1000
	pendingTxs, err := e.clientCtx.Client.UnconfirmedTxs(e.ctx, &limit)
	if err != nil {
		return nil, err
	}

	transactions := make([]*Transaction, pendingTxs.Count)
	for _, tx := range pendingTxs.Txs {
		ethTx, err := RawTxToEthTx(e.clientCtx, tx)
		if err != nil {
			return nil, err
		}

		// * Should check signer and reference against accounts the node manages in future
		rpcTx, err := NewTransaction(ethTx, common.BytesToHash(tx.Hash()), common.Hash{}, 0, 0)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, rpcTx)
	}

	return transactions, nil
}

// GetLogs returns all the logs from all the ethereum transactions in a block.
func (e *EthermintBackend) GetLogs(blockHash common.Hash) ([][]*ethtypes.Log, error) {
	// NOTE: we query the state in case the tx result logs are not persisted after an upgrade.
	req := &evmtypes.QueryBlockLogsRequest{
		Hash: blockHash.String(),
	}

	res, err := e.queryClient.BlockLogs(e.ctx, req)
	if err != nil {
		return nil, err
	}

	var blockLogs = [][]*ethtypes.Log{}
	for _, txLog := range res.TxLogs {
		blockLogs = append(blockLogs, txLog.EthLogs())
	}

	return blockLogs, nil
}

// BloomStatus returns the BloomBitsBlocks and the number of processed sections maintained
// by the chain indexer.
func (e *EthermintBackend) BloomStatus() (uint64, uint64) {
	return 4096, 0
}
