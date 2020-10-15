package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/ethermint/x/orders/internal/types"
)

// GETTERS SECTION

// GetOrderFillRequest gets an OrderFillRequest by order hash and particular Tx hash.
func (k Keeper) GetOrderFillRequest(
	ctx sdk.Context,
	orderHash, txHash common.Hash,
) *types.OrderFillRequest {
	store := ctx.KVStore(k.storeKey)

	orderFillRequestKey := append(types.OrderFillRequestsStorePrefix(orderHash), txHash.Bytes()...)
	bz := store.Get(orderFillRequestKey)
	if bz == nil {
		return nil
	}

	var req types.OrderFillRequest
	k.cdc.MustUnmarshalBinaryBare(bz, &req)

	return &req
}

// GetOrderSoftCancelRequest gets an OrderSoftCancelRequest by order hash and particular Tx hash.
func (k Keeper) GetOrderSoftCancelRequest(
	ctx sdk.Context,
	orderHash, txHash common.Hash,
) *types.OrderSoftCancelRequest {
	store := ctx.KVStore(k.storeKey)

	orderSoftCancelRequestKey := append(types.OrderSoftCancelRequestsStorePrefix(orderHash), txHash.Bytes()...)
	bz := store.Get(orderSoftCancelRequestKey)
	if bz == nil {
		return nil
	}

	var req types.OrderSoftCancelRequest
	k.cdc.MustUnmarshalBinaryBare(bz, &req)

	return &req
}

// ListOrderFillRequests returns a list of OrderFillRequests related to specific order hash.
func (k Keeper) ListOrderFillRequests(
	ctx sdk.Context,
	orderHash common.Hash,
) (results []*types.OrderFillRequest) {
	k.IterateOrderFillRequestsByOrderHash(ctx, orderHash,
		func(fillReq *types.OrderFillRequest) (stop bool) {

			results = append(results, fillReq)

			return false
		})

	return results
}

// ListOrderFillRequestsByTxHash returns a list of OrderFillRequests related to specific tx hash.
func (k Keeper) ListOrderFillRequestsByTxHash(
	ctx sdk.Context,
	txHash common.Hash,
) (results []*types.OrderFillRequest) {
	k.IterateOrderFillRequestsByTxHash(ctx, txHash,
		func(fillReq *types.OrderFillRequest) (stop bool) {

			results = append(results, fillReq)

			return false
		})

	return results
}

// ListOrderSoftCancelRequests returns a list of OrderSoftCancelRequests related to specific order hash.
func (k Keeper) ListOrderSoftCancelRequests(
	ctx sdk.Context,
	orderHash common.Hash,
) (results []*types.OrderSoftCancelRequest) {
	k.IterateOrderSoftCancelRequestsByOrderHash(ctx, orderHash,
		func(cancelReq *types.OrderSoftCancelRequest) (stop bool) {

			results = append(results, cancelReq)

			return false
		})

	return results
}

func (k Keeper) FindAllSoftCancelledOrders(
	ctx sdk.Context,
	orderHashes []common.Hash,
) (cancelled []common.Hash) {
	for _, hash := range orderHashes {
		k.IterateOrderSoftCancelRequestsByOrderHash(ctx, hash,
			func(cancelReq *types.OrderSoftCancelRequest) (stop bool) {
				if cancelReq != nil {
					cancelled = append(cancelled, hash)
				}

				return false
			})
	}

	return
}

// ListOrderSoftCancelRequestsByTxHash returns a list of OrderSoftCancelRequests related to specific tx hash.
func (k Keeper) ListOrderSoftCancelRequestsByTxHash(
	ctx sdk.Context,
	txHash common.Hash,
) (results []*types.OrderSoftCancelRequest) {
	k.IterateOrderSoftCancelRequestsByTxHash(ctx, txHash,
		func(cancelReq *types.OrderSoftCancelRequest) (stop bool) {

			results = append(results, cancelReq)

			return false
		})

	return results
}

// GetZeroExTransaction gets an ZeroExTransaction by Tx hash.
func (k Keeper) GetZeroExTransaction(
	ctx sdk.Context,
	txHash common.Hash,
) *types.ZeroExTransaction {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.ZeroExTransactionsStoreKey(txHash))
	if bz == nil {
		return nil
	}

	var req types.ZeroExTransaction
	k.cdc.MustUnmarshalBinaryBare(bz, &req)

	return &req
}

// SETTERS SECTION

// Stores OrderFillRequest object in keeper by order hash and tx hash.
func (k Keeper) SetOrderFillRequest(
	ctx sdk.Context,
	txHash common.Hash,
	req *types.OrderFillRequest,
) {
	store := ctx.KVStore(k.storeKey)
	keyPrefix := types.OrderFillRequestsStorePrefix(req.OrderHash)

	bz := k.cdc.MustMarshalBinaryBare(req)

	key := append(keyPrefix, txHash.Bytes()...)
	store.Set(key, bz)
}

// Stores OrderSoftCancelRequest object in keeper by order hash and tx hash.
func (k Keeper) SetOrderSoftCancelRequest(
	ctx sdk.Context,
	txHash common.Hash,
	req *types.OrderSoftCancelRequest,
) {
	store := ctx.KVStore(k.storeKey)
	keyPrefix := types.OrderSoftCancelRequestsStorePrefix(req.OrderHash)

	bz := k.cdc.MustMarshalBinaryBare(req)

	key := append(keyPrefix, txHash.Bytes()...)
	store.Set(key, bz)
}

// SetZeroExTransaction stores ZeroExTransaction object in keeper.
func (k Keeper) SetZeroExTransaction(
	ctx sdk.Context,
	txHash common.Hash,
	tx *types.ZeroExTransaction,
) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(tx)
	store.Set(types.ZeroExTransactionsStoreKey(txHash), bz)
}

// ITERATORS SECTION

// Iterates over OrderFillRequests of particular orderHash.
func (k Keeper) IterateOrderFillRequestsByOrderHash(
	ctx sdk.Context,
	orderHash common.Hash,
	process func(*types.OrderFillRequest) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)
	prefix := types.OrderFillRequestsStorePrefix(orderHash)
	iter := sdk.KVStorePrefixIterator(store, prefix)
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}

		bz := iter.Value()
		var req types.OrderFillRequest
		k.cdc.MustUnmarshalBinaryBare(bz, &req)

		if process(&req) {
			return
		}

		iter.Next()
	}
}

// IterateOrderFillRequestsByTxHash iterates over OrderFillRequests of particular txHash.
func (k Keeper) IterateOrderFillRequestsByTxHash(
	ctx sdk.Context,
	txHash common.Hash,
	process func(*types.OrderFillRequest) (stop bool),
) {
	tx := k.GetZeroExTransaction(ctx, txHash)
	if tx.Type != types.ZeroExOrderFillRequestTx {
		return
	}

	for _, orderHash := range tx.Orders {
		req := k.GetOrderFillRequest(ctx, orderHash, txHash)
		if req != nil {
			if process(req) {
				return
			}
		}
	}
}

// Iterates over OrderSoftCancelRequests of particular orderHash.
func (k Keeper) IterateOrderSoftCancelRequestsByOrderHash(
	ctx sdk.Context,
	orderHash common.Hash,
	process func(*types.OrderSoftCancelRequest) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)
	prefix := types.OrderSoftCancelRequestsStorePrefix(orderHash)
	iter := sdk.KVStorePrefixIterator(store, prefix)
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}

		bz := iter.Value()
		var req types.OrderSoftCancelRequest
		k.cdc.MustUnmarshalBinaryBare(bz, &req)

		if process(&req) {
			return
		}

		iter.Next()
	}
}

// IterateOrderSoftCancelRequestsByTxHash iterates over OrderSoftCancelRequests of particular txHash.
func (k Keeper) IterateOrderSoftCancelRequestsByTxHash(
	ctx sdk.Context,
	txHash common.Hash,
	process func(*types.OrderSoftCancelRequest) (stop bool),
) {
	tx := k.GetZeroExTransaction(ctx, txHash)
	if tx.Type != types.ZeroExOrderSoftCancelRequestTx {
		return
	}

	for _, orderHash := range tx.Orders {
		req := k.GetOrderSoftCancelRequest(ctx, orderHash, txHash)
		if req != nil {
			if process(req) {
				return
			}
		}
	}
}
