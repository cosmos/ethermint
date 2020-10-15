package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/ethermint/x/orders/types"
	"github.com/ethereum/go-ethereum/common"
)

// Returns Derivative Market from hash
func (k Keeper) GetDerivativeMarket(ctx sdk.Context, hash string) *types.DerivativeMarket {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DerivativeMarketStoreKey(hash))
	if bz == nil {
		return nil
	}
	var market types.DerivativeMarket
	k.cdc.MustUnmarshalBinaryBare(bz, &market)
	return &market
}

// SetDerivativeMarket saves derivative market in keeper.
func (k Keeper) SetDerivativeMarket(ctx sdk.Context, market *types.DerivativeMarket) {
	hash, err := market.Hash()
	if err != nil {
		k.Logger(ctx).Error("failed to compute tradePair hash:", "error", err.Error())
		return
	}
	market.MarketId = hash.Hex()
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(market)
	store.Set(types.DerivativeMarketStoreKey(market.MarketId), bz)
}

// SetBaseCurrencyInDerivativeMarket changes base currency address in a derivative market.
func (k Keeper) SetBaseCurrencyInDerivativeMarket(ctx sdk.Context, hash common.Hash, baseCurrency common.Address) {
	m := k.GetDerivativeMarket(ctx, hash.Hex())
	if m == nil {
		return
	}

	m.BaseCurrency = baseCurrency.Hex()
	k.SetDerivativeMarket(ctx, m)
}

// SetOracleInDerivativeMarket changes oracle address in a derivative market.
func (k Keeper) SetOracleInDerivativeMarket(ctx sdk.Context, hash common.Hash, oracle common.Address) {
	m := k.GetDerivativeMarket(ctx, hash.Hex())
	if m == nil {
		return
	}

	m.Oracle = oracle.Hex()
	k.SetDerivativeMarket(ctx, m)
}

// Iterates over derivative markets calling process on each market
func (k Keeper) IterateDerivativeMarkets(ctx sdk.Context, process func(*types.DerivativeMarket) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.DerivativeMarketStoreKeyPrefix)
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		bz := iter.Value()
		var market types.DerivativeMarket
		k.cdc.MustUnmarshalBinaryBare(bz, &market)
		if process(&market) {
			return
		}
		iter.Next()
	}
}

// Returns all derivative markets
func (k Keeper) GetAllDerivativeMarkets(ctx sdk.Context) []*types.DerivativeMarket {
	markets := []*types.DerivativeMarket{}
	appendMarket := func(p *types.DerivativeMarket) (stop bool) {
		markets = append(markets, p)
		return false
	}
	k.IterateDerivativeMarkets(ctx, appendMarket)
	return markets
}

// Sets Derivative Market status to Enabled in keeper
func (k Keeper) SetMarketEnabled(ctx sdk.Context, hash common.Hash, enabled bool) {
	m := k.GetDerivativeMarket(ctx, hash.Hex())
	if m == nil {
		k.Logger(ctx).Error("derivative market not found", "marketID", hash.String())
		return
	} else if m.Enabled == enabled {
		return
	}
	m.Enabled = enabled
	k.SetDerivativeMarket(ctx, m)
}
