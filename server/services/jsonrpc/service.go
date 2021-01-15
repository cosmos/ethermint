package jsonrpc

import (
	"fmt"
	"net"
	"time"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/cosmos/ethermint/server/config"
)

type Service struct {
	rpcServer *rpc.Server
	apis      []rpc.API
}

// NewService creates a new gRPC server instance with a defined listener address.
func NewService(apis []rpc.API) *Service {
	return &Service{
		rpcServer: rpc.NewServer(),
		apis:      apis,
	}
}

// Name returns the JSON-RPC service name
func (Service) Name() string {
	return "JSON-RPC"
}

// RegisterRoutes registers the JSON-RPC server to the application. It fails if any of the
// API names fail to register.
func (s *Service) RegisterRoutes() error {
	for _, api := range s.apis {
		if err := s.rpcServer.RegisterName(api.Namespace, api.Service); err != nil {
			return err
		}
	}

	return nil
}

// Start starts the JSON-RPC server on the address defined on the configuration.
func (s *Service) Start(cfg config.JSONRPCConfig) error {
	if !cfg.Enable {
		return nil
	}

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}

	errCh := make(chan error)
	go func() {
		err = s.rpcServer.ServeListener(listener)
		if err != nil {
			errCh <- fmt.Errorf("failed to serve: %w", err)
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-time.After(5 * time.Second): // assume server started successfully
		return nil
	}
}

// Stop stops the JSON-RPC service by no longer reading new requests, waits for
// stopPendingRequestTimeout to allow pending requests to finish, then closes all codecs which will
// cancel pending requests and subscriptions.
func (s *Service) Stop() error {
	s.rpcServer.Stop()
	return nil
}
