package evm

import (
	"github.com/cosmos/cosmos-sdk/codec"
	version "github.com/cosmos/ethermint/version"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// Supported endpoints
const (
	QueryBlockNumber = "blockNumber"
)



// TODO: Implement querier to route RPC methods. Unable to access RPC otherwise
// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryBlockNumber:
			return queryBlockNumber(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}

func queryProtocolVersion(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	vers := version.ProtocolVersion

	res, err := codec.MarshalJSONIndent(keeper.cdc, vers)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

//func querySyncing(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper ) ([]byte, error) {
//	// TODO: Implement
//	status := true
//
//	res, err := codec.MarshalJSONIndent(keeper.cdc, status)
//	if err != nil {
//		panic("could not marshal result to JSON")
//	}
//
//	return res, nil
//}
//
//func queryCoinbase(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper ) ([]byte, error) {
//	// TODO: Implement
//	status := true
//
//	res, err := codec.MarshalJSONIndent(keeper.cdc, status)
//	if err != nil {
//		panic("could not marshal result to JSON")
//	}
//
//	return res, nil
//}

func queryBlockNumber(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	block := ctx.BlockHeight()

	res, err := codec.MarshalJSONIndent(keeper.cdc, block)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}
