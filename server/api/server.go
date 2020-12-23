package api

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/client"
	_ "github.com/cosmos/cosmos-sdk/client/docs/statik"
	"github.com/cosmos/cosmos-sdk/server/api"
	sdkconfig "github.com/cosmos/cosmos-sdk/server/config"

	gethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/cosmos/ethermint/rpc/websockets"
	"github.com/cosmos/ethermint/server/config"
)

var _ api.Server = &Server{}

// Server defines the server's API interface.
type Server struct {
	Base      *api.BaseServer
	JSONRPC   *gethrpc.Server
	Websocket *websockets.Server
}

// New creates a new Server instance.
func New(clientCtx client.Context, logger log.Logger, cfg *config.Config) *Server {
	baseSvr := api.New(clientCtx, logger)
	return &Server{
		Base:      baseSvr,
		JSONRPC:   gethrpc.NewServer(),
		Websocket: websockets.NewServer(clientCtx, logger, baseSvr.Router, cfg.API.Address, cfg.Ethereum.WebsocketAddress),
	}
}

// BaseServer implements the api.Server interface
func (s *Server) BaseServer() *api.BaseServer { return s.Base }

// Start starts the API server. Internally, the API server leverages Tendermint's
// JSON RPC server. Configuration options are provided via config.APIConfig
// and are delegated to the Tendermint JSON RPC server. The process is
// non-blocking, so an external signal handler must be used.
func (s *Server) Start(cfg sdkconfig.ServerConfig) error {
	ethermintCfg, ok := cfg.(*config.Config)
	if !ok {
		return fmt.Errorf("invalid config type, expected %T got %T", &config.Config{}, cfg)
	}

	if err := s.Base.Start(cfg.GetSDKConfig()); err != nil {
		return err
	}

	// NOTE: JSON-RPC APIs are already registered by the application
	if ethermintCfg.Ethereum.EnableJSONRPC {
		// Web3 RPC API route
		s.Base.Router.HandleFunc("/", s.JSONRPC.ServeHTTP).Methods("POST", "OPTIONS")
	}

	if ethermintCfg.Ethereum.EnableWebsocket {
		if err := s.Websocket.Start(); err != nil {
			return err
		}
	}

	return nil
}

// Close closes all the connections with the servers
func (s *Server) Close() error {
	if err := s.Base.Close(); err != nil {
		return err
	}

	s.JSONRPC.Stop()

	return nil
}
