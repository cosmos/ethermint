package evm

// import (
// 	"math/big"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/suite"

// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	"github.com/cosmos/ethermint/app"
// 	"github.com/cosmos/ethermint/crypto"
// 	"github.com/cosmos/ethermint/x/evm/keeper"

// 	ethcmn "github.com/ethereum/go-ethereum/common"
// 	ethtypes "github.com/ethereum/go-ethereum/core/types"

// 	abci "github.com/tendermint/tendermint/abci/types"
// )

// pragma solidity ^0.5.1;

// contract Test {
//     event Hello(uint256 indexed world);

//     constructor() public {
//         emit Hello(17);
//     }
// }

// {
// 	"linkReferences": {},
// 	"object": "6080604052348015600f57600080fd5b5060117f775a94827b8fd9b519d36cd827093c664f93347070a554f65e4a6f56cd73889860405160405180910390a2603580604b6000396000f3fe6080604052600080fdfea165627a7a723058206cab665f0f557620554bb45adf266708d2bd349b8a4314bdff205ee8440e3c240029",
// 	"opcodes": "PUSH1 0x80 PUSH1 0x40 MSTORE CALLVALUE DUP1 ISZERO PUSH1 0xF JUMPI PUSH1 0x0 DUP1 REVERT JUMPDEST POP PUSH1 0x11 PUSH32 0x775A94827B8FD9B519D36CD827093C664F93347070A554F65E4A6F56CD738898 PUSH1 0x40 MLOAD PUSH1 0x40 MLOAD DUP1 SWAP2 SUB SWAP1 LOG2 PUSH1 0x35 DUP1 PUSH1 0x4B PUSH1 0x0 CODECOPY PUSH1 0x0 RETURN INVALID PUSH1 0x80 PUSH1 0x40 MSTORE PUSH1 0x0 DUP1 REVERT INVALID LOG1 PUSH6 0x627A7A723058 KECCAK256 PUSH13 0xAB665F0F557620554BB45ADF26 PUSH8 0x8D2BD349B8A4314 0xbd SELFDESTRUCT KECCAK256 0x5e 0xe8 DIFFICULTY 0xe EXTCODECOPY 0x24 STOP 0x29 ",
// 	"sourceMap": "25:119:0:-;;;90:52;8:9:-1;5:2;;;30:1;27;20:12;5:2;90:52:0;132:2;126:9;;;;;;;;;;25:119;;;;;;"
// }

// type EvmTestSuite struct {
// 	suite.Suite

// 	ctx     sdk.Context
// 	querier sdk.Querier
// 	app     *app.EthermintApp
// }

// func (suite *EvmTestSuite) SetupTest() {
// 	checkTx := false

// 	suite.app = app.Setup(checkTx)
// 	suite.ctx = suite.app.BaseApp.NewContext(checkTx, abci.Header{Height: 1, ChainID: "3", Time: time.Now().UTC()})
// 	suite.querier = keeper.NewQuerier(suite.app.EvmKeeper)
// }

// func TestEvmTestSuite(t *testing.T) {
// 	suite.Run(t, new(EvmTestSuite))
// }

// func (suite *EvmTestSuite) TestHandler_Logs() {
// 	// Perform state transitions
// 	// suite.app.EvmKeeper.CreateAccount(suite.ctx, address)
// 	// suite.app.EvmKeeper.SetBalance(suite.ctx, address, big.NewInt(5))
// 	// suite.app.EvmKeeper.SetNonce(suite.ctx, address, 4)

// 	priv1, _ := crypto.GenerateKey()
// 	bytecode := common.FromHex("0x6080604052348015600f57600080fd5b5060117f775a94827b8fd9b519d36cd827093c664f93347070a554f65e4a6f56cd73889860405160405180910390a2603580604b6000396000f3fe6080604052600080fdfea165627a7a723058206cab665f0f557620554bb45adf266708d2bd349b8a4314bdff205ee8440e3c240029")
// 	tx := evmtypes.NewEthereumTxMsg(1, nil, big.NewInt(0), gasLimit, gasPrice, bytecode)
// 	tx.Sign(big.NewInt(1), priv1.ToECDSA())

// 	result := handleETHTxMsg(ctx, suite.App.EvmKeeper, *tx)
// 	resultData, err := evmtypes.DecodeResultData(result.Data)
// 	suite.Require().NoError(err, "failed to decode result data")

// 	suite.Require().Equal(len(resultData.Logs), 1)
// 	suite.Require().Equal(len(resultData.Logs[0].Topics), 2)

// 	hash := []byte{1}
// 	err = ek.SetTransactionLogs(ctx, resultData.Logs, hash)
// 	suite.Require().NoError(err, "failed to set logs")

// 	logs, err := ek.GetTransactionLogs(ctx, hash)
// 	suite.Require().NoError(err, "failed to get logs")

// 	suite.Require().Equal(logs, resultData.Logs)
// }

// func TestHandler_Logs(t *testing.T) {
// 	// create logger, codec and root multi-store
// 	cdc := newTestCodec()

// 	// The ParamsKeeper handles parameter storage for the application
// 	keyParams := sdk.NewKVStoreKey(params.StoreKey)
// 	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
// 	paramsKeeper := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
// 	// Set specific supspaces
// 	authSubspace := paramsKeeper.Subspace(auth.DefaultParamspace)
// 	ak := auth.NewAccountKeeper(cdc, accKey, authSubspace, eminttypes.ProtoBaseAccount)
// 	ek := NewKeeper(ak, storageKey, codeKey, blockKey, cdc)

// 	gasLimit := uint64(100000)
// 	gasPrice := big.NewInt(1000000)

// 	priv1, _ := crypto.GenerateKey()

// 	bytecode := common.FromHex("0x6080604052348015600f57600080fd5b5060117f775a94827b8fd9b519d36cd827093c664f93347070a554f65e4a6f56cd73889860405160405180910390a2603580604b6000396000f3fe6080604052600080fdfea165627a7a723058206cab665f0f557620554bb45adf266708d2bd349b8a4314bdff205ee8440e3c240029")

// 	tx := evmtypes.NewEthereumTxMsg(1, nil, big.NewInt(0), gasLimit, gasPrice, bytecode)
// 	tx.Sign(big.NewInt(1), priv1.ToECDSA())

// 	db := dbm.NewMemDB()
// 	cms := store.NewCommitMultiStore(db)
// 	keys := []*sdk.KVStoreKey{accKey, storageKey, codeKey, blockKey}
// 	for _, key := range keys {
// 		cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, nil)
// 	}

// 	err := cms.LoadLatestVersion()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	ms := cms.CacheMultiStore()
// 	ctx := sdk.NewContext(ms, abci.Header{}, false, logger)
// 	ctx = ctx.WithBlockHeight(1).WithChainID("1")

// 	result := handleETHTxMsg(ctx, ek, *tx)
// 	resultData, err := evmtypes.DecodeResultData(result.Data)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if len(resultData.Logs) != 1 {
// 		t.Fatal("Fail: expected 1 log")
// 	}

// 	if len(resultData.Logs[0].Topics) != 2 {
// 		t.Fatal("Fail: expected 2 topics")
// 	}

// 	hash := []byte{1}
// 	err = ek.SetTransactionLogs(ctx, resultData.Logs, hash)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	logs, err := ek.GetTransactionLogs(ctx, hash)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if !reflect.DeepEqual(logs, resultData.Logs) {
// 		t.Fatalf("Fail: got %v expected %v", logs, resultData.Logs)
// 	}
// }
