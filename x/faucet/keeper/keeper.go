package keeper

import (
	"fmt"
	"time"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/ethermint/x/faucet/types"
)

// Keeper defines the faucet Keeper.
type Keeper struct {
	cdc          *codec.Codec
	storeKey     sdk.StoreKey
	supplyKeeper types.SupplyKeeper

	// History of users and their funding timeouts. They are reset if the app is reinitialized.
	timeouts map[string]time.Time
}

// NewKeeper creates a new faucet Keeper instance.
func NewKeeper(
	cdc *codec.Codec, storeKey sdk.StoreKey, supplyKeeper types.SupplyKeeper,
) Keeper {
	return Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		supplyKeeper: supplyKeeper,
		timeouts:     make(map[string]time.Time),
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Fund checks for timeout and max thresholds and then mints coins and transfers
// coins to the recipient.
func (k Keeper) Fund(ctx sdk.Context, amount sdk.Coins, recipient sdk.AccAddress) error {
	if err := k.rateLimit(ctx, recipient.String()); err != nil {
		return err
	}

	totalRequested := sdk.ZeroInt()
	for _, coin := range amount {
		totalRequested = totalRequested.Add(coin.Amount)
	}

	// TODO: check max caps

	if err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, amount); err != nil {
		return err
	}

	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, amount); err != nil {
		return err
	}

	k.Logger(ctx).Info(fmt.Sprintf("funded %s to %s", amount, recipient))
	return nil
}

func (k Keeper) GetTimeout(ctx sdk.Context) time.Duration {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TimeoutKey)
	if len(bz) == 0 {
		return time.Duration(0)
	}

	var timeout time.Duration
	k.cdc.MustUnmarshalBinaryBare(bz, &timeout)

	return timeout
}

func (k Keeper) SetTimout(ctx sdk.Context, timout time.Duration) {
}

func (k Keeper) rateLimit(ctx sdk.Context, address string) error {
	// first time requester, can send request
	lastRequest, ok := k.timeouts[address]
	if !ok {
		k.timeouts[address] = time.Now().UTC()
		return nil
	}

	defaultTimeout := k.GetTimeout(ctx)
	sinceLastRequest := time.Since(lastRequest)

	if defaultTimeout > sinceLastRequest {
		wait := defaultTimeout - sinceLastRequest
		return fmt.Errorf("%s has requested funds within the last %s, wait %s before trying again", address, defaultTimeout.String(), wait.String())
	}

	// user able to send funds since they have waited for period
	k.timeouts[address] = time.Now().UTC()
	return nil
}
