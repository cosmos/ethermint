package evm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ethermint/x/evm/keeper"
	"github.com/cosmos/ethermint/x/evm/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis initializes genesis state based on exported genesis
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) []abci.ValidatorUpdate {
	for _, record := range data.Accounts {
		k.SetCode(ctx, record.Address, record.Code)
		k.CreateGenesisAccount(ctx, record)
	}
	return []abci.ValidatorUpdate{}
}

// ExportGenesis exports genesis state
func ExportGenesis(ctx sdk.Context, _ keeper.Keeper) types.GenesisState {
	return types.GenesisState{Accounts: nil}
}
