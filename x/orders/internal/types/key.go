package types

import (
	"github.com/ethereum/go-ethereum/common"
)

const (
	// module name
	ModuleName = "orders"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
)

var (
	// ActiveOrdersStoreKeyPrefix for order-by-hash store of active orders.
	ActiveOrdersStoreKeyPrefix = []byte{0x01}
	// ArchiveOrdersStoreKeyPrefix for order-by-hash store of archive orders.
	ArchiveOrdersStoreKeyPrefix = []byte{0x02}

	OrderFillRequestsStoreKeyPrefix       = []byte{0x03}
	OrderSoftCancelRequestsStoreKeyPrefix = []byte{0x04}
	ZeroExTransactionsStoreKeyPrefix      = []byte{0x05}

	OrderFillEventStoreKeyPrefix   = []byte{0x06}
	OrderCancelEventStoreKeyPrefix = []byte{0x07}

	// TradePairsStoreKeyPrefix for pair-by-hash store of trade pairs with asset data.
	TradePairsStoreKeyPrefix       = []byte{0x08}
	DerivativeMarketStoreKeyPrefix = []byte{0x09}

	FuturesPositionFillEventStoreKeyPrefix   = []byte{0x10}
	FuturesPositionCancelEventStoreKeyPrefix = []byte{0x11}

	// ModuleStatePrefix defines a collection of module state-related entries.
	ModuleStatePrefix = []byte{0xff}
)

// ActiveOrdersStoreKey turns a order hash to key used to get it from the store.
func ActiveOrdersStoreKey(orderHash common.Hash) []byte {
	return append(ActiveOrdersStoreKeyPrefix, orderHash.Bytes()...)
}

// ArchiveOrdersStoreKey turns an order hash to key used to get it from the store.
func ArchiveOrdersStoreKey(orderHash common.Hash) []byte {
	return append(ArchiveOrdersStoreKeyPrefix, orderHash.Bytes()...)
}

// ZeroExTransactionsStorePrefix allows to obtain key for zeroex transaction by its hash.
func ZeroExTransactionsStoreKey(txHash common.Hash) []byte {
	return append(ZeroExTransactionsStoreKeyPrefix, txHash.Bytes()...)
}

// OrderFillRequestsStorePrefix allows to obtain prefix of fill requests against particular order hash.
func OrderFillRequestsStorePrefix(orderHash common.Hash) []byte {
	return append(OrderFillRequestsStoreKeyPrefix, orderHash.Bytes()...)
}

// OrderSoftCancelRequestsStorePrefix allows to obtain prefix of soft cancel requests against particular order hash.
func OrderSoftCancelRequestsStorePrefix(orderHash common.Hash) []byte {
	return append(OrderSoftCancelRequestsStoreKeyPrefix, orderHash.Bytes()...)
}

// OrderFillEventStoreKey turns a hash to key used to get it from the store.
func OrderFillEventStoreKey(orderHash, eventHash common.Hash) []byte {
	return append(append(OrderFillEventStoreKeyPrefix, orderHash.Bytes()...), eventHash.Bytes()...)
}

// OrderCancelEventStoreKey turns a hash to key used to get it from the store.
func OrderCancelEventStoreKey(orderHash, eventHash common.Hash) []byte {
	return append(append(OrderCancelEventStoreKeyPrefix, orderHash.Bytes()...), eventHash.Bytes()...)
}

// FuturesPositionFillEventStoreKey turns a hash to key used to get it from the store.
func FuturesPositionFillEventStoreKey(orderHash, eventHash common.Hash, isLong bool) []byte {
	return append(append(FuturesPositionFillEventStoreKeyPrefix, orderHash.Bytes()...), eventHash.Bytes()...)
}

// FuturesPositionCancelEventStoreKey turns a hash to key used to get it from the store.
func FuturesPositionCancelEventStoreKey(orderHash, eventHash common.Hash) []byte {
	return append(append(FuturesPositionCancelEventStoreKeyPrefix, orderHash.Bytes()...), eventHash.Bytes()...)
}

// TradePairsStoreKey turns a pair hash to key used to get it from the store.
func TradePairsStoreKey(hash common.Hash) []byte {
	return append(TradePairsStoreKeyPrefix, hash.Bytes()...)
}

// EvmSyncStatusKey returns a key for EVM sync status (spot markets) data entry.
func EvmSyncStatusKey() []byte {
	return append(ModuleStatePrefix, []byte("evm_sync_status")...)
}

// EvmFuturesSyncStatusKey returns a key for EVM sync status (futures markets) data entry.
func EvmFuturesSyncStatusKey() []byte {
	return append(ModuleStatePrefix, []byte("evm_futures_sync_status")...)
}

// DerivativeMarketStoreKey turns a pair hash to key used to get it from the store.
func DerivativeMarketStoreKey(hash common.Hash) []byte {
	return append(DerivativeMarketStoreKeyPrefix, hash.Bytes()...)
}
