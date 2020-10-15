package orders

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/cosmos/ethermint/x/orders/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/InjectiveLabs/zeroex-go"
	"github.com/InjectiveLabs/zeroex-go/wrappers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	log "github.com/xlab/suplog"

	"github.com/cosmos/ethermint/ethereum/provider"
	"github.com/cosmos/ethermint/ethereum/registry"
	"github.com/cosmos/ethermint/eventdb"
	"github.com/cosmos/ethermint/metrics"
)

type OrderMsgHandler interface {
	Handler() sdk.Handler
}

func NewOrderMsgHandler(
	keeper Keeper,
	isExportOnly bool,
	ethOrderEventDB eventdb.OrderEventDB,
	ethFuturesPositionEventDB eventdb.FuturesPositionEventDB,
	ethProvider func() provider.EVMProvider,
	ethContracts registry.ContractDiscoverer,
) sdk.Handler {

	h := &orderMsgHandler{
		svcTags: metrics.Tags{
			"svc": "orders_h",
		},
		keeper:                    keeper,
		ethOrderEventDB:           ethOrderEventDB,
		ethFuturesPositionEventDB: ethFuturesPositionEventDB,
		ethProvider:               ethProvider,
		ethContracts:              ethContracts,
	}

	if !isExportOnly {
		go h.postInit()
	}

	return h.Handler()
}

type orderMsgHandler struct {
	svcTags metrics.Tags

	keeper Keeper

	ethContracts             registry.ContractDiscoverer
	devUtilsContractCaller   *wrappers.DevUtilsCaller
	exchangeContractFilterer *wrappers.ExchangeFilterer
	futuresContractFilterer  *wrappers.FuturesFilterer

	ethOrderEventDB           eventdb.OrderEventDB
	ethFuturesPositionEventDB eventdb.FuturesPositionEventDB

	ethProvider func() provider.EVMProvider
}

func (h *orderMsgHandler) postInit() {
	set := h.ethContracts.GetContracts()

	devUtilsContractCaller, err := wrappers.NewDevUtilsCaller(set.DevUtilsContract, h.ethProvider())
	if err != nil && (set.DevUtilsContract != common.Address{}) {
		err = errors.Wrap(err, "failed to init devutils caller")
		log.Fatalln(err)
	}

	exchangeContractFilterer, err := wrappers.NewExchangeFilterer(set.ExchangeContract, h.ethProvider())
	if err != nil && (set.ExchangeContract != common.Address{}) {
		err = errors.Wrap(err, "failed to init exchange events filterer")
		log.Fatalln(err)
	}

	futuresContractFilterer, err := wrappers.NewFuturesFilterer(set.FuturesContract, h.ethProvider())
	if err != nil && (set.FuturesContract != common.Address{}) {
		err = errors.Wrap(err, "failed to init exchange events filterer")
		log.Fatalln(err)
	}

	h.devUtilsContractCaller = devUtilsContractCaller
	h.exchangeContractFilterer = exchangeContractFilterer
	h.futuresContractFilterer = futuresContractFilterer
}

// Handler returns a handler for "orders" type messages.
func (h *orderMsgHandler) Handler() sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (res *sdk.Result, err error) {
		defer Recover(&err)

		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgRegisterDerivativeMarket:
			return h.handleMsgRegisterDerivativeMarket(ctx, msg)
		case *types.MsgSuspendDerivativeMarket:
			return h.handleMsgSuspendDerivativeMarket(ctx, msg)
		case *types.MsgResumeDerivativeMarket:
			return h.handleMsgResumeDerivativeMarket(ctx, msg)
		case *types.MsgCreateDerivativeOrder:
			return h.handleMsgCreateDerivativeOrder(ctx, msg)
		case *types.MsgFilledDerivativeOrder:
			return h.handleMsgFilledDerivativeOrder(ctx, msg)
		case *types.MsgCancelledDerivativeOrder:
			return h.handleMsgCancelledDerivativeOrder(ctx, msg)
		case *types.MsgRegisterSpotMarket:
			return h.handleMsgRegisterSpotMarket(ctx, msg)
		case *types.MsgSuspendSpotMarket:
			return h.handleMsgSuspendSpotMarket(ctx, msg)
		case *types.MsgResumeSpotMarket:
			return h.handleMsgResumeSpotMarket(ctx, msg)
		case *types.MsgCreateSpotOrder:
			return h.handleMsgCreateSpotOrder(ctx, msg)
		case *types.MsgFilledSpotOrder:
			return h.handleMsgFilledSpotOrder(ctx, msg)
		case *types.MsgCancelledSpotOrder:
			return h.handleMsgCancelledSpotOrder(ctx, msg)
		case *types.MsgRequestFillSpotOrder:
			return h.handleMsgRequestFillSpotOrder(ctx, msg)
		case *types.MsgRequestSoftCancelSpotOrder:
			return h.handleMsgRequestSoftCancelSpotOrder(ctx, msg)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,
				fmt.Sprintf("Unrecognized orders Msg type: %v", msg.Type()))
		}
	}
}

func Recover(err *error) {
	if r := recover(); r != nil {
		*err = sdkerrors.Wrapf(sdkerrors.ErrPanic, "%v", r)

		if e, ok := r.(error); ok {
			log.WithError(e).Errorln("orders msg handler panicked with an error")
		} else {
			log.Errorln(r)
		}
	}
}

// Registers the Derivative Market in the keeper, enabling trades of this Derivative Market
func (h *orderMsgHandler) handleMsgRegisterDerivativeMarket(
	ctx sdk.Context,
	msg *types.MsgRegisterDerivativeMarket,
) (*sdk.Result, error) {
	metrics.ReportFuncCall(h.svcTags)
	doneFn := metrics.ReportFuncTiming(h.svcTags)
	defer doneFn()

	logger := ctx.Logger().With(
		"module", "orders",
		"handler", "MsgRegisterDerivativeMarket",
	)

	market := msg.Market

	hash, err := market.Hash()
	if err != nil {
		logger.Error("market hash failed", "error", err.Error())
		metrics.ReportFuncError(h.svcTags)
		return nil, sdkerrors.Wrap(err, "market hash error")
	}
	if hash != common.HexToHash(market.MarketId) {
		logger.Error("The MarketID provided does not match the MarketID computed", "error", err)
		metrics.ReportFuncError(h.svcTags)
		return nil, sdkerrors.Wrap(types.ErrMarketInvalid, "The MarketID provided does not match the MarketID computed")
	}
	if m := h.keeper.GetDerivativeMarket(ctx, hash.Hex()); m != nil {
		logger.Error("derivative market exists already", "marketID", market.MarketId)
		metrics.ReportFuncError(h.svcTags)
		return nil, sdkerrors.Wrap(types.ErrMarketExists, market.Ticker)
	}

	h.keeper.SetDerivativeMarket(ctx, market)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// Stores a new order in the orderbook, it will be published and can be filled.
// Requires the TradePair of order to exist in the keeper and be enabled.
func (h *orderMsgHandler) handleMsgCreateSpotOrder(ctx sdk.Context, msg *types.MsgCreateSpotOrder) (*sdk.Result, error) {
	metrics.ReportFuncCall(h.svcTags)
	doneFn := metrics.ReportFuncTiming(h.svcTags)
	defer doneFn()

	logger := ctx.Logger().With(
		"module", "orders",
		"handler", "MsgCreateSpotOrder",
	)

	tradePairHash, _ := (&types.TradePair{
		MakerAssetData: msg.Order.MakerAssetData,
		TakerAssetData: msg.Order.TakerAssetData,
	}).ComputeHash()

	tradePair := h.keeper.GetTradePair(ctx, tradePairHash)
	if tradePair == nil {
		logger.Error("trade pair doesn't exist", "hash", tradePairHash.Hex())
		metrics.ReportFuncError(h.svcTags)
		return nil, sdkerrors.Wrap(types.ErrPairNotFound, "trade pair doesn't exist: "+tradePairHash.Hex())
	}

	if !tradePair.Enabled {
		logger.With("name", tradePair.Name).Error("trade pair is suspended", "name", tradePair.Name)
		metrics.ReportFuncError(h.svcTags)
		return nil, sdkerrors.Wrap(types.ErrPairSuspended, "trade pair is suspended: "+tradePair.Name)
	}
	order := &types.Order{
		Order:         msg.Order,
		TradePairHash: tradePairHash.Hex(),
		Status:        int64(types.StatusUnfilled),
	}

	h.keeper.SetOrder(ctx, order)

	signedOrder := msg.Order.ToSignedOrder()
	json, _ := signedOrder.MarshalJSON()
	orderString := string(json)
	hash, _ := signedOrder.ComputeOrderHash()

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeNewOrder,
			sdk.NewAttribute(types.AttributeKeyOrderHash, hash.String()),
			sdk.NewAttribute(types.AttributeKeyTradePairHash, tradePairHash.String()),
			sdk.NewAttribute(types.AttributeKeySignedOrder, orderString),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// Stores a new derivative order in the orderbook, it will be published and can be filled.
// Requires the market of the order to exist in the keeper and be enabled.
func (h *orderMsgHandler) handleMsgCreateDerivativeOrder(ctx sdk.Context, msg *types.MsgCreateDerivativeOrder) (*sdk.Result, error) {
	metrics.ReportFuncCall(h.svcTags)
	doneFn := metrics.ReportFuncTiming(h.svcTags)
	defer doneFn()

	logger := ctx.Logger().With(
		"module", "orders",
		"handler", "MsgCreateDerivativeOrder",
	)

	marketID := msg.Order.TakerAssetData
	// TODO: Do stricter validation in validate basic
	if marketID == "0000000000000000000000000000000000000000000000000000000000000000" {
		marketID = msg.Order.MakerAssetData
	}
	market := h.keeper.GetDerivativeMarket(ctx, marketID)
	if market == nil {
		logger.Error("Derivative market doesn't exist", "id", msg.Order.TakerAssetData)
		metrics.ReportFuncError(h.svcTags)
		return nil, sdkerrors.Wrap(types.ErrMarketNotFound, "Derivative market doesn't exist: "+marketID)
	}

	if !market.Enabled {
		logger.With("ticker", market.Ticker).Error("Derivative market is suspended", "marketID", market.String())
		metrics.ReportFuncError(h.svcTags)
		return nil, sdkerrors.Wrap(types.ErrMarketSuspended, "Derivative market is suspended: "+marketID)
	}
	order := &types.Order{
		Order:         msg.Order,
		TradePairHash: market.MarketId,
		Status:        int64(types.StatusUnfilled),
		FilledAmount:  msg.InitialQuantityMatched,
	}

	h.keeper.SetOrder(ctx, order)

	signedOrder := msg.Order.ToSignedOrder()
	json, _ := signedOrder.MarshalJSON()
	orderString := string(json)
	hash, _ := signedOrder.ComputeOrderHash()

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeNewDerivativeOrder,
			sdk.NewAttribute(types.AttributeKeyOrderHash, hash.String()),
			sdk.NewAttribute(types.AttributeKeyMarketID, market.MarketId),
			sdk.NewAttribute(types.AttributeKeySignedOrder, orderString),
			sdk.NewAttribute(types.AttributeKeyFilledAmount, msg.InitialQuantityMatched),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func (h *orderMsgHandler) handleMsgRequestFillSpotOrder(ctx sdk.Context, msg *types.MsgRequestFillSpotOrder) (*sdk.Result, error) {
	metrics.ReportFuncCall(h.svcTags)
	doneFn := metrics.ReportFuncTiming(h.svcTags)
	defer doneFn()

	logger := ctx.Logger().With(
		"module", "orders",
		"handler", "MsgRequestFillSpotOrder",
	)

	tx := &zeroex.SignedTransaction{
		Transaction: zeroex.Transaction{
			Salt:                  types.BigNum(msg.SignedTransaction.Salt).Int(),
			SignerAddress:         common.HexToAddress(msg.SignedTransaction.SignerAddress),
			Data:                  common.FromHex(msg.SignedTransaction.Data),
			ExpirationTimeSeconds: types.BigNum(msg.SignedTransaction.ExpirationTimeSeconds).Int(),
			GasPrice:              types.BigNum(msg.SignedTransaction.GasPrice).Int(),
		},
		Signature: common.Hex2Bytes(msg.SignedTransaction.Signature),
	}
	tx.Domain.VerifyingContract = common.HexToAddress(msg.SignedTransaction.Domain.VerifyingContract)
	tx.Domain.ChainID = types.BigNum(msg.SignedTransaction.Domain.ChainId).Int()

	txData, err := tx.DecodeTransactionData()
	if err != nil {
		logger.Error("failed to decode tx data", "error", err)
		metrics.ReportFuncError(h.svcTags)
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "failed to decode tx data: "+err.Error())
	}
	if txData.FunctionName == zeroex.BatchMatchOrdersWithMaximalFill {
		// TODO: refactor
		txData.Orders = append(txData.LeftOrders, txData.RightOrders[0])
	}
	if h.hasCancelledOrders(ctx, txData.Orders) {
		err = errors.New("transaction contains soft-cancelled orders")
		metrics.ReportFuncError(h.svcTags)
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "transaction contains soft-cancelled orders: "+err.Error())
	}

	txHash, err := tx.ComputeTransactionHash()
	if err != nil {
		logger.Error("failed to compute zeroex tx hash", "error", err)
		metrics.ReportFuncError(h.svcTags)
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "failed to compute zeroex tx hash: "+err.Error())
	}

	approval := &zeroex.CoordinatorApproval{
		TxOrigin:             common.HexToAddress(msg.TxOrigin),
		TransactionHash:      txHash,
		TransactionSignature: tx.Signature,
		Domain: zeroex.EIP712Domain{
			VerifyingContract: h.ethContracts.GetContracts().CoordinatorContract,
			ChainID:           tx.Domain.ChainID,
		},
	}

	approvalHash, _ := approval.ComputeApprovalHash()

	_, err = h.addressFromSignature(approvalHash.Bytes(), common.FromHex(msg.ApprovalSignature))
	if err != nil {
		err = errors.New("unable to get address from approval sig")
		logger.Error("rejecting fill request", "error", err)
		metrics.ReportFuncError(h.svcTags)
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "rejecting fill request: "+err.Error())
	}
	// TODO: (@albertchon) re-introduce coordinator account logic
	//if !h.isActiveStaker(ctx, coordinatorAddr) {
	//	err = errors.Errorf("coordinator is not found in active stakers")
	//	logger.Error("rejecting fill request", "error", err, "address", coordinatorAddr.Hex())
	//	metrics.ReportFuncError(h.svcTags)
	//	return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "rejecting fill request: "+err.Error())
	//}

	if err = txData.ValidateAssetFillAmounts(); err != nil {
		logger.Error("ValidateAssetFillAmounts rejected orders", "error", err)
		metrics.ReportFuncError(h.svcTags)
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "rejecting fill request: "+err.Error())
	}

	fillRequests := make(map[common.Hash]*types.OrderFillRequest, len(txData.Orders))
	for _, order := range txData.Orders {
		orderHash, _ := order.ComputeOrderHash()
		req, ok := fillRequests[orderHash]
		if !ok {
			req = &types.OrderFillRequest{
				OrderHash: orderHash.String(),
				ApprovalSignatures: []string{
					msg.ApprovalSignature,
				},
				ExpiryAt:             tx.ExpirationTimeSeconds.Int64(),
				TakerAssetFillAmount: order.TakerAssetAmount.String(),
			}
		} else {
			logger.Error("seen order multiple times, a different fee payer?", "orderHash", orderHash.Hex())
			// req.ApprovalSignatures = append(req.ApprovalSignatures, txData.Signatures[orderIdx])
		}
		fillRequests[orderHash] = req
	}

	fillRequestsSorted := make([]*types.OrderFillRequest, 0, len(fillRequests))
	for _, req := range fillRequests {
		fillRequestsSorted = append(fillRequestsSorted, req)
	}

	// sort fill requests after mapping to their corresponding order hash
	sort.Slice(fillRequestsSorted, func(i, j int) bool {
		return strings.Compare(fillRequestsSorted[i].OrderHash, fillRequestsSorted[j].OrderHash) < 0
	})

	orderHashes := make([]string, 0, len(fillRequestsSorted))
	// fillRequestsSorted at this point have grouped-per-order signatures
	for _, fillRequest := range fillRequestsSorted {
		h.keeper.SetOrderFillRequest(ctx, txHash, fillRequest)
		orderHashes = append(orderHashes, fillRequest.OrderHash)
	}
	h.keeper.SetZeroExTransaction(ctx, txHash, &types.ZeroExTransaction{
		ZeroExTransactionType: int64(types.ZeroExOrderFillRequestTx),
		Orders:                orderHashes,
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func (h *orderMsgHandler) addressFromSignature(message, sig []byte) (address common.Address, err error) {
	if len(sig) < 65 {
		err = errors.New("signature is too short")
		return
	}

	digestHash, _ := textAndHash(message)

	ecSignature := make([]byte, 65)
	copy(ecSignature[:32], sig[1:33])    // R
	copy(ecSignature[32:64], sig[33:65]) // S
	ecSignature[64] = sig[0] - 27        // V (0 or 1)

	var pubKey *ecdsa.PublicKey

	if pubKey, err = ethcrypto.SigToPub(digestHash, ecSignature); err != nil {
		log.WithError(err).Errorln("failed to EC recover from sig")
		return common.Address{}, err
	}

	address = ethcrypto.PubkeyToAddress(*pubKey)
	return address, nil
}

// textAndHash is a helper function that calculates a hash for the given message that can be
// safely used to calculate a signature from.
//
// The hash is calulcated as
//   keccak256("\x19Ethereum Signed Message:\n"${message length}${message}).
//
// This gives context to the signed message and prevents signing of transactions.
func textAndHash(data []byte) ([]byte, string) {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), string(data))
	hasher := sha3.NewLegacyKeccak256()
	// Note: Write will never return an error here. We added placeholders in order
	// to satisfy the linter.
	_, _ = hasher.Write([]byte(msg))
	return hasher.Sum(nil), msg
}

func (h *orderMsgHandler) hasCancelledOrders(ctx sdk.Context, orders []*zeroex.Order) bool {
	orderHashes := make([]common.Hash, 0, len(orders))
	for _, order := range orders {
		orderHash, _ := order.ComputeOrderHash()
		orderHashes = append(orderHashes, orderHash)
	}

	cancelledOrders := h.keeper.FindAllSoftCancelledOrders(ctx, orderHashes)

	return len(cancelledOrders) > 0
}

func (h *orderMsgHandler) handleMsgRequestSoftCancelSpotOrder(
	ctx sdk.Context,
	msg *types.MsgRequestSoftCancelSpotOrder,
) (*sdk.Result, error) {
	metrics.ReportFuncCall(h.svcTags)
	doneFn := metrics.ReportFuncTiming(h.svcTags)
	defer doneFn()

	logger := ctx.Logger().With(
		"module", "orders",
		"handler", "MsgRequestSoftCancelSpotOrder",
	)

	tx := &zeroex.SignedTransaction{
		Transaction: zeroex.Transaction{
			Salt:                  types.BigNum(msg.SignedTransaction.Salt).Int(),
			SignerAddress:         common.HexToAddress(msg.SignedTransaction.SignerAddress),
			Data:                  common.FromHex(msg.SignedTransaction.Data),
			ExpirationTimeSeconds: types.BigNum(msg.SignedTransaction.ExpirationTimeSeconds).Int(),
			GasPrice:              types.BigNum(msg.SignedTransaction.GasPrice).Int(),
		},
		Signature: common.FromHex(msg.SignedTransaction.Signature),
	}
	tx.Domain.VerifyingContract = common.HexToAddress(msg.SignedTransaction.Domain.VerifyingContract)
	tx.Domain.ChainID = types.BigNum(msg.SignedTransaction.Domain.ChainId).Int()

	txData, err := tx.DecodeTransactionData()
	if err != nil {
		logger.Error("failed to decode tx data", "error", err)
		metrics.ReportFuncError(h.svcTags)
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "failed to decode tx data: "+err.Error())
	}

	txHash, err := tx.ComputeTransactionHash()
	if err != nil {
		logger.Error("failed to compute zeroex tx hash", "error", err)
		metrics.ReportFuncError(h.svcTags)
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "failed to compute zeroex tx hash: "+err.Error())
	}

	cancelRequests := make(map[common.Hash]*types.OrderSoftCancelRequest, len(txData.Orders))
	for _, order := range txData.Orders {
		orderHash, _ := order.ComputeOrderHash()

		req, ok := cancelRequests[orderHash]
		if !ok {
			req = &types.OrderSoftCancelRequest{
				TxHash:             txHash.String(),
				OrderHash:          orderHash.String(),
				ApprovalSignatures: []string{},
			}

			orderObj := h.keeper.GetOrder(ctx, orderHash)
			if orderObj == nil {
				return nil, sdkerrors.Wrap(types.ErrOrderNotFound, "order cannot be canceled because not found")
			}

			signedOrder := orderObj.Order.ToSignedOrder()
			json, _ := signedOrder.MarshalJSON()
			orderString := string(json)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeSoftCancelOrder,
					sdk.NewAttribute(types.AttributeKeyOrderHash, orderHash.String()),
					sdk.NewAttribute(types.AttributeKeySignedOrder, orderString),
				),
			)
		} else {
			logger.Error("seen order multiple times, different fee payer?", "orderHash", orderHash.Hex())
			// req.ApprovalSignatures = append(req.ApprovalSignatures, signedApproval.Signature)
		}
		cancelRequests[orderHash] = req
	}
	cancelRequestsSorted := make([]*types.OrderSoftCancelRequest, 0, len(cancelRequests))
	for _, req := range cancelRequests {
		cancelRequestsSorted = append(cancelRequestsSorted, req)
	}
	// sort canсel requests after mapping to their corresponding order hash
	sort.Slice(cancelRequestsSorted, func(i, j int) bool {
		return strings.Compare(cancelRequestsSorted[i].OrderHash, cancelRequestsSorted[j].OrderHash) < 0
	})

	orderHashes := make([]string, 0, len(cancelRequestsSorted))
	// cancelRequests at this point have grouped-per-order signatures
	for _, cancelRequest := range cancelRequestsSorted {
		h.keeper.SetOrderSoftCancelRequest(ctx, txHash, cancelRequest)
		h.keeper.SetActiveOrderStatus(ctx, common.HexToHash(cancelRequest.OrderHash), types.StatusSoftCancelled)
		orderHashes = append(orderHashes, cancelRequest.OrderHash)
	}

	h.keeper.SetZeroExTransaction(ctx, txHash, &types.ZeroExTransaction{
		ZeroExTransactionType: int64(types.ZeroExOrderSoftCancelRequestTx),
		Orders:                orderHashes,
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

const defaultOnlineLookupTimeout = 250 * time.Millisecond

func (h *orderMsgHandler) handleMsgFilledSpotOrder(
	ctx sdk.Context,
	msg *types.MsgFilledSpotOrder,
) (*sdk.Result, error) {
	metrics.ReportFuncCall(h.svcTags)
	doneFn := metrics.ReportFuncTiming(h.svcTags)
	defer doneFn()

	logger := ctx.Logger().With(
		"module", "orders",
		"handler", "MsgFilledSpotOrder",
		"blockNum", msg.BlockNum,
		"txHash", msg.TxHash,
		"orderHash", msg.OrderHash,
	)

	order := h.keeper.GetActiveOrder(ctx, common.HexToHash(msg.OrderHash))
	if order == nil {
		logger.Info("no active order found for hash")
		// no active order, the event is irrelevant
		return &sdk.Result{}, nil
	}

	msgEvent := &eventdb.OrderEvent{
		Type:       eventdb.OrderUpdateFilled,
		BlockNum:   uint64(msg.BlockNum),
		TxHash:     common.HexToHash(msg.TxHash),
		OrderHash:  common.HexToHash(msg.OrderHash),
		FillAmount: types.BigNum(msg.AmountFilled).Int(),
	}

	// validate against local cache of events
	ev, ok := h.ethOrderEventDB.GetFillEvent(uint64(msg.BlockNum), common.HexToHash(msg.TxHash), common.HexToHash(msg.OrderHash))
	if !ok {
		// if not found, lookup online with a timeout
		logger.Info("order event not found in DB")
		onlineCtx, cancelFn := context.WithTimeout(context.Background(), defaultOnlineLookupTimeout)
		ev, ok = h.getOrderFillEventFromNode(onlineCtx, uint64(msg.BlockNum), common.HexToHash(msg.TxHash), common.HexToHash(msg.OrderHash))
		cancelFn()

		if !ok {
			// if not found online, bail out with error
			logger.Error("order event not found in the node logs")
			metrics.ReportFuncError(h.svcTags)

			return nil, types.ErrBadUpdateEvent
		} else if ev == nil {
			// a special case when online lookup failed due to timeout or other technical reasons
			// just trust the original message
			logger.Error("couldn't validate the order event, trusting the source")

			ev = msgEvent
		}
	}

	if !ev.Equals(msgEvent) {
		logger.Error("couldn't validate the order event, locally stored data is different")
		metrics.ReportFuncError(h.svcTags)

		return nil, types.ErrBadUpdateEvent
	}

	prevAmount := types.BigNum(order.FilledAmount).Int()
	newAmount := prevAmount.Add(prevAmount, ev.FillAmount)

	logger.With(
		"prevAmount", prevAmount.String(),
		"newAmount", newAmount.String(),
	).Info("updating FilledTakerAssetAmount")

	order.FilledAmount = newAmount.String()
	if newAmount.Cmp(types.BigNum(order.Order.TakerAssetAmount).Int()) != -1 { // >=
		order.Status = int64(types.StatusFilled)
		h.keeper.SetActiveOrderStatus(ctx, common.HexToHash(msg.OrderHash), types.StatusFilled)

		logger.Info("order is fully filled")
	} else {
		order.Status = int64(types.StatusPartialFilled)

		logger.Info("order is partially filled")
	}

	h.keeper.SetOrder(ctx, order)

	syncStatus := h.keeper.GetEvmSyncStatus(ctx)
	if syncStatus != nil {
		if syncStatus.LatestBlockSynced < msg.BlockNum {
			logger.With(
				"prevLatestBlockSynced", syncStatus.LatestBlockSynced,
				"newLatestBlockSynced", msg.BlockNum,
			).Info("updating EvmSyncStatus")

			syncStatus.LatestBlockSynced = msg.BlockNum
			h.keeper.SetEvmSyncStatus(ctx, syncStatus)
		}
	} else {
		logger.With("latestBlockSynced", msg.BlockNum).Info("saving new EvmSyncStatus")

		h.keeper.SetEvmSyncStatus(ctx, &types.EvmSyncStatus{
			LatestBlockSynced: msg.BlockNum,
		})
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func (h *orderMsgHandler) handleMsgCancelledSpotOrder(ctx sdk.Context, msg *types.MsgCancelledSpotOrder) (*sdk.Result, error) {
	metrics.ReportFuncCall(h.svcTags)
	doneFn := metrics.ReportFuncTiming(h.svcTags)
	defer doneFn()

	logger := ctx.Logger().With(
		"module", "orders",
		"handler", "MsgCancelledSpotOrder",
		"blockNum", msg.BlockNum,
		"txHash", msg.TxHash,
		"orderHash", msg.OrderHash,
	)

	order := h.keeper.GetActiveOrder(ctx, common.HexToHash(msg.OrderHash))
	if order == nil {
		// no active order, the event is irrelevant
		logger.Info("no active order found for hash")

		order = h.keeper.GetArchiveOrder(ctx, common.HexToHash(msg.OrderHash))
		if order == nil {
			logger.Info("no archived order found for hash")
			return &sdk.Result{}, nil
		} else if types.OrderStatus(order.Status) != types.StatusSoftCancelled {
			logger.With("orderStatus", order.Status).Info("refusing to hard-cancel order that is not soft-cancelled")
			return &sdk.Result{}, nil
		}
	}

	msgEvent := &eventdb.OrderEvent{
		Type:      eventdb.OrderUpdateHardCancelled,
		BlockNum:  uint64(msg.BlockNum),
		TxHash:    common.HexToHash(msg.TxHash),
		OrderHash: common.HexToHash(msg.OrderHash),
	}

	// validate against local cache of events
	ev, ok := h.ethOrderEventDB.GetCancelEvent(uint64(msg.BlockNum), common.HexToHash(msg.TxHash), common.HexToHash(msg.OrderHash))
	if !ok {
		logger.Info("order event not found in DB")

		// if not found, lookup online with a timeout
		onlineCtx, cancelFn := context.WithTimeout(context.Background(), defaultOnlineLookupTimeout)
		ev, ok = h.getOrderCancelEventFromNode(onlineCtx, uint64(msg.BlockNum), common.HexToHash(msg.TxHash), common.HexToHash(msg.OrderHash))
		cancelFn()

		if !ok {
			// if not found online, bail out with error
			logger.Error("order event not found in the node logs")
			metrics.ReportFuncError(h.svcTags)

			return nil, types.ErrBadUpdateEvent
		} else if ev == nil {
			// a special case when online lookup failed due to timeout or other technical reasons
			// just trust the original message
			logger.Error("couldn't validate the order event, trusting the source")

			ev = msgEvent
		}
	}

	if !ev.Equals(msgEvent) {
		logger.Error("couldn't validate the order event, locally stored data is different")
		metrics.ReportFuncError(h.svcTags)

		return nil, types.ErrBadUpdateEvent
	}

	if types.OrderStatus(order.Status) == types.StatusSoftCancelled {
		order.Status = int64(types.StatusHardCancelled)
		h.keeper.SetArchiveOrderStatus(ctx, common.HexToHash(msg.OrderHash), types.StatusHardCancelled)
	} else {
		order.Status = int64(types.StatusHardCancelled)
		h.keeper.SetActiveOrderStatus(ctx, common.HexToHash(msg.OrderHash), types.StatusHardCancelled)
	}

	syncStatus := h.keeper.GetEvmSyncStatus(ctx)
	if syncStatus != nil {
		if syncStatus.LatestBlockSynced < msg.BlockNum {
			logger.With(
				"prevLatestBlockSynced", syncStatus.LatestBlockSynced,
				"newLatestBlockSynced", msg.BlockNum,
			).Info("updating EvmSyncStatus")

			syncStatus.LatestBlockSynced = msg.BlockNum
			h.keeper.SetEvmSyncStatus(ctx, syncStatus)
		}
	} else {
		logger.With("latestBlockSynced", msg.BlockNum).Info("saving new EvmSyncStatus")

		h.keeper.SetEvmSyncStatus(ctx, &types.EvmSyncStatus{
			LatestBlockSynced: msg.BlockNum,
		})
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func (h *orderMsgHandler) getOrderFillEventFromNode(ctx context.Context, blockNum uint64, txHash, orderHash common.Hash) (*eventdb.OrderEvent, bool) {
	metrics.ReportFuncCall(h.svcTags)
	doneFn := metrics.ReportFuncTiming(h.svcTags)
	defer doneFn()

	receipt, err := h.ethProvider().TransactionReceiptByHash(ctx, txHash)
	if err != nil {
		// TODO: detect "not found" error and return false
		// because faking txHash is not allowed

		log.WithError(err).Errorln("failed to get transaction receipt from node")
		metrics.ReportFuncError(h.svcTags)
		return nil, true // technical error
	}

	if uint64(receipt.BlockNumber) != blockNum {
		err = fmt.Errorf("block num mismatch: %d != %d", uint64(receipt.BlockNumber), blockNum)
		log.WithError(err).Errorln("failed verify transaction log")
		metrics.ReportFuncError(h.svcTags)
		return nil, false
	}

	var eventFound bool
	var amountFilled *big.Int
	for _, log := range receipt.Logs {
		fillEvent, err := h.exchangeContractFilterer.ParseFill(*log)
		if err != nil {
			continue
		}
		if common.Hash(fillEvent.OrderHash) != orderHash {
			continue
		}

		eventFound = true
		amountFilled = fillEvent.TakerAssetFilledAmount
	}

	if !eventFound {
		metrics.ReportFuncError(h.svcTags)
		return nil, false
	}

	orderEvent := &eventdb.OrderEvent{
		Type:       eventdb.OrderUpdateFilled,
		BlockNum:   blockNum,
		TxHash:     txHash,
		OrderHash:  orderHash,
		FillAmount: amountFilled,
	}

	return orderEvent, true
}

func (h *orderMsgHandler) getOrderCancelEventFromNode(ctx context.Context, blockNum uint64, txHash, orderHash common.Hash) (*eventdb.OrderEvent, bool) {
	metrics.ReportFuncCall(h.svcTags)
	doneFn := metrics.ReportFuncTiming(h.svcTags)
	defer doneFn()

	receipt, err := h.ethProvider().TransactionReceiptByHash(ctx, txHash)
	if err != nil {
		// TODO: detect "not found" error and return false
		// because faking txHash is not allowed

		log.WithError(err).Errorln("failed to get transaction receipt from node")
		metrics.ReportFuncError(h.svcTags)
		return nil, true // technical error
	}

	if uint64(receipt.BlockNumber) != blockNum {
		err = fmt.Errorf("block num mismatch: %d != %d", uint64(receipt.BlockNumber), blockNum)
		log.WithError(err).Errorln("failed verify transaction log")
		metrics.ReportFuncError(h.svcTags)
		return nil, false
	}

	var eventFound bool
	for _, log := range receipt.Logs {
		cancelEvent, err := h.exchangeContractFilterer.ParseCancel(*log)
		if err != nil {
			continue
		}
		if common.Hash(cancelEvent.OrderHash) != orderHash {
			continue
		}

		eventFound = true
	}

	if !eventFound {
		metrics.ReportFuncError(h.svcTags)
		return nil, false
	}

	orderEvent := &eventdb.OrderEvent{
		Type:      eventdb.OrderUpdateHardCancelled,
		BlockNum:  blockNum,
		TxHash:    txHash,
		OrderHash: orderHash,
	}

	return orderEvent, true
}

func (h *orderMsgHandler) handleMsgFilledDerivativeOrder(
	ctx sdk.Context,
	msg *types.MsgFilledDerivativeOrder,
) (*sdk.Result, error) {
	metrics.ReportFuncCall(h.svcTags)
	doneFn := metrics.ReportFuncTiming(h.svcTags)
	defer doneFn()

	logger := ctx.Logger().With(
		"module", "orders",
		"handler", "MsgFilledDerivativeOrder",
		"blockNum", msg.BlockNum,
		"txHash", msg.TxHash,
		"orderHash", msg.OrderHash,
	)

	order := h.keeper.GetActiveOrder(ctx, common.HexToHash(msg.OrderHash))
	if order == nil {
		// no active order, the event is irrelevant
		return &sdk.Result{}, nil
	}

	msgEvent := &eventdb.FuturesPositionEvent{
		Type:           eventdb.FuturesPositionUpdateFilled,
		BlockNum:       uint64(msg.BlockNum),
		TxHash:         common.HexToHash(msg.TxHash),
		MakerAddress:   common.HexToAddress(msg.MakerAddress),
		OrderHash:      common.HexToHash(msg.OrderHash),
		MarketID:       common.HexToHash(msg.MarketId),
		QuantityFilled: types.BigNum(msg.QuantityFilled).Int(),
		ContractPrice:  types.BigNum(msg.ContractPrice).Int(),
		PositionID:     types.BigNum(msg.PositionId).Int(),
		IsLong:         msg.IsLong,
	}

	// validate against local cache of events
	ev, ok := h.ethFuturesPositionEventDB.GetFillEvent(uint64(msg.BlockNum), common.HexToHash(msg.TxHash), common.HexToHash(msg.OrderHash), msg.IsLong)
	if !ok {
		// if not found, lookup online with a timeout
		onlineCtx, cancelFn := context.WithTimeout(context.Background(), defaultOnlineLookupTimeout)
		ev, ok = h.getFuturesPositionFillEventFromNode(onlineCtx, uint64(msg.BlockNum), common.HexToHash(msg.TxHash), common.HexToHash(msg.OrderHash), msg.IsLong)
		cancelFn()

		if !ok {
			// if not found online, bail out with error
			logger.Error("furtures position event not found in the node logs")
			metrics.ReportFuncError(h.svcTags)

			return &sdk.Result{}, nil
		} else if ev == nil {
			// a special case when online lookup failed due to timeout or other technical reasons
			// just trust the original message
			logger.Error("couldn't validate the furtures position event, trusting the source")

			ev = msgEvent
		}
	}

	if !ev.Equals(msgEvent) {
		logger.Error("couldn't validate the furtures position event, locally stored data is different")
		metrics.ReportFuncError(h.svcTags)
		return nil, types.ErrBadUpdateEvent
	}

	prevAmount := types.BigNum(order.FilledAmount).Int()
	// if pre-filled from a match, just update the state
	if prevAmount.Cmp(ev.QuantityFilled) == 0 && types.OrderStatus(order.Status) == types.StatusUnfilled {
		order.Status = int64(types.StatusPartialFilled)
	} else {
		newAmount := prevAmount.Add(prevAmount, ev.QuantityFilled)

		logger.With(
			"prevAmount", prevAmount.String(),
			"newAmount", newAmount.String(),
		).Info("updating FilledTakerAssetAmount")

		order.FilledAmount = newAmount.String()
		if newAmount.Cmp(types.BigNum(order.Order.TakerAssetAmount).Int()) != -1 { // >=
			order.Status = int64(types.StatusFilled)
			h.keeper.SetActiveOrderStatus(ctx, common.HexToHash(msg.OrderHash), types.StatusFilled)

			logger.Info("position is fully filled")
		} else {
			order.Status = int64(types.StatusPartialFilled)
			logger.Info("position is partially filled")
		}
	}

	h.keeper.SetOrder(ctx, order)

	syncStatus := h.keeper.GetFuturesEvmSyncStatus(ctx)
	if syncStatus != nil {
		if syncStatus.LatestBlockSynced < msg.BlockNum {

			logger.With(
				"prevLatestBlockSynced", syncStatus.LatestBlockSynced,
				"newLatestBlockSynced", msg.BlockNum,
			).Info("updating FuturesEvmSyncStatus")

			syncStatus.LatestBlockSynced = msg.BlockNum
			h.keeper.SetFuturesEvmSyncStatus(ctx, syncStatus)
		}
	} else {
		logger.With("latestBlockSynced", msg.BlockNum).Info("saving new FuturesEvmSyncStatus")
		h.keeper.SetFuturesEvmSyncStatus(ctx, &types.EvmSyncStatus{
			LatestBlockSynced: msg.BlockNum,
		})
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func (h *orderMsgHandler) handleMsgCancelledDerivativeOrder(ctx sdk.Context, msg *types.MsgCancelledDerivativeOrder) (*sdk.Result, error) {
	metrics.ReportFuncCall(h.svcTags)
	doneFn := metrics.ReportFuncTiming(h.svcTags)
	defer doneFn()

	logger := ctx.Logger().With(
		"module", "orders",
		"handler", "MsgCancelledDerivativeOrder",
		"blockNum", msg.BlockNum,
		"txHash", msg.TxHash,
		"orderHash", msg.OrderHash,
	)

	order := h.keeper.GetActiveOrder(ctx, common.HexToHash(msg.OrderHash))
	if order == nil {
		// no active order, the event is irrelevant

		order = h.keeper.GetArchiveOrder(ctx, common.HexToHash(msg.OrderHash))
		if order == nil {
			logger.Info("no archived order found for hash")
			return &sdk.Result{}, nil
		} else if types.OrderStatus(order.Status) != types.StatusSoftCancelled {
			logger.With("orderStatus", order.Status).Info("refusing to hard-cancel order that is not soft-cancelled")
			return &sdk.Result{}, nil
		}
	}

	msgEvent := &eventdb.FuturesPositionEvent{
		Type:         eventdb.FuturesPositionUpdateHardCancelled,
		BlockNum:     uint64(msg.BlockNum),
		TxHash:       common.HexToHash(msg.TxHash),
		MakerAddress: common.HexToAddress(msg.MakerAddress),
		MarketID:     common.HexToHash(msg.MarketId),
		OrderHash:    common.HexToHash(msg.OrderHash),
		PositionID:   types.BigNum(msg.PositionId).Int(),
	}

	// validate against local cache of events
	ev, ok := h.ethFuturesPositionEventDB.GetCancelEvent(uint64(msg.BlockNum), common.HexToHash(msg.TxHash), common.HexToHash(msg.OrderHash))
	if !ok {
		// if not found, lookup online with a timeout
		onlineCtx, cancelFn := context.WithTimeout(context.Background(), defaultOnlineLookupTimeout)
		ev, ok = h.getFuturesPositionCancelEventFromNode(onlineCtx, uint64(msg.BlockNum), common.HexToHash(msg.TxHash), common.HexToHash(msg.OrderHash))
		cancelFn()

		if !ok {
			// if not found online, bail out with error
			logger.Error("futures position event not found in the node logs")
			metrics.ReportFuncError(h.svcTags)

			return nil, types.ErrBadUpdateEvent
		} else if ev == nil {
			// a special case when online lookup failed due to timeout or other technical reasons
			// just trust the original message
			logger.Error("couldn't validate the futures position event, trusting the source")

			ev = msgEvent
		}
	}

	if !ev.Equals(msgEvent) {
		logger.Error("couldn't validate the futures position event, locally stored data is different")

		metrics.ReportFuncError(h.svcTags)
		return nil, types.ErrBadUpdateEvent
	}

	if types.OrderStatus(order.Status) == types.StatusSoftCancelled {
		order.Status = int64(types.StatusHardCancelled)
		h.keeper.SetArchiveOrderStatus(ctx, common.HexToHash(msg.OrderHash), types.StatusHardCancelled)
	} else {
		order.Status = int64(types.StatusHardCancelled)
		h.keeper.SetActiveOrderStatus(ctx, common.HexToHash(msg.OrderHash), types.StatusHardCancelled)
	}

	syncStatus := h.keeper.GetFuturesEvmSyncStatus(ctx)
	if syncStatus != nil {
		if syncStatus.LatestBlockSynced < msg.BlockNum {
			logger.With(
				"prevLatestBlockSynced", syncStatus.LatestBlockSynced,
				"newLatestBlockSynced", msg.BlockNum,
			).Info("updating FuturesEvmSyncStatus")

			syncStatus.LatestBlockSynced = msg.BlockNum
			h.keeper.SetFuturesEvmSyncStatus(ctx, syncStatus)
		}
	} else {
		logger.With("latestBlockSynced", msg.BlockNum).Info("saving new FuturesEvmSyncStatus")

		h.keeper.SetFuturesEvmSyncStatus(ctx, &types.EvmSyncStatus{
			LatestBlockSynced: msg.BlockNum,
		})
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func (h *orderMsgHandler) getFuturesPositionFillEventFromNode(
	ctx context.Context,
	blockNum uint64,
	txHash, orderHash common.Hash, isLong bool,
) (*eventdb.FuturesPositionEvent, bool) {
	metrics.ReportFuncCall(h.svcTags)
	doneFn := metrics.ReportFuncTiming(h.svcTags)
	defer doneFn()

	receipt, err := h.ethProvider().TransactionReceiptByHash(ctx, txHash)
	if err != nil {
		// TODO: detect "not found" error and return false
		// because faking txHash is not allowed

		log.WithError(err).Errorln("failed to get transaction receipt from node")
		metrics.ReportFuncError(h.svcTags)
		return nil, true // technical error
	}

	if uint64(receipt.BlockNumber) != blockNum {
		err = fmt.Errorf("block num mismatch: %d != %d", uint64(receipt.BlockNumber), blockNum)
		log.WithError(err).Errorln("failed verify transaction log")
		metrics.ReportFuncError(h.svcTags)
		return nil, false
	}

	var eventFound bool
	var quantityFilled *big.Int
	for _, log := range receipt.Logs {
		fillEvent, err := h.futuresContractFilterer.ParseFuturesPosition(*log)
		if err != nil {
			continue
		}
		if common.Hash(fillEvent.OrderHash) != orderHash {
			continue
		} else if fillEvent.IsLong != isLong {
			continue
		}

		eventFound = true
		quantityFilled = fillEvent.QuantityFilled
	}

	if !eventFound {
		metrics.ReportFuncError(h.svcTags)
		return nil, false
	}

	positionEvent := &eventdb.FuturesPositionEvent{
		Type:           eventdb.FuturesPositionUpdateFilled,
		BlockNum:       blockNum,
		TxHash:         txHash,
		OrderHash:      orderHash,
		QuantityFilled: quantityFilled,
	}

	return positionEvent, true
}

func (h *orderMsgHandler) getFuturesPositionCancelEventFromNode(
	ctx context.Context,
	blockNum uint64,
	txHash, orderHash common.Hash,
) (*eventdb.FuturesPositionEvent, bool) {
	metrics.ReportFuncCall(h.svcTags)
	doneFn := metrics.ReportFuncTiming(h.svcTags)
	defer doneFn()

	receipt, err := h.ethProvider().TransactionReceiptByHash(ctx, txHash)
	if err != nil {
		// TODO: detect "not found" error and return false
		// because faking txHash is not allowed

		log.WithError(err).Errorln("failed to get transaction receipt from node")
		metrics.ReportFuncError(h.svcTags)
		return nil, true // technical error
	}

	if uint64(receipt.BlockNumber) != blockNum {
		err = fmt.Errorf("block num mismatch: %d != %d", uint64(receipt.BlockNumber), blockNum)
		log.WithError(err).Errorln("failed verify transaction log")
		metrics.ReportFuncError(h.svcTags)
		return nil, false
	}

	var eventFound bool
	for _, log := range receipt.Logs {
		cancelEvent, err := h.futuresContractFilterer.ParseFuturesCancel(*log)
		if err != nil {
			continue
		}
		if common.Hash(cancelEvent.OrderHash) != orderHash {
			continue
		}

		eventFound = true
	}

	if !eventFound {
		metrics.ReportFuncError(h.svcTags)
		return nil, false
	}

	positionEvent := &eventdb.FuturesPositionEvent{
		Type:      eventdb.FuturesPositionUpdateHardCancelled,
		BlockNum:  blockNum,
		TxHash:    txHash,
		OrderHash: orderHash,
	}

	return positionEvent, true
}

// Registers the TradePair in the keeper, enabling trades of this TradePair
func (h *orderMsgHandler) handleMsgRegisterSpotMarket(
	ctx sdk.Context,
	msg *types.MsgRegisterSpotMarket,
) (*sdk.Result, error) {
	metrics.ReportFuncCall(h.svcTags)
	doneFn := metrics.ReportFuncTiming(h.svcTags)
	defer doneFn()

	logger := ctx.Logger().With(
		"module", "orders",
		"handler", "MsgRegisterSpotMarket",
	)

	tradePair := &types.TradePair{
		Name:           msg.Name,
		MakerAssetData: msg.MakerAssetData,
		TakerAssetData: msg.TakerAssetData,
		Enabled:        msg.Enabled,
	}
	hash, err := tradePair.ComputeHash()
	if err != nil {
		logger.Error("trade pair hash failed", "error", err.Error())
		metrics.ReportFuncError(h.svcTags)
		return nil, sdkerrors.Wrap(err, "trade pair hash error")
	}
	if pair := h.keeper.GetTradePair(ctx, hash); pair != nil {
		logger.Error("trade pair exists already", "name", msg.Name)
		metrics.ReportFuncError(h.svcTags)
		return nil, sdkerrors.Wrap(types.ErrPairExists, msg.Name)
	}

	h.keeper.SetTradePair(ctx, tradePair)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// Disables trading of the TradePair
func (h *orderMsgHandler) handleMsgSuspendSpotMarket(
	ctx sdk.Context,
	msg *types.MsgSuspendSpotMarket,
) (*sdk.Result, error) {
	metrics.ReportFuncCall(h.svcTags)
	doneFn := metrics.ReportFuncTiming(h.svcTags)
	defer doneFn()

	logger := ctx.Logger().With(
		"module", "orders",
		"handler", "MsgSuspendSpotMarket",
	)

	var tradePair *types.TradePair
	if len(msg.MakerAssetData) > 0 && len(msg.TakerAssetData) > 0 {
		hash, _ := (&types.TradePair{
			MakerAssetData: msg.MakerAssetData,
			TakerAssetData: msg.TakerAssetData,
		}).ComputeHash()
		tradePair = h.keeper.GetTradePair(ctx, hash)
	} else if len(msg.Name) > 0 {
		tradePair = h.keeper.GetTradePairByName(ctx, msg.Name)
	}
	if tradePair == nil {
		logger.Error("trade pair doesn't exist", "name", msg.Name)
		metrics.ReportFuncError(h.svcTags)
		return nil, sdkerrors.Wrap(types.ErrPairNotFound, msg.Name)
	}

	hash, _ := tradePair.ComputeHash()

	h.keeper.SetTradePairEnabled(ctx, hash, false)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// Enables trading of the TradePair
func (h *orderMsgHandler) handleMsgResumeSpotMarket(
	ctx sdk.Context,
	msg *types.MsgResumeSpotMarket,
) (*sdk.Result, error) {
	metrics.ReportFuncCall(h.svcTags)
	doneFn := metrics.ReportFuncTiming(h.svcTags)
	defer doneFn()

	logger := ctx.Logger().With(
		"module", "orders",
		"handler", "MsgResumeSpotMarket",
		"pairName", msg.Name,
	)

	var tradePair *types.TradePair
	if len(msg.MakerAssetData) > 0 && len(msg.TakerAssetData) > 0 {
		hash, _ := (&types.TradePair{
			MakerAssetData: msg.MakerAssetData,
			TakerAssetData: msg.TakerAssetData,
		}).ComputeHash()
		tradePair = h.keeper.GetTradePair(ctx, hash)
	} else if len(msg.Name) > 0 {
		tradePair = h.keeper.GetTradePairByName(ctx, msg.Name)
	}
	if tradePair == nil {
		logger.Error("trade pair doesn't exist")
		metrics.ReportFuncError(h.svcTags)
		return nil, sdkerrors.Wrap(types.ErrPairNotFound, msg.Name)
	}

	hash, _ := tradePair.ComputeHash()

	h.keeper.SetTradePairEnabled(ctx, hash, true)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// Disables trading of the derivative market
func (h *orderMsgHandler) handleMsgSuspendDerivativeMarket(
	ctx sdk.Context,
	msg *types.MsgSuspendDerivativeMarket,
) (*sdk.Result, error) {
	metrics.ReportFuncCall(h.svcTags)
	doneFn := metrics.ReportFuncTiming(h.svcTags)
	defer doneFn()

	logger := ctx.Logger().With(
		"module", "orders",
		"handler", "MsgSuspendDerivativeMarket",
	)

	oldDerivativeMarket := h.keeper.GetDerivativeMarket(ctx, msg.MarketId)
	if oldDerivativeMarket == nil {
		logger.Error("derivative market doesn't exist")
		metrics.ReportFuncError(h.svcTags)
		return nil, sdkerrors.Wrap(types.ErrMarketNotFound, msg.MarketId)
	}

	h.keeper.SetMarketEnabled(ctx, common.HexToHash(msg.MarketId), false)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

// Re-enables trading of the derivative market
func (h *orderMsgHandler) handleMsgResumeDerivativeMarket(
	ctx sdk.Context,
	msg *types.MsgResumeDerivativeMarket,
) (*sdk.Result, error) {
	metrics.ReportFuncCall(h.svcTags)
	doneFn := metrics.ReportFuncTiming(h.svcTags)
	defer doneFn()

	logger := ctx.Logger().With(
		"module", "orders",
		"handler", "MsgResumeDerivativeMarket",
	)

	oldDerivativeMarket := h.keeper.GetDerivativeMarket(ctx, msg.MarketId)
	if oldDerivativeMarket == nil {
		logger.Error("derivative market doesn't exist")
		metrics.ReportFuncError(h.svcTags)
		return nil, sdkerrors.Wrap(types.ErrMarketNotFound, msg.MarketId)
	} else if !oldDerivativeMarket.Enabled {
		h.keeper.SetMarketEnabled(ctx, common.HexToHash(msg.MarketId), true)
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
