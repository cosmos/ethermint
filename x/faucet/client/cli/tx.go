package cli

import (
	"bufio"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/crypto/keys"

	"github.com/cosmos/ethermint/crypto"
	"github.com/cosmos/ethermint/x/faucet/types"
)

// GetTxCmd return faucet sub-command for tx
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	faucetTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "faucet transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	faucetTxCmd.AddCommand(flags.PostCommands(
		GetCmdRequest(cdc),
	)...)

	return faucetTxCmd
}

// GetCmdRequest is the CLI command to fund an address with the requested coins
func GetCmdRequest(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "request [amount] [other-recipient (optional)]",
		Short: "request an address with the requested coins",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authclient.GetTxEncoder(cdc))

			var err error
			cliCtx.Keybase, err = keys.NewKeyring(
				sdk.KeyringServiceName(),
				viper.GetString(flags.FlagKeyringBackend),
				viper.GetString(flags.FlagHome),
				cliCtx.Input,
				crypto.EthSecp256k1Options()...,
			)
			if err != nil {
				fmt.Println("failed to create keybase")
				return err
			}

			amount, err := sdk.ParseCoins(args[0])
			if err != nil {
				return err
			}

			var recipient sdk.AccAddress
			if len(args) == 1 {
				recipient = cliCtx.GetFromAddress()
			} else {
				recipient, err = sdk.AccAddressFromBech32(args[1])
			}

			if err != nil {
				fmt.Println("failed to create acc address")
				return err
			}

			msg := types.NewMsgFund(amount, cliCtx.GetFromAddress(), recipient)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			accRet := auth.NewAccountRetriever(cliCtx)
			err = accRet.EnsureExists(recipient)
			if err != nil {
				// account doesn't exist
				return fmt.Errorf("nonexistent account %s: %s", recipient, err)
			}

			err = authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
			if err != nil {
				fmt.Println("failed to run request command:", err)
			}
			return err
		},
	}
}
