package validator

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// OrderStatus encodes order status according to LibOrder.OrderStatus
type OrderStatus uint8

const (
	// OrderInvalid is the default value
	OrderInvalid OrderStatus = 0
	// OrderInsufficientMarginForContractPrice when Order does not have enough margin for contract price
	OrderInsufficientMarginForContractPrice OrderStatus = 1
	// OrderInsufficientMarginIndexPrice when Order does not have enough margin for index price
	OrderInsufficientMarginIndexPrice OrderStatus = 2
	// OrderFillable when order is fillable
	OrderFillable OrderStatus = 3
	// OrderExpired when order has already expired
	OrderExpired OrderStatus = 4
	// OrderFullyFilled when order is fully filled
	OrderFullyFilled OrderStatus = 5
	// OrderCancelled when order has been cancelled
	OrderCancelled OrderStatus = 6
	// Maker of the order does not have sufficient funds deposited to be filled.
	OrderUnfunded OrderStatus = 7
	// Index Price has not been triggered
	OrderUntriggered OrderStatus = 8
)

func (o OrderStatus) String() string {
	switch o {
	case OrderInvalid:
		return "invalid"
	case OrderInsufficientMarginForContractPrice:
		return "insufficientMarginForContractPrice"
	case OrderInsufficientMarginIndexPrice:
		return "insufficientMarginIndexPrice"
	case OrderFillable:
		return "fillable"
	case OrderExpired:
		return "expired"
	case OrderFullyFilled:
		return "fullyFilled"
	case OrderCancelled:
		return "cancelled"
	case OrderUnfunded:
		return "unfunded"
	case OrderUntriggered:
		return "untriggered"
	}

	return ""
}

// OrderState is a representation of the state returned by DevUtils ABI
type OrderState struct {
	Status                   OrderStatus `json:"status"`
	Hash                     common.Hash `json:"hash"`
	TakerAssetFilledAmount   *big.Int    `json:"takerAssetFilledAmount"`
	FillableTakerAssetAmount *big.Int    `json:"fillableTakerAssetAmount"`
	IsValidSignature         bool        `json:"isValidSignature"`
}

func (o OrderState) String() string {
	return fmt.Sprintf("%s [%s]", o.Hash.Hex(), o.Status.String())
}
