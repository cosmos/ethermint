package cli

import (
	"bufio"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/cosmos/ethermint/x/faucet/types"
)

// GetTxCmd return faucet sub-command for tx
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	faucetTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "faucet transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	faucetTxCmd.AddCommand(flags.PostCommands(
		GetCmdFund(cdc),
	)...)

	return faucetTxCmd
}

// GetCmdFund is the CLI command to fund an address with the requested coins
func GetCmdFund(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "fund [amount] [[address]]",
		Short: "fund an address with the requested coins",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))


			amount, err := sdk.ParseCoins(args[0])
			if err != nil {
				return err
			}

			var recipient sdk.AccAddress
			if len(args) == 1 {
				recipient = cliCtx.GetFromAddress()
			} else {
				recipient, err := sdk.AccAddressFromBech32(args[1])
			}

			if err != nil {
				return err
			}

			msg := types.NewMsgFund(cliCtx.GetFromAddress(), recipient, time.Now().Unix()
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}


