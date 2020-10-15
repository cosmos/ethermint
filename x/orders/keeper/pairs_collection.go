package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/ethermint/x/orders/types"
	"github.com/ethereum/go-ethereum/common"
)

// Returns TradePair from hash
func (k Keeper) GetTradePair(ctx sdk.Context, hash common.Hash) *types.TradePair {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TradePairsStoreKey(hash))
	var pair types.TradePair
	k.cdc.MustUnmarshalBinaryBare(bz, &pair)
	return &pair
}

// Returns TradePair from name
func (k Keeper) GetTradePairByName(ctx sdk.Context, name string) *types.TradePair {
	var tradePair *types.TradePair
	matchPair := func(p *types.TradePair) (stop bool) {
		if p.Name == name {
			tradePair = p
			return true
		}
		return false
	}
	k.IterateTradePairs(ctx, matchPair)
	return tradePair
}

// Returns all TradePairs
func (k Keeper) GetAllTradePairs(ctx sdk.Context) []*types.TradePair {
	tradePairs := []*types.TradePair{}
	appendPair := func(p *types.TradePair) (stop bool) {
		tradePairs = append(tradePairs, p)
		return false
	}
	k.IterateTradePairs(ctx, appendPair)
	return tradePairs
}

// Sets TradePair in keeper
func (k Keeper) SetTradePair(ctx sdk.Context, tradePair *types.TradePair) {
	hash, err := tradePair.ComputeHash()
	if err != nil {
		k.Logger(ctx).Error("failed to compute tradePair hash:", "error", err.Error())
		return
	}
	tradePair.Hash = hash.Hex()
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(tradePair)
	store.Set(types.TradePairsStoreKey(hash), bz)
}

// Deletes TradePair from keeper (needed for moving to another hash)
func (k Keeper) DeleteTradePair(ctx sdk.Context, hash common.Hash) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.TradePairsStoreKey(hash))
}

// Sets TradePair status to Enabled in keeper
func (k Keeper) SetTradePairEnabled(ctx sdk.Context, hash common.Hash, enabled bool) {
	tradePair := k.GetTradePair(ctx, hash)
	if tradePair == nil {
		k.Logger(ctx).Error("trade pair not found", "hash", hash.String())
		return
	} else if tradePair.Enabled == enabled {
		return
	}
	tradePair.Enabled = enabled
	k.SetTradePair(ctx, tradePair)
}

// Iterates over TradePairs calling process on each pair
func (k Keeper) IterateTradePairs(ctx sdk.Context, process func(*types.TradePair) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.TradePairsStoreKeyPrefix)
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		bz := iter.Value()
		var tradePair types.TradePair
		k.cdc.MustUnmarshalBinaryBare(bz, &tradePair)
		if process(&tradePair) {
			return
		}
		iter.Next()
	}
}
