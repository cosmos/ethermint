package main

import (
	"encoding/json"
	"io"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clientkeys "github.com/cosmos/cosmos-sdk/client/keys"
	clientrpc "github.com/cosmos/cosmos-sdk/client/rpc"
	sdkcodec "github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/ethermint/app"
	"github.com/cosmos/ethermint/client"
	"github.com/cosmos/ethermint/codec"
	"github.com/cosmos/ethermint/crypto"
	"github.com/cosmos/ethermint/rpc"
	ethermint "github.com/cosmos/ethermint/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	tmamino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

const flagInvCheckPeriod = "inv-check-period"

var (
	cdc            = codec.MakeCodec(app.ModuleBasics)
	invCheckPeriod uint
)

func main() {
	cobra.EnableCommandSorting = false

	cdc := codec.MakeCodec(app.ModuleBasics)
	appCodec := codec.NewAppCodec(cdc)

	tmamino.RegisterKeyType(crypto.PubKeySecp256k1{}, crypto.PubKeyAminoName)
	tmamino.RegisterKeyType(crypto.PrivKeySecp256k1{}, crypto.PrivKeyAminoName)

	keyring.CryptoCdc = cdc
	genutil.ModuleCdc = cdc
	genutiltypes.ModuleCdc = cdc
	clientkeys.KeysCdc = cdc

	config := sdk.GetConfig()
	ethermint.SetBech32Prefixes(config)
	config.Seal()

	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:   "ethermintd",
		Short: "Command line interface for interacting with ethermintd",
	}

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := client.InitConfig(rootCmd); err != nil {
			return err
		}

		return server.PersistentPreRunEFn(ctx)(cmd, args)
	}

	// CLI commands to initialize the chain
	rootCmd.AddCommand(
		client.ValidateChainID(
			genutilcli.InitCmd(ctx, cdc, app.ModuleBasics, app.DefaultNodeHome),
		),
		genutilcli.CollectGenTxsCmd(ctx, cdc, bank.GenesisBalancesIterator{}, app.DefaultNodeHome),
		genutilcli.MigrateGenesisCmd(ctx, cdc),
		genutilcli.GenTxCmd(
			ctx, cdc, app.ModuleBasics, staking.AppModuleBasic{}, bank.GenesisBalancesIterator{},
			app.DefaultNodeHome, app.DefaultCLIHome,
		),
		genutilcli.ValidateGenesisCmd(ctx, cdc, app.ModuleBasics),
		client.TestnetCmd(ctx, cdc, app.ModuleBasics, bank.GenesisBalancesIterator{}),
		// AddGenesisAccountCmd allows users to add accounts to the genesis file
		AddGenesisAccountCmd(ctx, cdc, appCodec, app.DefaultNodeHome, app.DefaultCLIHome),
		flags.NewCompletionCmd(rootCmd, true),
	)

	// Tendermint node base commands
	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	// Construct Root Command
	rootCmd.AddCommand(
		clientrpc.StatusCommand(),
		sdkclient.ConfigCmd(app.DefaultNodeHome),
		queryCmd(cdc),
		txCmd(cdc),
		rpc.EmintServeCmd(cdc),
		client.KeyCommands(),
		version.Cmd,
	)

	rootCmd.PersistentFlags().String(flags.FlagChainID, "", "Chain ID of tendermint node")
	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "EM", app.DefaultNodeHome)
	rootCmd.PersistentFlags().UintVar(&invCheckPeriod, flagInvCheckPeriod,
		0, "Assert registered invariants every N blocks")
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return app.NewEthermintApp(logger, db, traceStore, true, 0,
		baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))))
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {

	if height != -1 {
		emintApp := app.NewEthermintApp(logger, db, traceStore, true, 0)
		err := emintApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return emintApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}

	emintApp := app.NewEthermintApp(logger, db, traceStore, true, 0)

	return emintApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}

func queryCmd(cdc *sdkcodec.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		authcmd.GetAccountCmd(cdc),
		flags.LineBreak,
		authcmd.QueryTxsByEventsCmd(cdc),
		authcmd.QueryTxCmd(cdc),
		flags.LineBreak,
	)

	// add modules' query commands
	app.ModuleBasics.AddQueryCommands(queryCmd, cdc)

	return queryCmd
}

func txCmd(cdc *sdkcodec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(
		bankcmd.SendTxCmd(cdc),
		flags.LineBreak,
		authcmd.GetSignCommand(cdc),
		authcmd.GetMultiSignCommand(cdc),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(cdc),
		authcmd.GetEncodeCommand(cdc),
		authcmd.GetDecodeCommand(cdc),
		flags.LineBreak,
	)

	// add modules' tx commands
	app.ModuleBasics.AddTxCommands(txCmd, cdc)

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
