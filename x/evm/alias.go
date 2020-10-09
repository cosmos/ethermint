package evm

import (
	"github.com/cosmos/ethermint/x/evm/keeper"
	"github.com/cosmos/ethermint/x/evm/types"
)

// nolint
const (
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	RouterKey         = types.RouterKey
	DefaultParamspace = types.DefaultParamspace
)

// nolint
var (
	NewKeeper = keeper.NewKeeper
)

//nolint
type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
)
