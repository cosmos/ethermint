package crypto

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	tmcrypto "github.com/tendermint/tendermint/crypto"
)

// EthermintKeygenFunc is the key generation function to generate secp256k1 ToECDSA
// from ethereum.
func EthermintKeygenFunc(bz []byte, algo keys.SigningAlgo) (tmcrypto.PrivKey, error) {
	return PrivKeySecp256k1(bz), nil
}
