package keeper

import (
	"context"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ethermint/x/orders/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) QueryOrder(c context.Context, req *types.QueryOrderRequest) (*types.QueryOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	order := k.GetOrder(ctx, common.HexToHash(req.OrderHash))
	return &types.QueryOrderResponse{Order: order}, nil
}

func (k Keeper) QueryDerivativeMarket(c context.Context, req *types.QueryDerivativeMarketRequest) (*types.QueryDerivativeMarketResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	market := k.GetDerivativeMarket(ctx, req.MarketId)
	return &types.QueryDerivativeMarketResponse{Market: market}, nil
}

func (k Keeper) QueryDerivativeMarkets(c context.Context, req *types.QueryDerivativeMarketsRequest) (*types.QueryDerivativeMarketsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	markets := k.GetAllDerivativeMarkets(ctx)
	return &types.QueryDerivativeMarketsResponse{Markets: markets}, nil
}

func (k Keeper) QueryTradePairs(c context.Context, req *types.QueryTradePairsRequest) (*types.QueryTradePairsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	pairs := k.GetAllTradePairs(ctx)
	return &types.QueryTradePairsResponse{Records: pairs}, nil
}

func (k Keeper) QueryDerivativeOrders(c context.Context, req *types.QueryDerivativeOrdersRequest) (*types.QueryDerivativeOrdersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	resp := &types.QueryDerivativeOrdersResponse{}

	orderStatus := types.OrderStatusFromString(req.Status)
	tradePairHash := common.HexToHash(req.TradePairHash)

	resp.Records = k.queryOrders(ctx, &orderStatus, &tradePairHash, req.Filters, req.Collection)
	return resp, nil
}

func (k Keeper) queryOrders(ctx sdk.Context, byStatus *types.OrderStatus, byTradePair *common.Hash, byOrderFilters *types.OrderFilters, collection string) []*types.Order {
	var records []*types.Order
	orderFilterPredicate := func(
		byStatus *types.OrderStatus,
		byTradePair *common.Hash,
		byOrderFilters *types.OrderFilters,
	) func(order *types.Order) bool {
		return func(order *types.Order) (stop bool) {
			if byStatus != nil {
				if types.OrderStatus(order.Status) != *byStatus {
					return false
				}
			}

			if byTradePair != nil {
				if order.TradePairHash != byTradePair.String() {
					return false
				}
			}

			if byOrderFilters != nil {
				if !matchToOrderFilters(order.Order, byOrderFilters) {
					return false
				}
			}

			records = append(records, order)
			return false
		}
	}

	predicateFn := orderFilterPredicate(
		byStatus,
		byTradePair,
		byOrderFilters,
	)

	switch types.OrderCollectionType(collection) {
	case types.OrderCollectionActive:
		k.IterateActiveOrders(ctx, predicateFn)
	case types.OrderCollectionArchive:
		k.IterateArchiveOrders(ctx, predicateFn)
	default:
		k.IterateActiveOrders(ctx, predicateFn)
		k.IterateArchiveOrders(ctx, predicateFn)
	}
	return records
}

func (k Keeper) QuerySoftCancelledOrders(c context.Context, req *types.QuerySoftCancelledOrdersRequest) (*types.QuerySoftCancelledOrdersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	resp := &types.QuerySoftCancelledOrdersResponse{}
	for _, hash := range req.OrderHashes {
		k.IterateOrderSoftCancelRequestsByOrderHash(ctx, common.HexToHash(hash),
			func(cancelReq *types.OrderSoftCancelRequest) (stop bool) {
				if cancelReq != nil {
					resp.OrderHashes = append(resp.OrderHashes, hash)
				}

				return false
			})
	}
	return resp, nil
}

func (k Keeper) QueryActiveOrder(c context.Context, req *types.QueryActiveOrderRequest) (*types.QueryActiveOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	order := k.GetActiveOrder(ctx, common.HexToHash(req.OrderHash))
	return &types.QueryActiveOrderResponse{Order: order}, nil
}

func (k Keeper) QueryArchiveOrder(c context.Context, req *types.QueryArchiveOrderRequest) (*types.QueryArchiveOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	order := k.GetArchiveOrder(ctx, common.HexToHash(req.OrderHash))
	return &types.QueryArchiveOrderResponse{Order: order}, nil
}

func (k Keeper) QuerySpotOrders(c context.Context, req *types.QuerySpotOrdersRequest) (*types.QuerySpotOrdersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	resp := &types.QuerySpotOrdersResponse{}

	orderStatus := types.OrderStatusFromString(req.Status)
	tradePairHash := common.HexToHash(req.TradePairHash)

	resp.Records = k.queryOrders(ctx, &orderStatus, &tradePairHash, req.Filters, req.Collection)
	return resp, nil
}

func (k Keeper) QueryFillRequests(c context.Context, req *types.QueryFillRequestsRequest) (*types.QueryFillRequestsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	var results []*types.OrderFillRequest
	k.IterateOrderFillRequestsByOrderHash(ctx, common.HexToHash(req.OrderHash),
		func(fillReq *types.OrderFillRequest) (stop bool) {

			results = append(results, fillReq)

			return false
		},
	)

	return &types.QueryFillRequestsResponse{FillRequests: results}, nil
}

func (k Keeper) QueryTradePair(c context.Context, req *types.QueryTradePairRequest) (*types.QueryTradePairResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	var pair *types.TradePair
	if req.Hash != "" {
		pair = k.GetTradePair(ctx, common.HexToHash(req.Hash))
	} else if len(req.MakerAssetData) > 0 && len(req.TakerAssetData) > 0 {
		hash, err2 := (&types.TradePair{
			MakerAssetData: req.MakerAssetData,
			TakerAssetData: req.TakerAssetData,
		}).ComputeHash()
		if err2 != nil {
			return nil, sdkerrors.Wrap(types.ErrPairNotFound, "failed to compute pair hash: "+err2.Error())
		}
		pair = k.GetTradePair(ctx, hash)
	} else if len(req.Name) > 0 {
		pair = k.GetTradePairByName(ctx, req.Name)
	} else {
		return nil, sdkerrors.Wrap(types.ErrPairNotFound, req.Name)
	}

	return &types.QueryTradePairResponse{
		Name:           pair.Name,
		MakerAssetData: pair.MakerAssetData,
		TakerAssetData: pair.TakerAssetData,
		Hash:           pair.Hash,
		Enabled:        pair.Enabled,
	}, nil
}

func (k Keeper) QueryZeroExTransaction(c context.Context, req *types.QueryZeroExTransactionRequest) (*types.QueryZeroExTransactionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	txInfo := k.GetZeroExTransaction(ctx, common.HexToHash(req.TxHash))

	if txInfo == nil {
		return nil, nil
	}

	txResp := &types.QueryZeroExTransactionResponse{
		TxType: txInfo.ZeroExTransactionType,
	}

	if types.ZeroExTransactionType(txInfo.ZeroExTransactionType) == types.ZeroExOrderFillRequestTx {
		txResp.FillRequests = k.ListOrderFillRequestsByTxHash(ctx, common.HexToHash(req.TxHash))
	} else if types.ZeroExTransactionType(txInfo.ZeroExTransactionType) == types.ZeroExOrderSoftCancelRequestTx {
		txResp.SoftCancelRequests = k.ListOrderSoftCancelRequestsByTxHash(ctx, common.HexToHash(req.TxHash))
	}
	return txResp, nil
}

func (k Keeper) QueryOutstandingFillRequests(c context.Context, req *types.QueryOutstandingFillRequestsRequest) (*types.QueryOutstandingFillRequestsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	fillReqs := k.ListOrderFillRequestsByTxHash(ctx, common.HexToHash(req.TxHash))
	return &types.QueryOutstandingFillRequestsResponse{FillRequests: fillReqs}, nil
}

func (k Keeper) QueryEvmSyncStatus(c context.Context, req *types.QueryEvmSyncStatusRequest) (*types.QueryEvmSyncStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	syncStatus := k.GetEvmSyncStatus(ctx)
	return &types.QueryEvmSyncStatusResponse{SyncStatus: syncStatus}, nil
}

func isZeroAssetData(b string) bool {
	return b == "0x000000000000000000000000000000000000000000000000000000000000000000000000"
}

func matchToOrderFilters(order *types.BaseOrder, filters *types.OrderFilters) bool {

	if filters.NotExpired == true {
		if types.BigNum(order.ExpirationTimeSeconds).Int().Cmp(big.NewInt(time.Now().Unix())) < 1 {
			return false
		}
	}

	// derivative order matching
	if filters.ContractPriceBound != "" && filters.MarketId != "" {
		contractPriceBound := types.BigNum(filters.ContractPriceBound).Int()
		contractPrice := types.BigNum(order.MakerAssetAmount).Int()
		// if my order is long,
		if filters.IsLong {
			// filter out other longs or other orders whose marketID's dont match
			if !isZeroAssetData(order.MakerAssetData) || order.TakerAssetData != filters.MarketId {
				return false
			}
			// filter out short orders whose price is > my long order price
			if contractPrice.Cmp(contractPriceBound) > 0 {
				return false
			}
		} else {
			// my order is short
			// filter out other shorts or other orders whose marketID's dont match
			if !isZeroAssetData(order.TakerAssetData) || (order.MakerAssetData != filters.MarketId) {
				return false
			}
			// filter out long orders whose price is < my short order price
			if contractPrice.Cmp(contractPriceBound) < 0 {
				return false
			}
		}
		// return early for efficiency
		return true
	}

	if filters.MakerAssetData == "" {
		if order.MakerAssetData != filters.MakerAssetData {
			return false
		}
	}

	if filters.TakerAssetData != "" {
		if order.TakerAssetData != filters.TakerAssetData {
			return false
		}
	}

	if filters.MakerAssetAddress != "" {
		if len(order.MakerAssetData) < 36 {
			return false
		}

		makerAssetAddress := (order.MakerAssetData[4 : 4+32])
		if makerAssetAddress != filters.MakerAssetAddress {
			return false
		}
	}

	if filters.TakerAssetAddress != "" {
		if len(order.TakerAssetData) < 36 {
			return false
		}

		takerAssetAddress := (order.TakerAssetData[4 : 4+32])
		if takerAssetAddress != filters.TakerAssetAddress {
			return false
		}
	}

	if filters.ExchangeAddress != "" {
		if order.ExchangeAddress != filters.ExchangeAddress {
			return false
		}
	}

	if filters.SenderAddress != "" {
		if order.SenderAddress != filters.SenderAddress {
			return false
		}
	}

	if filters.MakerFeeAssetData != "" {
		if order.MakerFeeAssetData != filters.MakerFeeAssetData {
			return false
		}
	}

	if filters.TakerFeeAssetData != "" {
		if order.TakerFeeAssetData != filters.TakerFeeAssetData {
			return false
		}
	}

	if filters.MakerAddress != "" {
		if order.MakerAddress != filters.MakerAddress {
			return false
		}
	}

	if filters.NotMakerAddress != "" {
		if order.MakerAddress != filters.NotMakerAddress {
			return false
		}
	}

	if filters.TakerAddress != "" {
		if order.TakerAddress != filters.TakerAddress {
			return false
		}
	}

	if filters.TraderAddress != "" {
		if (order.MakerAddress != filters.TraderAddress) &&
			(order.TakerAddress != filters.TraderAddress) {
			return false
		}
	}

	if filters.FeeRecipientAddress != "" {
		if order.FeeRecipientAddress != filters.FeeRecipientAddress {
			return false
		}
	}

	if filters.MakerAssetAmount != "" && filters.TakerAssetAmount != "" {
		// maker is offering his 1 filters.TakerAssetAmount for 1 filters.MakerAssetAmount
		// maker would accept his 1 filters.TakerAssetAmount for >1 filters.MakerAssetAmount
		// so let's say this order's giving 2 order.TakerAssetAmount for 3 order.MakerAssetAmount
		// maker would accept since (filters.TakerAssetAmount / filters.MakerAssetAmount) >= (order.TakerAssetAmount / order.MakerAssetAmount)

		// maker is offering his 1 filters.TakerAssetAmount for 2 filters.MakerAssetAmount
		// maker would accept his 1 filters.TakerAssetAmount for >2 filters.MakerAssetAmount
		// so let's say this order's giving 2 order.TakerAssetAmount for 3 order.MakerAssetAmount
		// maker would NOT accept since (1 / 2) is not >= (2 / 3)

		// expressed in just multiplication, the condition we need to satisfy is
		// filters.TakerAssetAmount * order.MakerAssetAmount >= order.TakerAssetAmount * filters.MakerAssetAmount
		//a := types.BigNum(takerAmount)
		a := types.BigNum(filters.TakerAssetAmount).Int()
		b := types.BigNum(order.MakerAssetAmount).Int()
		lh := a.Mul(a, b)
		c := types.BigNum(filters.MakerAssetAmount).Int()
		d := types.BigNum(order.TakerAssetAmount).Int()
		rh := c.Mul(c, d)
		if lh.Cmp(rh) == -1 {
			return false
		}
	}

	return true
}

//func (k Keeper) QueryContractSet(c context.Context, req *types.QueryContractSetRequest) (*types.QueryContractSetResponse, error) {
//	if req == nil {
//		return nil, status.Error(codes.InvalidArgument, "empty request")
//	}
//	ctx := sdk.UnwrapSDKContext(c)
//	contractSet := k.QueryContractSet(ctx)
//	return &types.QueryContractSetResponse{ContractSet:}, nil
//}
