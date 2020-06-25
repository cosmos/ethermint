package types_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/cosmos/ethermint/app"
	"github.com/cosmos/ethermint/crypto"
	"github.com/cosmos/ethermint/x/evm/keeper"

	abci "github.com/tendermint/tendermint/abci/types"
)

// nolint: unused
type StateDBTestSuite struct {
	suite.Suite

	ctx     sdk.Context
	querier sdk.Querier
	app     *app.EthermintApp
}

func TestStateDBTestSuite(t *testing.T) {
	suite.Run(t, new(StateDBTestSuite))
}

func (suite *StateDBTestSuite) SetupTest() {
	checkTx := false

	suite.app = app.Setup(checkTx)
	suite.ctx = suite.app.BaseApp.NewContext(checkTx, abci.Header{Height: 1})
	suite.querier = keeper.NewQuerier(suite.app.EvmKeeper)

	stateDB := suite.app.EvmKeeper.CommitStateDB
	resp := stateDB.WithContext(suite.ctx)
	suite.Require().NotNil(suite.T(), resp)
}

func (suite *StateDBTestSuite) TestBloomFilter() {
	stateDB := suite.app.EvmKeeper.CommitStateDB

	// Prepare db for logs
	tHash := ethcmn.BytesToHash([]byte{0x1})
	stateDB.Prepare(tHash, ethcmn.Hash{}, 0)

	contractAddress := ethcmn.BigToAddress(big.NewInt(1))

	// Generate and add a log to test
	log := ethtypes.Log{Address: contractAddress}
	stateDB.AddLog(&log)

	// Get log from db
	logs, err := stateDB.GetLogs(tHash)
	suite.Require().NoError(err)
	suite.Require().Len(logs, 1)
	suite.Require().Equal(log, *logs[0])

	// get logs bloom from the log
	bloomInt := ethtypes.LogsBloom(logs)
	bloomFilter := ethtypes.BytesToBloom(bloomInt.Bytes())

	// Check to make sure bloom filter will succeed on
	suite.Require().True(ethtypes.BloomLookup(bloomFilter, contractAddress))
	suite.Require().False(ethtypes.BloomLookup(bloomFilter, ethcmn.BigToAddress(big.NewInt(2))))
}

func (suite *StateDBTestSuite) TestStateDBBalance() {
	stateDB := suite.app.EvmKeeper.CommitStateDB

	priv, err := crypto.GenerateKey()
	suite.Require().NoError(err)

	addr := ethcrypto.PubkeyToAddress(priv.ToECDSA().PublicKey)
	value := big.NewInt(100)
	stateDB.SetBalance(addr, value)
	suite.Require().Equal(value, stateDB.GetBalance(addr))

	stateDB.SubBalance(addr, value)
	suite.Require().Equal(big.NewInt(0), stateDB.GetBalance(addr))

	stateDB.AddBalance(addr, value)
	suite.Require().Equal(value, stateDB.GetBalance(addr))
}

func (suite *StateDBTestSuite) TestStateDBNonce() {
	stateDB := suite.app.EvmKeeper.CommitStateDB

	priv, err := crypto.GenerateKey()
	suite.Require().NoError(err)
	addr := ethcrypto.PubkeyToAddress(priv.ToECDSA().PublicKey)

	nonce := uint64(123)
	stateDB.SetNonce(addr, nonce)

	suite.Require().Equal(nonce, stateDB.GetNonce(addr))
}

func (suite *StateDBTestSuite) TestStateDBState() {
	stateDB := suite.app.EvmKeeper.CommitStateDB

	priv, err := crypto.GenerateKey()
	suite.Require().NoError(err)

	addr := ethcrypto.PubkeyToAddress(priv.ToECDSA().PublicKey)
	key := ethcmn.BytesToHash([]byte("foo"))
	val := ethcmn.BytesToHash([]byte("bar"))

	stateDB.SetState(addr, key, val)

	suite.Require().Equal(val, stateDB.GetState(addr, key))
}

func (suite *StateDBTestSuite) TestStateDBCode() {
	stateDB := suite.app.EvmKeeper.CommitStateDB

	priv, err := crypto.GenerateKey()
	suite.Require().NoError(err)

	addr := ethcrypto.PubkeyToAddress(priv.ToECDSA().PublicKey)
	code := []byte("foobar")

	stateDB.SetCode(addr, code)

	suite.Require().Equal(code, stateDB.GetCode(addr))

	codelen := len(code)
	suite.Require().Equal(codelen, stateDB.GetCodeSize(addr))
}

func (suite *StateDBTestSuite) TestStateDBLogs() {
	stateDB := suite.app.EvmKeeper.CommitStateDB

	priv, err := crypto.GenerateKey()
	suite.Require().NoError(err)

	addr := ethcrypto.PubkeyToAddress(priv.ToECDSA().PublicKey)

	hash := ethcmn.BytesToHash([]byte("hash"))
	log := ethtypes.Log{
		Address:     addr,
		Topics:      []common.Hash{ethcmn.BytesToHash([]byte("topic"))},
		Data:        []byte("data"),
		BlockNumber: 1,
		TxHash:      common.Hash{},
		TxIndex:     1,
		BlockHash:   common.Hash{},
		Index:       1,
		Removed:     false,
	}
	logs := []*ethtypes.Log{&log}

	err = stateDB.SetLogs(hash, logs)
	suite.Require().NoError(err)
	dbLogs, err := stateDB.GetLogs(hash)
	suite.Require().NoError(err)
	suite.Require().Equal(logs, dbLogs)

	stateDB.DeleteLogs(hash)
	dbLogs, err = stateDB.GetLogs(hash)
	suite.Require().NoError(err)
	suite.Require().Empty(dbLogs)

	stateDB.AddLog(&log)
	suite.Require().Equal(logs, stateDB.AllLogs())

	//resets state but checking to see if storekey still persists.
	stateDB.Reset(hash)
	suite.Require().Equal(logs, stateDB.AllLogs())
}

func (suite *StateDBTestSuite) TestStateDBPreimage() {
	stateDB := suite.app.EvmKeeper.CommitStateDB

	hash := ethcmn.BytesToHash([]byte("hash"))
	preimage := []byte("preimage")

	stateDB.AddPreimage(hash, preimage)

	suite.Require().Equal(preimage, stateDB.Preimages()[hash])
}

func (suite *StateDBTestSuite) TestStateDBRefund() {
	stateDB := suite.app.EvmKeeper.CommitStateDB

	value := uint64(100)

	stateDB.AddRefund(value)
	suite.Require().Equal(value, stateDB.GetRefund())

	stateDB.SubRefund(value)
	suite.Require().Equal(uint64(0), stateDB.GetRefund())
}

func (suite *StateDBTestSuite) TestStateDBCreateAcct() {
	stateDB := suite.app.EvmKeeper.CommitStateDB

	priv, err := crypto.GenerateKey()
	suite.Require().NoError(err)

	addr := ethcrypto.PubkeyToAddress(priv.ToECDSA().PublicKey)

	stateDB.CreateAccount(addr)
	suite.Require().True(stateDB.Exist(addr))

	value := big.NewInt(100)
	stateDB.AddBalance(addr, value)

	stateDB.CreateAccount(addr)
	suite.Require().Equal(value, stateDB.GetBalance(addr))
}

func (suite *StateDBTestSuite) TestStateDBClearStateOjb() {
	stateDB := suite.app.EvmKeeper.CommitStateDB

	priv, err := crypto.GenerateKey()
	suite.Require().NoError(err)

	addr := ethcrypto.PubkeyToAddress(priv.ToECDSA().PublicKey)

	stateDB.CreateAccount(addr)
	suite.Require().True(stateDB.Exist(addr))

	stateDB.ClearStateObjects()
	suite.Require().False(stateDB.Exist(addr))
}

func (suite *StateDBTestSuite) TestStateDBReset() {
	stateDB := suite.app.EvmKeeper.CommitStateDB

	priv, err := crypto.GenerateKey()
	suite.Require().NoError(err)

	addr := ethcrypto.PubkeyToAddress(priv.ToECDSA().PublicKey)

	hash := ethcmn.BytesToHash([]byte("hash"))

	stateDB.CreateAccount(addr)
	suite.Require().True(stateDB.Exist(addr))

	stateDB.Reset(hash)
	suite.Require().False(stateDB.Exist(addr))
}

func (suite *StateDBTestSuite) TestStateDBUpdateAcct() {

}

func (suite *StateDBTestSuite) TestSuiteDBPrepare() {
	stateDB := suite.app.EvmKeeper.CommitStateDB

	thash := ethcmn.BytesToHash([]byte("thash"))
	bhash := ethcmn.BytesToHash([]byte("bhash"))
	txi := 1

	stateDB.Prepare(thash, bhash, txi)

	suite.Require().Equal(txi, stateDB.TxIndex())
	suite.Require().Equal(bhash, stateDB.BlockHash())
}

func (suite *StateDBTestSuite) TestSuiteDBCopyState() {
	stateDB := suite.app.EvmKeeper.CommitStateDB

	priv, err := crypto.GenerateKey()
	suite.Require().NoError(err)

	addr := ethcrypto.PubkeyToAddress(priv.ToECDSA().PublicKey)

	hash := ethcmn.BytesToHash([]byte("hash"))
	log := ethtypes.Log{
		Address:     addr,
		Topics:      []common.Hash{ethcmn.BytesToHash([]byte("topic"))},
		Data:        []byte("data"),
		BlockNumber: 1,
		TxHash:      common.Hash{},
		TxIndex:     1,
		BlockHash:   common.Hash{},
		Index:       1,
		Removed:     false,
	}
	logs := []*ethtypes.Log{&log}

	err = stateDB.SetLogs(hash, logs)
	suite.Require().NoError(err)

	copyDB := stateDB.Copy()

	copiedDBLogs, err := copyDB.GetLogs(hash)
	suite.Require().NoError(err)
	suite.Require().Equal(logs, copiedDBLogs)
	suite.Require().Equal(stateDB.Exist(addr), copyDB.Exist(addr))
}

func (suite *StateDBTestSuite) TestSuiteDBEmpty() {
	stateDB := suite.app.EvmKeeper.CommitStateDB

	priv, err := crypto.GenerateKey()
	suite.Require().NoError(err)

	addr := ethcrypto.PubkeyToAddress(priv.ToECDSA().PublicKey)

	suite.Require().True(stateDB.Empty(addr))

	stateDB.SetBalance(addr, big.NewInt(100))

	suite.Require().False(stateDB.Empty(addr))
}

func (suite *StateDBTestSuite) TestSuiteDBSuicide() {
	stateDB := suite.app.EvmKeeper.CommitStateDB

	priv, err := crypto.GenerateKey()
	suite.Require().NoError(err)

	addr := ethcrypto.PubkeyToAddress(priv.ToECDSA().PublicKey)

	suicide := stateDB.Suicide(addr)
	suite.Require().False(suicide)
	suite.Require().False(stateDB.HasSuicided(addr))

	//Suicide only works for an account with non-zero balance/nonce
	stateDB.SetBalance(addr, big.NewInt(100))
	suicide = stateDB.Suicide(addr)

	suite.Require().True(suicide)
	suite.Require().True(stateDB.HasSuicided(addr))

	delete := true
	stateDB.Commit(delete)
	suite.Require().False(stateDB.Exist(addr))
}
