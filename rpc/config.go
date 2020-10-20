package rpc

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/cosmos/ethermint/server/config"

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

	ws := newWebsocketsServer(clientCtx, apiConfig.Address, apiConfig.WebsocketAddress)
	ws.start()
}

func (s *websocketsServer) start() {
	ws := mux.NewRouter()
	ws.Handle("/", s)

	errCh := make(chan error)
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%s", s.wsAddr), ws)
		if err != nil {
			errCh <- fmt.Errorf("failed to serve: %w", err)
		}
	}()

	select {
	case err := <-errCh:
		return nil, err
	case <-time.After(5 * time.Second): // assume server started successfully
		return grpcSrv, nil
	}
}
