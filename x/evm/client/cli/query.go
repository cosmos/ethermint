package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/ethermint/x/evm/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd(moduleName string, cdc *codec.Codec) *cobra.Command {
	evmQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the evm module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	evmQueryCmd.AddCommand(client.GetCommands(
		GetCmdGetBlockNumber(moduleName, cdc),
	)...)
	return evmQueryCmd
}

// GetCmdResolveName queries information about a name
func GetCmdGetBlockNumber(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "block-number",
		Short: "Gets block number (block height)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/blockNumber", queryRoute), nil)
			if err != nil {
				fmt.Printf("could not resolve: %s\n", err)
				return nil
			}

			var out types.QueryResBlockNumber
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
