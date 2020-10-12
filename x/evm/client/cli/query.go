package cli

import (
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
		Use:   "storage [account] [key]",
		Short: "Gets storage for an account at a given key",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			// TODO: gRPC
			// queryClient := types.NewQueryClient(clientCtx)

			// account, err := accountToHex(args[0])
			// if err != nil {
			// 	return errors.Wrap(err, "could not parse account address")
			// }

			// key := formatKeyToHash(args[1])

			// res, _, err := clientCtx.Query(
			// 	fmt.Sprintf("custom/%s/storage/%s/%s", queryRoute, account, key))

			// if err != nil {
			// 	return fmt.Errorf("could not resolve: %s", err)
			// }
			// var out types.QueryResStorage
			// cdc.MustUnmarshalJSON(res, &out)
			// return clientCtx.PrintOutput(out)
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCodeCmd queries the code field of a given address
func GetCodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "code [account]",
		Short: "Gets code from an account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			// account, err := accountToHex(args[0])
			// if err != nil {
			// 	return errors.Wrap(err, "could not parse account address")
			// }

			// res, _, err := clientCtx.Query(
			// 	fmt.Sprintf("custom/%s/code/%s", queryRoute, account))

			// if err != nil {
			// 	return fmt.Errorf("could not resolve: %s", err)
			// }

			// var out types.QueryResCode
			// cdc.MustUnmarshalJSON(res, &out)
			// return clientCtx.PrintOutput(out)

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
