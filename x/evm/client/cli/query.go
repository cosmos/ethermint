package cli

import (
	"context"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

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
		Use:   "storage [account] [key] [height]",
		Short: "Gets storage for an account with a given key and height",
		Long:  "Gets storage for an account with a given key and height. If the height is not provided, it will use the latest height from context.",
		Args:  cobra.RangeArgs(2, 3),
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

			var height int64
			if len(args) == 3 {
				height, err = strconv.ParseInt(args[2], 10, 64)
				if err != nil {
					return err
				}
			}

			req := &types.QueryStorageRequest{
				Address: address,
				Key:     key,
				Height:  height,
			}

			res, err := queryClient.Storage(context.Background(), req)
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
		Use:   "code [account] [height]",
		Short: "Gets code from an account",
		Long:  "Gets code from an account. If the height is not provided, it will use the latest height from context.",
		Args:  cobra.RangeArgs(1, 2),
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

			var height int64
			if len(args) == 2 {
				height, err = strconv.ParseInt(args[1], 10, 64)
				if err != nil {
					return err
				}
			}

			req := &types.QueryCodeRequest{
				Address: address,
				Height:  height,
			}

			res, err := queryClient.Code(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
