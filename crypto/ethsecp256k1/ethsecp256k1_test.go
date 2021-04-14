package ethsecp256k1

import (
	"testing"

	"github.com/stretchr/testify/require"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	ethsecp256k1 "github.com/ethereum/go-ethereum/crypto/secp256k1"

	tmcrypto "github.com/tendermint/tendermint/crypto"
)

func TestPrivKeyPrivKey(t *testing.T) {
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
	expectedSig, _ := ethsecp256k1.Sign(sigHash.Bytes(), privKey)

	sig, err := privKey.Sign(msg)
	require.NoError(t, err)
	require.Equal(t, expectedSig, sig)
}

func TestPrivKeyPubKey(t *testing.T) {
	privKey, err := GenerateKey()
	require.NoError(t, err)

	// validate type and equality
	pubKey := privKey.PubKey().(PubKey)
	require.Implements(t, (*tmcrypto.PubKey)(nil), pubKey)

	// validate inequality
	privKey2, err := GenerateKey()
	require.NoError(t, err)
	require.False(t, pubKey.Equals(privKey2.PubKey()))

	// validate signature
	msg := []byte("hello world")
	sig, err := privKey.Sign(msg)
	require.NoError(t, err)

	res := pubKey.VerifyBytes(msg, sig)
	require.True(t, res)
}
