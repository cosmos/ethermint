package evm_test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/status-im/keycard-go/hexutils"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/cosmos/ethermint/app"
	"github.com/cosmos/ethermint/x/evm"
	"github.com/cosmos/ethermint/x/evm/keeper"
	"github.com/cosmos/ethermint/x/evm/types"

	abci "github.com/tendermint/tendermint/abci/types"
)

type EvmTestSuite struct {
	suite.Suite

	ctx     sdk.Context
	handler sdk.Handler
	querier sdk.Querier
	app     *app.EthermintApp
	codec   *codec.Codec
}

func (suite *EvmTestSuite) SetupTest() {
	checkTx := false

	suite.app = app.Setup(checkTx)
	suite.ctx = suite.app.BaseApp.NewContext(checkTx, abci.Header{Height: 1, ChainID: "3", Time: time.Now().UTC()})
	suite.handler = evm.NewHandler(suite.app.EvmKeeper)
	suite.querier = keeper.NewQuerier(suite.app.EvmKeeper)
	suite.codec = codec.New()
}

func TestEvmTestSuite(t *testing.T) {
	suite.Run(t, new(EvmTestSuite))
}

func (suite *EvmTestSuite) TestHandler_Logs() {
	// Test contract:

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

	gasLimit := uint64(100000)
	gasPrice := big.NewInt(1000000)

	priv, err := crypto.GenerateKey()
	suite.Require().NoError(err, "failed to create key")

	bytecode := common.FromHex("0x6080604052348015600f57600080fd5b5060117f775a94827b8fd9b519d36cd827093c664f93347070a554f65e4a6f56cd73889860405160405180910390a2603580604b6000396000f3fe6080604052600080fdfea165627a7a723058206cab665f0f557620554bb45adf266708d2bd349b8a4314bdff205ee8440e3c240029")
	tx := types.NewMsgEthereumTx(1, nil, big.NewInt(0), gasLimit, gasPrice, bytecode)
	tx.Sign(big.NewInt(3), priv)

	result, err := suite.handler(suite.ctx, tx)
	suite.Require().NoError(err, "failed to handle eth tx msg")

	resultData, err := types.DecodeResultData(result.Data)
	suite.Require().NoError(err, "failed to decode result data")

	suite.Require().Equal(len(resultData.Logs), 1)
	suite.Require().Equal(len(resultData.Logs[0].Topics), 2)

	hash := []byte{1}
	err = suite.app.EvmKeeper.SetTransactionLogs(suite.ctx, resultData.Logs, hash)
	suite.Require().NoError(err, "failed to set logs")

	logs, err := suite.app.EvmKeeper.GetTransactionLogs(suite.ctx, hash)
	suite.Require().NoError(err, "failed to get logs")

	suite.Require().Equal(logs, resultData.Logs)
}

func (suite *EvmTestSuite) TestQueryTxLogs() {
	gasLimit := uint64(100000)
	gasPrice := big.NewInt(1000000)

	priv, err := crypto.GenerateKey()
	suite.Require().NoError(err, "failed to create key")

	// send contract deployment transaction with an event in the constructor
	bytecode := common.FromHex("0x6080604052348015600f57600080fd5b5060117f775a94827b8fd9b519d36cd827093c664f93347070a554f65e4a6f56cd73889860405160405180910390a2603580604b6000396000f3fe6080604052600080fdfea165627a7a723058206cab665f0f557620554bb45adf266708d2bd349b8a4314bdff205ee8440e3c240029")
	tx := types.NewMsgEthereumTx(1, nil, big.NewInt(0), gasLimit, gasPrice, bytecode)
	tx.Sign(big.NewInt(3), priv)

	// result, err := evm.HandleEthTxMsg(suite.ctx, suite.app.EvmKeeper, tx)
	// suite.Require().NoError(err, "failed to handle eth tx msg")

	result, err := suite.handler(suite.ctx, tx)
	suite.Require().NoError(err)
	suite.Require().NotNil(result)

	resultData, err := types.DecodeResultData(result.Data)
	suite.Require().NoError(err, "failed to decode result data")

	suite.Require().Equal(len(resultData.Logs), 1)
	suite.Require().Equal(len(resultData.Logs[0].Topics), 2)

	// get logs by tx hash
	hash := resultData.TxHash.Bytes()

	logs, err := suite.app.EvmKeeper.GetTransactionLogs(suite.ctx, hash)
	suite.Require().NoError(err, "failed to get logs")

	suite.Require().Equal(logs, resultData.Logs)

	// query tx logs
	path := []string{"txLogs", fmt.Sprintf("0x%x", hash)}
	res, err := suite.querier(suite.ctx, path, abci.RequestQuery{})
	suite.Require().NoError(err, "failed to query txLogs")

	var txLogs types.QueryETHLogs
	suite.codec.MustUnmarshalJSON(res, &txLogs)

	// amino decodes an empty byte array as nil, whereas JSON decodes it as []byte{} causing a discrepancy
	resultData.Logs[0].Data = []byte{}
	suite.Require().Equal(txLogs.Logs[0], resultData.Logs[0])
}

func (suite *EvmTestSuite) TestDeployAndCallContract() {
	// Test contract:
	//http://remix.ethereum.org/#optimize=false&evmVersion=istanbul&version=soljson-v0.5.15+commit.6a57276f.js
	//2_Owner.sol
	//
	//pragma solidity >=0.4.22 <0.7.0;
	//
	///**
	// * @title Owner
	// * @dev Set & change owner
	// */
	//contract Owner {
	//
	//	address private owner;
	//
	//	// event for EVM logging
	//	event OwnerSet(address indexed oldOwner, address indexed newOwner);
	//
	//	// modifier to check if caller is owner
	//	modifier isOwner() {
	//	// If the first argument of 'require' evaluates to 'false', execution terminates and all
	//	// changes to the state and to Ether balances are reverted.
	//	// This used to consume all gas in old EVM versions, but not anymore.
	//	// It is often a good idea to use 'require' to check if functions are called correctly.
	//	// As a second argument, you can also provide an explanation about what went wrong.
	//	require(msg.sender == owner, "Caller is not owner");
	//	_;
	//}
	//
	//	/**
	//	 * @dev Set contract deployer as owner
	//	 */
	//	constructor() public {
	//	owner = msg.sender; // 'msg.sender' is sender of current call, contract deployer for a constructor
	//	emit OwnerSet(address(0), owner);
	//}
	//
	//	/**
	//	 * @dev Change owner
	//	 * @param newOwner address of new owner
	//	 */
	//	function changeOwner(address newOwner) public isOwner {
	//	emit OwnerSet(owner, newOwner);
	//	owner = newOwner;
	//}
	//
	//	/**
	//	 * @dev Return owner address
	//	 * @return address of owner
	//	 */
	//	function getOwner() external view returns (address) {
	//	return owner;
	//}
	//}

	// Deploy contract - Owner.sol
	gasLimit := uint64(100000000)
	gasPrice := big.NewInt(10000)

	priv, err := crypto.GenerateKey()
	suite.Require().NoError(err, "failed to create key")

	bytecode := common.FromHex("0x608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167f342827c97908e5e2f71151c08502a66d44b6f758e3ac2f1de95f02eb95f0a73560405160405180910390a36102c4806100dc6000396000f3fe608060405234801561001057600080fd5b5060043610610053576000357c010000000000000000000000000000000000000000000000000000000090048063893d20e814610058578063a6f9dae1146100a2575b600080fd5b6100606100e6565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6100e4600480360360208110156100b857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061010f565b005b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146101d1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260138152602001807f43616c6c6572206973206e6f74206f776e65720000000000000000000000000081525060200191505060405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f342827c97908e5e2f71151c08502a66d44b6f758e3ac2f1de95f02eb95f0a73560405160405180910390a3806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505056fea265627a7a72315820f397f2733a89198bc7fed0764083694c5b828791f39ebcbc9e414bccef14b48064736f6c63430005100032")
	tx := types.NewMsgEthereumTx(1, nil, big.NewInt(0), gasLimit, gasPrice, bytecode)
	tx.Sign(big.NewInt(3), priv)
	suite.Require().NoError(err)

	result, err := suite.handler(suite.ctx, tx)
	suite.Require().NoError(err, "failed to handle eth tx msg")

	resultData, err := types.DecodeResultData(result.Data)
	suite.Require().NoError(err, "failed to decode result data")

	// store - changeOwner
	gasLimit = uint64(100000000000)
	gasPrice = big.NewInt(100)
	receiver := common.HexToAddress(resultData.Address.String())

	storeAddr := "0xa6f9dae10000000000000000000000006a82e4a67715c8412a9114fbd2cbaefbc8181424"
	bytecode = common.FromHex(storeAddr)
	tx = types.NewMsgEthereumTx(2, &receiver, big.NewInt(0), gasLimit, gasPrice, bytecode)
	tx.Sign(big.NewInt(3), priv)
	suite.Require().NoError(err)

	result, err = suite.handler(suite.ctx, tx)
	suite.Require().NoError(err, "failed to handle eth tx msg")

	resultData, err = types.DecodeResultData(result.Data)
	suite.Require().NoError(err, "failed to decode result data")

	// query - getOwner
	bytecode = common.FromHex("0x893d20e8")
	tx = types.NewMsgEthereumTx(2, &receiver, big.NewInt(0), gasLimit, gasPrice, bytecode)
	tx.Sign(big.NewInt(3), priv)
	suite.Require().NoError(err)

	result, err = suite.handler(suite.ctx, tx)
	suite.Require().NoError(err, "failed to handle eth tx msg")

	resultData, err = types.DecodeResultData(result.Data)
	suite.Require().NoError(err, "failed to decode result data")

	getAddr := strings.ToLower(hexutils.BytesToHex(resultData.Ret))
	suite.Require().Equal(true, strings.HasSuffix(storeAddr, getAddr), "Fail to query the address")
}
