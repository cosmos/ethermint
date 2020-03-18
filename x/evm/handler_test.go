package evm_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	ethcmn "github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/ethermint/app"
	"github.com/cosmos/ethermint/crypto"
	"github.com/cosmos/ethermint/x/evm"
	"github.com/cosmos/ethermint/x/evm/types"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type EvmTestSuite struct {
	suite.Suite

	ctx     sdk.Context
	handler sdk.Handler
	app     *app.EthermintApp
}

func (suite *EvmTestSuite) SetupTest() {
	checkTx := false

	suite.app = app.Setup(checkTx)
	suite.ctx = suite.app.BaseApp.NewContext(checkTx, abci.Header{Height: 1, ChainID: "3", Time: time.Now().UTC()})
	suite.handler = evm.NewHandler(suite.app.EvmKeeper)
}

func TestEvmTestSuite(t *testing.T) {
	suite.Run(t, new(EvmTestSuite))
}

func (suite *EvmTestSuite) TestHandleMsgEthereumTx() {
	privkey, err := crypto.GenerateKey()
	suite.Require().NoError(err)
	sender := ethcmn.HexToAddress(privkey.PubKey().Address().String())

	var (
		tx      types.MsgEthereumTx
		chainID *big.Int
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"passed",
			func() {
				suite.app.EvmKeeper.SetBalance(suite.ctx, sender, big.NewInt(100))
				tx = types.NewMsgEthereumTx(0, &sender, big.NewInt(100), 0, big.NewInt(10000), nil)

				// parse context chain ID to big.Int
				var ok bool
				chainID, ok = new(big.Int).SetString(suite.ctx.ChainID(), 10)
				suite.Require().True(ok)

				// sign transaction
				err = tx.Sign(chainID, privkey.ToECDSA())
				suite.Require().NoError(err)
			},
			true,
		},
		{
			"insufficient balance",
			func() {
				tx = types.NewMsgEthereumTxContract(0, big.NewInt(100), 0, big.NewInt(10000), nil)

				// parse context chain ID to big.Int
				var ok bool
				chainID, ok = new(big.Int).SetString(suite.ctx.ChainID(), 10)
				suite.Require().True(ok)

				// sign transaction
				err = tx.Sign(chainID, privkey.ToECDSA())
				suite.Require().NoError(err)
			},
			false,
		},
		{
			"tx encoding failed",
			func() {
				tx = types.NewMsgEthereumTxContract(0, big.NewInt(100), 0, big.NewInt(10000), nil)
			},
			false,
		},
		{
			"invalid chain ID",
			func() {
				suite.ctx = suite.ctx.WithChainID("chainID")
			},
			false,
		},
		{
			"VerifySig failed",
			func() {
				tx = types.NewMsgEthereumTxContract(0, big.NewInt(100), 0, big.NewInt(10000), nil)
			},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run("", func() {
			suite.SetupTest() // reset
			tc.malleate()

			res, err := suite.handler(suite.ctx, tx)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}

func (suite *EvmTestSuite) TestMsgEthermint() {
	var (
		tx   types.MsgEthermint
		from = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
		to   = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"passed",
			func() {
				tx = types.NewMsgEthermint(0, &to, sdk.NewInt(1), 100000, sdk.NewInt(2), []byte("test"), from)
				suite.app.EvmKeeper.SetBalance(suite.ctx, ethcmn.BytesToAddress(from.Bytes()), big.NewInt(100))
			},
			true,
		},
		{
			"invalid state transition",
			func() {
				tx = types.NewMsgEthermint(0, &to, sdk.NewInt(1), 100000, sdk.NewInt(2), []byte("test"), from)
			},
			false,
		},
		{
			"invalid chain ID",
			func() {
				suite.ctx = suite.ctx.WithChainID("chainID")
			},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run("", func() {
			suite.SetupTest() // reset
			tc.malleate()

			res, err := suite.handler(suite.ctx, tx)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}
