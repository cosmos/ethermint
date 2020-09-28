package ethsecp256k1

import (
	"bytes"
	"crypto/ecdsa"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	ethsecp256k1 "github.com/ethereum/go-ethereum/crypto/secp256k1"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	tmcrypto "github.com/tendermint/tendermint/crypto"
)

const (
	PrivKeySize = 32
	keyType     = "eth_secp256k1"
	PrivKeyName = "ethermint/PrivKeyEthSecp256k1"
	PubKeyName  = "ethermint/PubKeyEthSecp256k1"
)


// ----------------------------------------------------------------------------
// secp256k1 Private Key

var _ cryptotypes.PrivKey = &PrivKey{}
var _ codec.AminoMarshaler = &PrivKey{}

// PrivKey defines a type alias for an ecdsa.PrivateKey that implements
// Tendermint's PrivateKey interface.
type PrivKey []byte

// GenerateKey generates a new random private key. It returns an error upon
// failure.
func GenerateKey() (*PrivKey, error) {
	priv, err := ethcrypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	return &PrivKey{
			Key: ethcrypto.FromECDSA(priv),
		}, nil
}

// Bytes returns the byte representation of the ECDSA Private Key.
func (privKey *PrivKey) Bytes() []byte {
	return privKey.Key
}

// PubKey returns the ECDSA private key's public key.
func (privkey PrivKey) PubKey() tmcrypto.PubKey {
	ecdsaPrivKey := privkey.ToECDSA()
	return &PubKey{
		Key: ethcrypto.CompressPubkey(&ecdsaPrivKey.PublicKey),
	}
}

// Equals returns true if two ECDSA private keys are equal and false otherwise.
func (privKey *PrivKey) Equals(other crypto.PrivKey) bool {
	return privKey.Type() == other.Type() && subtle.ConstantTimeCompare(privKey.Bytes(), other.Bytes()) == 1
}

// Type returns eth_secp256k1
func (privKey *PrivKey) Type() string {
	return keyType
}

// MarshalAmino overrides Amino binary marshalling.
func (privKey PrivKey) MarshalAmino() ([]byte, error) {
	return privKey.Key, nil
}

// UnmarshalAmino overrides Amino binary marshalling.
func (privKey *PrivKey) UnmarshalAmino(bz []byte) error {
	if len(bz) != PrivKeySize {
		return fmt.Errorf("invalid privkey size, expected %d got %d", PrivKeySize, len(bz))
	}
	privKey.Key = bz

	return nil
}

// MarshalAminoJSON overrides Amino JSON marshalling.
func (privKey PrivKey) MarshalAminoJSON() ([]byte, error) {
	// When we marshal to Amino JSON, we don't marshal the "key" field itself,
	// just its contents (i.e. the key bytes).
	return privKey.MarshalAmino()
}

// UnmarshalAminoJSON overrides Amino JSON marshalling.
func (privKey *PrivKey) UnmarshalAminoJSON(bz []byte) error {
	return privKey.UnmarshalAmino(bz)
}

// Sign creates a recoverable ECDSA signature on the secp256k1 curve over the
// Keccak256 hash of the provided message. The produced signature is 65 bytes
// where the last byte contains the recovery ID.
func (privkey PrivKey) Sign(msg []byte) ([]byte, error) {
	return ethcrypto.Sign(ethcrypto.Keccak256Hash(msg).Bytes(), privkey.ToECDSA())
}


// ToECDSA returns the ECDSA private key as a reference to ecdsa.PrivateKey type.
// The function will panic if the private key is invalid.
func (privkey PrivKey) ToECDSA() *ecdsa.PrivateKey {
	key, err := ethcrypto.ToECDSA(privkey.Bytes())
	if err != nil {
		panic(err)
	}
	return key
}

// ----------------------------------------------------------------------------
// secp256k1 Public Key

var _ cryptotypes.PubKey = &PubKey{}
var _ codec.AminoMarshaler = &PubKey{}

// PubKeySize is comprised of 32 bytes for one field element
// (the x-coordinate), plus one byte for the parity of the y-coordinate.
const PubKeySize = 33

// PubKey defines a type alias for an ecdsa.PublicKey that implements Tendermint's PubKey
// interface. It represents the 33-byte compressed public key format.
type PubKey []byte

// Address returns the address of the ECDSA public key.
// The function will panic if the public key is invalid.
func (pubKey PubKey) Address() tmcrypto.Address {
	if len(pubKey.Key) != PubKeySize {
		panic("length of pubkey is incorrect")
	}

	pubk, err := ethcrypto.DecompressPubkey(pubKey.Key)
	if err != nil {
		panic(err)
	}

	return tmcrypto.Address(ethcrypto.PubkeyToAddress(*pubk).Bytes())
}

// Bytes returns the raw bytes of the ECDSA public key.
func (pubKey PubKey) Bytes() []byte {
	return pubKey.Key
}

// String implements the fmt.Stringer interface.
func (pubKey *PubKey) String() string {
	return fmt.Sprintf("EthPubKeySecp256k1{%X}", pubKey.Key)
}

// Type returns eth_secp256k1
func (pubKey *PubKey) Type() string {
	return keyType
}

// Equals returns true if the pubkey type is the same and their bytes are deeply equal.
func (pubKey *PubKey) Equals(other tmcrypto.PubKey) bool {
	return pubKey.Type() == other.Type() && bytes.Equal(pubKey.Bytes(), other.Bytes())
}

// MarshalAmino overrides Amino binary marshalling.
func (pubKey PubKey) MarshalAmino() ([]byte, error) {
	return pubKey.Key, nil
}

// UnmarshalAmino overrides Amino binary marshalling.
func (pubKey *PubKey) UnmarshalAmino(bz []byte) error {
	if len(bz) != PubKeySize {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, "invalid pubkey size, expected %d, got %d", PubKeySize, len(bz))
	}
	pubKey.Key = bz

	return nil
}

// MarshalAminoJSON overrides Amino JSON marshalling.
func (pubKey PubKey) MarshalAminoJSON() ([]byte, error) {
	// When we marshal to Amino JSON, we don't marshal the "key" field itself,
	// just its contents (i.e. the key bytes).
	return pubKey.MarshalAmino()
}

// UnmarshalAminoJSON overrides Amino JSON marshalling.
func (pubKey *PubKey) UnmarshalAminoJSON(bz []byte) error {
	return pubKey.UnmarshalAmino(bz)
}


// VerifyBytes verifies that the ECDSA public key created a given signature over
// the provided message. It will calculate the Keccak256 hash of the message
// prior to verification.
func (pubKeykey PubKey) VerifyBytes(msg []byte, sig []byte) bool {
	if len(sig) == 65 {
		// remove recovery ID if contained in the signature
		sig = sig[:len(sig)-1]
	}

	// the signature needs to be in [R || S] format when provided to VerifySignature
	return ethsecp256k1.VerifySignature(pubKey.Key, ethcrypto.Keccak256Hash(msg).Bytes(), sig)
}