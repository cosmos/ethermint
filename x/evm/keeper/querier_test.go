package keeper_test

import (
	"github.com/cosmos/ethermint/x/evm/types"

	abci "github.com/tendermint/tendermint/abci/types"
)

func (suite *KeeperTestSuite) TestQuerier() {

	testCases := []struct {
		msg     string
		path    []string
		expPass bool
	}{
		{"protocol version", []string{types.QueryProtocolVersion}, true},
		{"balance", []string{types.QueryBalance, "0x0"}, true},
		{"block number", []string{types.QueryBlockNumber, "0x0"}, true},
		{"storage", []string{types.QueryStorage, "0x0", "0x0"}, true},
		{"code", []string{types.QueryCode, "0x0"}, true},
		{"nonce", []string{types.QueryNonce, "0x0"}, true},
		{"hash to height", []string{types.QueryHashToHeight, "0x0"}, true},
		{"tx logs", []string{types.QueryTxLogs, "0x0"}, true},
		{"logs bloom", []string{types.QueryLogsBloom, "0x0"}, true},
		{"logs", []string{types.QueryLogs, "0x0"}, true},
		{"account", []string{types.QueryAccount, "0x0"}, true},
		{"unknown request", []string{"other"}, false},
	}

	for i, tc := range testCases {
		tc := tc
		bz, err := suite.querier(suite.ctx, tc.path, abci.RequestQuery{})
		if tc.expPass {
			suite.Require().NoError(err, "valid test %d failed: %s", i, tc.msg)
			suite.Require().NotZero(len(bz))
			// err = json.Unmarshal(bz)
			// suite.Require().NoError(err)
		} else {
			suite.Require().Error(err, "invalid test %d passed: %s", i, tc.msg)
		}
	}
}
