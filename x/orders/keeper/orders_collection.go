package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/ethermint/x/orders/types"
	"github.com/ethereum/go-ethereum/common"
)

// Returns Order from hash
func (k Keeper) GetOrder(ctx sdk.Context, hash common.Hash) *types.Order {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.ActiveOrdersStoreKey(hash))
	if bz == nil {
		bz = store.Get(types.ArchiveOrdersStoreKey(hash))
		if bz == nil {
			return nil
		}
	}

	var order types.Order
	k.cdc.MustUnmarshalBinaryBare(bz, &order)

	return &order
}

// Returns active Order from hash
func (k Keeper) GetActiveOrder(ctx sdk.Context, hash common.Hash) *types.Order {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.ActiveOrdersStoreKey(hash))
	if bz == nil {
		return nil
	}

	var order types.Order
	k.cdc.MustUnmarshalBinaryBare(bz, &order)

	return &order
}

// Returns archive Order from hash
func (k Keeper) GetArchiveOrder(ctx sdk.Context, hash common.Hash) *types.Order {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.ArchiveOrdersStoreKey(hash))
	if bz == nil {
		return nil
	}

	var order types.Order
	k.cdc.MustUnmarshalBinaryBare(bz, &order)

	return &order
}

// Returns array of all orders
func (k Keeper) GetAllOrders(ctx sdk.Context) []*types.Order {
	orders := []*types.Order{}
	appendOrder := func(order *types.Order) (stop bool) {
		orders = append(orders, order)
		return false
	}
	k.IterateActiveOrders(ctx, appendOrder)
	k.IterateArchiveOrders(ctx, appendOrder)
	return orders
}

// Returns array of all active orders
func (k Keeper) GetActiveOrders(ctx sdk.Context) []*types.Order {
	orders := []*types.Order{}
	appendOrder := func(order *types.Order) (stop bool) {
		orders = append(orders, order)
		return false
	}
	k.IterateActiveOrders(ctx, appendOrder)
	return orders
}

// Returns array of all archive orders
func (k Keeper) GetArchiveOrders(ctx sdk.Context) []*types.Order {
	orders := []*types.Order{}
	appendOrder := func(order *types.Order) (stop bool) {
		orders = append(orders, order)
		return false
	}
	k.IterateArchiveOrders(ctx, appendOrder)
	return orders
}

// Stores Order order in keeper
func (k Keeper) SetOrder(ctx sdk.Context, order *types.Order) {

	hash, err := order.Order.ToSignedOrder().ComputeOrderHash()
	if err != nil {
		k.Logger(ctx).Error("failed to compute order hash:", "error", err.Error())
		return
	}

	var key []byte
	if isActiveStatus(types.OrderStatus(order.Status)) {
		key = types.ActiveOrdersStoreKey(hash)
	} else {
		key = types.ArchiveOrdersStoreKey(hash)
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(order)
	store.Set(key, bz)
}

// Update OrderStatus of order from the provided order hash in the active collection
func (k Keeper) SetActiveOrderStatus(ctx sdk.Context, hash common.Hash, status types.OrderStatus) {
	order := k.GetActiveOrder(ctx, hash)
	if order == nil {
		k.Logger(ctx).Error("no active order found", "hash", hash.String())
		return
	}

	order.Status = int64(status)
	if !isActiveStatus(status) {
		// status is not active anymore
		store := ctx.KVStore(k.storeKey)
		store.Delete(types.ActiveOrdersStoreKey(hash))
	}

	k.SetOrder(ctx, order)
}

// Update OrderStatus of order from the provided order hash in the archive collection
func (k Keeper) SetArchiveOrderStatus(ctx sdk.Context, hash common.Hash, status types.OrderStatus) {
	order := k.GetArchiveOrder(ctx, hash)
	if order == nil {
		k.Logger(ctx).Error("order not found", "hash", hash.String())
		return
	} else if isActiveStatus(status) {
		k.Logger(ctx).Error("incorrect status transition for archived order", "status", status.String())
		return
	}

	order.Status = int64(status)
	k.SetOrder(ctx, order)
}

// Returns true if status is an active status, false otherwise
func isActiveStatus(status types.OrderStatus) bool {
	switch status {
	case types.StatusUnfilled, types.StatusPartialFilled:
		return true
	default:
		return false
	}
}

// Iterates over active Orders calling process on each one.
func (k Keeper) IterateActiveOrders(ctx sdk.Context, process func(*types.Order) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ActiveOrdersStoreKeyPrefix)
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}

		bz := iter.Value()
		var order types.Order
		k.cdc.MustUnmarshalBinaryBare(bz, &order)

		if process(&order) {
			return
		}

		iter.Next()
	}
}

// Iterates over archive Orders calling process on each one.
func (k Keeper) IterateArchiveOrders(ctx sdk.Context, process func(*types.Order) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ArchiveOrdersStoreKeyPrefix)
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}

		bz := iter.Value()
		var order types.Order
		k.cdc.MustUnmarshalBinaryBare(bz, &order)

		if process(&order) {
			return
		}

		iter.Next()
	}
}
