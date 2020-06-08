package crypto

import (
	"crypto/hmac"
	"crypto/sha512"

	"github.com/tyler-smith/go-bip39"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	tmcrypto "github.com/tendermint/tendermint/crypto"
)

var _ keyring.SignatureAlgo = ethSecp256k1{}

// EthSeckp256k1Option defines a keyring option for the ethereum Secp256k1 curve.
func EthSeckp256k1Option(options *keyring.Options) {
	options.SupportedAlgos = append(options.SupportedAlgos, Secp256k1)
	options.SupportedAlgosLedger = append(options.SupportedAlgosLedger, Secp256k1)
}

// Secp256k1 represents the Secp256k1 curve used in Ethereum.
var Secp256k1 = ethSecp256k1{}

type ethSecp256k1 struct{}

func (s ethSecp256k1) Name() hd.PubKeyType {
	return hd.PubKeyType("ethsecp256k1")
}

// Derive derives and returns the secp256k1 private key for the given seed and HD path.
func (s ethSecp256k1) Derive() hd.DeriveFn {
	return func(mnemonic string, bip39Passphrase, hdPath string) ([]byte, error) {
		seed, err := bip39.NewSeedWithErrorChecking(mnemonic, bip39Passphrase)
		if err != nil {
			return nil, err
		}

		// HMAC the seed to produce the private key and chain code
		mac := hmac.New(sha512.New, []byte("Bitcoin seed"))
		mac.Write(seed)
		seed = mac.Sum(nil)

		priv, err := ethcrypto.ToECDSA(seed[:32])
		if err != nil {
			return nil, err
		}

		derivedKey := PrivKeySecp256k1(ethcrypto.FromECDSA(priv))

		return derivedKey, nil
	}
}

func (ethSecp256k1) Generate() hd.GenerateFn {
	return func(bz []byte) tmcrypto.PrivKey {
		var bzArr []byte
		copy(bzArr[:], bz)
		return PrivKeySecp256k1(bzArr)
	}
}
