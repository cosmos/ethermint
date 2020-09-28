package hd

import (
	"fmt"

	"github.com/pkg/errors"

	"crypto/hmac"
	"crypto/sha512"

	"github.com/tyler-smith/go-bip39"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	ethaccounts "github.com/ethereum/go-ethereum/accounts"

	tmcrypto "github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
)

const (
	// EthSecp256k1 defines the ECDSA secp256k1 curve used on Ethereum
	EthSecp256k1Type = hd.PubKeyType("eth_secp256k1")
)

var (
	_ keyring.SignatureAlgo = EthSecp256k1{}

	// EthSecp256k1 defines the signature algorithm type used on Ethermint
	EthSecp256k1 = ethSecp256k1Algo{}

	// SupportedAlgorithms defines the list of signing algorithms used on Ethermint:
	//  - eth_secp256k1 (Ethereum)
	//  - secp256k1 (Tendermint)
	SupportedAlgorithms = keyring.SigningAlgoList{EthSecp256k1, hd.Secp256k1}

	// SupportedAlgorithmsLedger defines the list of signing algorithms used on Ethermint with a Ledger device:
	//  - eth_secp256k1 (Ethereum)
	//  - secp256k1 (Tendermint)
	SupportedAlgorithmsLedger = keyring.SigningAlgoList{EthSecp256k1, hd.Secp256k1}
)


// EthSecp256k1Options defines a keys options for the ethereum Secp256k1 curve.
func EthSecp256k1Options(options *keyring.Options) {
	options.SupportedAlgos = SupportedAlgorithms
	options.WithSupportedAlgosLedger = SupportedAlgorithmsLedger
}

// Name returns eth_secp256k1
func (s ethSecp256k1Algo) Name() hd.PubKeyType {
	return EthSecp256k1Type
}

// Derive derives and returns the eth_secp256k1 private key for the given seed and HD path.
func (s ethSecp256k1Algo) Derive() hd.DeriveFn {
	return func(mnemonic string, bip39Passphrase, path string) ([]byte, error) {
		hdpath, err := ethaccounts.ParseDerivationPath(hdpath)
		if err != nil {
			return nil, err 
		}

		seed, err := bip39.NewSeedWithErrorChecking(mnemonic, bip39Passphrase)
		if err != nil {
			return nil, err
		}

		// HMAC the seed to produce the private key and chain code
		mac := hmac.New(sha512.New, []byte("Bitcoin seed"))
		_, err = mac.Write(seed)
		if err != nil {
			return nil, err
		}

		seed = mac.Sum(nil)

		masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
		if err != nil {
			return nil, err
		}

		for _, n := range hdpath {
			masterKey, err = masterKey.Child(n)
			if err != nil {
				return nil, err
			}
		}

		privateKey, err := key.ECPrivKey()
		if err != nil {
			return nil, err
		}

		privateKeyECDSA := privateKey.ToECDSA()
		derivedKey := ethcrypto.FromECDSA(privateKeyECDSA)
		return derivedKey, nil
	}
}

// Generate generates a eth_secp256k1 private key from the given bytes.
func (ethSecp256k1Algo) Generate() hd.GenerateFn {
	return func(bz []byte) tmcrypto.PrivKey {
		var bzArr = make([]byte, ethsecp256k1.PrivKeySize)
		copy(bzArr, bz)

		return &ethsecp256k1.PrivKey{Key: bzArr}
	}
}