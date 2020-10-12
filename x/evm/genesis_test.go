package evm_test

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/cosmos/ethermint/crypto/ethsecp256k1"
	"github.com/cosmos/ethermint/x/evm"
	"github.com/cosmos/ethermint/x/evm/types"

	"github.com/ethereum/go-ethereum/common"
)

func (suite *EvmTestSuite) TestExportImport() {
	var genState types.GenesisState
	suite.Require().NotPanics(func() {
		genState = evm.ExportGenesis(suite.ctx, suite.app.EvmKeeper, suite.app.AccountKeeper)
	})

	_ = evm.InitGenesis(suite.ctx, suite.app.EvmKeeper, genState)
}
