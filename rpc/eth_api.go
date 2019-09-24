package rpc

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/cosmos/cosmos-sdk/client/context"
	authutils "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	emintcrypto "github.com/cosmos/ethermint/crypto"
	emintkeys "github.com/cosmos/ethermint/keys"
	"github.com/cosmos/ethermint/rpc/args"
	"github.com/cosmos/ethermint/version"
	"github.com/cosmos/ethermint/x/evm/types"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/viper"
)

// PublicEthAPI is the eth_ prefixed set of APIs in the Web3 JSON-RPC spec.
type PublicEthAPI struct {
	cliCtx    context.CLIContext
	key       emintcrypto.PrivKeySecp256k1
	nonceLock *AddrLocker
}

// NewPublicEthAPI creates an instance of the public ETH Web3 API.
func NewPublicEthAPI(cliCtx context.CLIContext, nonceLock *AddrLocker,
	key emintcrypto.PrivKeySecp256k1) *PublicEthAPI {

	return &PublicEthAPI{
		cliCtx:    cliCtx,
		key:       key,
		nonceLock: nonceLock,
	}
}

// ProtocolVersion returns the supported Ethereum protocol version.
func (e *PublicEthAPI) ProtocolVersion() hexutil.Uint {
	return hexutil.Uint(version.ProtocolVersion)
}

// Syncing returns whether or not the current node is syncing with other peers. Returns false if not, or a struct
// outlining the state of the sync if it is.
func (e *PublicEthAPI) Syncing() interface{} {
	return false
}

// Coinbase returns this node's coinbase address. Not used in Ethermint.
func (e *PublicEthAPI) Coinbase() (addr common.Address) {
	return
}

// Mining returns whether or not this node is currently mining. Always false.
func (e *PublicEthAPI) Mining() bool {
	return false
}

// Hashrate returns the current node's hashrate. Always 0.
func (e *PublicEthAPI) Hashrate() hexutil.Uint64 {
	return 0
}

// GasPrice returns the current gas price based on Ethermint's gas price oracle.
func (e *PublicEthAPI) GasPrice() *hexutil.Big {
	out := big.NewInt(0)
	return (*hexutil.Big)(out)
}

// Accounts returns the list of accounts available to this node.
func (e *PublicEthAPI) Accounts() ([]common.Address, error) {
	addresses := make([]common.Address, 0) // return [] instead of nil if empty
	keybase, err := emintkeys.NewKeyBaseFromHomeFlag()
	if err != nil {
		return addresses, err
	}

	infos, err := keybase.List()
	if err != nil {
		return addresses, err
	}

	for _, info := range infos {
		addressBytes := info.GetPubKey().Address().Bytes()
		addresses = append(addresses, common.BytesToAddress(addressBytes))
	}

	return addresses, nil
}

// BlockNumber returns the current block number.
func (e *PublicEthAPI) BlockNumber() (hexutil.Uint64, error) {
	res, _, err := e.cliCtx.QueryWithData(fmt.Sprintf("custom/%s/blockNumber", types.ModuleName), nil)
	if err != nil {
		return hexutil.Uint64(0), err
	}

	var qRes uint64
	e.cliCtx.Codec.MustUnmarshalJSON(res, &qRes)

	return hexutil.Uint64(qRes), nil
}

// GetBalance returns the provided account's balance up to the provided block number.
func (e *PublicEthAPI) GetBalance(address common.Address, blockNum rpc.BlockNumber) (*hexutil.Big, error) {
	ctx := e.cliCtx.WithHeight(blockNum.Int64())
	res, _, err := ctx.QueryWithData(fmt.Sprintf("custom/%s/balance/%s", types.ModuleName, address), nil)
	if err != nil {
		return nil, err
	}

	var out *big.Int
	e.cliCtx.Codec.MustUnmarshalJSON(res, &out)

	return (*hexutil.Big)(out), nil
}

// GetStorageAt returns the contract storage at the given address, block number, and key.
func (e *PublicEthAPI) GetStorageAt(address common.Address, key string, blockNum rpc.BlockNumber) (hexutil.Bytes, error) {
	ctx := e.cliCtx.WithHeight(blockNum.Int64())
	res, _, err := ctx.QueryWithData(fmt.Sprintf("custom/%s/storage/%s/%s", types.ModuleName, address, key), nil)
	if err != nil {
		return nil, err
	}

	var out []byte
	e.cliCtx.Codec.MustUnmarshalJSON(res, &out)
	return out, nil
}

// GetTransactionCount returns the number of transactions at the given address up to the given block number.
func (e *PublicEthAPI) GetTransactionCount(address common.Address, blockNum rpc.BlockNumber) (*hexutil.Uint64, error) {
	ctx := e.cliCtx.WithHeight(blockNum.Int64())
	res, _, err := ctx.QueryWithData(fmt.Sprintf("custom/%s/nonce/%s", types.ModuleName, address), nil)
	if err != nil {
		return nil, err
	}

	var out *hexutil.Uint64
	e.cliCtx.Codec.MustUnmarshalJSON(res, out)
	return out, nil
}

// GetBlockTransactionCountByHash returns the number of transactions in the block identified by hash.
func (e *PublicEthAPI) GetBlockTransactionCountByHash(hash common.Hash) hexutil.Uint {
	return 0
}

// GetBlockTransactionCountByNumber returns the number of transactions in the block identified by number.
func (e *PublicEthAPI) GetBlockTransactionCountByNumber(blockNum rpc.BlockNumber) (hexutil.Uint, error) {
	node, err := e.cliCtx.GetNode()
	if err != nil {
		return 0, err
	}

	height := blockNum.Int64()
	block, err := node.Block(&height)
	if err != nil {
		return 0, err
	}

	return hexutil.Uint(block.Block.NumTxs), nil
}

// GetUncleCountByBlockHash returns the number of uncles in the block idenfied by hash. Always zero.
func (e *PublicEthAPI) GetUncleCountByBlockHash(hash common.Hash) hexutil.Uint {
	return 0
}

// GetUncleCountByBlockNumber returns the number of uncles in the block idenfied by number. Always zero.
func (e *PublicEthAPI) GetUncleCountByBlockNumber(blockNum rpc.BlockNumber) hexutil.Uint {
	return 0
}

// GetCode returns the contract code at the given address and block number.
func (e *PublicEthAPI) GetCode(address common.Address, blockNumber rpc.BlockNumber) (hexutil.Bytes, error) {
	ctx := e.cliCtx.WithHeight(blockNumber.Int64())
	res, _, err := ctx.QueryWithData(fmt.Sprintf("custom/%s/code/%s", types.ModuleName, address), nil)
	if err != nil {
		return nil, err
	}

	var out []byte
	e.cliCtx.Codec.MustUnmarshalJSON(res, &out)
	return out, nil
}

// Sign signs the provided data using the private key of address via Geth's signature standard.
func (e *PublicEthAPI) Sign(address common.Address, data hexutil.Bytes) (hexutil.Bytes, error) {
	// TODO: Change this functionality to find an unlocked account by address
	if e.key == nil || !bytes.Equal(e.key.PubKey().Address().Bytes(), address.Bytes()) {
		return nil, keystore.ErrLocked
	}

	// Sign the requested hash with the wallet
	signature, err := e.key.Sign(data)
	if err == nil {
		signature[64] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	}

	return signature, err
}

// SendTransaction sends an Ethereum transaction.
func (e *PublicEthAPI) SendTransaction(args args.SendTxArgs) (common.Hash, error) {
	// TODO: Change this functionality to find an unlocked account by address
	if e.key == nil || !bytes.Equal(e.key.PubKey().Address().Bytes(), args.From.Bytes()) {
		return common.Hash{}, keystore.ErrLocked
	}

	// Mutex lock the address' nonce to avoid assigning it to multiple requests
	if args.Nonce == nil {
		e.nonceLock.LockAddr(args.From)
		defer e.nonceLock.UnlockAddr(args.From)
	}

	// Assemble transaction from fields
	tx, err := types.GenerateFromArgs(args, e.cliCtx)
	if err != nil {
		return common.Hash{}, err
	}

	// ChainID must be set as flag to send transaction
	chainID := viper.GetString(flags.FlagChainID)
	// parse the chainID from a string to a base-10 integer
	intChainID, ok := new(big.Int).SetString(chainID, 10)
	if !ok {
		return common.Hash{}, fmt.Errorf(
			fmt.Sprintf("Invalid chainID: %s, must be integer format", chainID))
	}

	// Sign transaction
	tx.Sign(intChainID, e.key.ToECDSA())

	// Encode transaction by default Tx encoder
	txEncoder := authutils.GetTxEncoder(e.cliCtx.Codec)
	txBytes, err := txEncoder(tx)
	if err != nil {
		return common.Hash{}, err
	}

	// Broadcast transaction
	res, err := e.cliCtx.BroadcastTx(txBytes)
	// If error is encountered on the node, the broadcast will not return an error
	// TODO: Remove res log
	fmt.Println(res.RawLog)
	if err != nil {
		return common.Hash{}, err
	}

	// Return RLP encoded bytes
	return tx.Hash(), nil
}

// SendRawTransaction send a raw Ethereum transaction.
func (e *PublicEthAPI) SendRawTransaction(data hexutil.Bytes) (common.Hash, error) {
	tx := new(types.EthereumTxMsg)

	// RLP decode raw transaction bytes
	if err := rlp.DecodeBytes(data, tx); err != nil {
		// Return nil is for when gasLimit overflows uint64
		return common.Hash{}, nil
	}

	// Encode transaction by default Tx encoder
	txEncoder := authutils.GetTxEncoder(e.cliCtx.Codec)
	txBytes, err := txEncoder(tx)
	if err != nil {
		return common.Hash{}, err
	}

	// TODO: Possibly log the contract creation address (if recipient address is nil) or tx data
	res, err := e.cliCtx.BroadcastTx(txBytes)
	// If error is encountered on the node, the broadcast will not return an error
	// TODO: Remove res log
	fmt.Println(res.RawLog)
	if err != nil {
		return common.Hash{}, err
	}

	// Return RLP encoded bytes
	return tx.Hash(), nil
}

// CallArgs represents arguments to a smart contract call as provided by RPC clients.
type CallArgs struct {
	From     common.Address `json:"from"`
	To       common.Address `json:"to"`
	Gas      hexutil.Uint64 `json:"gas"`
	GasPrice hexutil.Big    `json:"gasPrice"`
	Value    hexutil.Big    `json:"value"`
	Data     hexutil.Bytes  `json:"data"`
}

// Call performs a raw contract call.
func (e *PublicEthAPI) Call(args CallArgs, blockNum rpc.BlockNumber) hexutil.Bytes {
	return nil
}

// EstimateGas estimates gas usage for the given smart contract call.
func (e *PublicEthAPI) EstimateGas(args CallArgs, blockNum rpc.BlockNumber) hexutil.Uint64 {
	return 0
}

// GetBlockByHash returns the block identified by hash.
func (e *PublicEthAPI) GetBlockByHash(hash common.Hash, fullTx bool) map[string]interface{} {
	return nil
}

// GetBlockByNumber returns the block identified by number.
func (e *PublicEthAPI) GetBlockByNumber(blockNum rpc.BlockNumber, fullTx bool) map[string]interface{} {
	return nil
}

// Transaction represents a transaction returned to RPC clients.
type Transaction struct {
	BlockHash        common.Hash     `json:"blockHash"`
	BlockNumber      *hexutil.Big    `json:"blockNumber"`
	From             common.Address  `json:"from"`
	Gas              hexutil.Uint64  `json:"gas"`
	GasPrice         *hexutil.Big    `json:"gasPrice"`
	Hash             common.Hash     `json:"hash"`
	Input            hexutil.Bytes   `json:"input"`
	Nonce            hexutil.Uint64  `json:"nonce"`
	To               *common.Address `json:"to"`
	TransactionIndex hexutil.Uint    `json:"transactionIndex"`
	Value            *hexutil.Big    `json:"value"`
	V                *hexutil.Big    `json:"v"`
	R                *hexutil.Big    `json:"r"`
	S                *hexutil.Big    `json:"s"`
}

// GetTransactionByHash returns the transaction identified by hash.
func (e *PublicEthAPI) GetTransactionByHash(hash common.Hash) *Transaction {
	return nil
}

// GetTransactionByBlockHashAndIndex returns the transaction identified by hash and index.
func (e *PublicEthAPI) GetTransactionByBlockHashAndIndex(hash common.Hash, idx hexutil.Uint) *Transaction {
	return nil
}

// GetTransactionByBlockNumberAndIndex returns the transaction identified by number and index.
func (e *PublicEthAPI) GetTransactionByBlockNumberAndIndex(blockNumber rpc.BlockNumber, idx hexutil.Uint) *Transaction {
	return nil
}

// GetTransactionReceipt returns the transaction receipt identified by hash.
func (e *PublicEthAPI) GetTransactionReceipt(hash common.Hash) map[string]interface{} {
	return nil
}

// GetUncleByBlockHashAndIndex returns the uncle identified by hash and index. Always returns nil.
func (e *PublicEthAPI) GetUncleByBlockHashAndIndex(hash common.Hash, idx hexutil.Uint) map[string]interface{} {
	return nil
}

// GetUncleByBlockNumberAndIndex returns the uncle identified by number and index. Always returns nil.
func (e *PublicEthAPI) GetUncleByBlockNumberAndIndex(number hexutil.Uint, idx hexutil.Uint) map[string]interface{} {
	return nil
}
