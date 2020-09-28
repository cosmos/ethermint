package faucet

import (
	"github.com/cosmos/ethermint/x/faucet/keeper"
	"github.com/cosmos/ethermint/x/faucet/types"
)

const (
	ModuleName   = types.ModuleName
	RouterKey    = types.RouterKey
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute
)

var (
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier
	ModuleCdc     = types.ModuleCdc
	RegisterLegacyAminoCodec = types.RegisterLegacyAminoCodec
)

type (
	Keeper = keeper.Keeper
)
