package codec

import (
	tmcrypto "github.com/tendermint/tendermint/crypto"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/cosmos/ethermint/crypto/ethsecp256k1"
)

// RegisterInterfaces registers the sdk.Tx interface.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterInterface("tendermint.crypto.PubKey", (*tmcrypto.PubKey)(nil))
	registry.RegisterImplementations((*tmcrypto.PubKey)(nil), &ed25519.PubKey{})
	registry.RegisterImplementations((*tmcrypto.PubKey)(nil), &secp256k1.PubKey{})
	registry.RegisterImplementations((*tmcrypto.PubKey)(nil), &multisig.LegacyAminoPubKey{})
	registry.RegisterImplementations((*tmcrypto.PubKey)(nil), &ethsecp256k1.PubKey{})

	registry.RegisterInterface("cosmos.crypto.PubKey", (*cryptotypes.PubKey)(nil))
	registry.RegisterImplementations((*cryptotypes.PubKey)(nil), &ed25519.PubKey{})
	registry.RegisterImplementations((*cryptotypes.PubKey)(nil), &secp256k1.PubKey{})
	registry.RegisterImplementations((*cryptotypes.PubKey)(nil), &multisig.LegacyAminoPubKey{})
	registry.RegisterImplementations((*cryptotypes.PubKey)(nil), &ethsecp256k1.PubKey{})
}
