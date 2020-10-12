package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/ethermint/x/evm/types"
)

// NewTxCmd returns a root CLI command handler for all x/evm transaction commands.
func NewTxCmd() *cobra.Command {
	evmTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "EVM transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	evmTxCmd.AddCommand(
		NewSendTxCmd(),
		NewCreateContractCmd(),
	)

	return evmTxCmd
}

// NewSendTxCmd generates an Ethermint transaction (excludes create operations)
func NewSendTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [to_address] [amount (in aphotons)] [<data>]",
		Short: "send transaction to address (call operations included)",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			toAddr, err := cosmosAddressFromArg(args[0])
			if err != nil {
				return errors.Wrap(err, "must provide a valid Bech32 address for to_address")
			}

			// Ambiguously decode amount from any base
			amount, err := strconv.ParseInt(args[1], 0, 64)
			if err != nil {
				return err
			}

			var data []byte
			if len(args) > 2 {
				payload := args[2]
				if !strings.HasPrefix(payload, "0x") {
					payload = "0x" + payload
				}

				data, err = hexutil.Decode(payload)
				if err != nil {
					return err
				}
			}

			from := clientCtx.GetFromAddress()

			_, seq, err := clientCtx.AccountRetriever.GetAccountNumberSequence(clientCtx, from)
			if err != nil {
				return errors.Wrap(err, "Could not retrieve account sequence")
			}

			txFactory := tx.NewFactoryCLI(clientCtx, cmd.Flags())
			gasPrice := txFactory.GasPrices()[0].Amount.TruncateInt()

			msg := types.NewMsgEthermint(
				seq,
				&toAddr,
				sdk.NewInt(amount),
				txFactory.Gas(),
				gasPrice,
				data,
				from,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewCreateContractCmd generates an Ethermint transaction (excludes create operations)
func NewCreateContractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [contract bytecode] [<amount (in aphotons)>]",
		Short: "create contract through the evm using compiled bytecode",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			payload := args[0]
			if !strings.HasPrefix(payload, "0x") {
				payload = "0x" + payload
			}

			data, err := hexutil.Decode(payload)
			if err != nil {
				return err
			}

			var amount int64
			if len(args) > 1 {
				// Ambiguously decode amount from any base
				amount, err = strconv.ParseInt(args[1], 0, 64)
				if err != nil {
					return errors.Wrap(err, "invalid amount")
				}
			}

			from := clientCtx.GetFromAddress()

			_, seq, err := clientCtx.AccountRetriever.GetAccountNumberSequence(clientCtx, from)
			if err != nil {
				return errors.Wrap(err, "Could not retrieve account sequence")
			}

			txFactory := tx.NewFactoryCLI(clientCtx, cmd.Flags())
			gasPrice := txFactory.GasPrices()[0].Amount.TruncateInt()

			msg := types.NewMsgEthermint(
				seq,
				nil,
				sdk.NewInt(amount),
				txFactory.Gas(),
				gasPrice,
				data,
				from,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			if err := tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg); err != nil {
				return err
			}

			contractAddr := ethcrypto.CreateAddress(common.BytesToAddress(from.Bytes()), seq)
			fmt.Printf(
				"Contract will be deployed to: \nHex: %s\nCosmos Address: %s\n",
				contractAddr.Hex(),
				sdk.AccAddress(contractAddr.Bytes()),
			)
			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
