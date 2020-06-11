package rpc

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/cosmos/cosmos-sdk/client"
)

// PersonalEthAPI is the eth_ prefixed set of APIs in the Web3 JSON-RPC spec.
type PersonalEthAPI struct {
	clientCtx client.Context
	nonceLock *AddrLocker
}

// NewPersonalEthAPI creates an instance of the public ETH Web3 API.
func NewPersonalEthAPI(clientCtx client.Context, nonceLock *AddrLocker) *PersonalEthAPI {
	return &PersonalEthAPI{
		clientCtx: clientCtx,
		nonceLock: nonceLock,
	}
}

// Sign calculates an Ethereum ECDSA signature for:
// keccack256("\x19Ethereum Signed Message:\n" + len(message) + message))
//
// Note, the produced signature conforms to the secp256k1 curve R, S and V values,
// where the V value will be 27 or 28 for legacy reasons.
//
// The key used to calculate the signature is decrypted with the given password.
//
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_sign
func (e *PersonalEthAPI) Sign(ctx context.Context, data hexutil.Bytes, addr common.Address, passwd string) (hexutil.Bytes, error) {
	return nil, nil
}
