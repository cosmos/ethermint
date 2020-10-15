package orders

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ethermint/x/orders/types"
)

func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	for _, tradePair := range data.TradePairs {
		keeper.SetTradePair(ctx, tradePair)
	}
	for _, market := range data.DerivativeMarkets {
		keeper.SetDerivativeMarket(ctx, market)
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) *types.GenesisState {
	return &types.GenesisState{
		TradePairs:        k.GetAllTradePairs(ctx),
		DerivativeMarkets: k.GetAllDerivativeMarkets(ctx),
	}
}
