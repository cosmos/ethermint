package types

// Orders module event types
const (
	EventTypeNewOrder           = "new_order"
	EventTypeNewDerivativeOrder = "new_derivative_order"
	EventTypeSoftCancelOrder    = "soft_cancel"

	AttributeKeyOrderHash     = "order_hash"
	AttributeKeyMarketID      = "market_id"
	AttributeKeyTradePairHash = "trade_pair_hash"
	AttributeKeySignedOrder   = "signed_order"
	AttributeKeyFilledAmount  = "filled_amount"
)
