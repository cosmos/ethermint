package validator

import (
	"context"
	"math/big"

	"github.com/InjectiveLabs/zeroex-go"
	"github.com/InjectiveLabs/zeroex-go/wrappers"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/cosmos/ethermint/ethereum/provider"
	"github.com/cosmos/ethermint/metrics"
)

// OrderValidator allows to interact with 0x DevUtil smartcontract and perform on-chain orders validation.
// Also computes remaining fill amounts, useful for precise matching.
type OrderValidator interface {
	// OrderState gets relevant order state from DevUtils for a single order.
	OrderState(ctx context.Context, order *zeroex.SignedOrder) (*OrderState, error)
	// OrderStates gets relevant order states from DevUtils for a list of orders.
	OrderStates(ctx context.Context, orders []*zeroex.SignedOrder) ([]*OrderState, error)
	// FuturesPositionState gets relevant futures position state from Futures contract for a single order.
	FuturesPositionState(ctx context.Context, order *zeroex.SignedOrder) (*OrderState, error)
	// FuturesPositionStates gets relevant futures position states from Futures for a list of orders.
	FuturesPositionStates(ctx context.Context, orders []*zeroex.SignedOrder) ([]*OrderState, error)
}

// NewOrderValidator returns a new instance of validator that interacts with devUtilsContract via ethProvider.
func NewOrderValidator(
	ethProvider provider.EVMProvider,
	devUtilsContract, futuresContract common.Address,
) (OrderValidator, error) {
	ethDevUtils, err := wrappers.NewDevUtilsCaller(devUtilsContract, ethProvider)
	if err != nil {
		err := errors.Wrap(err, "failed to init DevUtils caller")
		return nil, err
	}

	ethFutures, err := wrappers.NewFuturesCaller(futuresContract, ethProvider)
	if err != nil {
		err := errors.Wrap(err, "failed to init Futures caller")
		return nil, err
	}

	v := &orderValidator{
		svcTags: metrics.Tags{
			"module": "order_validator",
		},

		ethProvider:      ethProvider,
		ethDevUtils:      ethDevUtils,
		ethFutures:       ethFutures,
		devUtilsContract: devUtilsContract,
		futuresContract:  futuresContract,
	}

	return v, nil
}

type orderValidator struct {
	ethProvider      provider.EVMProvider
	ethDevUtils      *wrappers.DevUtilsCaller
	ethFutures       *wrappers.FuturesCaller
	devUtilsContract common.Address
	futuresContract  common.Address

	svcTags metrics.Tags
}

func (o *orderValidator) OrderState(ctx context.Context, order *zeroex.SignedOrder) (*OrderState, error) {
	metrics.ReportFuncCall(o.svcTags)
	doneFn := metrics.ReportFuncTiming(o.svcTags)
	defer doneFn()

	states, err := o.orderStates(ctx, []*zeroex.SignedOrder{
		order,
	})
	if err != nil {
		return nil, err
	}

	return states[0], nil
}

func (o *orderValidator) OrderStates(ctx context.Context, orders []*zeroex.SignedOrder) ([]*OrderState, error) {
	metrics.ReportFuncCall(o.svcTags)
	doneFn := metrics.ReportFuncTiming(o.svcTags)
	defer doneFn()

	return o.orderStates(ctx, orders)
}

func (o *orderValidator) orderStates(ctx context.Context, orders []*zeroex.SignedOrder) ([]*OrderState, error) {
	wrappedOrders := make([]wrappers.Order, 0, len(orders))
	signatures := make([][]byte, 0, len(orders))
	for _, signedOrder := range orders {
		wrappedOrders = append(wrappedOrders, so2wo(signedOrder))
		signatures = append(signatures, signedOrder.Signature)
	}

	opts := &bind.CallOpts{
		Context: ctx,
		// From is set to some known address only to avoid bug in Ganache
		From: o.devUtilsContract,
	}

	wrappedOrderStates, err := o.ethDevUtils.GetOrderRelevantStates(opts, wrappedOrders, signatures)
	if err != nil {
		err = errors.Wrapf(err, "failed to get order relevant states")
		metrics.ReportFuncError(o.svcTags)
		return nil, err
	}

	orderStates := repackOrderStates(wrappedOrderStates)
	return orderStates, nil
}

func (o *orderValidator) FuturesPositionState(ctx context.Context, order *zeroex.SignedOrder) (*OrderState, error) {
	metrics.ReportFuncCall(o.svcTags)
	doneFn := metrics.ReportFuncTiming(o.svcTags)
	defer doneFn()

	states, err := o.futuresPositionStates(ctx, []*zeroex.SignedOrder{
		order,
	})
	if err != nil {
		return nil, err
	}

	return states[0], nil
}

func (o *orderValidator) FuturesPositionStates(ctx context.Context, orders []*zeroex.SignedOrder) ([]*OrderState, error) {
	metrics.ReportFuncCall(o.svcTags)
	doneFn := metrics.ReportFuncTiming(o.svcTags)
	defer doneFn()

	return o.futuresPositionStates(ctx, orders)
}

func (o *orderValidator) futuresPositionStates(ctx context.Context, orders []*zeroex.SignedOrder) ([]*OrderState, error) {
	wrappedOrders := make([]wrappers.Order, 0, len(orders))
	signatures := make([][]byte, 0, len(orders))
	for _, signedOrder := range orders {
		wrappedOrders = append(wrappedOrders, so2wo(signedOrder))
		signatures = append(signatures, signedOrder.Signature)
	}

	opts := &bind.CallOpts{
		Context: ctx,
		// From is set to some known address only to avoid bug in Ganache
		From: o.devUtilsContract,
	}

	wrappedOrderStates, err := o.ethFutures.GetOrderRelevantStates(opts, wrappedOrders, signatures)
	if err != nil {
		err = errors.Wrapf(err, "failed to get futures position relevant states")
		metrics.ReportFuncError(o.svcTags)
		return nil, err
	}

	orderStates := repackOrderStates(wrappedOrderStates)
	return orderStates, nil
}

type WrappedOrderStates = struct {
	OrdersInfo                []wrappers.OrderInfo
	FillableTakerAssetAmounts []*big.Int
	IsValidSignature          []bool
}

func repackOrderStates(wrapped WrappedOrderStates) []*OrderState {
	states := make([]*OrderState, len(wrapped.OrdersInfo))

	for idx := range states {
		states[idx] = &OrderState{
			Status:                   OrderStatus(wrapped.OrdersInfo[idx].OrderStatus),
			Hash:                     wrapped.OrdersInfo[idx].OrderHash,
			TakerAssetFilledAmount:   wrapped.OrdersInfo[idx].OrderTakerAssetFilledAmount,
			FillableTakerAssetAmount: wrapped.FillableTakerAssetAmounts[idx],
			IsValidSignature:         wrapped.IsValidSignature[idx],
		}
	}

	return states
}

func so2wo(order *zeroex.SignedOrder) wrappers.Order {
	wrappedOrder := wrappers.Order{
		MakerAddress:          order.MakerAddress,
		TakerAddress:          order.TakerAddress,
		FeeRecipientAddress:   order.FeeRecipientAddress,
		SenderAddress:         order.SenderAddress,
		MakerAssetAmount:      order.MakerAssetAmount,
		TakerAssetAmount:      order.TakerAssetAmount,
		MakerFee:              order.MakerFee,
		TakerFee:              order.TakerFee,
		ExpirationTimeSeconds: order.ExpirationTimeSeconds,
		Salt:                  order.Salt,
		MakerAssetData:        order.MakerAssetData,
		TakerAssetData:        order.TakerAssetData,
		MakerFeeAssetData:     order.MakerFeeAssetData,
		TakerFeeAssetData:     order.TakerFeeAssetData,
	}
	return wrappedOrder
}
