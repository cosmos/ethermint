package client

import (
	"bufio"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/ethermint/crypto/hd"
)

const (
	flagDryRun = "dry-run"
)

// KeyCommands registers a sub-tree of commands to interact with
// local private key storage.
func KeyCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keys",
		Short: "Add or view local private keys",
		Long: `Keys allows you to manage your local keystore for tendermint.

    These keys may be in any format supported by go-crypto and can be
    used by light-clients, full nodes, or any other application that
    needs to sign with a private key.`,
	}

	// support adding Ethereum supported keys
	addCmd := keys.AddKeyCommand()

	// update the default signing algorithm value to "eth_secp256k1"
	algoFlag := addCmd.Flag("algo")
	algoFlag.DefValue = string(hd.EthSecp256k1Type)
	err := algoFlag.Value.Set(string(hd.EthSecp256k1Type))
	if err != nil {
		panic(err)
	}
	addCmd.RunE = runAddCmd

	cmd.AddCommand(
		keys.MnemonicKeyCommand(),
		keys.AddKeyCommand(),
		keys.ExportKeyCommand(),
		keys.ImportKeyCommand(),
		keys.ListKeysCmd(),
		keys.ShowKeysCmd(),
		flags.LineBreak,
		keys.DeleteKeyCommand(),
		keys.ParseKeyStringCommand(),
		keys.MigrateCommand(),
		flags.LineBreak,
		UnsafeExportEthKeyCommand(),
	)
	return cmd
}

func runAddCmd(cmd *cobra.Command, args []string) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())
	kb, err := getKeybase(viper.GetBool(flagDryRun), inBuf)
	if err != nil {
		return err
	}

	return keys.RunAddCmd(cmd, args, kb, inBuf)
}

func getKeybase(transient bool, buf io.Reader) (keyring.Keyring, error) {
	if transient {
		return keyring.NewInMemory(
			hd.EthSecp256k1Option(),
		), nil
	}

	return keyring.New(
		sdk.KeyringServiceName(),
		viper.GetString(flags.FlagKeyringBackend),
		viper.GetString(flags.FlagHome),
		buf,
		hd.EthSecp256k1Option(),
	)
}
