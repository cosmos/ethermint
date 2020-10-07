package cmd

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/version"
)

// add server commands
func AddCommands(rootCmd *cobra.Command, defaultNodeHome string, appCreator types.AppCreator, appExport types.AppExporter) {
	tendermintCmd := &cobra.Command{
		Use:   "tendermint",
		Short: "Tendermint subcommands",
	}

	tendermintCmd.AddCommand(
		ShowNodeIDCmd(),
		ShowValidatorCmd(),
		ShowAddressCmd(),
		VersionCmd(),
	)

	rootCmd.AddCommand(
		StartCmd(appCreator, defaultNodeHome),
		UnsafeResetAllCmd(),
		flags.LineBreak,
		tendermintCmd,
		ExportCmd(appExport, defaultNodeHome),
		flags.LineBreak,
		version.NewVersionCommand(),
	)
}
