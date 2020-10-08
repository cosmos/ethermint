package ethsecp256k1

import (
	"testing"

	"github.com/stretchr/testify/require"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"

	tmcrypto "github.com/tendermint/tendermint/crypto"
)

func TestPrivKey(t *testing.T) {
	// validate type and equality
	privKey, err := GenerateKey()
	require.NoError(t, err)
	require.True(t, privKey.Equals(privKey))
	require.Implements(t, (*tmcrypto.PrivKey)(nil), privKey)

	// validate inequality
	privKey2, err := GenerateKey()
	require.NoError(t, err)
	require.False(t, privKey.Equals(privKey2))

	// validate Ethereum address equality
	addr := privKey.PubKey().Address()
	expectedAddr := ethcrypto.PubkeyToAddress(privKey.ToECDSA().PublicKey)
	require.Equal(t, expectedAddr.Bytes(), addr.Bytes())

	// validate we can sign some bytes
	msg := []byte("hello world")
	sigHash := ethcrypto.Keccak256Hash(msg)
	expectedSig, err := secp256k1.Sign(sigHash.Bytes(), privKey.Bytes())
	require.NoError(t, err)

	sig, err := privKey.Sign(msg)
	require.NoError(t, err)
	require.Equal(t, expectedSig, sig)
}

func TestPrivKey_PubKey(t *testing.T) {
	privKey, err := GenerateKey()
	require.NoError(t, err)

	// validate type and equality
	pubKey := &PubKey{
		Key: privKey.PubKey().Bytes(),
	}
	require.Implements(t, (*tmcrypto.PubKey)(nil), pubKey)

	// validate inequality
	privKey2, err := GenerateKey()
	require.NoError(t, err)
	require.False(t, pubKey.Equals(privKey2.PubKey()))

	// validate signature
	msg := []byte("hello world")
	sig, err := privKey.Sign(msg)
	require.NoError(t, err)

	res := pubKey.VerifySignature(msg, sig)
	require.True(t, res)
}
