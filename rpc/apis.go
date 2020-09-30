// Package rpc contains RPC handler methods and utilities to start
// Ethermint's Web3-compatibly JSON-RPC server.
package rpc

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	emintcrypto "github.com/cosmos/ethermint/crypto"
	"github.com/ethereum/go-ethereum/rpc"
)

// RPC namespaces and API version
const (
	Web3Namespace     = "web3"
	EthNamespace      = "eth"
	PersonalNamespace = "personal"
	NetNamespace      = "net"

	apiVersion = "1.0"
)

// GetRPCAPIs returns the list of all APIs
func GetRPCAPIs(cliCtx context.CLIContext, keys []emintcrypto.PrivKeySecp256k1) []rpc.API {
	nonceLock := new(AddrLocker)
	backend := NewEthermintBackend(cliCtx)
	ethAPI := NewPublicEthAPI(cliCtx, backend, nonceLock, keys)

	return []rpc.API{
		{
			Namespace: Web3Namespace,
			Version:   apiVersion,
			Service:   NewPublicWeb3API(),
			Public:    true,
		},
		{
			Namespace: EthNamespace,
			Version:   apiVersion,
			Service:   ethAPI,
			Public:    true,
		},
		{
			Namespace: PersonalNamespace,
			Version:   apiVersion,
			Service:   NewPersonalEthAPI(ethAPI),
			Public:    false,
		},
		{
			Namespace: EthNamespace,
			Version:   apiVersion,
			Service:   NewPublicFilterAPI(cliCtx, backend),
			Public:    true,
		},
		{
			Namespace: NetNamespace,
			Version:   apiVersion,
			Service:   NewPublicNetAPI(cliCtx),
			Public:    true,
		},
	}
}
