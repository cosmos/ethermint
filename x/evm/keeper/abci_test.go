package keeper_test

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestBeginBlock() {
	req := abci.RequestBeginBlock{
		Header: abci.Header{
			LastBlockId: abci.BlockID{
				Hash: []byte("hash"),
			},
			Height: 10,
		},
	}

	// set gas limit to 1 to ensure no gas is consumed during the operation
	ctx := suite.ctx.WithBlockGasMeter(sdk.NewGasMeter(1))
	initialConsumed := suite.ctx.GasMeter().GasConsumed()
	suite.app.EvmKeeper.Bloom.SetInt64(10)
	suite.app.EvmKeeper.TxCount = 10

	suite.app.EvmKeeper.BeginBlock(ctx, abci.RequestBeginBlock{})
	suite.Require().NotZero(suite.app.EvmKeeper.Bloom.Int64())
	suite.Require().NotZero(suite.app.EvmKeeper.TxCount)

	suite.Require().Zero(ctx.BlockGasMeter().GasConsumed())
	suite.Require().Equal(int64(initialConsumed), int64(suite.ctx.GasMeter().GasConsumed()))

	suite.app.EvmKeeper.BeginBlock(ctx, req)
	suite.Require().Zero(suite.app.EvmKeeper.Bloom.Int64())
	suite.Require().Zero(suite.app.EvmKeeper.TxCount)

	suite.Require().Zero(ctx.BlockGasMeter().GasConsumed())
	suite.Require().Equal(int64(initialConsumed), int64(suite.ctx.GasMeter().GasConsumed()))

	// query using the original context
	lastHeight, found := suite.app.EvmKeeper.GetBlockHash(suite.ctx, req.Header.LastBlockId.Hash)
	suite.Require().True(found)
	suite.Require().Equal(int64(9), lastHeight)
}

func (suite *KeeperTestSuite) TestEndBlock() {
	suite.app.EvmKeeper.Bloom.SetInt64(10)

	// set gas limit to 1 to ensure no gas is consumed during the operation
	ctx := suite.ctx.WithBlockGasMeter(sdk.NewGasMeter(1))
	initialConsumed := ctx.GasMeter().GasConsumed()

	_ = suite.app.EvmKeeper.EndBlock(ctx, abci.RequestEndBlock{Height: 10})

	suite.Require().Zero(ctx.BlockGasMeter().GasConsumed())
	suite.Require().Equal(int64(initialConsumed), int64(ctx.GasMeter().GasConsumed()))

	bloom, found := suite.app.EvmKeeper.GetBlockBloom(suite.ctx, 100)
	suite.Require().True(found)
	suite.Require().Equal(int64(10), bloom.Big().Int64())

}
