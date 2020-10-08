package rpc

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"os"
	"sync"

	"github.com/spf13/viper"

	"github.com/cosmos/ethermint/crypto/ethsecp256k1"
	"github.com/cosmos/ethermint/crypto/hd"
	params "github.com/cosmos/ethermint/rpc/args"
	ethermint "github.com/cosmos/ethermint/types"
	"github.com/cosmos/ethermint/utils"
	"github.com/cosmos/ethermint/version"
	evmtypes "github.com/cosmos/ethermint/x/evm/types"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/merkle"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/rpc/client"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// PublicEthAPI is the eth_ prefixed set of APIs in the Web3 JSON-RPC spec.
type PublicEthAPI struct {
	clientCtx    client.Context
	chainIDEpoch *big.Int
	logger       log.Logger
	backend      Backend
	keys         []ethsecp256k1.PrivKey // unlocked keys
	nonceLock    *AddrLocker
	keybaseLock  sync.Mutex
}

// NewPublicEthAPI creates an instance of the public ETH Web3 API.
func NewPublicEthAPI(clientCtx client.Context, backend Backend, nonceLock *AddrLocker,
	key []ethsecp256k1.PrivKey) *PublicEthAPI {

	epoch, err := ethermint.ParseChainID(clientCtx.ChainID)
	if err != nil {
		panic(err)
	}

	api := &PublicEthAPI{
		clientCtx:    clientCtx,
		chainIDEpoch: epoch,
		logger:       log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "json-rpc"),
		backend:      backend,
		keys:         key,
		nonceLock:    nonceLock,
	}

	if err := api.getKeybaseInfo(); err != nil {
		api.logger.Error("failed to get keybase info", "error", err)
	}

	return api
}

func (e *PublicEthAPI) getKeybaseInfo() error {
	e.keybaseLock.Lock()
	defer e.keybaseLock.Unlock()

	if e.clientCtx.Keybase == nil {
		keybase, err := keyring.New(
			sdk.KeyringServiceName(),
			viper.GetString(flags.FlagKeyringBackend),
			viper.GetString(flags.FlagHome),
			e.clientCtx.Input,
			hd.EthSecp256k1Option(),
		)
		if err != nil {
			return err
		}

		e.clientCtx.Keybase = keybase
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

	status, err := e.clientCtx.Client.Status()
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

	status, err := node.Status()
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
	e.keybaseLock.Lock()

	addresses := make([]common.Address, 0) // return [] instead of nil if empty

	keybase, err := keyring.New(
		sdk.KeyringServiceName(),
		viper.GetString(flags.FlagKeyringBackend),
		viper.GetString(flags.FlagHome),
		e.clientCtx.Input,
		hd.EthSecp256k1Option(),
	)
	if err != nil {
		return addresses, err
	}

	infos, err := keybase.List()
	if err != nil {
		return addresses, err
	}

	e.keybaseLock.Unlock()

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
func (e *PublicEthAPI) GetBalance(address common.Address, blockNum BlockNumber) (*hexutil.Big, error) {
	e.logger.Debug("eth_getBalance", "address", address, "block number", blockNum)
	ctx := e.clientCtx.WithHeight(blockNum.Int64())
	res, _, err := ctx.QueryWithData(fmt.Sprintf("custom/%s/balance/%s", evmtypes.ModuleName, address.Hex()), nil)
	if err != nil {
		return nil, err
	}

	var out evmtypes.QueryResBalance
	e.clientCtx.Codec.MustUnmarshalJSON(res, &out)
	val, err := utils.UnmarshalBigInt(out.Balance)
	if err != nil {
		return nil, err
	}

	return (*hexutil.Big)(val), nil
}

// GetStorageAt returns the contract storage at the given address, block number, and key.
func (e *PublicEthAPI) GetStorageAt(address common.Address, key string, blockNum BlockNumber) (hexutil.Bytes, error) {
	e.logger.Debug("eth_getStorageAt", "address", address, "key", key, "block number", blockNum)
	ctx := e.clientCtx.WithHeight(blockNum.Int64())
	res, _, err := ctx.QueryWithData(fmt.Sprintf("custom/%s/storage/%s/%s", evmtypes.ModuleName, address.Hex(), key), nil)
	if err != nil {
		return nil, err
	}

	var out evmtypes.QueryResStorage
	e.clientCtx.Codec.MustUnmarshalJSON(res, &out)
	return out.Value, nil
}

// GetTransactionCount returns the number of transactions at the given address up to the given block number.
func (e *PublicEthAPI) GetTransactionCount(address common.Address, blockNum BlockNumber) (*hexutil.Uint64, error) {
	e.logger.Debug("eth_getTransactionCount", "address", address, "block number", blockNum)
	ctx := e.clientCtx.WithHeight(blockNum.Int64())

	// Get nonce (sequence) from account
	from := sdk.AccAddress(address.Bytes())
	accRet := authtypes.NewAccountRetriever(ctx)

	err := accRet.EnsureExists(from)
	if err != nil {
		// account doesn't exist yet, return 0
		n := hexutil.Uint64(0)
		return &n, nil
	}

	_, nonce, err := accRet.GetAccountNumberSequence(from)
	if err != nil {
		return nil, err
	}

	n := hexutil.Uint64(nonce)
	return &n, nil
}

// GetBlockTransactionCountByHash returns the number of transactions in the block identified by hash.
func (e *PublicEthAPI) GetBlockTransactionCountByHash(hash common.Hash) *hexutil.Uint {
	e.logger.Debug("eth_getBlockTransactionCountByHash", "hash", hash)
	res, _, err := e.clientCtx.Query(fmt.Sprintf("custom/%s/%s/%s", evmtypes.ModuleName, evmtypes.QueryHashToHeight, hash.Hex()))
	if err != nil {
		// Return nil if block does not exist
		return nil
	}

	var out evmtypes.QueryResBlockNumber
	e.clientCtx.Codec.MustUnmarshalJSON(res, &out)
	return e.getBlockTransactionCountByNumber(out.Number)
}

// GetBlockTransactionCountByNumber returns the number of transactions in the block identified by number.
func (e *PublicEthAPI) GetBlockTransactionCountByNumber(blockNum BlockNumber) *hexutil.Uint {
	e.logger.Debug("eth_getBlockTransactionCountByNumber", "block number", blockNum)
	height := blockNum.Int64()
	return e.getBlockTransactionCountByNumber(height)
}

func (e *PublicEthAPI) getBlockTransactionCountByNumber(number int64) *hexutil.Uint {
	block, err := e.clientCtx.Client.Block(&number)
	if err != nil {
		// Return nil if block doesn't exist
		return nil
	}

	n := hexutil.Uint(len(block.Block.Txs))
	return &n
}

// GetUncleCountByBlockHash returns the number of uncles in the block idenfied by hash. Always zero.
func (e *PublicEthAPI) GetUncleCountByBlockHash(hash common.Hash) hexutil.Uint {
	return 0
}

// GetUncleCountByBlockNumber returns the number of uncles in the block idenfied by number. Always zero.
func (e *PublicEthAPI) GetUncleCountByBlockNumber(blockNum BlockNumber) hexutil.Uint {
	return 0
}

// GetCode returns the contract code at the given address and block number.
func (e *PublicEthAPI) GetCode(address common.Address, blockNumber BlockNumber) (hexutil.Bytes, error) {
	e.logger.Debug("eth_getCode", "address", address, "block number", blockNumber)
	ctx := e.clientCtx.WithHeight(blockNumber.Int64())
	res, _, err := ctx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", evmtypes.ModuleName, evmtypes.QueryCode, address.Hex()), nil)
	if err != nil {
		return nil, err
	}

	var out evmtypes.QueryResCode
	e.clientCtx.Codec.MustUnmarshalJSON(res, &out)
	return out.Code, nil
}

// GetTransactionLogs returns the logs given a transaction hash.
func (e *PublicEthAPI) GetTransactionLogs(txHash common.Hash) ([]*ethtypes.Log, error) {
	e.logger.Debug("eth_getTransactionLogs", "hash", txHash)
	return e.backend.GetTransactionLogs(txHash)
}

// ExportAccount exports an account's balance, code, and storage at the given block number
// TODO: deprecate this once the export genesis command works
func (e *PublicEthAPI) ExportAccount(address common.Address, blockNumber BlockNumber) (string, error) {
	e.logger.Debug("eth_exportAccount", "address", address, "block number", blockNumber)
	ctx := e.clientCtx.WithHeight(blockNumber.Int64())

	res, _, err := ctx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", evmtypes.ModuleName, evmtypes.QueryExportAccount, address.Hex()), nil)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func checkKeyInKeyring(keys []ethsecp256k1.PrivKey, address common.Address) (key ethsecp256k1.PrivKey, exist bool) {
	if len(keys) > 0 {
		for _, key := range keys {
			if bytes.Equal(key.PubKey().Address().Bytes(), address.Bytes()) {
				return key, true
			}
		}
	}
	return nil, false
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
	if err == nil {
		signature[64] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	}

	return signature, err
}

// SendTransaction sends an Ethereum transaction.
func (e *PublicEthAPI) SendTransaction(args params.SendTxArgs) (common.Hash, error) {
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
	txEncoder := authclient.GetTxEncoder(e.clientCtx.Codec)
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
	txEncoder := authclient.GetTxEncoder(e.clientCtx.Codec)
	txBytes, err := txEncoder(tx)
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

// CallArgs represents the arguments for a call.
type CallArgs struct {
	From     *common.Address `json:"from"`
	To       *common.Address `json:"to"`
	Gas      *hexutil.Uint64 `json:"gas"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
	Value    *hexutil.Big    `json:"value"`
	Data     *hexutil.Bytes  `json:"data"`
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

// account indicates the overriding fields of account during the execution of
// a message call.
// Note, state and stateDiff can't be specified at the same time. If state is
// set, message execution will only use the data in the given state. Otherwise
// if statDiff is set, all diff will be applied first and then execute the call
// message.
type account struct {
	Nonce     *hexutil.Uint64              `json:"nonce"`
	Code      *hexutil.Bytes               `json:"code"`
	Balance   **hexutil.Big                `json:"balance"`
	State     *map[common.Hash]common.Hash `json:"state"`
	StateDiff *map[common.Hash]common.Hash `json:"stateDiff"`
}

// DoCall performs a simulated call operation through the evmtypes. It returns the
// estimated gas used on the operation or an error if fails.
func (e *PublicEthAPI) doCall(
	args CallArgs, blockNr BlockNumber, globalGasCap *big.Int,
) (*sdk.SimulationResponse, error) {
	// Set height for historical queries
	ctx := e.clientCtx

	if blockNr.Int64() != 0 {
		ctx = e.clientCtx.WithHeight(blockNr.Int64())
	}

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
	var toAddr sdk.AccAddress
	if args.To != nil {
		toAddr = sdk.AccAddress(args.To.Bytes())
	}

	// Create new call message
	msg := evmtypes.NewMsgEthermint(0, &toAddr, sdk.NewIntFromBigInt(value), gas,
		sdk.NewIntFromBigInt(gasPrice), data, sdk.AccAddress(addr.Bytes()))

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Generate tx to be used to simulate (signature isn't needed)
	var stdSig authtypes.StdSignature
	tx := authtypes.NewStdTx([]sdk.Msg{msg}, authtypes.StdFee{}, []authtypes.StdSignature{stdSig}, "")

	txEncoder := authclient.GetTxEncoder(ctx.Codec)
	txBytes, err := txEncoder(tx)
	if err != nil {
		return nil, err
	}

	// Transaction simulation through query
	res, _, err := ctx.QueryWithData("app/simulate", txBytes)
	if err != nil {
		return nil, err
	}

	var simResponse sdk.SimulationResponse
	if err := ctx.Codec.UnmarshalBinaryBare(res, &simResponse); err != nil {
		return nil, err
	}

	return &simResponse, nil
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

func formatBlock(
	header tmtypes.Header, size int, gasLimit int64,
	gasUsed *big.Int, transactions interface{}, bloom ethtypes.Bloom,
) map[string]interface{} {
	if bytes.Equal(header.DataHash, []byte{}) {
		header.DataHash = tmbytes.HexBytes(common.Hash{}.Bytes())
	}

	return map[string]interface{}{
		"number":           hexutil.Uint64(header.Height),
		"hash":             hexutil.Bytes(header.Hash()),
		"parentHash":       hexutil.Bytes(header.LastBlockID.Hash),
		"nonce":            hexutil.Uint64(0), // PoW specific
		"sha3Uncles":       common.Hash{},     // No uncles in Tendermint
		"logsBloom":        bloom,
		"transactionsRoot": hexutil.Bytes(header.DataHash),
		"stateRoot":        hexutil.Bytes(header.AppHash),
		"miner":            common.Address{},
		"mixHash":          common.Hash{},
		"difficulty":       0,
		"totalDifficulty":  0,
		"extraData":        hexutil.Uint64(0),
		"size":             hexutil.Uint64(size),
		"gasLimit":         hexutil.Uint64(gasLimit), // Static gas limit
		"gasUsed":          (*hexutil.Big)(gasUsed),
		"timestamp":        hexutil.Uint64(header.Time.Unix()),
		"transactions":     transactions.([]common.Hash),
		"uncles":           []string{},
		"receiptsRoot":     common.Hash{},
	}
}

func convertTransactionsToRPC(clientCtx client.Context, txs []tmtypes.Tx, blockHash common.Hash, height uint64) ([]common.Hash, *big.Int, error) {
	transactions := make([]common.Hash, len(txs))
	gasUsed := big.NewInt(0)

	for i, tx := range txs {
		ethTx, err := bytesToEthTx(clientCtx, tx)
		if err != nil {
			// continue to next transaction in case it's not a MsgEthereumTx
			continue
		}
		// TODO: Remove gas usage calculation if saving gasUsed per block
		gasUsed.Add(gasUsed, ethTx.Fee())
		tx, err := newRPCTransaction(*ethTx, common.BytesToHash(tx.Hash()), blockHash, &height, uint64(i))
		if err != nil {
			return nil, nil, err
		}
		transactions[i] = tx.Hash
	}

	return transactions, gasUsed, nil
}

// Transaction represents a transaction returned to RPC clients.
type Transaction struct {
	BlockHash        *common.Hash    `json:"blockHash"`
	BlockNumber      *hexutil.Big    `json:"blockNumber"`
	From             common.Address  `json:"from"`
	Gas              hexutil.Uint64  `json:"gas"`
	GasPrice         *hexutil.Big    `json:"gasPrice"`
	Hash             common.Hash     `json:"hash"`
	Input            hexutil.Bytes   `json:"input"`
	Nonce            hexutil.Uint64  `json:"nonce"`
	To               *common.Address `json:"to"`
	TransactionIndex *hexutil.Uint64 `json:"transactionIndex"`
	Value            *hexutil.Big    `json:"value"`
	V                *hexutil.Big    `json:"v"`
	R                *hexutil.Big    `json:"r"`
	S                *hexutil.Big    `json:"s"`
}

func bytesToEthTx(clientCtx client.Context, bz []byte) (*evmtypes.MsgEthereumTx, error) {
	var stdTx sdk.Tx
	// TODO: switch to UnmarshalBinaryBare on SDK v0.40.0
	err := clientCtx.Codec.UnmarshalBinaryLengthPrefixed(bz, &stdTx)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	ethTx, ok := stdTx.(evmtypes.MsgEthereumTx)
	if !ok {
		return nil, fmt.Errorf("invalid transaction type %T, expected MsgEthereumTx", stdTx)
	}
	return &ethTx, nil
}

// newRPCTransaction returns a transaction that will serialize to the RPC
// representation, with the given location metadata set (if available).
func newRPCTransaction(tx evmtypes.MsgEthereumTx, txHash, blockHash common.Hash, blockNumber *uint64, index uint64) (*Transaction, error) {
	// Verify signature and retrieve sender address
	from, err := tx.VerifySig(tx.ChainID())
	if err != nil {
		return nil, err
	}

	result := Transaction{
		From:     from,
		Gas:      hexutil.Uint64(tx.Data.GasLimit),
		GasPrice: (*hexutil.Big)(tx.Data.Price),
		Hash:     txHash,
		Input:    hexutil.Bytes(tx.Data.Payload),
		Nonce:    hexutil.Uint64(tx.Data.AccountNonce),
		To:       tx.To(),
		Value:    (*hexutil.Big)(tx.Data.Amount),
		V:        (*hexutil.Big)(tx.Data.V),
		R:        (*hexutil.Big)(tx.Data.R),
		S:        (*hexutil.Big)(tx.Data.S),
	}

	if blockHash != (common.Hash{}) {
		result.BlockHash = &blockHash
		result.BlockNumber = (*hexutil.Big)(new(big.Int).SetUint64(*blockNumber))
		result.TransactionIndex = (*hexutil.Uint64)(&index)
	}

	return &result, nil
}

// GetTransactionByHash returns the transaction identified by hash.
func (e *PublicEthAPI) GetTransactionByHash(hash common.Hash) (*Transaction, error) {
	e.logger.Debug("eth_getTransactionByHash", "hash", hash)
	tx, err := e.clientCtx.Client.Tx(hash.Bytes(), false)
	if err != nil {
		// Return nil for transaction when not found
		return nil, nil
	}

	// Can either cache or just leave this out if not necessary
	block, err := e.clientCtx.Client.Block(&tx.Height)
	if err != nil {
		return nil, err
	}
	blockHash := common.BytesToHash(block.Block.Header.Hash())

	ethTx, err := bytesToEthTx(e.clientCtx, tx.Tx)
	if err != nil {
		return nil, err
	}

	height := uint64(tx.Height)
	return newRPCTransaction(*ethTx, common.BytesToHash(tx.Tx.Hash()), blockHash, &height, uint64(tx.Index))
}

// GetTransactionByBlockHashAndIndex returns the transaction identified by hash and index.
func (e *PublicEthAPI) GetTransactionByBlockHashAndIndex(hash common.Hash, idx hexutil.Uint) (*Transaction, error) {
	e.logger.Debug("eth_getTransactionByHashAndIndex", "hash", hash, "index", idx)
	res, _, err := e.clientCtx.Query(fmt.Sprintf("custom/%s/%s/%s", evmtypes.ModuleName, evmtypes.QueryHashToHeight, hash.Hex()))
	if err != nil {
		return nil, err
	}

	var out evmtypes.QueryResBlockNumber
	e.clientCtx.Codec.MustUnmarshalJSON(res, &out)
	return e.getTransactionByBlockNumberAndIndex(out.Number, idx)
}

// GetTransactionByBlockNumberAndIndex returns the transaction identified by number and index.
func (e *PublicEthAPI) GetTransactionByBlockNumberAndIndex(blockNum BlockNumber, idx hexutil.Uint) (*Transaction, error) {
	e.logger.Debug("eth_getTransactionByBlockNumberAndIndex", "number", blockNum, "index", idx)
	value := blockNum.Int64()
	return e.getTransactionByBlockNumberAndIndex(value, idx)
}

func (e *PublicEthAPI) getTransactionByBlockNumberAndIndex(number int64, idx hexutil.Uint) (*Transaction, error) {
	block, err := e.clientCtx.Client.Block(&number)
	if err != nil {
		return nil, err
	}
	header := block.Block.Header

	txs := block.Block.Txs
	if uint64(idx) >= uint64(len(txs)) {
		return nil, nil
	}
	ethTx, err := bytesToEthTx(e.clientCtx, txs[idx])
	if err != nil {
		return nil, err
	}

	height := uint64(header.Height)
	return newRPCTransaction(*ethTx, common.BytesToHash(txs[idx].Hash()), common.BytesToHash(header.Hash()), &height, uint64(idx))
}

// GetTransactionReceipt returns the transaction receipt identified by hash.
func (e *PublicEthAPI) GetTransactionReceipt(hash common.Hash) (map[string]interface{}, error) {
	e.logger.Debug("eth_getTransactionReceipt", "hash", hash)
	tx, err := e.clientCtx.Client.Tx(hash.Bytes(), false)
	if err != nil {
		// Return nil for transaction when not found
		return nil, nil
	}

	// Query block for consensus hash
	block, err := e.clientCtx.Client.Block(&tx.Height)
	if err != nil {
		return nil, err
	}
	blockHash := common.BytesToHash(block.Block.Header.Hash())

	// Convert tx bytes to eth transaction
	ethTx, err := bytesToEthTx(e.clientCtx, tx.Tx)
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

	if data.Logs == nil {
		data.Logs = []*ethtypes.Log{}
	}

	receipt := map[string]interface{}{
		// Consensus fields: These fields are defined by the Yellow Paper
		"status":            status,
		"cumulativeGasUsed": nil, // ignore until needed
		"logsBloom":         data.Bloom,
		"logs":              data.Logs,

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

// Copied the Account and StorageResult types since they are registered under an
// internal pkg on geth.

// AccountResult struct for account proof
type AccountResult struct {
	Address      common.Address  `json:"address"`
	AccountProof []string        `json:"accountProof"`
	Balance      *hexutil.Big    `json:"balance"`
	CodeHash     common.Hash     `json:"codeHash"`
	Nonce        hexutil.Uint64  `json:"nonce"`
	StorageHash  common.Hash     `json:"storageHash"`
	StorageProof []StorageResult `json:"storageProof"`
}

// StorageResult defines the format for storage proof return
type StorageResult struct {
	Key   string       `json:"key"`
	Value *hexutil.Big `json:"value"`
	Proof []string     `json:"proof"`
}

// GetProof returns an account object with proof and any storage proofs
func (e *PublicEthAPI) GetProof(address common.Address, storageKeys []string, block BlockNumber) (*AccountResult, error) {
	e.logger.Debug("eth_getProof", "address", address, "keys", storageKeys, "number", block)
	e.clientCtx = e.clientCtx.WithHeight(int64(block))
	path := fmt.Sprintf("custom/%s/%s/%s", evmtypes.ModuleName, evmtypes.QueryAccount, address.Hex())

	// query eth account at block height
	resBz, _, err := e.clientCtx.Query(path)
	if err != nil {
		return nil, err
	}

	var account evmtypes.QueryResAccount
	e.clientCtx.Codec.MustUnmarshalJSON(resBz, &account)

	storageProofs := make([]StorageResult, len(storageKeys))
	opts := client.ABCIQueryOptions{Height: int64(block), Prove: true}
	for i, k := range storageKeys {
		// Get value for key
		vPath := fmt.Sprintf("custom/%s/%s/%s/%s", evmtypes.ModuleName, evmtypes.QueryStorage, address, k)
		vRes, err := e.clientCtx.Client.ABCIQueryWithOptions(vPath, nil, opts)
		if err != nil {
			return nil, err
		}

		var value evmtypes.QueryResStorage
		e.clientCtx.Codec.MustUnmarshalJSON(vRes.Response.GetValue(), &value)

		// check for proof
		proof := vRes.Response.GetProof()
		proofStr := new(merkle.Proof).String()
		if proof != nil {
			proofStr = proof.String()
		}

		storageProofs[i] = StorageResult{
			Key:   k,
			Value: (*hexutil.Big)(common.BytesToHash(value.Value).Big()),
			Proof: []string{proofStr},
		}
	}

	req := abci.RequestQuery{
		Path:   fmt.Sprintf("store/%s/key", auth.StoreKey),
		Data:   auth.AddressStoreKey(sdk.AccAddress(address.Bytes())),
		Height: int64(block),
		Prove:  true,
	}

	res, err := e.clientCtx.QueryABCI(req)
	if err != nil {
		return nil, err
	}

	// check for proof
	accountProof := res.GetProof()
	accProofStr := new(merkle.Proof).String()
	if accountProof != nil {
		accProofStr = accountProof.String()
	}

	return &AccountResult{
		Address:      address,
		AccountProof: []string{accProofStr},
		Balance:      (*hexutil.Big)(utils.MustUnmarshalBigInt(account.Balance)),
		CodeHash:     common.BytesToHash(account.CodeHash),
		Nonce:        hexutil.Uint64(account.Nonce),
		StorageHash:  common.Hash{}, // Ethermint doesn't have a storage hash
		StorageProof: storageProofs,
	}, nil
}

// generateFromArgs populates tx message with args (used in RPC API)
func (e *PublicEthAPI) generateFromArgs(args params.SendTxArgs) (*evmtypes.MsgEthereumTx, error) {
	var (
		nonce    uint64
		gasLimit uint64
		err      error
	)

	amount := (*big.Int)(args.Value)
	gasPrice := (*big.Int)(args.GasPrice)

	if args.GasPrice == nil {

		// Set default gas price
		// TODO: Change to min gas price from context once available through server/daemon
		gasPrice = big.NewInt(ethermint.DefaultGasPrice)
	}

	if args.Nonce == nil {
		// Get nonce (sequence) from account
		from := sdk.AccAddress(args.From.Bytes())
		accRet := authtypes.NewAccountRetriever(e.clientCtx)

		if e.clientCtx.Keybase == nil {
			return nil, fmt.Errorf("clientCtx.Keybase is nil")
		}

		err = accRet.EnsureExists(from)
		if err != nil {
			// account doesn't exist
			return nil, fmt.Errorf("nonexistent account %s: %s", args.From.Hex(), err)
		}

		_, nonce, err = accRet.GetAccountNumberSequence(from)
		if err != nil {
			return nil, err
		}
	} else {
		nonce = (uint64)(*args.Nonce)
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

	return &msg, nil
}
