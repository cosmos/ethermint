// Package rpc contains RPC handler methods and utilities to start
// Ethermint's Web3-compatibly JSON-RPC server.
package rpc

import (
	"github.com/ethereum/go-ethereum/rpc"
)

// GetRPCAPIs returns the master list of public APIs for use with
// StartHTTPEndpoint.
func GetRPCAPIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "web3",
			Version:   "1.0",
			Service:   NewPublicWeb3API(),
		},
		{
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicEthAPI(),
		},
	}
}
