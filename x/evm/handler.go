package evm

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ethermint/x/evm/types"
)

// NewHandler returns a handler for ethermint type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.EthereumTxMsg:
			return handleETHTxMsg(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized ethermint Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle an etherem specific tx
func handleETHTxMsg(ctx sdk.Context, keeper Keeper, msg types.EthereumTxMsg) sdk.Result {
	// TODO: Implement transaction logic
	return sdk.Result{}
}
