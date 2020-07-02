package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/ethermint/x/faucet/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewQuerier is the module level router for state queries
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryFunded:
			return queryFunded(ctx, req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown query endpoint")
		}
	}
}

func queryFunded(ctx sdk.Context, _ abci.RequestQuery, k Keeper) ([]byte, error) {
	funded := k.GetFunded(ctx)

	bz, err := codec.MarshalJSONIndent(k.cdc, funded)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
