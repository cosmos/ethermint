package rpc

import (
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/cobra"
	"log"
)

// Config contains configuration fields that determine the
// behavior of the RPC HTTP server.
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

func Web3RpcCmd(cdc *codec.Codec) *cobra.Command {
	return lcd.ServeCommand(cdc, registerRoutes)
}

func registerRoutes(rs *lcd.RestServer) {
	s := rpc.NewServer()
	apis := GetRPCAPIs()

	// TODO: Read modules from config and whitelist
	//modules := cliCtx.Modules
	//whitelist := make(map[string]bool)
	//for _, module := range modules {
	//	whitelist[module] = true
	//}

	// Register all the APIs exposed by the services
	for _, api := range apis {
		//if whitelist[api.Namespace] || (len(whitelist) == 0 && api.Public) {
		if err := s.RegisterName(api.Namespace, api.Service); err != nil {
			log.Println(err)
			return
		}
		//}
	}

	rs.Mux.HandleFunc("/rpc", s.ServeHTTP).Methods("POST")
}
