package cli

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	emint "github.com/cosmos/ethermint/types"
	"github.com/cosmos/ethermint/x/evm/types"
)

// GetTxCmd defines the CLI commands regarding evm module transactions
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	evmTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "EVM transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	evmTxCmd.AddCommand(client.PostCommands(
		GetCmdGenTx(cdc),
	)...)

	return evmTxCmd
}

// GetCmdGenTx generates an Emint transaction (excludes create operations)
func GetCmdGenTx(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "send [to_address] [amount (in photons)] [<data>]",
		Short: "send transaction to address (call operations included)",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// Ambiguously decode amount from any base
			amount, err := strconv.ParseInt(args[1], 0, 64)
			if err != nil {
				return err
			}

			var data []byte
			if len(args) > 2 {
				payload := args[2]
				data, err = hexutil.Decode(payload)
				if err != nil {
					fmt.Println(err)
				}
			}

			toAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return errors.Wrap(err, "must provide a valid Bech32 address for to_address")
			}

			from := cliCtx.GetFromAddress()

			_, seq, err := authtypes.NewAccountRetriever(cliCtx).GetAccountNumberSequence(from)
			if err != nil {
				return errors.Wrap(err, "Could not retrieve account sequence")
			}

			// TODO: Potentially allow overriding of gas price and gas limit
			msg := types.NewEmintMsg(seq, &toAddr, sdk.NewInt(amount), txBldr.Gas(),
				sdk.NewInt(emint.DefaultGasPrice), data, from)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
