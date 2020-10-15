package rpc

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server/config"

	"github.com/ethereum/go-ethereum/rpc"
)

// RegisterRoutes creates a new ethereum JSON-RPC server and recreates a CLI command to start Cosmos REST server with web3 RPC API and
// Cosmos rest-server endpoints
func RegisterRoutes(clientCtx client.Context, r *mux.Router, _ config.APIConfig) {
	server := rpc.NewServer()
	r.HandleFunc("/", server.ServeHTTP).Methods("POST", "OPTIONS")

	// accounts := strings.Split(apiConfig.Accounts, ",")

	// privkeys := []ethsecp256k1.PrivKey{}
	// if len(accounts) != 0 {
	// 	privkeys = addKeys(accounts)
	// }

	apis := GetRPCAPIs(clientCtx)

	// Register all the APIs exposed by the namespace services
	// TODO: handle allowlist and private APIs
	for _, api := range apis {
		if err := server.RegisterName(api.Namespace, api.Service); err != nil {
			panic(err)
		}
	}
}

// StartEthereumWebsocket starts the Filter api websocket
func StartEthereumWebsocket(clientCtx client.Context, apiConfig config.APIConfig) {
	websocketAddr := apiConfig.Address
	ws := newWebsocketsServer(clientCtx, websocketAddr)
	ws.start()
}

// // registerRoutes creates a new server and registers the `/rpc` endpoint.
// // Rpc calls are enabled based on their associated namespace (eg. "eth").
// func registerRoutes() {
// 	s := rpc.NewServer()
// 	accountName := viper.GetString(flagUnlockKey)
// 	accountNames := strings.Split(accountName, ",")

// 	var privkeys []ethsecp256k1.PrivKey
// 	if len(accountName) > 0 {
// 		var err error
// 		inBuf := bufio.NewReader(os.Stdin)

// 		keyringBackend := viper.GetString(flags.FlagKeyringBackend)
// 		passphrase := ""
// 		switch keyringBackend {
// 		case keyring.BackendOS:
// 			break
// 		case keyring.BackendFile:
// 			passphrase, err = input.GetPassword(
// 				"Enter password to unlock key for RPC API: ",
// 				inBuf)
// 			if err != nil {
// 				panic(err)
// 			}
// 		}

// 		privkeys, err = unlockKeyFromNameAndPassphrase(accountNames, passphrase)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}

// 	apis := GetRPCAPIs(rs.clientCtx, privkeys)

// 	// TODO: Allow cli to configure modules https://github.com/cosmos/ethermint/issues/74
// 	whitelist := make(map[string]bool)

// 	// Register all the APIs exposed by the services
// 	for _, api := range apis {
// 		if whitelist[api.Namespace] || (len(whitelist) == 0 && api.Public) {
// 			if err := s.RegisterName(api.Namespace, api.Service); err != nil {
// 				panic(err)
// 			}
// 		} else if !api.Public { // TODO: how to handle private apis? should only accept local calls
// 			if err := s.RegisterName(api.Namespace, api.Service); err != nil {
// 				panic(err)
// 			}
// 		}
// 	}

// 	// Web3 RPC API route
// 	rs.Mux.HandleFunc("/", s.ServeHTTP).Methods("POST", "OPTIONS")

// 	// start websockets server
// 	websocketAddr := viper.GetString(flagWebsocket)
// 	ws := newWebsocketsServer(rs.clientCtx, websocketAddr)
// 	ws.start()
// }

// func unlockKeyFromNameAndPassphrase(accountNames []string, passphrase string) ([]*ethsecp256k1.PrivKey, error) {
// 	kr, err := keyring.New(
// 		sdk.KeyringServiceName(),
// 		viper.GetString(flags.FlagKeyringBackend),
// 		viper.GetString(flags.FlagHome),
// 		os.Stdin,
// 		hd.EthSecp256k1Option(),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// try the for loop with array []string accountNames
// 	// run through the bottom code inside the for loop

// 	keys := make([]*ethsecp256k1.PrivKey, len(accountNames))
// 	for i, acc := range accountNames {
// 		// With keyring keybase, password is not required as it is pulled from the OS prompt
// 		armor, err := kr.ExportPrivKeyArmor(acc, passphrase)
// 		if err != nil {
// 			return nil, err
// 		}

// 		privKey, algo, err := crypto.UnarmorDecryptPrivKey(armor, passphrase)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if algo != ethsecp256k1.KeyType {
// 			return nil, fmt.Errorf("invalid key algorithm, got %s, expected %s", algo, ethsecp256k1.KeyType)
// 		}

// 		var ok bool
// 		keys[i], ok = privKey.(*ethsecp256k1.PrivKey)
// 		if !ok {
// 			return nil, fmt.Errorf("invalid private key type %T, expected %T", privKey, &ethsecp256k1.PrivKey{})
// 		}
// 	}

// 	return keys, nil
// }
