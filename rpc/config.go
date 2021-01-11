package rpc

import (
<<<<<<< HEAD
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/cosmos/ethermint/rpc/websockets"
	"github.com/cosmos/ethermint/server/config"
=======
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ethermint/app"
	"github.com/cosmos/ethermint/crypto/ethsecp256k1"
	"github.com/cosmos/ethermint/crypto/hd"
	"github.com/cosmos/ethermint/rpc/websockets"
	evmrest "github.com/cosmos/ethermint/x/evm/client/rest"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/viper"
)
>>>>>>> 9ecd264ae0e8c6ace10abe3bc45ab8d3c29dc10f

	"github.com/ethereum/go-ethereum/rpc"
)

// RegisterEthereum creates a new ethereum JSON-RPC server and recreates a CLI command to start Cosmos REST server with web3 RPC API and
// Cosmos rest-server endpoints
func RegisterEthereum(clientCtx client.Context, r *mux.Router) {
	server := rpc.NewServer()
	r.HandleFunc("/", server.ServeHTTP).Methods("POST", "OPTIONS")

	apis := GetAPIs(clientCtx)

	// Register all the APIs exposed by the namespace services
	// TODO: handle allowlist and private APIs
	for _, api := range apis {
		if err := server.RegisterName(api.Namespace, api.Service); err != nil {
			panic(err)
		}
	}
<<<<<<< HEAD
=======

	// Web3 RPC API route
	rs.Mux.HandleFunc("/", server.ServeHTTP).Methods("POST", "OPTIONS")

	// Register all other Cosmos routes
	client.RegisterRoutes(rs.CliCtx, rs.Mux)
	evmrest.RegisterRoutes(rs.CliCtx, rs.Mux)
	app.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)

	// start websockets server
	websocketAddr := viper.GetString(flagWebsocket)
	ws := websockets.NewServer(rs.CliCtx, websocketAddr)
	ws.Start()
>>>>>>> 9ecd264ae0e8c6ace10abe3bc45ab8d3c29dc10f
}

// StartEthereumWebsocket starts the Filter api websocket
func StartEthereumWebsocket(clientCtx client.Context, apiConfig config.EthereumConfig) {
	ws := websockets.NewServer(clientCtx, apiConfig.WebsocketAddress)
	ws.Start()
}
