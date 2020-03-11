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
func EncodeResultData(addr ethcmn.Address, bloom ethtypes.Bloom, logs []*ethtypes.Log, evmRet []byte) []byte {
	// Append address, bloom, evm return bytes in that order
	returnData := append(addr.Bytes(), bloom.Bytes()...)
	returnData = append(returnData, byte(len(logs)))
	return append(returnData, evmRet...)
}

// DecodeReturnData decodes the byte slice of values to their respective types
func DecodeResultData(bytes []byte) (addr ethcmn.Address, bloom ethtypes.Bloom, ret []byte, err error) {
	if len(bytes) >= returnIdx {
		addr = ethcmn.BytesToAddress(bytes[:bloomIdx])
		bloom = ethtypes.BytesToBloom(bytes[bloomIdx : bloomIdx+ethtypes.BloomByteLength])
		ret = bytes[returnIdx:]
	} else {
		err = fmt.Errorf("Invalid format for encoded data, message must be an EVM state transition")
	}

	return
}

// func decodeResultData(r io.Reader) (ethcmn.Address, ethtypes.Bloom, []*ethtypes.Log, []byte, error) {
// 	addr := make([]byte, ethcmn.AddressLength)
// 	n, err := r.Read(addr)
// 	if err != nil {
// 		return nil, nil, nil, nil, err
// 	} else if n != ethcmn.AddressLength {
// 		return nil, nil, nil, nil, fmt.Errorf("could not read address")
// 	}

// 	bloom := make([]byte, ethtypes.BloomByteLength)
// 	n, err = r.Read(bloom)
// 	if err != nil {
// 		return nil, nil, nil, nil, err
// 	} else if n != ethcmn.BloomByteLength {
// 		return nil, nil, nil, nil, fmt.Errorf("could not read bloom")
// 	}

// 	logLen := make([]byte, 1)
// 	n, err = r.Read(logLen)
// 	if err != nil {
// 		return nil, nil, nil, nil, err
// 	}

// 	topics := [][]byte{}

// 	for i := range logLen[0] {
// 		topic := make([]byte, 32)
// 		n, err = r.Read(topic)
// 		if err != nil {
// 			return nil, nil, nil, nil, err
// 		}

// 		topics = append(topics, topic)
// 	}

// 	evmRet := []byte{}
// 	for {
// 		buf := make([]byte, 1)
// 		n, err = r.Read(topic)
// 		if err != nil || n == 0 {
// 			break
// 		}
// 		evmRet = append(evmRet, buf[0])
// 	}

// 	return
// }
