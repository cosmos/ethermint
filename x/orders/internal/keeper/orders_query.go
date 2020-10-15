package keeper

import (
	"bytes"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"
	"math/big"
	"time"

	"github.com/cosmos/ethermint/x/orders/internal/types"
)

// queryOrder queries the keeper with an order filtering params and returns marshalled order if found.
func queryOrder(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryOrderParams
	if err := keeper.cdc.UnmarshalBinaryBare(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, "failed to parse params: "+err.Error())
	}

	if order := keeper.GetOrder(ctx, params.Hash); order != nil {
		bz, err := keeper.cdc.MarshalBinaryBare(&types.QueryOrderResponse{
			Order: order,
		})
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, "could not marshal result to JSON: "+err.Error())
		}
		return bz, nil
	}

	err := sdkerrors.Wrap(types.ErrOrderNotFound, "order does not exist: "+params.Hash.String())
	return nil, err
}

// queryActiveOrder queries the keeper with the order filtering params and returns marshalled active order if found.
func queryActiveOrder(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryActiveOrderParams
	if err := keeper.cdc.UnmarshalBinaryBare(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, "failed to parse params:"+err.Error())
	}

	if order := keeper.GetActiveOrder(ctx, params.Hash); order != nil {
		bz, err := keeper.cdc.MarshalBinaryBare(&types.QueryActiveOrderResponse{
			Order: order,
		})
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, "could not marshal result to JSON: "+err.Error())
		}
		return bz, nil
	}

	err := sdkerrors.Wrap(types.ErrOrderNotFound, "active order does not exist: "+params.Hash.String())
	return nil, err
}

// queryArchiveOrder queries the keeper with the order filtering params and returns marshalled archive order if found.
func queryArchiveOrder(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryArchiveOrderParams
	if err := keeper.cdc.UnmarshalBinaryBare(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, "failed to parse params: "+err.Error())
	}

	if order := keeper.GetArchiveOrder(ctx, params.Hash); order != nil {
		bz, err := keeper.cdc.MarshalBinaryBare(&types.QueryArchiveOrderResponse{
			Order: order,
		})
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, "could not marshal result to JSON: "+err.Error())
		}
		return bz, nil
	}

	err := sdkerrors.Wrap(types.ErrOrderNotFound, fmt.Sprintf("archive order %s does not exist", params.Hash.String()))
	return nil, err
}

// queryOrdersList queries the keeper using order filtering params
func queryOrdersList(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryOrdersListParams
	if err := keeper.cdc.UnmarshalBinaryBare(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, fmt.Sprintf("failed to parse params: %s", err))
	}

	resp := &types.QueryOrdersListResponse{}

	orderFilterPredicate := func(
		byStatus *types.OrderStatus,
		byTradePair *common.Hash,
		byOrderFilters *types.OrderFilters,
	) func(order *types.Order) bool {
		return func(order *types.Order) (stop bool) {
			if byStatus != nil {
				if order.Status != *byStatus {
					return false
				}
			}

			if byTradePair != nil {
				if order.TradePairHash != *byTradePair {
					return false
				}
			}

			if byOrderFilters != nil {
				if !matchToOrderFilters(order.Order, byOrderFilters) {
					return false
				}
			}

			resp.Orders = append(resp.Orders, order)
			return false
		}
	}

	predicateFn := orderFilterPredicate(
		params.ByStatus,
		params.ByTradePairHash,
		params.ByOrderFilters,
	)

	switch params.ByCollection {
	case types.OrderCollectionActive:
		keeper.IterateActiveOrders(ctx, predicateFn)
	case types.OrderCollectionArchive:
		keeper.IterateArchiveOrders(ctx, predicateFn)
	default:
		keeper.IterateActiveOrders(ctx, predicateFn)
		keeper.IterateArchiveOrders(ctx, predicateFn)
	}

	if params.Limit > 0 {
		if params.Limit < len(resp.Orders) {
			resp.Orders = resp.Orders[:params.Limit]
		}
	}

	bz, err := keeper.cdc.MarshalBinaryBare(resp)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, "could not marshal result to JSON "+err.Error())
	}

	return bz, nil
}

func isZeroAddr(addr common.Address) bool {
	return addr == common.Address{}
}

func isZeroAssetData(b []byte) bool {
	return bytes.Equal(b, common.FromHex("0x000000000000000000000000000000000000000000000000000000000000000000000000"))
}

func isEqualAddr(addr1, addr2 common.Address) bool {
	return bytes.Equal(addr1.Bytes(), addr2.Bytes())
}

func isEqualBytes(bytes1, bytes2 []byte) bool {
	return bytes.Equal(bytes1, bytes2)
}

func matchToOrderFilters(order *types.SafeSignedOrder, filters *types.OrderFilters) bool {

	if filters.NotExpired == true {
		if order.ExpirationTimeSeconds.Int().Cmp(big.NewInt(time.Now().Unix())) < 1 {
			return false
		}
	}

	// derivative order matching
	if filters.ContractPriceBound != nil && filters.MarketID != nil {
		contractPriceBound := types.BigNum(*filters.ContractPriceBound).Int()
		contractPrice := order.MakerAssetAmount.Int()
		// if my order is long,
		if filters.IsLong {
			// filter out other longs or other orders whose marketID's dont match
			if !isZeroAssetData(order.MakerAssetData) || !isEqualBytes(order.TakerAssetData, *filters.MarketID) {
				return false
			}
			// filter out short orders whose price is > my long order price
			if contractPrice.Cmp(contractPriceBound) > 0 {
				return false
			}
		} else {
			// my order is short
			// filter out other shorts or other orders whose marketID's dont match
			if !isZeroAssetData(order.TakerAssetData) || !isEqualBytes(order.MakerAssetData, *filters.MarketID) {
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

	if filters.MakerAssetData != nil {
		if !isEqualBytes(order.MakerAssetData, *filters.MakerAssetData) {
			return false
		}
	}

	if filters.TakerAssetData != nil {
		if !isEqualBytes(order.TakerAssetData, *filters.TakerAssetData) {
			return false
		}
	}

	if filters.MakerAssetProxyID != nil {
		if len(order.MakerAssetData) < 4 {
			return false
		}

		makerAssetProxyID := order.MakerAssetData[:4]
		if !isEqualBytes(makerAssetProxyID, *filters.MakerAssetProxyID) {
			return false
		}
	}

	if filters.TakerAssetProxyID != nil {
		if len(order.TakerAssetData) < 4 {
			return false
		}

		takerAssetProxyID := order.TakerAssetData[:4]
		if !isEqualBytes(takerAssetProxyID, *filters.TakerAssetProxyID) {
			return false
		}
	}

	if filters.MakerAssetAddress != nil {
		if len(order.MakerAssetData) < 36 {
			return false
		}

		makerAssetAddress := common.BytesToAddress(order.MakerAssetData[4 : 4+32])
		if !isEqualAddr(makerAssetAddress, *filters.MakerAssetAddress) {
			return false
		}
	}

	if filters.TakerAssetAddress != nil {
		if len(order.TakerAssetData) < 36 {
			return false
		}

		takerAssetAddress := common.BytesToAddress(order.TakerAssetData[4 : 4+32])
		if !isEqualAddr(takerAssetAddress, *filters.TakerAssetAddress) {
			return false
		}
	}

	if filters.ExchangeAddress != nil {
		if !isEqualAddr(order.ExchangeAddress.Address, *filters.ExchangeAddress) {
			return false
		}
	}

	if filters.SenderAddress != nil {
		if !isEqualAddr(order.SenderAddress.Address, *filters.SenderAddress) {
			return false
		}
	}

	if filters.MakerFeeAssetData != nil {
		if !isEqualBytes(order.MakerFeeAssetData, *filters.MakerFeeAssetData) {
			return false
		}
	}

	if filters.TakerFeeAssetData != nil {
		if !isEqualBytes(order.TakerFeeAssetData, *filters.TakerFeeAssetData) {
			return false
		}
	}

	if filters.TraderAssetData != nil {
		if !isEqualBytes(order.MakerAssetData, *filters.TraderAssetData) &&
			!isEqualBytes(order.TakerAssetData, *filters.TraderAssetData) {
			return false
		}
	}

	if filters.MakerAddress != nil {
		if !isEqualAddr(order.MakerAddress.Address, *filters.MakerAddress) {
			return false
		}
	}

	if filters.NotMakerAddress != nil {
		if isEqualAddr(order.MakerAddress.Address, *filters.NotMakerAddress) {
			return false
		}
	}

	if filters.TakerAddress != nil {
		if !isEqualAddr(order.TakerAddress.Address, *filters.TakerAddress) {
			return false
		}
	}

	if filters.TraderAddress != nil {
		if !isEqualAddr(order.MakerAddress.Address, *filters.TraderAddress) &&
			!isEqualAddr(order.TakerAddress.Address, *filters.TraderAddress) {
			return false
		}
	}

	if filters.FeeRecipientAddress != nil {
		if !isEqualAddr(order.FeeRecipientAddress.Address, *filters.FeeRecipientAddress) {
			return false
		}
	}

	if filters.MakerAssetAmount != nil && filters.TakerAssetAmount != nil {
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
		a := types.BigNum(*filters.TakerAssetAmount).Int()
		b := order.MakerAssetAmount.Int()
		lh := a.Mul(a, b)
		c := types.BigNum(*filters.MakerAssetAmount).Int()
		d := order.TakerAssetAmount.Int()
		rh := c.Mul(c, d)
		if lh.Cmp(rh) == -1 {
			return false
		}
	}

	return true
}
