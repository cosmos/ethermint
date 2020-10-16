package orders

import (
	"github.com/cosmos/ethermint/x/orders/keeper"
	"github.com/cosmos/ethermint/x/orders/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
)

var (
	NewKeeper = keeper.NewKeeper
)

type (
	Keeper = keeper.Keeper
)
