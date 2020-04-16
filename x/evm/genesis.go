package evm

import (
	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"

	emint "github.com/cosmos/ethermint/types"
	"github.com/cosmos/ethermint/x/evm/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis initializes genesis state based on exported genesis
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, record := range data.Accounts {
		k.CreateGenesisAccount(ctx, record)
	}
	return []abci.ValidatorUpdate{}
}

// ExportGenesis exports genesis state
func ExportGenesis(ctx sdk.Context, k Keeper, ak types.AccountKeeper) GenesisState {
	var ethGenAccounts []types.GenesisAccount

	accounts := ak.GetAllAccounts(ctx)

	for _, account := range accounts {
		ethAccount, ok := account.(emint.Account)
		if !ok {
			continue
		}

		addr := common.BytesToAddress(ethAccount.GetAddress().Bytes())

		genAccount := types.GenesisAccount{
			Address: addr,
			Balance: k.GetBalance(ctx, addr),
			Code:    k.GetCode(ctx, addr),
			// Storage: k.GetStorage(ctx, addr), TODO:
		}

		ethGenAccounts = append(ethGenAccounts, genAccount)

	}

	return GenesisState{Accounts: ethGenAccounts}
}
