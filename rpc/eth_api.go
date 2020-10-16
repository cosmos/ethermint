package rpc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"sync"

	"github.com/spf13/viper"

	"github.com/cosmos/ethermint/crypto/ethsecp256k1"
	"github.com/cosmos/ethermint/crypto/hd"
	ethermint "github.com/cosmos/ethermint/types"
	"github.com/cosmos/ethermint/version"
	evmtypes "github.com/cosmos/ethermint/x/evm/types"

	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/grpc/simulate"
	"github.com/cosmos/cosmos-sdk/client/tx"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// PublicEthAPI is the eth_ prefixed set of APIs in the Web3 JSON-RPC spec.
type PublicEthAPI struct {
	ctx          context.Context
	clientCtx    client.Context
	queryClient  *QueryClient // gRPC query client
	chainIDEpoch *big.Int
	logger       log.Logger
	backend      Backend
	keys         []ethsecp256k1.PrivKey // unlocked keys
	nonceLock    *AddrLocker
	keyringLock  sync.Mutex
}

// NewPublicEthAPI creates an instance of the public ETH Web3 API.
func NewPublicEthAPI(clientCtx client.Context, backend Backend, nonceLock *AddrLocker,
	keys ...ethsecp256k1.PrivKey) *PublicEthAPI {

	epoch, err := ethermint.ParseChainID(clientCtx.ChainID)
	if err != nil {
		panic(err)
	}

	api := &PublicEthAPI{
		ctx:          context.Background(),
		clientCtx:    clientCtx,
		queryClient:  NewQueryClient(clientCtx),
		chainIDEpoch: epoch,
		logger:       log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "json-rpc"),
		backend:      backend,
		keys:         keys,
		nonceLock:    nonceLock,
	}

	if err := api.getKeyringInfo(); err != nil {
		api.logger.Error("failed to get keyring info", "error", err)
	}

	return api
}

func (e *PublicEthAPI) getKeyringInfo() error {
	e.keyringLock.Lock()
	defer e.keyringLock.Unlock()

	if e.clientCtx.Keyring == nil {
		keyring, err := keyring.New(
			sdk.KeyringServiceName(),
			viper.GetString(flags.FlagKeyringBackend),
			viper.GetString(flags.FlagHome),
			e.clientCtx.Input,
			hd.EthSecp256k1Option(),
		)
		if err != nil {
			return err
		}

		e.clientCtx.Keyring = keyring
	}

	return nil
}

// ProtocolVersion returns the supported Ethereum protocol version.
func (e *PublicEthAPI) ProtocolVersion() hexutil.Uint {
	e.logger.Debug("eth_protocolVersion")
	return hexutil.Uint(version.ProtocolVersion)
}

// ChainId returns the chain's identifier in hex format
func (e *PublicEthAPI) ChainId() (hexutil.Uint, error) { // nolint
	e.logger.Debug("eth_chainId")
	return hexutil.Uint(uint(e.chainIDEpoch.Uint64())), nil
}

// Syncing returns whether or not the current node is syncing with other peers. Returns false if not, or a struct
// outlining the state of the sync if it is.
func (e *PublicEthAPI) Syncing() (interface{}, error) {
	e.logger.Debug("eth_syncing")

	status, err := e.clientCtx.Client.Status(e.ctx)
	if err != nil {
		return false, err
	}

	if !status.SyncInfo.CatchingUp {
		return false, nil
	}

	return map[string]interface{}{
		// "startingBlock": nil, // NA
		"currentBlock": hexutil.Uint64(status.SyncInfo.LatestBlockHeight),
		// "highestBlock":  nil, // NA
		// "pulledStates":  nil, // NA
		// "knownStates":   nil, // NA
	}, nil
}

// Coinbase is the address that staking rewards will be send to (alias for Etherbase).
func (e *PublicEthAPI) Coinbase() (common.Address, error) {
	e.logger.Debug("eth_coinbase")

	node, err := e.clientCtx.GetNode()
	if err != nil {
		return common.Address{}, err
	}

	status, err := node.Status(e.ctx)
	if err != nil {
		return common.Address{}, err
	}

	return common.BytesToAddress(status.ValidatorInfo.Address.Bytes()), nil
}

// Mining returns whether or not this node is currently mining. Always false.
func (e *PublicEthAPI) Mining() bool {
	e.logger.Debug("eth_mining")
	return false
}

// Hashrate returns the current node's hashrate. Always 0.
func (e *PublicEthAPI) Hashrate() hexutil.Uint64 {
	e.logger.Debug("eth_hashrate")
	return 0
}

// GasPrice returns the current gas price based on Ethermint's gas price oracle.
func (e *PublicEthAPI) GasPrice() *hexutil.Big {
	e.logger.Debug("eth_gasPrice")
	out := big.NewInt(0)
	return (*hexutil.Big)(out)
}

// Accounts returns the list of accounts available to this node.
func (e *PublicEthAPI) Accounts() ([]common.Address, error) {
	e.logger.Debug("eth_accounts")
	e.keyringLock.Lock()

	addresses := make([]common.Address, 0) // return [] instead of nil if empty

	infos, err := e.clientCtx.Keyring.List()
	if err != nil {
		return addresses, err
	}

	e.keyringLock.Unlock()

	for _, info := range infos {
		addressBytes := info.GetPubKey().Address().Bytes()
		addresses = append(addresses, common.BytesToAddress(addressBytes))
	}

	return addresses, nil
}

// BlockNumber returns the current block number.
func (e *PublicEthAPI) BlockNumber() (hexutil.Uint64, error) {
	e.logger.Debug("eth_blockNumber")
	return e.backend.BlockNumber()
}

// GetBalance returns the provided account's balance up to the provided block number.
func (e *PublicEthAPI) GetBalance(address common.Address, blockNum BlockNumber) (*hexutil.Big, error) { // nolint: interfacer
	e.logger.Debug("eth_getBalance", "address", address, "block number", blockNum)

	req := &evmtypes.QueryBalanceRequest{
		Address: address.String(),
	}

	res, err := e.queryClient.Balance(ContextWithHeight(blockNum.Int64()), req)
	if err != nil {
		return nil, err
	}

	val, err := ethermint.UnmarshalBigInt(res.Balance)
	if err != nil {
		return nil, err
	}

	return (*hexutil.Big)(val), nil
}

// GetStorageAt returns the contract storage at the given address, block number, and key.
func (e *PublicEthAPI) GetStorageAt(address common.Address, key string, blockNum BlockNumber) (hexutil.Bytes, error) { // nolint: interfacer
	e.logger.Debug("eth_getStorageAt", "address", address, "key", key, "block number", blockNum)

	req := &evmtypes.QueryStorageRequest{
		Address: address.String(),
		Key:     key,
	}

	res, err := e.queryClient.Storage(ContextWithHeight(blockNum.Int64()), req)
	if err != nil {
		return nil, err
	}

	value := common.HexToHash(res.Value)
	return value.Bytes(), nil
}

// GetTransactionCount returns the number of transactions at the given address up to the given block number.
func (e *PublicEthAPI) GetTransactionCount(address common.Address, blockNum BlockNumber) (*hexutil.Uint64, error) {
	e.logger.Debug("eth_getTransactionCount", "address", address, "block number", blockNum)

	// Get nonce (sequence) from account
	from := sdk.AccAddress(address.Bytes())
	accRet := e.clientCtx.AccountRetriever

	err := accRet.EnsureExists(e.clientCtx, from)
	if err != nil {
		// account doesn't exist yet, return 0
		n := hexutil.Uint64(0)
		return &n, nil
	}

	_, nonce, err := accRet.GetAccountNumberSequence(e.clientCtx, from)
	if err != nil {
		return nil, err
	}

	n := hexutil.Uint64(nonce)
	return &n, nil
}

// GetBlockTransactionCountByHash returns the number of transactions in the block identified by hash.
func (e *PublicEthAPI) GetBlockTransactionCountByHash(hash common.Hash) *hexutil.Uint {
	e.logger.Debug("eth_getBlockTransactionCountByHash", "hash", hash)

	resBlock, err := e.clientCtx.Client.BlockByHash(e.ctx, hash.Bytes())
	if err != nil {
		return nil
	}

	n := hexutil.Uint(len(resBlock.Block.Txs))
	return &n
}

// GetBlockTransactionCountByNumber returns the number of transactions in the block identified by number.
func (e *PublicEthAPI) GetBlockTransactionCountByNumber(blockNum BlockNumber) *hexutil.Uint {
	e.logger.Debug("eth_getBlockTransactionCountByNumber", "block number", blockNum)
	resBlock, err := e.clientCtx.Client.Block(e.ctx, blockNum.TmHeight())
	if err != nil {
		return nil
	}

	n := hexutil.Uint(len(resBlock.Block.Txs))
	return &n
}

// GetUncleCountByBlockHash returns the number of uncles in the block identified by hash. Always zero.
func (e *PublicEthAPI) GetUncleCountByBlockHash(hash common.Hash) hexutil.Uint {
	return 0
}

// GetUncleCountByBlockNumber returns the number of uncles in the block identified by number. Always zero.
func (e *PublicEthAPI) GetUncleCountByBlockNumber(blockNum BlockNumber) hexutil.Uint {
	return 0
}

// GetCode returns the contract code at the given address and block number.
func (e *PublicEthAPI) GetCode(address common.Address, blockNumber BlockNumber) (hexutil.Bytes, error) { // nolint: interfacer
	e.logger.Debug("eth_getCode", "address", address, "block number", blockNumber)

	req := &evmtypes.QueryCodeRequest{
		Address: address.String(),
	}

	res, err := e.queryClient.Code(ContextWithHeight(blockNumber.Int64()), req)
	if err != nil {
		return nil, err
	}

	return res.Code, nil
}

// GetTransactionLogs returns the logs given a transaction hash.
func (e *PublicEthAPI) GetTransactionLogs(txHash common.Hash) ([]*ethtypes.Log, error) {
	e.logger.Debug("eth_getTransactionLogs", "hash", txHash)
	return e.backend.GetTransactionLogs(txHash)
}

// Sign signs the provided data using the private key of address via Geth's signature standard.
func (e *PublicEthAPI) Sign(address common.Address, data hexutil.Bytes) (hexutil.Bytes, error) {
	e.logger.Debug("eth_sign", "address", address, "data", data)
	// TODO: Change this functionality to find an unlocked account by address

	key, exist := checkKeyInKeyring(e.keys, address)
	if !exist {
		return nil, keystore.ErrLocked
	}

	// Sign the requested hash with the wallet
	signature, err := key.Sign(data)
	if err != nil {
		return nil, err
	}

	signature[64] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	return signature, nil
}

// SendTransaction sends an Ethereum transaction.
func (e *PublicEthAPI) SendTransaction(args SendTxArgs) (common.Hash, error) {
	e.logger.Debug("eth_sendTransaction", "args", args)
	// TODO: Change this functionality to find an unlocked account by address

	for _, key := range e.keys {
		e.logger.Debug("eth_sendTransaction", "key", fmt.Sprintf("0x%x", key.PubKey().Address().Bytes()))
	}

	key, exist := checkKeyInKeyring(e.keys, args.From)
	if !exist {
		e.logger.Debug("failed to find key in keyring", "key", args.From)
		return common.Hash{}, keystore.ErrLocked
	}

	// Mutex lock the address' nonce to avoid assigning it to multiple requests
	if args.Nonce == nil {
		e.nonceLock.LockAddr(args.From)
		defer e.nonceLock.UnlockAddr(args.From)
	}

	// Assemble transaction from fields
	tx, err := e.generateFromArgs(args)
	if err != nil {
		e.logger.Debug("failed to generate tx", "error", err)
		return common.Hash{}, err
	}

	// ChainID must be set as flag to send transaction
	chainID := viper.GetString(flags.FlagChainID)
	// parse the chainID from a string to a base-10 integer
	chainIDEpoch, err := ethermint.ParseChainID(chainID)
	if err != nil {
		return common.Hash{}, err
	}

	// Sign transaction
	if err := tx.Sign(chainIDEpoch, key.ToECDSA()); err != nil {
		e.logger.Debug("failed to sign tx", "error", err)
		return common.Hash{}, err
	}

	// Encode transaction by default Tx encoder
	txEncoder := e.clientCtx.TxConfig.TxEncoder()
	txBytes, err := txEncoder(tx)
	if err != nil {
		return common.Hash{}, err
	}

	// Broadcast transaction in sync mode (default)
	res, err := e.clientCtx.BroadcastTx(txBytes)
	// If error is encountered on the node, the broadcast will not return an error
	if err != nil {
		return common.Hash{}, err
	}

	// Return transaction hash
	return common.HexToHash(res.TxHash), nil
}

// SendRawTransaction send a raw Ethereum transaction.
func (e *PublicEthAPI) SendRawTransaction(data hexutil.Bytes) (common.Hash, error) {
	e.logger.Debug("eth_sendRawTransaction", "data", data)
	tx := new(evmtypes.MsgEthereumTx)

	// RLP decode raw transaction bytes
	if err := rlp.DecodeBytes(data, tx); err != nil {
		// Return nil is for when gasLimit overflows uint64
		return common.Hash{}, nil
	}

	// Encode transaction by default Tx encoder
	txBytes, err := e.clientCtx.TxConfig.TxEncoder()(tx)
	if err != nil {
		return common.Hash{}, err
	}

	// TODO: Possibly log the contract creation address (if recipient address is nil) or tx data
	// If error is encountered on the node, the broadcast will not return an error
	res, err := e.clientCtx.BroadcastTx(txBytes)
	if err != nil {
		return common.Hash{}, err
	}

	// Return transaction hash
	return common.HexToHash(res.TxHash), nil
}

// Call performs a raw contract call.
func (e *PublicEthAPI) Call(args CallArgs, blockNr BlockNumber, _ *map[common.Address]account) (hexutil.Bytes, error) {
	e.logger.Debug("eth_call", "args", args, "block number", blockNr)
	simRes, err := e.doCall(args, blockNr, big.NewInt(ethermint.DefaultRPCGasLimit))
	if err != nil {
		return []byte{}, err
	}

	data, err := evmtypes.DecodeResultData(simRes.Result.Data)
	if err != nil {
		return []byte{}, err
	}

	return (hexutil.Bytes)(data.Ret), nil
}

// DoCall performs a simulated call operation through the evmtypes. It returns the
// estimated gas used on the operation or an error if fails.
func (e *PublicEthAPI) doCall(
	args CallArgs, blockNr BlockNumber, globalGasCap *big.Int,
) (*simulate.SimulateResponse, error) {
	// Set sender address or use a default if none specified
	var addr common.Address

	if args.From == nil {
		addrs, err := e.Accounts()
		if err == nil && len(addrs) > 0 {
			addr = addrs[0]
		}
	} else {
		addr = *args.From
	}

	// Set default gas & gas price if none were set
	// Change this to uint64(math.MaxUint64 / 2) if gas cap can be configured
	gas := uint64(ethermint.DefaultRPCGasLimit)
	if args.Gas != nil {
		gas = uint64(*args.Gas)
	}
	if globalGasCap != nil && globalGasCap.Uint64() < gas {
		e.logger.Debug("Caller gas above allowance, capping", "requested", gas, "cap", globalGasCap)
		gas = globalGasCap.Uint64()
	}

	// Set gas price using default or parameter if passed in
	gasPrice := new(big.Int).SetUint64(ethermint.DefaultGasPrice)
	if args.GasPrice != nil {
		gasPrice = args.GasPrice.ToInt()
	}

	// Set value for transaction
	value := new(big.Int)
	if args.Value != nil {
		value = args.Value.ToInt()
	}

	// Set Data if provided
	var data []byte
	if args.Data != nil {
		data = []byte(*args.Data)
	}

	// Set destination address for call
	var fromAddr sdk.AccAddress
	if args.From != nil {
		fromAddr = sdk.AccAddress(args.From.Bytes())
	}

	accNum, seq, err := e.clientCtx.AccountRetriever.GetAccountNumberSequence(e.clientCtx, fromAddr)
	if err != nil {
		return nil, err
	}

	// Create new call message
	msg := evmtypes.NewMsgEthereumTx(seq, args.To, value, gas, gasPrice, data)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	privKey, exists := checkKeyInKeyring(e.keys, addr)
	if !exists {
		return nil, fmt.Errorf("account with address %s does not exist in keyring", addr.String())
	}

	fees := sdk.NewCoins(ethermint.NewPhotonCoin(sdk.NewIntFromBigInt(msg.Fee())))
	signMode := e.clientCtx.TxConfig.SignModeHandler().DefaultMode()
	signerData := authsigning.SignerData{ChainID: e.clientCtx.ChainID, AccountNumber: accNum, Sequence: seq}

	// Create a TxBuilder
	txBuilder := e.clientCtx.TxConfig.NewTxBuilder()
	if err := txBuilder.SetMsgs(msg); err != nil {
		return nil, err
	}
	txBuilder.SetFeeAmount(fees)
	txBuilder.SetGasLimit(gas)

	// TODO: use tx.Factory

	// sign with the private key
	sigV2, err := tx.SignWithPrivKey(
		signMode, signerData,
		txBuilder, privKey, e.clientCtx.TxConfig, seq,
	)

	if err != nil {
		return nil, err
	}

	if err := txBuilder.SetSignatures(sigV2); err != nil {
		return nil, err
	}

	tx, ok := txBuilder.(codectypes.IntoAny).AsAny().GetCachedValue().(*txtypes.Tx)
	if !ok {
		return nil, errors.New("cannot cast to tx")
	}

	req := &simulate.SimulateRequest{
		Tx: tx,
	}

	simResponse, err := e.queryClient.Simulate(ContextWithHeight(blockNr.Int64()), req)
	if err != nil {
		return nil, err
	}

	return simResponse, nil
}

// EstimateGas returns an estimate of gas usage for the given smart contract call.
// It adds 1,000 gas to the returned value instead of using the gas adjustment
// param from the SDK.
func (e *PublicEthAPI) EstimateGas(args CallArgs) (hexutil.Uint64, error) {
	e.logger.Debug("eth_estimateGas", "args", args)
	simResponse, err := e.doCall(args, 0, big.NewInt(ethermint.DefaultRPCGasLimit))
	if err != nil {
		return 0, err
	}

	// TODO: change 1000 buffer for more accurate buffer (eg: SDK's gasAdjusted)
	estimatedGas := simResponse.GasInfo.GasUsed
	gas := estimatedGas + 1000

	return hexutil.Uint64(gas), nil
}

// GetBlockByHash returns the block identified by hash.
func (e *PublicEthAPI) GetBlockByHash(hash common.Hash, fullTx bool) (map[string]interface{}, error) {
	e.logger.Debug("eth_getBlockByHash", "hash", hash, "full", fullTx)
	return e.backend.GetBlockByHash(hash, fullTx)
}

// GetBlockByNumber returns the block identified by number.
func (e *PublicEthAPI) GetBlockByNumber(blockNum BlockNumber, fullTx bool) (map[string]interface{}, error) {
	e.logger.Debug("eth_getBlockByNumber", "number", blockNum, "full", fullTx)
	return e.backend.GetBlockByNumber(blockNum, fullTx)
}

// GetTransactionByHash returns the transaction identified by hash.
func (e *PublicEthAPI) GetTransactionByHash(hash common.Hash) (*Transaction, error) {
	e.logger.Debug("eth_getTransactionByHash", "hash", hash)
	tx, err := e.clientCtx.Client.Tx(e.ctx, hash.Bytes(), false)
	if err != nil {
		// Return nil for transaction when not found
		return nil, nil
	}

	// Can either cache or just leave this out if not necessary
	block, err := e.clientCtx.Client.Block(e.ctx, &tx.Height)
	if err != nil {
		return nil, err
	}
	blockHash := common.BytesToHash(block.Block.Header.Hash())

	ethTx, err := RawTxToEthTx(e.clientCtx, tx.Tx)
	if err != nil {
		return nil, err
	}

	height := uint64(tx.Height)
	return NewTransaction(ethTx, common.BytesToHash(tx.Tx.Hash()), blockHash, height, uint64(tx.Index))
}

// GetTransactionByBlockHashAndIndex returns the transaction identified by hash and index.
func (e *PublicEthAPI) GetTransactionByBlockHashAndIndex(hash common.Hash, idx hexutil.Uint) (*Transaction, error) {
	e.logger.Debug("eth_getTransactionByHashAndIndex", "hash", hash, "index", idx)

	resBlock, err := e.clientCtx.Client.BlockByHash(e.ctx, hash.Bytes())
	if err != nil {
		return nil, err
	}
	return e.getTransactionByBlockAndIndex(resBlock.Block, idx)
}

// GetTransactionByBlockNumberAndIndex returns the transaction identified by number and index.
func (e *PublicEthAPI) GetTransactionByBlockNumberAndIndex(blockNum BlockNumber, idx hexutil.Uint) (*Transaction, error) {
	e.logger.Debug("eth_getTransactionByBlockNumberAndIndex", "number", blockNum, "index", idx)
	height := blockNum.Int64()

	resBlock, err := e.clientCtx.Client.Block(e.ctx, &height)
	if err != nil {
		return nil, err
	}

	return e.getTransactionByBlockAndIndex(resBlock.Block, idx)
}

func (e *PublicEthAPI) getTransactionByBlockAndIndex(block *tmtypes.Block, idx hexutil.Uint) (*Transaction, error) {
	// return if index out of bounds
	if uint64(idx) >= uint64(len(block.Txs)) {
		return nil, nil
	}

	ethTx, err := RawTxToEthTx(e.clientCtx, block.Txs[idx])
	if err != nil {
		// return nil error if the transaction is not a MsgEthereumTx
		return nil, nil
	}

	height := uint64(block.Height)
	txHash := common.BytesToHash(block.Txs[idx].Hash())
	blockHash := common.BytesToHash(block.Header.Hash())
	return NewTransaction(ethTx, txHash, blockHash, height, uint64(idx))
}

// GetTransactionReceipt returns the transaction receipt identified by hash.
func (e *PublicEthAPI) GetTransactionReceipt(hash common.Hash) (map[string]interface{}, error) {
	e.logger.Debug("eth_getTransactionReceipt", "hash", hash)
	tx, err := e.clientCtx.Client.Tx(e.ctx, hash.Bytes(), false)
	if err != nil {
		// Return nil for transaction when not found
		return nil, nil
	}

	// Query block for consensus hash
	block, err := e.clientCtx.Client.Block(e.ctx, &tx.Height)
	if err != nil {
		return nil, err
	}

	blockHash := common.BytesToHash(block.Block.Header.Hash())

	// Convert tx bytes to eth transaction
	ethTx, err := RawTxToEthTx(e.clientCtx, tx.Tx)
	if err != nil {
		return nil, err
	}

	from, err := ethTx.VerifySig(ethTx.ChainID())
	if err != nil {
		return nil, err
	}

	// Set status codes based on tx result
	var status hexutil.Uint
	if tx.TxResult.IsOK() {
		status = hexutil.Uint(1)
	} else {
		status = hexutil.Uint(0)
	}

	txData := tx.TxResult.GetData()

	data, err := evmtypes.DecodeResultData(txData)
	if err != nil {
		status = 0 // transaction failed
	}

	if len(data.TxLogs.Logs) == 0 {
		data.TxLogs.Logs = []*evmtypes.Log{}
	}

	receipt := map[string]interface{}{
		// Consensus fields: These fields are defined by the Yellow Paper
		"status":            status,
		"cumulativeGasUsed": nil, // ignore until needed
		"logsBloom":         data.Bloom,
		"logs":              data.TxLogs.EthLogs(),

		// Implementation fields: These fields are added by geth when processing a transaction.
		// They are stored in the chain database.
		"transactionHash": hash,
		"contractAddress": data.ContractAddress,
		"gasUsed":         hexutil.Uint64(tx.TxResult.GasUsed),

		// Inclusion information: These fields provide information about the inclusion of the
		// transaction corresponding to this receipt.
		"blockHash":        blockHash,
		"blockNumber":      hexutil.Uint64(tx.Height),
		"transactionIndex": hexutil.Uint64(tx.Index),

		// sender and receiver (contract or EOA) addreses
		"from": from,
		"to":   ethTx.To(),
	}

	return receipt, nil
}

// PendingTransactions returns the transactions that are in the transaction pool
// and have a from address that is one of the accounts this node manages.
func (e *PublicEthAPI) PendingTransactions() ([]*Transaction, error) {
	e.logger.Debug("eth_getPendingTransactions")
	return e.backend.PendingTransactions()
}

// GetUncleByBlockHashAndIndex returns the uncle identified by hash and index. Always returns nil.
func (e *PublicEthAPI) GetUncleByBlockHashAndIndex(hash common.Hash, idx hexutil.Uint) map[string]interface{} {
	return nil
}

// GetUncleByBlockNumberAndIndex returns the uncle identified by number and index. Always returns nil.
func (e *PublicEthAPI) GetUncleByBlockNumberAndIndex(number hexutil.Uint, idx hexutil.Uint) map[string]interface{} {
	return nil
}

// GetProof returns an account object with proof and any storage proofs
func (e *PublicEthAPI) GetProof(address common.Address, storageKeys []string, blockNumber BlockNumber) (*AccountResult, error) {
	height := blockNumber.Int64()
	e.logger.Debug("eth_getProof", "address", address, "keys", storageKeys, "number", height)

	ctx := ContextWithHeight(height)
	clientCtx := e.clientCtx.WithHeight(height)

	// query storage proofs
	storageProofs := make([]StorageResult, len(storageKeys))
	for i, key := range storageKeys {
		hexKey := common.HexToHash(key)
		valueBz, proof, err := e.queryClient.GetProof(clientCtx, evmtypes.StoreKey, evmtypes.StateKey(address, hexKey.Bytes()))
		if err != nil {
			return nil, err
		}

		// check for proof
		var proofStr string
		if proof != nil {
			proofStr = proof.String()
		}

		storageProofs[i] = StorageResult{
			Key:   key,
			Value: (*hexutil.Big)(new(big.Int).SetBytes(valueBz)),
			Proof: []string{proofStr},
		}
	}

	// query EVM account
	req := &evmtypes.QueryAccountRequest{
		Address: address.String(),
	}

	res, err := e.queryClient.Account(ctx, req)
	if err != nil {
		return nil, err
	}

	// query account proofs
	accountKey := authtypes.AddressStoreKey(sdk.AccAddress(address.Bytes()))
	_, proof, err := e.queryClient.GetProof(clientCtx, authtypes.StoreKey, accountKey)
	if err != nil {
		return nil, err
	}

	// check for proof
	var accProofStr string
	if proof != nil {
		accProofStr = proof.String()
	}

	balance, err := ethermint.UnmarshalBigInt(res.Balance)
	if err != nil {
		return nil, err
	}

	return &AccountResult{
		Address:      address,
		AccountProof: []string{accProofStr},
		Balance:      (*hexutil.Big)(balance),
		CodeHash:     common.BytesToHash(res.CodeHash),
		Nonce:        hexutil.Uint64(res.Nonce),
		StorageHash:  common.Hash{}, // NOTE: Ethermint doesn't have a storage hash. TODO: implement?
		StorageProof: storageProofs,
	}, nil
}

// generateFromArgs populates tx message with args (used in RPC API)
func (e *PublicEthAPI) generateFromArgs(args SendTxArgs) (*evmtypes.MsgEthereumTx, error) {
	var (
		nonce    uint64
		gasLimit uint64
		err      error
	)

	amount := args.Value.ToInt()
	gasPrice := args.GasPrice.ToInt()

	if args.GasPrice == nil {

		// Set default gas price
		// TODO: Change to min gas price from context once available through server/daemon
		gasPrice = big.NewInt(ethermint.DefaultGasPrice)
	}

	if args.Nonce == nil {
		// Get nonce (sequence) from account
		from := sdk.AccAddress(args.From.Bytes())
		accRet := e.clientCtx.AccountRetriever

		if e.clientCtx.Keyring == nil {
			return nil, fmt.Errorf("clientCtx.Keyring is nil")
		}

		err = accRet.EnsureExists(e.clientCtx, from)
		if err != nil {
			return nil, fmt.Errorf("nonexistent account %s: %s", args.From.String(), err)
		}

		_, nonce, err = accRet.GetAccountNumberSequence(e.clientCtx, from)
		if err != nil {
			return nil, err
		}
	} else {
		nonce = uint64(*args.Nonce)
	}

	if args.Data != nil && args.Input != nil && !bytes.Equal(*args.Data, *args.Input) {
		return nil, errors.New(`both "data" and "input" are set and not equal. Please use "input" to pass transaction call data`)
	}

	// Sets input to either Input or Data, if both are set and not equal error above returns
	var input []byte
	if args.Input != nil {
		input = *args.Input
	} else if args.Data != nil {
		input = *args.Data
	}

	if args.To == nil && len(input) == 0 {
		// Contract creation
		return nil, fmt.Errorf("contract creation without any data provided")
	}

	if args.Gas == nil {
		callArgs := CallArgs{
			From:     &args.From,
			To:       args.To,
			Gas:      args.Gas,
			GasPrice: args.GasPrice,
			Value:    args.Value,
			Data:     args.Data,
		}
		gl, err := e.EstimateGas(callArgs)
		if err != nil {
			return nil, err
		}
		gasLimit = uint64(gl)
	} else {
		gasLimit = (uint64)(*args.Gas)
	}

	msg := evmtypes.NewMsgEthereumTx(nonce, args.To, amount, gasLimit, gasPrice, input)

	return msg, nil
}
