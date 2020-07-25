package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/ethermint/app"
	"github.com/cosmos/ethermint/x/evm/keeper"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	abci "github.com/tendermint/tendermint/abci/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx     sdk.Context
	querier sdk.Querier
	app     *app.EthermintApp
	address ethcmn.Address
}

func (suite *KeeperTestSuite) SetupTest() {
	checkTx := false

	suite.app = app.Setup(checkTx)
	suite.ctx = suite.app.BaseApp.NewContext(checkTx, abci.Header{Height: 1, ChainID: "3", Time: time.Now().UTC()})
	suite.querier = keeper.NewQuerier(suite.app.EvmKeeper)
	suite.address = ethcmn.HexToAddress("0x756F45E3FA69347A9A973A725E3C98bC4db0b4c1")
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestTransactionLogs() {
	hash := ethcmn.FromHex("0x0d87a3a5f73140f46aac1bf419263e4e94e87c292f25007700ab7f2060e2af68")
	ethHash := ethcmn.BytesToHash(hash)
	log := &ethtypes.Log{
		Address:     suite.address,
		Data:        []byte("log"),
		BlockNumber: 10,
	}
	log2 := &ethtypes.Log{
		Address:     suite.address,
		Data:        []byte("log2"),
		BlockNumber: 11,
	}
	expLogs := []*ethtypes.Log{log}

	err := suite.app.EvmKeeper.SetLogs(suite.ctx, ethHash, expLogs)
	suite.Require().NoError(err)

	logs, err := suite.app.EvmKeeper.GetLogs(suite.ctx, ethHash)
	suite.Require().NoError(err)
	suite.Require().Equal(expLogs, logs)

	expLogs = []*ethtypes.Log{log2, log}

	// add another log under the zero hash
	suite.app.EvmKeeper.AddLog(suite.ctx, log2)
	logs = suite.app.EvmKeeper.AllLogs(suite.ctx)
	suite.Require().Equal(expLogs, logs)

	// add another log under the zero hash
	log3 := &ethtypes.Log{
		Address:     suite.address,
		Data:        []byte("log3"),
		BlockNumber: 10,
	}
	suite.app.EvmKeeper.AddLog(suite.ctx, log3)

	txLogs := suite.app.EvmKeeper.GetAllTxLogs(suite.ctx)
	suite.Require().Equal(2, len(txLogs))

	suite.Require().Equal(ethcmn.Hash{}.String(), txLogs[0].Hash.String())
	suite.Require().Equal([]*ethtypes.Log{log2, log3}, txLogs[0].Logs)

	suite.Require().Equal(ethHash.String(), txLogs[1].Hash.String())
	suite.Require().Equal([]*ethtypes.Log{log}, txLogs[1].Logs)
}

func (suite *KeeperTestSuite) TestBlockHash() {
	testCase := []struct {
		name    string
		hash    []byte
		expPass bool
	}{
		{
			"valid hash",
			[]byte{0x43, 0x32},
			true,
		},
		{
			"invalid hash",
			[]byte{0x3, 0x2},
			false,
		},
	}

	for _, tc := range testCase {
		if tc.expPass {
			suite.app.EvmKeeper.SetBlockHash(suite.ctx, tc.hash, 7)
			height, found := suite.app.EvmKeeper.GetBlockHash(suite.ctx, tc.hash)
			suite.Require().True(found, tc.name)
			suite.Require().Equal(int64(7), height, tc.name)
		} else {
			height, found := suite.app.EvmKeeper.GetBlockHash(suite.ctx, tc.hash)
			suite.Require().False(found, tc.name)
			suite.Require().Equal(int64(0), height, tc.name)
		}
	}
}

func (suite *KeeperTestSuite) TestBlockBloom() {
	testCase := []struct {
		name    string
		height  int64
		expPass bool
	}{
		{
			"found bloom",
			4,
			true,
		},
		{
			"missing bloom",
			5,
			false,
		},
	}

	for _, tc := range testCase {
		if tc.expPass {
			testBloom := ethtypes.BytesToBloom([]byte{0x1, 0x3})
			suite.app.EvmKeeper.SetBlockBloom(suite.ctx, tc.height, testBloom)
			bloom, found := suite.app.EvmKeeper.GetBlockBloom(suite.ctx, tc.height)
			suite.Require().True(found, tc.name)
			suite.Require().Equal(testBloom, bloom, tc.name)
		} else {
			bloom, found := suite.app.EvmKeeper.GetBlockBloom(suite.ctx, tc.height)
			suite.Require().False(found, tc.name)
			suite.Require().Equal(ethtypes.Bloom{}, bloom, tc.name)
		}
	}
}
