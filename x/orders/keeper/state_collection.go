package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/ethermint/x/orders/types"
	log "github.com/xlab/suplog"
)

// Returns EvmSyncStatus from keeper for spot market
func (k Keeper) GetEvmSyncStatus(ctx sdk.Context) *types.EvmSyncStatus {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.EvmSyncStatusKey())
	if bz == nil {
		return &types.EvmSyncStatus{}
	}

	var status types.EvmSyncStatus
	k.cdc.MustUnmarshalBinaryBare(bz, &status)

	return &status
}

// Returns EvmSyncStatus from keeper for the futures market
func (k Keeper) GetFuturesEvmSyncStatus(ctx sdk.Context) *types.EvmSyncStatus {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.EvmFuturesSyncStatusKey())
	if bz == nil {
		return &types.EvmSyncStatus{}
	}

	var status types.EvmSyncStatus
	k.cdc.MustUnmarshalBinaryBare(bz, &status)

	return &status
}

// Sets EvmSyncStatus in keeper for spot market
func (k Keeper) SetEvmSyncStatus(ctx sdk.Context, status *types.EvmSyncStatus) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(status)
	if bz == nil {
		log.Errorln("expected to marshal EvmSyncStatus but got nil")
		return
	}

	store.Set(types.EvmSyncStatusKey(), bz)
}

// Sets EvmSyncStatus in keeper for the futures market
func (k Keeper) SetFuturesEvmSyncStatus(ctx sdk.Context, status *types.EvmSyncStatus) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(status)
	if bz == nil {
		log.Errorln("expected to marshal EvmSyncStatus but got nil")
		return
	}

	store.Set(types.EvmFuturesSyncStatusKey(), bz)
}
