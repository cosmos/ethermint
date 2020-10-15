package orders

import (
	"github.com/cosmos/ethermint/x/orders/internal/client/cli"
	"github.com/cosmos/ethermint/x/orders/internal/keeper"
	"github.com/cosmos/ethermint/x/orders/internal/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey

	QueryOrder        = types.QueryOrder
	QueryActiveOrder  = types.QueryActiveOrder
	QueryArchiveOrder = types.QueryArchiveOrder
	QueryOrdersList   = types.QueryOrdersList
	QueryPair         = types.QueryPair
	QueryPairsList    = types.QueryPairsList
	QueryMarketsList  = types.QueryMarketsList

	QueryZeroExTransaction       = types.QueryZeroExTransaction
	QuerySoftCancelledOrders     = types.QuerySoftCancelledOrders
	QueryOutstandingFillRequests = types.QueryOutstandingFillRequests
	QueryOrderFillRequests       = types.QueryOrderFillRequests

	StatusUnfilled      = types.StatusUnfilled
	StatusSoftCancelled = types.StatusSoftCancelled
	StatusPartialFilled = types.StatusPartialFilled
	StatusFilled        = types.StatusFilled
	StatusExpired       = types.StatusExpired

	ZeroExOrderFillRequestTx       = types.ZeroExOrderFillRequestTx
	ZeroExOrderSoftCancelRequestTx = types.ZeroExOrderSoftCancelRequestTx
)

var (
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier
	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc
	NewTxCmd      = cli.NewTxCmd

	ActiveOrdersStoreKey               = types.ActiveOrdersStoreKey
	ArchiveOrdersStoreKey              = types.ArchiveOrdersStoreKey
	OrderFillRequestsStorePrefix       = types.OrderFillRequestsStorePrefix
	OrderSoftCancelRequestsStorePrefix = types.OrderSoftCancelRequestsStorePrefix
	ZeroExTransactionsStoreKey         = types.ZeroExTransactionsStoreKey
	OrderFillEventStoreKey             = types.OrderFillEventStoreKey
	OrderCancelEventStoreKey           = types.OrderCancelEventStoreKey
	FuturesPositionFillEventStoreKey   = types.FuturesPositionFillEventStoreKey
	FuturesPositionCancelEventStoreKey = types.FuturesPositionCancelEventStoreKey
	TradePairsStoreKey                 = types.TradePairsStoreKey
	DerivativeMarketStoreKey           = types.DerivativeMarketStoreKey

	OrderStatusFromString  = types.OrderStatusFromString
	NewMsgSignedOrder      = types.NewSafeSignedOrder
	OrderCollectionAny     = types.OrderCollectionAny
	OrderCollectionActive  = types.OrderCollectionActive
	OrderCollectionArchive = types.OrderCollectionArchive
)

type (
	Keeper = keeper.Keeper

	MsgCreateSpotOrder            = types.MsgCreateSpotOrder
	MsgCreateDerivativeOrder      = types.MsgCreateDerivativeOrder
	MsgRequestFillSpotOrder       = types.MsgRequestFillSpotOrder
	MsgRequestSoftCancelSpotOrder = types.MsgRequestSoftCancelSpotOrder
	MsgFilledSpotOrder            = types.MsgFilledSpotOrder
	MsgCancelledSpotOrder         = types.MsgCancelledSpotOrder
	MsgFilledDerivativeOrder      = types.MsgFilledDerivativeOrder
	MsgCancelledDerivativeOrder   = types.MsgCancelledDerivativeOrder

	MsgRegisterSpotMarket       = types.MsgRegisterSpotMarket
	MsgSuspendSpotMarket        = types.MsgSuspendSpotMarket
	MsgResumeSpotMarket         = types.MsgResumeSpotMarket
	MsgRegisterDerivativeMarket = types.MsgRegisterDerivativeMarket
	MsgSuspendDerivativeMarket  = types.MsgSuspendDerivativeMarket
	MsgResumeDerivativeMarket   = types.MsgResumeDerivativeMarket

	QueryOrderParams        = types.QueryOrderParams
	QueryActiveOrderParams  = types.QueryActiveOrderParams
	QueryArchiveOrderParams = types.QueryArchiveOrderParams
	QueryOrdersListParams   = types.QueryOrdersListParams
	QueryPairParams         = types.QueryPairParams
	QueryPairsListParams    = types.QueryPairsListParams
	QueryMarketParams       = types.QueryPairParams
	QueryMarketsListParams  = types.QueryMarketsListParams

	QueryZeroExTransactionParams       = types.QueryZeroExTransactionParams
	QuerySoftCancelledOrdersParams     = types.QuerySoftCancelledOrdersParams
	QueryOutstandingFillRequestsParams = types.QueryOutstandingFillRequestsParams
	QueryOrderFillRequestsParams       = types.QueryOrderFillRequestsParams

	QueryOrderResponse        = types.QueryOrderResponse
	QueryActiveOrderResponse  = types.QueryActiveOrderResponse
	QueryArchiveOrderResponse = types.QueryArchiveOrderResponse
	QueryOrdersListResponse   = types.QueryOrdersListResponse
	QueryPairResponse         = types.QueryPairResponse
	QueryPairsListResponse    = types.QueryPairsListResponse
	QueryMarketsListResponse  = types.QueryMarketsListResponse

	QueryZeroExTransactionResponse       = types.QueryZeroExTransactionResponse
	QuerySoftCancelledOrdersResponse     = types.QuerySoftCancelledOrdersResponse
	QueryOutstandingFillRequestsResponse = types.QueryOutstandingFillRequestsResponse
	QueryOrderFillRequestsResponse       = types.QueryOrderFillRequestsResponse

	OrderStatus         = types.OrderStatus
	OrderCollectionType = types.OrderCollectionType

	Order                  = types.Order
	EvmSyncStatus          = types.EvmSyncStatus
	OrderFilters           = types.OrderFilters
	OrderFillRequest       = types.OrderFillRequest
	OrderSoftCancelRequest = types.OrderSoftCancelRequest
	SignedTransaction      = types.SignedTransaction
	ZeroExTransactionType  = types.ZeroExTransactionType
	CoordinatorDomain      = types.CoordinatorDomain
	SafeSignedOrder        = types.SafeSignedOrder
	TradePair              = types.TradePair
	DerivativeMarket       = types.DerivativeMarket
	Address                = types.Address
	HexBytes               = types.HexBytes
	BigNum                 = types.BigNum
	Hash                   = types.Hash
)
