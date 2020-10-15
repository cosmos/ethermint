package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/ethermint/x/orders/internal/types"
)

// queryPair queries keeper with trading pair params and returns found pair marshaled to bytes.
func queryPair(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryPairParams

	if err := keeper.cdc.UnmarshalBinaryBare(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, fmt.Sprintf("failed to parse params: %s", err))
	}

	var tradePair *types.TradePair

	if (params.Hash != common.Hash{}) {
		tradePair = keeper.GetTradePair(ctx, params.Hash)
	} else if len(params.MakerAssetData) > 0 && len(params.TakerAssetData) > 0 {
		hash, err := (&types.TradePair{
			MakerAssetData: params.MakerAssetData,
			TakerAssetData: params.TakerAssetData,
		}).Hash()
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrPairNotFound, "failed to compute pair hash: "+err.Error())
		}
		tradePair = keeper.GetTradePair(ctx, hash)
	} else if len(params.Name) > 0 {
		tradePair = keeper.GetTradePairByName(ctx, params.Name)
	}

	if tradePair == nil {
		return nil, sdkerrors.Wrap(types.ErrPairNotFound, params.Name)
	}

	bz, err := keeper.cdc.MarshalBinaryBare(&types.QueryPairResponse{
		Pair: tradePair,
	})

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, "could not marshal result to JSON: "+err.Error())
	}

	return bz, nil
}

// queryPairsList queries keeper with trading pair filter and returns found pairs list marshaled to bytes.
func queryPairsList(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryPairsListParams
	if err := keeper.cdc.UnmarshalBinaryBare(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, fmt.Sprintf("failed to parse params: %s", err))
	}

	resp := &types.QueryPairsListResponse{}
	filterTradePair := func(onlyEnabled bool) func(tradePair *types.TradePair) bool {
		return func(tradePair *types.TradePair) bool {
			if onlyEnabled && !tradePair.Enabled {
				return false
			}
			resp.Pairs = append(resp.Pairs, tradePair)
			return false
		}
	}

	keeper.IterateTradePairs(ctx, filterTradePair(!params.All))

	bz, err := keeper.cdc.MarshalBinaryBare(resp)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, "could not marshal result to JSON: "+err.Error())
	}

	return bz, nil
}
