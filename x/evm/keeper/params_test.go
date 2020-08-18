package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/ethermint/x/evm/types"
)

func (suite *KeeperTestSuite) TestParams() {
	params := suite.app.EvmKeeper.GetParams(suite.ctx)

	suite.Require().Equal(types.DefaultParams().IstanbulBlock.BigInt(), params.IstanbulBlock.BigInt())

	params.EIP150Block = sdk.NewInt(10000)
	suite.app.EvmKeeper.SetParams(suite.ctx, params)
}
