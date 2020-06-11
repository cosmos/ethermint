package main

import (
	"bufio"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/ethermint/app"
	emintcrypto "github.com/cosmos/ethermint/crypto"
	"github.com/cosmos/ethermint/rpc"

	tmamino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	clientkeys "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	clientrpc "github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
)

const flagUnlockKey = "unlock-key"

var (
	appCodec, cdc = app.MakeCodecs()
)

func main() {
	// Configure cobra to sort commands
	cobra.EnableCommandSorting = false

	tmamino.RegisterKeyType(emintcrypto.PubKeySecp256k1{}, emintcrypto.PubKeyAminoName)
	tmamino.RegisterKeyType(emintcrypto.PrivKeySecp256k1{}, emintcrypto.PrivKeyAminoName)

	keyring.CryptoCdc = cdc
	clientkeys.KeysCdc = cdc

	// Read in the configuration file for the sdk
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	rootCmd := &cobra.Command{
		Use:   "emintcli",
		Short: "Command line interface for interacting with emintd",
	}

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(flags.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// Construct Root Command
	rootCmd.AddCommand(
		clientrpc.StatusCommand(),
		client.ConfigCmd(app.DefaultCLIHome),
		queryCmd(cdc),
		txCmd(cdc),
		EmintServeCmd(cdc),
		flags.LineBreak,
		keyCommands(),
		flags.LineBreak,
		version.Cmd,
		flags.NewCompletionCmd(rootCmd, true),
	)

	// Add flags and prefix all env exposed with EM
	executor := cli.PrepareMainCmd(rootCmd, "EM", app.DefaultCLIHome)

	err := executor.Execute()
	if err != nil {
		panic(fmt.Errorf("failed executing CLI command: %w", err))
	}
}

func queryCmd(cdc *codec.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		authcmd.GetAccountCmd(cdc),
		flags.LineBreak,
		clientrpc.ValidatorCommand(cdc),
		clientrpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(cdc),
		authcmd.QueryTxCmd(cdc),
		flags.LineBreak,
	)

	// add modules' query commands
	clientCtx := client.Context{}
	clientCtx = clientCtx.
		WithJSONMarshaler(appCodec).
		WithCodec(cdc)
	app.ModuleBasics.AddQueryCommands(queryCmd, clientCtx)

	return queryCmd
}

func txCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	clientCtx := client.Context{}
	clientCtx = clientCtx.
		WithJSONMarshaler(appCodec).
		WithTxGenerator(authtypes.StdTxGenerator{Cdc: cdc}).
		WithAccountRetriever(authtypes.NewAccountRetriever(appCodec)).
		WithCodec(cdc)

	txCmd.AddCommand(
		bankcmd.NewSendTxCmd(clientCtx),
		flags.LineBreak,
		authcmd.GetSignCommand(cdc),
		authcmd.GetSignBatchCommand(cdc),
		authcmd.GetMultiSignCommand(cdc),
		authcmd.GetValidateSignaturesCommand(cdc),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(cdc),
		authcmd.GetEncodeCommand(cdc),
		authcmd.GetDecodeCommand(cdc),
		flags.LineBreak,
	)

	// add modules' tx commands
	app.ModuleBasics.AddTxCommands(txCmd, clientCtx)

	// remove auth and bank commands as they're mounted under the root tx command
	var cmdsToRemove []*cobra.Command

	for _, cmd := range txCmd.Commands() {
		if cmd.Use == auth.ModuleName || cmd.Use == bank.ModuleName {
			cmdsToRemove = append(cmdsToRemove, cmd)
		}
	}

	txCmd.RemoveCommand(cmdsToRemove...)

	return txCmd
}

// EmintServeCmd creates a CLI command to start Cosmos REST server with web3 RPC API and
// Cosmos rest-server endpoints
func EmintServeCmd(cdc *codec.Codec) *cobra.Command {
	cmd := lcd.ServeCommand(cdc, registerRoutes)
	cmd.Flags().String(flagUnlockKey, "", "Select a key to unlock on the RPC server")
	cmd.Flags().StringP(flags.FlagBroadcastMode, "b", flags.BroadcastSync, "Transaction broadcasting mode (sync|async|block)")
	return cmd
}

// registerRoutes creates a new server and registers the `/rpc` endpoint.
// Rpc calls are enabled based on their associated module (eg. "eth").
func registerRoutes(rs *lcd.RestServer) {
	s := rpc.NewServer()
	accountName := viper.GetString(flagUnlockKey)

	var emintKey emintcrypto.PrivKeySecp256k1
	if len(accountName) > 0 {
		var err error
		inBuf := bufio.NewReader(os.Stdin)

		keyringBackend := viper.GetString(flags.FlagKeyringBackend)
		passphrase := ""
		switch keyringBackend {
		case keyring.BackendOS:
			break
		case keyring.BackendFile:
			passphrase, err = input.GetPassword(
				"Enter password to unlock key for RPC API: ",
				inBuf)
			if err != nil {
				panic(err)
			}
		}

		emintKey, err = unlockKeyFromNameAndPassphrase(accountName, passphrase)
		if err != nil {
			panic(err)
		}
	}

	apis := GetRPCAPIs(rs.ClientCtx, emintKey)

	// TODO: Allow cli to configure modules https://github.com/ChainSafe/ethermint/issues/74
	whitelist := make(map[string]bool)

	// Register all the APIs exposed by the services
	for _, api := range apis {
		if whitelist[api.Namespace] || (len(whitelist) == 0 && api.Public) {
			if err := s.RegisterName(api.Namespace, api.Service); err != nil {
				panic(err)
			}
		}
	}

	// Web3 RPC API route
	rs.Mux.HandleFunc("/", s.ServeHTTP).Methods("POST", "OPTIONS")

	// Register all other Cosmos routes
	client.RegisterRoutes(rs.ClientCtx, rs.Mux)
	authrest.RegisterTxRoutes(rs.ClientCtx, rs.Mux)
	app.ModuleBasics.RegisterRESTRoutes(rs.ClientCtx, rs.Mux)
}

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}

	cfgFile := path.Join(home, "config", "config.toml")
	if _, err := os.Stat(cfgFile); err == nil {
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}
	if err := viper.BindPFlag(flags.FlagChainID, cmd.PersistentFlags().Lookup(flags.FlagChainID)); err != nil {
		return err
	}
	if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}
