package rpc

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"

	"github.com/cosmos/ethermint/app"
	"github.com/cosmos/ethermint/crypto/ethsecp256k1"
	"github.com/cosmos/ethermint/crypto/hd"
	"github.com/cosmos/ethermint/rpc/websockets"
	"github.com/ethereum/go-ethereum/rpc"
)

	"github.com/ethereum/go-ethereum/rpc"
)

// RegisterRoutes creates a new ethereum JSON-RPC server and recreates a CLI command to start Cosmos REST server with web3 RPC API and
// Cosmos rest-server endpoints
func RegisterEthereum(clientCtx client.Context, r *mux.Router, apiConfig config.EthereumConfig) {
	server := rpc.NewServer()
	r.HandleFunc("/", server.ServeHTTP).Methods("POST", "OPTIONS")

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
func StartEthereumWebsocket(clientCtx client.Context, apiConfig server.APIConfig) {
	ws := websockets.NewServer(clientCtx, apiConfig.Address, apiConfig.WebsocketAddress)
	ws.Start()
}
