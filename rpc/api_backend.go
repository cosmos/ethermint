package rpc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/core/bloombits"

	context2 "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/ethermint/x/evm"
	"github.com/cosmos/ethermint/x/evm/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
)

const (
	// bloomFilterThreads is the number of goroutines used locally per filter to
	// multiplex requests onto the global servicing goroutines.
	bloomFilterThreads = 3

	// bloomRetrievalBatch is the maximum number of bloom bit retrievals to service
	// in a single batch.
	bloomRetrievalBatch = 16

	// bloomRetrievalWait is the maximum time to wait for enough bloom bit requests
	// to accumulate request an entire batch (avoiding hysteresis).
	bloomRetrievalWait = time.Duration(0)
)

type EmintAPIBackend struct {
	bloomRequests chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	//bloomIndexer  *core.ChainIndexer             // Bloom indexer operating during block imports
}

// ServiceFilter stubbed right now
func (e *EmintAPIBackend) ServiceFilter(session *bloombits.MatcherSession) {
	for i := 0; i < bloomFilterThreads; i++ {
		go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, e.bloomRequests)
	}
}

// BloomStatus stubbed for right now
func (e *EmintAPIBackend) BloomStatus() (uint64, uint64) { return 4096, 0 }

// GetBlockNumber returns block number from canonical head
func (e *EmintAPIBackend) GetBlockNumber(cliCtx context2.CLIContext) (int64, error) {
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/blockNumber", types.ModuleName), nil)
	if err != nil {
		return int64(0), err
	}

	var out types.QueryResBlockNumber
	cliCtx.Codec.MustUnmarshalJSON(res, &out)
	return out.Number, nil
}

// GetBlockNumberFromHash returns block number from provided block hash
func (e *EmintAPIBackend) GetBlockNumberFromHash(cliCtx context2.CLIContext, bhash common.Hash) (int64, error) {
	res, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s/%s", types.ModuleName, evm.QueryHashToHeight, bhash.Hex()))
	if err != nil {
		// Return nil if block does not exist
		return 0, err
	}

	var out types.QueryResBlockNumber
	cliCtx.Codec.MustUnmarshalJSON(res, &out)
	return out.Number, nil
}

// GetBloom returns bloomFilter from provided block
func (e *EmintAPIBackend) GetBloom(cliCtx context2.CLIContext, block *core_types.ResultBlock) (ethtypes.Bloom, error) {
	res, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s/%s", types.ModuleName, evm.QueryLogsBloom, strconv.FormatInt(block.Block.Height, 10)))
	if err != nil {
		return ethtypes.Bloom{}, err
	}

	var out types.QueryBloomFilter
	cliCtx.Codec.MustUnmarshalJSON(res, &out)
	return out.Bloom, nil
}

// GetReceipts returns a slice of receipts from provided block
//
// Used in checkMatches fn to check if we need to resolve full logs
func (e *EmintAPIBackend) GetReceipts(cliCtx context2.CLIContext, block *core_types.ResultBlock) []Receipt {
	var transactions []Receipt

	for _, tx := range block.Block.Txs {

		receipts, err := e.GetTransactionReceipt(cliCtx, common.BytesToHash(tx.Hash()))
		if err != nil {
			fmt.Println("err retrieving tx receipt", err)
		}
		transactions = append(transactions, *receipts)
	}
	return transactions
}

// Ethereum receipt struct
type Receipt struct {
	BlockHash         common.Hash     `json:"blockHash"`
	BlockNumber       uint64          `json:"blockNumber"`
	TransactionHash   common.Hash     `json:"transactionHash"`
	TransactionIndex  uint64          `json:"transactionIndex"`
	From              common.Address  `json:"from"`
	To                *common.Address `json:"to"`
	GasUsed           uint64          `json:"gasUsed"`
	CumulativeGasUsed uint64          `json:"cumulativeGasUsed"`
	ContractAddress   common.Address  `json:"contractAddress"`
	Logs              []*ethtypes.Log `json:"logs"`
	LogsBloom         ethtypes.Bloom  `json:"logsBloom"`
	Status            uint            `json:"status"`
}

// GetTransactionReceipt returns the transaction receipt identified by hash.
func (e *EmintAPIBackend) GetTransactionReceipt(cliCtx context2.CLIContext, hash common.Hash) (*Receipt, error) {
	tx, err := cliCtx.Client.Tx(hash.Bytes(), false)
	if err != nil {
		// Return nil for transaction when not found
		return nil, nil
	}

	// Query block for consensus hash
	block, err := cliCtx.Client.Block(&tx.Height)
	if err != nil {
		return nil, err
	}
	blockHash := common.BytesToHash(block.BlockMeta.Header.Hash())

	// Convert tx bytes to eth transaction
	ethTx, err := bytesToEthTx(cliCtx, tx.Tx)
	if err != nil {
		return nil, err
	}

	from, _ := ethTx.VerifySig(ethTx.ChainID())

	// Set status codes based on tx result
	var status uint
	if tx.TxResult.IsOK() {
		status = uint(1)
	} else {
		status = uint(0)
	}

	res, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s/%s", types.ModuleName, evm.QueryTxLogs, hash.Hex()))
	if err != nil {
		return nil, err
	}

	var logs types.QueryETHLogs
	cliCtx.Codec.MustUnmarshalJSON(res, &logs)

	// TODO: change hard coded indexing of bytes
	bloomFilter := ethtypes.BytesToBloom(tx.TxResult.GetData()[20:])

	rec := Receipt{
		BlockNumber:       uint64(tx.Height),
		BlockHash:         blockHash,
		TransactionHash:   hash,
		TransactionIndex:  uint64(tx.Index),
		From:              from,
		To:                ethTx.To(),
		GasUsed:           uint64(tx.TxResult.GasUsed),
		CumulativeGasUsed: 0,
		ContractAddress:   common.Address{},
		Logs:              logs.Logs,
		LogsBloom:         bloomFilter,
		Status:            status,
	}

	contractAddress := common.BytesToAddress(tx.TxResult.GetData()[:20])
	if contractAddress != (common.Address{}) {
		// TODO: change hard coded indexing of first 20 bytes
		rec.ContractAddress = contractAddress
	}

	return &rec, nil
}

// GetBlockLogs creates a slice of logs matching the given criteria.
func (e *EmintAPIBackend) GetBlockLogs(ctx context2.CLIContext, logs []*ethtypes.Log, bhash common.Hash, blockNumber int64) ([]*ethtypes.Log, error) {
	var ret []*ethtypes.Log
	var bn int64
	var blockHash common.Hash
	var err error
	// filter criteria specified block hash
	if bhash != (common.Hash{}) {
		blockHash = bhash
		bn, err = e.GetBlockNumberFromHash(ctx, blockHash)
		if err != nil {
			return nil, err
		}
	} else {
		// no block hash provided but a block range
		bn = blockNumber
		b, err := ctx.Client.Block(&blockNumber)
		if err != nil {
			return nil, err
		}
		hash := b.BlockMeta.Header.Hash()
		bhash := "0x" + common.Bytes2Hex(hash)
		blockHash = common.HexToHash(bhash)
	}

	for _, lg := range logs {
		// if block numbers match append correct blockhash for response
		if bn > 0 && uint64(bn) == lg.BlockNumber {
			lg.BlockHash = blockHash
			ret = append(ret, lg)
		}
	}

	return ret, nil
}

// GetLogs returns logs from committed state
func (e *EmintAPIBackend) GetLogs(cliCtx context2.CLIContext) (results []*ethtypes.Log) {
	var res types.QueryETHLogs
	l, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/logs", types.ModuleName), nil)
	if err != nil {
		fmt.Printf("error from querier %e ", err)
	}
	if err := json.Unmarshal(l, &res); err != nil {
		panic(err)
	}
	return res.Logs
}
