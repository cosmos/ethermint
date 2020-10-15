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