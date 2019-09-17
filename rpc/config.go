package rpc

import (
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/cobra"

	"log"
)

const (
	flagUnlockKey = "unlock-key"
)

// Config contains configuration fields that determine the behavior of the RPC HTTP server.
// TODO: These may become irrelevant if HTTP config is handled by the SDK
type Config struct {
	// EnableRPC defines whether or not to enable the RPC server
	EnableRPC bool
	// RPCAddr defines the IP address to listen on
	RPCAddr string
	// RPCPort defines the port to listen on
	RPCPort int
	// RPCCORSDomains defines list of domains to enable CORS headers for (used by browsers)
	RPCCORSDomains []string
	// RPCVhosts defines list of domains to listen on (useful if Tendermint is addressable via DNS)
	RPCVHosts []string
}

// Web3RpcCmd creates a CLI command to start RPC server
func Web3RpcCmd(cdc *codec.Codec) *cobra.Command {
	cmd := lcd.ServeCommand(cdc, registerRoutes)
	// Attach flag to cmd output to be handled in registerRoutes
	cmd.Flags().String(flagUnlockKey, "", "Select a key to unlock on the RPC server")
	return cmd
}

// registerRoutes creates a new server and registers the `/rpc` endpoint.
// Rpc calls are enabled based on their associated module (eg. "eth").
func registerRoutes(rs *lcd.RestServer) {
	s := rpc.NewServer()
	apis := GetRPCAPIs(rs.CliCtx)

	// TODO: Allow cli to configure modules https://github.com/ChainSafe/ethermint/issues/74
	whitelist := make(map[string]bool)

	// Register all the APIs exposed by the services
	for _, api := range apis {
		if whitelist[api.Namespace] || (len(whitelist) == 0 && api.Public) {
			if err := s.RegisterName(api.Namespace, api.Service); err != nil {
				log.Println(err)
				return
			}
		}
	}

	rs.Mux.HandleFunc("/rpc", s.ServeHTTP).Methods("POST")
}
