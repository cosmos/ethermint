package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ethermint/x/orders/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of this module maintains collections of orders.
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      codec.BinaryMarshaler
}

// NewKeeper creates new instances of the orders Keeper
func NewKeeper(
	storeKey sdk.StoreKey,
	cdc codec.BinaryMarshaler,
) Keeper {
	return Keeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("injective-chain/modules/%s", types.ModuleName))
}

//// NewQuerier is the module level router for state queries of orders
//func NewQuerier(keeper Keeper) sdk.Querier {
//	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
//		switch path[0] {
//		case types.QueryOrder:
//			return queryOrder(ctx, req, keeper)
//		case types.QueryActiveOrder:
//			return queryActiveOrder(ctx, req, keeper)
//		case types.QueryArchiveOrder:
//			return queryArchiveOrder(ctx, req, keeper)
//		case types.QueryOrdersList:
//			return queryOrdersList(ctx, req, keeper)
//		case types.QueryPair:
//			return queryPair(ctx, req, keeper)
//		case types.QueryPairsList:
//			return queryPairsList(ctx, req, keeper)
//		case types.QueryMarketsList:
//			return queryMarketsList(ctx, req, keeper)
//		case types.QueryZeroExTransaction:
//			return queryZeroExTransaction(ctx, req, keeper)
//		case types.QuerySoftCancelledOrders:
//			return querySoftCancelledOrders(ctx, req, keeper)
//		case types.QueryOutstandingFillRequests:
//			return queryOutstandingFillRequests(ctx, req, keeper)
//		case types.QueryOrderFillRequests:
//			return queryOrderFillRequests(ctx, req, keeper)
//		case types.QueryEvmSyncStatus:
//			return queryEvmSyncStatus(ctx, req, keeper)
//		default:
//			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown orders query endpoint")
//		}
//	}
//}
