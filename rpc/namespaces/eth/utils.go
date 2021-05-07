package eth

import (
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

func getStorageByAddressKey(addr common.Address, key []byte) common.Hash {
	prefix := addr.Bytes()
	compositeKey := make([]byte, len(prefix)+len(key))

	copy(compositeKey, prefix)
	copy(compositeKey[len(prefix):], key)

	return ethcrypto.Keccak256Hash(compositeKey)
}
