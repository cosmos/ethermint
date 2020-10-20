package server

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmcfg "github.com/tendermint/tendermint/config"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdkserver "github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
)

// DONTCOVER

// interceptConfigs parses and updates a Tendermint configuration file or
// creates a new one and saves it. It also parses and saves the application
// configuration file. The Tendermint configuration file is parsed given a root
// Viper object, whereas the application is parsed with the private package-aware
// viperCfg object.
func interceptConfigs(rootViper *viper.Viper) (*tmcfg.Config, error) {
	rootDir := rootViper.GetString(flags.FlagHome)
	configPath := filepath.Join(rootDir, "config")
	configFile := filepath.Join(configPath, "config.toml")

	conf := tmcfg.DefaultConfig()

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		tmcfg.EnsureRoot(rootDir)

		if err = conf.ValidateBasic(); err != nil {
			return nil, fmt.Errorf("error in config file: %v", err)
		}

		conf.RPC.PprofListenAddress = "localhost:6060"
		conf.P2P.RecvRate = 5120000
		conf.P2P.SendRate = 5120000
		conf.Consensus.TimeoutCommit = 5 * time.Second
		tmcfg.WriteConfigFile(configFile, conf)
	} else {
		rootViper.SetConfigType("toml")
		rootViper.SetConfigName("config")
		rootViper.AddConfigPath(configPath)
		if err := rootViper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read in app.toml: %w", err)
		}
	}

	// Read into the configuration whatever data the viper instance has for it
	// This may come from the configuration file above but also any of the other sources
	// viper uses
	if err := rootViper.Unmarshal(conf); err != nil {
		return nil, err
	}
	conf.SetRoot(rootDir)

	appConfigFilePath := filepath.Join(configPath, "app.toml")
	if _, err := os.Stat(appConfigFilePath); os.IsNotExist(err) {
		appConf, err := config.ParseConfig(rootViper)
		if err != nil {
			return nil, fmt.Errorf("failed to parse app.toml: %w", err)
		}

		config.WriteConfigFile(appConfigFilePath, appConf)
	}

	rootViper.SetConfigType("toml")
	rootViper.SetConfigName("app")
	rootViper.AddConfigPath(configPath)
	if err := rootViper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read in app.toml: %w", err)
	}

	return conf, nil
}

// add server commands
func AddCommands(rootCmd *cobra.Command, defaultNodeHome string, appCreator AppCreator, appExport types.AppExporter) {
	tendermintCmd := &cobra.Command{
		Use:   "tendermint",
		Short: "Tendermint subcommands",
	}

	tendermintCmd.AddCommand(
		sdkserver.ShowNodeIDCmd(),
		sdkserver.ShowValidatorCmd(),
		sdkserver.ShowAddressCmd(),
		sdkserver.VersionCmd(),
	)

	rootCmd.AddCommand(
		StartCmd(appCreator, defaultNodeHome),
		sdkserver.UnsafeResetAllCmd(),
		flags.LineBreak,
		tendermintCmd,
		sdkserver.ExportCmd(appExport, defaultNodeHome),
		flags.LineBreak,
		version.NewVersionCommand(),
	)
}

func openDB(rootDir string) (dbm.DB, error) {
	dataDir := filepath.Join(rootDir, "data")
	return sdk.NewLevelDB("application", dataDir)
}

func openTraceWriter(traceWriterFile string) (w io.Writer, err error) {
	if traceWriterFile == "" {
		return
	}
	return os.OpenFile(
		traceWriterFile,
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0666,
	)
}
