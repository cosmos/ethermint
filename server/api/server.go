package api

import (
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server/api"

	// unnamed import of statik for swagger UI support
	_ "github.com/cosmos/cosmos-sdk/client/docs/statik"

	"github.com/cosmos/ethermint/server/config"
)

// Server defines the server's API interface.
type Server struct {
	*api.Server
	// TODO: define
	// WebsocketServer
	// JSONRPCServer
}

// New creates a new Server instance.
func New(clientCtx client.Context, logger log.Logger) *Server {
	return &Server{
		Server: api.New(clientCtx, logger),
	}
}

// Start starts the API server. Internally, the API server leverages Tendermint's
// JSON RPC server. Configuration options are provided via config.APIConfig
// and are delegated to the Tendermint JSON RPC server. The process is
// non-blocking, so an external signal handler must be used.
func (s *Server) Start(cfg config.Config) error {
	if err := s.Server.Start(*cfg.Config); err != nil {
		return err
	}

	// TODO: start rpc servers
	return nil
}
