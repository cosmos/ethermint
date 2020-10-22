package evm

import (
	"fmt"
	"time"

	"github.com/cosmos/ethermint/x/evm/keeper"
	"github.com/cosmos/ethermint/x/evm/types"

	ethcmn "github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns a handler for Ethermint type messages.
func NewHandler(k Keeper) sdk.Handler {
	defer telemetry.MeasureSince(time.Now(), "evm", "state_transition")

	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgEthereumTx:
			// execute state transition
			res, err := msgServer.EthereumTx(sdk.WrapSDKContext(ctx), msg)
			result, err := sdk.WrapServiceResult(ctx, res, err)
			if err != nil {
				return nil, err
			}

			// log state transition result
			var recipientLog string
			if res.ContractAddress != "" {
				recipientLog = fmt.Sprintf("contract address %s", res.ContractAddress)
			} else {
				recipientLog = fmt.Sprintf("recipient address %s", msg.Data.Recipient)
			}

			sender := ethcmn.BytesToAddress(msg.GetFrom().Bytes())

			log := fmt.Sprintf(
				"executed EVM state transition; sender address %s; %s", sender, recipientLog,
			)

			k.Logger(ctx).Info(log)
			result.Log = log

			return result, nil

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}
