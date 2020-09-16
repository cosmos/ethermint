package rpc

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/ethermint/crypto"

	"github.com/tendermint/tendermint/libs/log"
)

type FaucetAPI struct {
	logger log.Logger
	cliCtx context.CLIContext
	keys   []crypto.PrivKeySecp256k1
}

func NewFaucetAPI(cliCtx context.CLIContext, keys []crypto.PrivKeySecp256k1) *FaucetAPI {
	return &FaucetAPI{
		logger: log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "json-rpc", "api", "faucet"),
		cliCtx: cliCtx,
		keys:   keys,
	}
}

func (api *FaucetAPI) RequestFunds(address common.Address, amount *hexutil.Big) error {
	return nil
}
