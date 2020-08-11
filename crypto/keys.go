package crypto

import (
	"fmt"

	tmcrypto "github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
)

const (
	// EthSecp256k1 defines the ECDSA secp256k1 used on Ethereum
	EthSecp256k1 = keyring.SigningAlgo("eth_secp256k1")
)

// SupportedAlgorithms defines the list of signing algorithms used on Ethermint:
//  - eth_secp256k1 (Ethereum)
//  - secp256k1 (Tendermint)
var SupportedAlgorithms = []keyring.SigningAlgo{EthSecp256k1, keyring.Secp256k1}

// EthermintKeygenFunc is the key generation function to generate secp256k1 ToECDSA
// from ethereum.
func EthermintKeygenFunc(bz []byte, algo keyring.SigningAlgo) (tmcrypto.PrivKey, error) {
	if algo != EthSecp256k1 {
		return nil, fmt.Errorf("signing algorithm must be %s, got %s", EthSecp256k1, algo)
	}

	return PrivKeySecp256k1(bz), nil
}
