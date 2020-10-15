package types

import "github.com/ethereum/go-ethereum/common"

const (
	QueryOrder         = "order"
	QueryPair          = "pair"
	QueryActiveOrder   = "activeOrder"
	QueryArchiveOrder  = "archiveOrder"
	QueryOrdersList    = "ordersList"
	QueryPairsList     = "pairsList"
	QueryMarketsList   = "marketsList"
	QueryEvmSyncStatus = "evmSyncStatus"

	QueryZeroExTransaction       = "queryZeroExTransaction"
	QuerySoftCancelledOrders     = "querySoftCancelledOrders"
	QueryOutstandingFillRequests = "queryOutstandingFillRequests"
	QueryOrderFillRequests       = "queryOrderFillRequests"
)

// QueryOrderParams defines the params for querying an order.
type QueryOrderParams struct {
	Hash common.Hash `json:"hash"`
}

type QueryOrderResponse struct {
	Order *Order `json:"order,omitempty"`
}

// QueryActiveOrderParams defines the params for querying an active order.
type QueryActiveOrderParams struct {
	Hash common.Hash `json:"hash"`
}

type QueryActiveOrderResponse struct {
	Order *Order `json:"order,omitempty"`
}

// QueryArchiveOrderParams defines the params for querying an archive order.
type QueryArchiveOrderParams struct {
	Hash common.Hash `json:"hash"`
}

type QueryArchiveOrderResponse struct {
	Order *Order `json:"order,omitempty"`
}

// QueryOrdersListParams defines the params for querying an order list.
type QueryOrdersListParams struct {
	ByStatus        *OrderStatus        `json:"status"`
	ByCollection    OrderCollectionType `json:"collection"`
	ByTradePairHash *common.Hash        `json:"pair_hash"`
	ByOrderFilters  *OrderFilters       `json:"filters"`
	Limit           int                 `json:"limit"`
}

type OrderFilters struct {
	// for derivatives
	ContractPriceBound *string `json:"contractPriceBound"`
	MarketID           *[]byte `json:"marketID"`
	IsLong             bool    `json:"isLong"`
	// for normal orders
	NotExpired          bool            `json:"notExpired"`
	MakerAssetProxyID   *[]byte         `json:"makerAssetProxyId"`
	TakerAssetProxyID   *[]byte         `json:"takerAssetProxyId"`
	MakerAssetAddress   *common.Address `json:"makerAssetAddress"`
	TakerAssetAddress   *common.Address `json:"takerAssetAddress"`
	ExchangeAddress     *common.Address `json:"exchangeAddress"`
	SenderAddress       *common.Address `json:"senderAddress"`
	MakerAssetData      *[]byte         `json:"makerAssetData"`
	TakerAssetData      *[]byte         `json:"takerAssetData"`
	MakerFeeAssetData   *[]byte         `json:"makerFeeAssetData"`
	TakerFeeAssetData   *[]byte         `json:"takerFeeAssetData"`
	TraderAssetData     *[]byte         `json:"traderAssetData"`
	MakerAssetAmount    *string         `json:"makerAssetAmount"`
	TakerAssetAmount    *string         `json:"takerAssetAmount"`
	MakerAddress        *common.Address `json:"makerAddress"`
	NotMakerAddress     *common.Address `json:"notMakerAddress"`
	TakerAddress        *common.Address `json:"takerAddress"`
	TraderAddress       *common.Address `json:"traderAddress"`
	FeeRecipientAddress *common.Address `json:"feeRecipientAddress"`
}

type QueryOrdersListResponse struct {
	Orders []*Order `json:"orders,omitempty"`
}

type OrderCollectionType string

const (
	OrderCollectionAny     OrderCollectionType = ""
	OrderCollectionActive  OrderCollectionType = "active"
	OrderCollectionArchive OrderCollectionType = "archive"
)

// TODO: use
//// QueryMarketParams defines the params for querying a derivative market.
//type QueryMarketParams struct {
//	Name           string      `json:"name"`
//	Hash           common.Hash `json:"hash"`
//	MakerAssetData []byte      `json:"makerAssetData"`
//	TakerAssetData []byte      `json:"takerAssetData"`
//}

// QueryPairsListParams defines the params for querying existing trade pairs.
type QueryMarketsListParams struct {
	All bool `json:"all"`
}

type QueryMarketsListResponse struct {
	Markets []*DerivativeMarket `json:"markets"`
}

// QueryPairsParams defines the params for querying a trade pair.
type QueryPairParams struct {
	Name           string      `json:"name"`
	Hash           common.Hash `json:"hash"`
	MakerAssetData []byte      `json:"makerAssetData"`
	TakerAssetData []byte      `json:"takerAssetData"`
}

type QueryPairResponse struct {
	Pair *TradePair `json:"pair"`
}

// QueryPairsListParams defines the params for querying existing trade pairs.
type QueryPairsListParams struct {
	All bool `json:"all"`
}

type QueryPairsListResponse struct {
	Pairs []*TradePair `json:"pairs"`
}

type QueryZeroExTransactionParams struct {
	TxHash common.Hash `json:"txHash"`
}

type QueryZeroExTransactionResponse struct {
	TxType             ZeroExTransactionType     `json:"txType"`
	FillRequests       []*OrderFillRequest       `json:"fillRequests,omitempty"`
	SoftCancelRequests []*OrderSoftCancelRequest `json:"softCancelRequests,omitempty"`
}

type QueryOutstandingFillRequestsParams struct {
	TxHash common.Hash `json:"txHash"`
}

type QueryOutstandingFillRequestsResponse struct {
	FillRequests []*OrderFillRequest `json:"fillRequests"`
}

type QuerySoftCancelledOrdersParams struct {
	OrderHashes []common.Hash `json:"orderHashes"`
}

type QuerySoftCancelledOrdersResponse struct {
	OrderHashes []common.Hash `json:"orderHashes"`
}

type QueryOrderFillRequestsParams struct {
	OrderHash common.Hash `json:"orderHash"`
}

type QueryOrderFillRequestsResponse struct {
	FillRequests []*OrderFillRequest `json:"fillRequests"`
}

type QueryEvmSyncStatusParams struct {
}

type QueryEvmSyncStatusResponse struct {
	LatestBlockSynced uint64 `json:"latestBlock"`
}
