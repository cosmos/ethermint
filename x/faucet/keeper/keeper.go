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
	cdc *codec.Codec

	supplyKeeper  types.SupplyKeeper
	stakingKeeper types.StakingKeeper

	// history of users and their funding timeouts
	timeouts map[string]time.Time
}

// NewKeeper creates a new faucet Keeper instance.
func NewKeeper(
	cdc *codec.Codec, storeKey sdk.StoreKey, supplyKeeper types.SupplyKeeper, stakingKeeper types.StakingKeeper,
) Keeper {
	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		supplyKeeper:  supplyKeeper,
		stakingKeeper: stakingKeeper,
		timeouts:      make(map[string]time.Time),
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// MintAndSend mint coins and send to minter.
func (k Keeper) MintAndSend(ctx sdk.Context, minter sdk.AccAddress, mintTime int64) error {

	mining := k.getMining(ctx, minter)

	// refuse mint in 24 hours
	if k.isPresent(ctx, minter) &&
		time.Unix(mining.LastTime, 0).Add(k.Limit).UTC().After(time.Unix(mintTime, 0)) {
		return types.ErrWithdrawTooOften
	}

	denom := k.StakingKeeper.BondDenom(ctx)
	newCoin := sdk.NewCoin(denom, sdk.NewInt(k.amount))
	mining.Total = mining.Total.Add(newCoin)
	mining.LastTime = mintTime
	k.setMining(ctx, minter, mining)

	k.Logger(ctx).Info("Mint coin: %s", newCoin)

	err := k.SupplyKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(newCoin))
	if err != nil {
		return err
	}
	err = k.SupplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, minter, sdk.NewCoins(newCoin))
	if err != nil {
		return err
	}
	return nil
}

// TODO:
func (k Keeper) GetTimout(ctx sdk.Context) time.Duration {
	return time.Second
}

func (k Keeper) SetTimout(ctx sdk.Context, timout time.Duration) {
}

func (k Keeper) getMining(ctx sdk.Context, minter sdk.AccAddress) types.Mining {
	store := ctx.KVStore(k.storeKey)
	if !k.isPresent(ctx, minter) {
		denom := k.StakingKeeper.BondDenom(ctx)
		return types.NewMining(minter, sdk.NewCoin(denom, sdk.NewInt(0)))
	}
	bz := store.Get(minter.Bytes())
	var mining types.Mining
	k.cdc.MustUnmarshalBinaryBare(bz, &mining)
	return mining
}

func (k Keeper) setMining(ctx sdk.Context, minter sdk.AccAddress, mining types.Mining) {
	if mining.Minter.Empty() {
		return
	}
	if !mining.Total.IsPositive() {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(minter.Bytes(), k.cdc.MustMarshalBinaryBare(mining))
}

func (k Keeper) rateLimit(address string) error {
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
		return fmt.Errorf("%s has requested funds within the last %s, wait %s before trying again", address, k.timeout.String(), wait.String())
	}

	// user able to send funds since they have waited for period
	k.timeouts[address] = time.Now().UTC()
	return nil
}
