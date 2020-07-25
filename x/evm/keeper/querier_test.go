package keeper_test

import (
	"math/big"

	"github.com/cosmos/ethermint/x/evm/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	abci "github.com/tendermint/tendermint/abci/types"
)

func (suite *KeeperTestSuite) TestQuerier() {
	hex := "0x0d87a3a5f73140f46aac1bf419263e4e94e87c292f25007700ab7f2060e2af68"
	hash := ethcmn.FromHex(hex)
	addrHex := "0x756F45E3FA69347A9A973A725E3C98bC4db0b4c1"

	testCases := []struct {
		msg      string
		path     []string
		malleate func()
		expPass  bool
	}{
		{"protocol version", []string{types.QueryProtocolVersion}, func() {}, true},
		{"balance", []string{types.QueryBalance, addrHex}, func() {
			suite.app.EvmKeeper.SetBalance(suite.ctx, suite.address, big.NewInt(5))
		}, true},
		{"block number", []string{types.QueryBlockNumber, "0x0"}, func() {}, true},
		{"storage", []string{types.QueryStorage, "0x0", "0x0"}, func() {}, true},
		{"code", []string{types.QueryCode, "0x0"}, func() {}, true},
		{"hash to height", []string{types.QueryHashToHeight, hex}, func() {
			suite.app.EvmKeeper.SetBlockHash(suite.ctx, hash, 8)
		}, true},
		{"tx logs", []string{types.QueryTransactionLogs, "0x0"}, func() {}, true},
		{"bloom", []string{types.QueryBloom, "4"}, func() {
			testBloom := ethtypes.BytesToBloom([]byte{0x1, 0x3})
			suite.app.EvmKeeper.SetBlockBloom(suite.ctx, 4, testBloom)
		}, true},
		{"logs", []string{types.QueryLogs, "0x0"}, func() {}, true},
		{"account", []string{types.QueryAccount, "0x0"}, func() {}, true},
		{"exportAccount", []string{types.QueryExportAccount, "0x0"}, func() {}, true},
		{"unknown request", []string{"other"}, func() {}, false},
	}

	for i, tc := range testCases {
		suite.Run("", func() {
			//nolint
			tc := tc
			suite.SetupTest() // reset
			//nolint
			tc.malleate()

			bz, err := suite.querier(suite.ctx, tc.path, abci.RequestQuery{})

			//nolint
			if tc.expPass {
				//nolint
				suite.Require().NoError(err, "valid test %d failed: %s", i, tc.msg)
				suite.Require().NotZero(len(bz))
			} else {
				//nolint
				suite.Require().Error(err, "invalid test %d passed: %s", i, tc.msg)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestProtocolVersion() {

}
