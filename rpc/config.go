package rpc

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/codec"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"

	"github.com/cosmos/ethermint/app"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	flagWebsocket = "wsport"
)

// EmintServeCmd creates a CLI command to start Cosmos REST server with web3 RPC API and
// Cosmos rest-server endpoints
func EmintServeCmd(cdc *codec.Codec) *cobra.Command {
	cmd := lcd.ServeCommand(cdc, registerRoutes)
	cmd.Flags().String(flagWebsocket, "8546", "websocket port to listen to")
	cmd.Flags().StringP(flags.FlagBroadcastMode, "b", flags.BroadcastSync, "Transaction broadcasting mode (sync|async|block)")
	return cmd
}

// registerRoutes creates a new server and registers the `/rpc` endpoint.
// Rpc calls are enabled based on their associated module (eg. "eth").
func registerRoutes(rs *lcd.RestServer) {
	s := rpc.NewServer()

	apis := GetRPCAPIs(rs.CliCtx)

	// Register all the APIs exposed by the services
	for _, api := range apis {
		if err := s.RegisterName(api.Namespace, api.Service); err != nil {
			panic(err)
		}
	}

	// Web3 RPC API route
	rs.Mux.HandleFunc("/", s.ServeHTTP).Methods("POST", "OPTIONS")

	// Register all other Cosmos routes
	client.RegisterRoutes(rs.CliCtx, rs.Mux)
	authrest.RegisterTxRoutes(rs.CliCtx, rs.Mux)
	app.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)

	// start websockets server
	websocketAddr := viper.GetString(flagWebsocket)
	ws := newWebsocketsServer(rs.CliCtx, websocketAddr)
	ws.start()
}
