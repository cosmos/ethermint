package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/ethermint/x/orders/internal/types"
)

func queryEvmSyncStatus(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryEvmSyncStatusParams

	if err := keeper.cdc.UnmarshalBinaryBare(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, fmt.Sprintf("failed to parse params: %s", err))
	}

	evmSyncStatus := keeper.GetEvmSyncStatus(ctx)
	if evmSyncStatus == nil {
		evmSyncStatus = &types.EvmSyncStatus{}
	}

	bz, err := keeper.cdc.MarshalBinaryBare(&types.QueryEvmSyncStatusResponse{
		LatestBlockSynced: evmSyncStatus.LatestBlockSynced,
	})

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, "could not marshal result to JSON: "+err.Error())
	}

	return bz, nil
}

func queryFuturesEvmSyncStatus(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryEvmSyncStatusParams

	if err := keeper.cdc.UnmarshalBinaryBare(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, fmt.Sprintf("failed to parse params: %s", err))
	}

	evmSyncStatus := keeper.GetFuturesEvmSyncStatus(ctx)
	if evmSyncStatus == nil {
		evmSyncStatus = &types.EvmSyncStatus{}
	}

	bz, err := keeper.cdc.MarshalBinaryBare(&types.QueryEvmSyncStatusResponse{
		LatestBlockSynced: evmSyncStatus.LatestBlockSynced,
	})

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, "could not marshal result to JSON: "+err.Error())
	}

	return bz, nil
}
