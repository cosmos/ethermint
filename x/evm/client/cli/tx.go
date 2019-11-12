package cli

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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

const (
	flagToAddress = "to-address"
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

// GetCmdGenTx generates an ethereum transaction wrapped in a Cosmos standard transaction
func GetCmdGenTx(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-tx [amount (in photons)] [<data>]",
		Short: "generate eth tx wrapped in a Cosmos Standard tx",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			// Ambiguously decode amount from any base
			amount, err := strconv.ParseInt(args[0], 0, 64)
			if err != nil {
				return err
			}

			var data []byte
			if len(args) > 1 {
				payload := args[1]
				data, err = hexutil.Decode(payload)
				if err != nil {
					fmt.Println(err)
				}
			}

			var toAddr *sdk.AccAddress
			toFlag := viper.GetString(flagToAddress)
			if toFlag != "" {
				addr, err := sdk.AccAddressFromBech32(toFlag)
				if err != nil {
					return errors.Wrap(err, "must provide a valid Bech32 address for ")
				}
				toAddr = &addr
			}

			from := cliCtx.GetFromAddress()

			_, seq, err := authtypes.NewAccountRetriever(cliCtx).GetAccountNumberSequence(from)
			if err != nil {
				return errors.Wrap(err, "Could not retrieve account sequence")
			}

			msg := types.NewEmintMsg(seq, toAddr, sdk.NewInt(amount), txBldr.Gas(),
				sdk.NewInt(emint.DefaultGasPrice), data, from)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	// Optional parameter flags
	cmd.Flags().String(flagToAddress, "", "set to address for transaction")

	return cmd
}
