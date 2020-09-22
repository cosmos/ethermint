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
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client/utils"

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
		Use:   "request [faucet-address] [recipient-address] [amount] ",
		Short: "request an address with the requested coins",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContextWithInputAndFrom(inBuf, args[0]).WithCodec(cdc)
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

			fmt.Println("-- Keyring! -- ")

			infos, err := cliCtx.Keybase.List()
			if err != nil {
				return err
			}

			for _, info := range infos {
				cmd.Println("name: ", info.GetName())
				fmt.Println("pubkey: ", info.GetPubKey())
				fmt.Println("addr: ", info.GetAddress().String())
				path, err := info.GetPath()
				fmt.Println("path: ", path)
				fmt.Println("err? ", err)
			}

			var faucet, recipient sdk.AccAddress

			faucet, err = sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			// if len(args[0]) < 42 {
			// 	for _, info := range infos {
			// 		if args[0] == info.GetName() {
			// 			faucet = info.GetAddress()
			// 		}
			// 	}
			// } else {
			// 	faucet, err = sdk.AccAddressFromBech32(args[0])
			// 	if err != nil {
			// 		return err
			// 	}
			// }

			if len(args[1]) < 42 {
				for _, info := range infos {
					if args[1] == info.GetName() {
						recipient = info.GetAddress()
					}
				}
			} else {
				recipient, err = sdk.AccAddressFromBech32(args[1])
				if err != nil {
					return err
				}
			}

			amount, err := sdk.ParseCoins(args[2])
			if err != nil {
				return err
			}

			fmt.Println("recipient: ", recipient)

			fmt.Println("faucet addr: ", faucet)

			accRet := auth.NewAccountRetriever(cliCtx)
			fmt.Println("accRet: ", accRet)

			err = accRet.EnsureExists(faucet)
			if err != nil {
				return fmt.Errorf("faucet account does not exist: %s", faucet)
			}

			msg := types.NewMsgFund(amount, faucet, recipient)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			fmt.Println("msg:", msg)

			err = authclient.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
			if err != nil {
				fmt.Println("failed to run request command:", err)
			}

			return err
		},
	}
}
