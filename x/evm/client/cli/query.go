package cli

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/ethermint/x/evm/types"
)

// GetQueryCmd defines evm module queries through the cli
func GetQueryCmd(moduleName string, cdc *codec.Codec) *cobra.Command {
	evmQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the evm module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	evmQueryCmd.AddCommand(client.GetCommands(
		GetCmdGetStorageAt(moduleName, cdc),
		GetCmdGetCode(moduleName, cdc),
	)...)
	return evmQueryCmd
}

// GetCmdGetStorageAt queries a key in an accounts storage
func GetCmdGetStorageAt(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "storage [account] [key]",
		Short: "Gets storage for an account at a given key",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			account, err := accountToHex(args[0])
			if err != nil {
				return errors.Wrap(err, "could not parse account address")
			}
			key := args[1]

			res, _, err := cliCtx.Query(
				fmt.Sprintf("custom/%s/storage/%s/%s", queryRoute, account, key))

			if err != nil {
				return fmt.Errorf("could not resolve: %s", err)
			}
			var out types.QueryResStorage
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdGetCode queries the code field of a given address
func GetCmdGetCode(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "code [account]",
		Short: "Gets code from an account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			account, err := accountToHex(args[0])
			if err != nil {
				return errors.Wrap(err, "could not parse account address")
			}

			res, _, err := cliCtx.Query(
				fmt.Sprintf("custom/%s/code/%s", queryRoute, account))

			if err != nil {
				return fmt.Errorf("could not resolve: %s", err)
			}

			var out types.QueryResCode
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func accountToHex(addr string) (string, error) {
	if strings.HasPrefix("cosmos", addr) {
		// Check to see if address is Cosmos bech32 formatted
		toAddr, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			return "", errors.Wrap(err, "must provide a valid Bech32 address")
		}
		ethAddr := common.BytesToAddress(toAddr.Bytes())
		return ethAddr.Hex(), nil
	}

	if !strings.HasPrefix(addr, "0x") {
		addr = "0x" + addr
	}

	valid := common.IsHexAddress(addr)
	if !valid {
		return "", fmt.Errorf("%s is not a valid Ethereum or Cosmos address", addr)
	}

	return addr, nil
}
