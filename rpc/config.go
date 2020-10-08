package rpc

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"

	"github.com/cosmos/ethermint/app"
	"github.com/cosmos/ethermint/crypto/ethsecp256k1"
	"github.com/cosmos/ethermint/crypto/hd"

	"github.com/ethereum/go-ethereum/rpc"
)

const (
	flagUnlockKey = "unlock-key"
	flagWebsocket = "wsport"
)

// EmintServeCmd creates a CLI command to start Cosmos REST server with web3 RPC API and
// Cosmos rest-server endpoints
func EmintServeCmd(cdc *codec.LegacyAmino) *cobra.Command {
	cmd := lcd.ServeCommand(cdc, registerRoutes)
	cmd.Flags().String(flagUnlockKey, "", "Select a key to unlock on the RPC server")
	cmd.Flags().String(flagWebsocket, "8546", "websocket port to listen to")
	cmd.Flags().StringP(flags.FlagBroadcastMode, "b", flags.BroadcastSync, "Transaction broadcasting mode (sync|async|block)")
	return cmd
}

// registerRoutes creates a new server and registers the `/rpc` endpoint.
// Rpc calls are enabled based on their associated module (eg. "eth").
func registerRoutes(rs *lcd.RestServer) {
	s := rpc.NewServer()
	accountName := viper.GetString(flagUnlockKey)
	accountNames := strings.Split(accountName, ",")

	var privkeys []ethsecp256k1.PrivKey
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

		privkeys, err = unlockKeyFromNameAndPassphrase(accountNames, passphrase)
		if err != nil {
			panic(err)
		}
	}

	apis := GetRPCAPIs(rs.clientCtx, privkeys)

	// TODO: Allow cli to configure modules https://github.com/ChainSafe/ethermint/issues/74
	whitelist := make(map[string]bool)

	// Register all the APIs exposed by the services
	for _, api := range apis {
		if whitelist[api.Namespace] || (len(whitelist) == 0 && api.Public) {
			if err := s.RegisterName(api.Namespace, api.Service); err != nil {
				panic(err)
			}
		} else if !api.Public { // TODO: how to handle private apis? should only accept local calls
			if err := s.RegisterName(api.Namespace, api.Service); err != nil {
				panic(err)
			}
		}
	}

	// Web3 RPC API route
	rs.Mux.HandleFunc("/", s.ServeHTTP).Methods("POST", "OPTIONS")

	// Register all other Cosmos routes
	client.RegisterRoutes(rs.clientCtx, rs.Mux)
	authrest.RegisterTxRoutes(rs.clientCtx, rs.Mux)
	app.ModuleBasics.RegisterRESTRoutes(rs.clientCtx, rs.Mux)

	// start websockets server
	websocketAddr := viper.GetString(flagWebsocket)
	ws := newWebsocketsServer(rs.clientCtx, websocketAddr)
	ws.start()
}

func unlockKeyFromNameAndPassphrase(accountNames []string, passphrase string) ([]ethsecp256k1.PrivKey, error) {
	kr, err := keyring.New(
		sdk.KeyringServiceName(),
		viper.GetString(flags.FlagKeyringBackend),
		viper.GetString(flags.FlagHome),
		os.Stdin,
		hd.EthSecp256k1Option(),
	)
	if err != nil {
		return []ethsecp256k1.PrivKey{}, err
	}

	// try the for loop with array []string accountNames
	// run through the bottom code inside the for loop

	keys := make([]ethsecp256k1.PrivKey, len(accountNames))
	for i, acc := range accountNames {
		// With keyring keybase, password is not required as it is pulled from the OS prompt
		armor, err := kr.ExportPrivKeyArmor(acc, passphrase)
		if err != nil {
			return err
		}

		privKey, algo, err := crypto.UnarmorDecryptPrivKey(armor, passphrase)
		if err != nil {
			return err
		}

		if algo != ethsecp256k1.KeyType {
			return []ethsecp256k1.PrivKey{}, fmt.Errorf("invalid key algorithm, got %s, expected %s", algo, ethsecp256k1.KeyType)
		}

		var ok bool
		keys[i], ok = privKey.(*ethsecp256k1.PrivKey)
		if !ok {
			return []ethsecp256k1.PrivKey{}, fmt.Errorf("invalid private key type %T, expected %T", privKey, &ethsecp256k1.PrivKey{})
		}
	}

	return keys, nil
}
