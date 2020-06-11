// Package rpc contains RPC handler methods and utilities to start
// Ethermint's Web3-compatibly JSON-RPC server.
package rpc

import (
	"github.com/cosmos/cosmos-sdk/client"
	emintcrypto "github.com/cosmos/ethermint/crypto"
	"github.com/ethereum/go-ethereum/rpc"
)

const Web3Namespace = "web3"
const EthNamespace = "eth"
const PersonalNamespace = "personal"
const NetNamespace = "net"

// GetRPCAPIs returns the list of all APIs
func GetRPCAPIs(clientCtx client.Context, key emintcrypto.PrivKeySecp256k1) []rpc.API {
	nonceLock := new(AddrLocker)
	backend := NewEthermintBackend(clientCtx)
	return []rpc.API{
		{
			Namespace: Web3Namespace,
			Version:   "1.0",
			Service:   NewPublicWeb3API(),
			Public:    true,
		},
		{
			Namespace: EthNamespace,
			Version:   "1.0",
			Service:   NewPublicEthAPI(clientCtx, backend, nonceLock, key),
			Public:    true,
		},
		{
			Namespace: PersonalNamespace,
			Version:   "1.0",
			Service:   NewPersonalEthAPI(clientCtx, nonceLock),
			Public:    false,
		},
		{
			Namespace: EthNamespace,
			Version:   "1.0",
			Service:   NewPublicFilterAPI(clientCtx, backend),
			Public:    true,
		},
		{
			Namespace: NetNamespace,
			Version:   "1.0",
			Service:   NewPublicNetAPI(clientCtx),
			Public:    true,
		},
	}
}
