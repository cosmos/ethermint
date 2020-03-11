package types

import (
	"fmt"

	"github.com/cosmos/ethermint/crypto"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"
)

const (
	bloomIdx  = ethcmn.AddressLength
	returnIdx = bloomIdx + ethtypes.BloomByteLength + 1
)

// GenerateEthAddress generates an Ethereum address.
func GenerateEthAddress() ethcmn.Address {
	priv, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}

	return ethcrypto.PubkeyToAddress(priv.ToECDSA().PublicKey)
}

// ValidateSigner attempts to validate a signer for a given slice of bytes over
// which a signature and signer is given. An error is returned if address
// derived from the signature and bytes signed does not match the given signer.
func ValidateSigner(signBytes, sig []byte, signer ethcmn.Address) error {
	pk, err := ethcrypto.SigToPub(signBytes, sig)

	if err != nil {
		return errors.Wrap(err, "failed to derive public key from signature")
	} else if ethcrypto.PubkeyToAddress(*pk) != signer {
		return fmt.Errorf("invalid signature for signer: %s", signer)
	}

	return nil
}

func rlpHash(x interface{}) (hash ethcmn.Hash) {
	hasher := sha3.NewLegacyKeccak256()
	//nolint:gosec,errcheck
	rlp.Encode(hasher, x)
	hasher.Sum(hash[:0])

	return hash
}

type ResultData struct {
	addr  ethcmn.Address
	bloom ethtypes.Bloom
	logs  []*ethtypes.Log
	ret   []byte
}

// EncodeReturnData takes all of the necessary data from the EVM execution
// and returns the data as a byte slice
func EncodeResultData(data *ResultData) ([]byte, error) {
	return rlp.EncodeToBytes(data)
}

// DecodeReturnData decodes the byte slice of values to their respective types
func DecodeResultData(in []byte) (*ResultData, error) {
	data := new(ResultData)
	err := rlp.DecodeBytes(in, data)
	return data, err
}
