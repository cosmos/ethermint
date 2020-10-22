package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/cosmos/ethermint/rpc"
	"github.com/cosmos/ethermint/x/evm/types"
)

// GetQueryCmd returns the parent command for all x/bank CLi query commands.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the evm module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetStorageCmd(),
		GetCodeCmd(),
	)
	return cmd
}

// GetStorageCmd queries a key in an accounts storage
func GetStorageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "storage [address] [key]",
		Short: "Gets storage for an account with a given key and height",
		Long:  "Gets storage for an account with a given key and height. If the height is not provided, it will use the latest height from context.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			address, err := accountToHex(args[0])
			if err != nil {
				return err
			}

			key := formatKeyToHash(args[1])

			req := &types.QueryStorageRequest{
				Address: address,
				Key:     key,
			}

			res, err := queryClient.Storage(rpc.ContextWithHeight(clientCtx.Height), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCodeCmd queries the code field of a given address
func GetCodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "code [address]",
		Short: "Gets code from an account",
		Long:  "Gets code from an account. If the height is not provided, it will use the latest height from context.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			address, err := accountToHex(args[0])
			if err != nil {
				return err
			}

			req := &types.QueryCodeRequest{
				Address: address,
			}

			res, err := queryClient.Code(rpc.ContextWithHeight(clientCtx.Height), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
