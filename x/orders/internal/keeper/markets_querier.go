package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/ethermint/x/orders/internal/types"
)

//// queryMarket queries keeper with derivative market params and returns found market marshaled to bytes.
//func queryMarket(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
//	var params types.QueryPairParams
//
//	if err := keeper.cdc.UnmarshalBinaryBare(req.Data, &params); err != nil {
//		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, fmt.Sprintf("failed to parse params: %s", err))
//	}
//
//	var tradePair *types.TradePair
//
//	if (params.Hash != common.Hash{}) {
//		tradePair = keeper.GetTradePair(ctx, params.Hash)
//	} else if len(params.MakerAssetData) > 0 && len(params.TakerAssetData) > 0 {
//		hash, err := (&types.TradePair{
//			MakerAssetData: params.MakerAssetData,
//			TakerAssetData: params.TakerAssetData,
//		}).Hash()
//		if err != nil {
//			return nil, sdkerrors.Wrap(types.ErrPairNotFound, "failed to compute pair hash: "+err.Error())
//		}
//		tradePair = keeper.GetTradePair(ctx, hash)
//	} else if len(params.Name) > 0 {
//		tradePair = keeper.GetTradePairByName(ctx, params.Name)
//	}
//
//	if tradePair == nil {
//		return nil, sdkerrors.Wrap(types.ErrPairNotFound, params.Name)
//	}
//
//	bz, err := keeper.cdc.MarshalBinaryBare(&types.QueryPairResponse{
//		Pair: tradePair,
//	})
//
//	if err != nil {
//		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, "could not marshal result to JSON: "+err.Error())
//	}
//
//	return bz, nil
//}

// queryMarketsList queries keeper with market filter and returns found markets list marshaled to bytes.
func queryMarketsList(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryMarketsListParams
	if err := keeper.cdc.UnmarshalBinaryBare(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, fmt.Sprintf("failed to parse params: %s", err))
	}
	// TODO: GET RID OF
	params.All = true

	resp := &types.QueryMarketsListResponse{}
	filterMarket := func(onlyEnabled bool) func(market *types.DerivativeMarket) bool {
		return func(market *types.DerivativeMarket) bool {
			if onlyEnabled && !market.Enabled {
				return false
			}
			resp.Markets = append(resp.Markets, market)
			return false
		}
	}

	keeper.IterateDerivativeMarkets(ctx, filterMarket(!params.All))

	bz, err := keeper.cdc.MarshalBinaryBare(resp)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, "could not marshal result to JSON: "+err.Error())
	}

	return bz, nil
}
