package codec

import (
	tmcrypto "github.com/tendermint/tendermint/crypto"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	
	"github.com/cosmos/ethermint/crypto/ethsecp256k1"
)

// RegisterInterfaces registers the sdk.Tx interface.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*tmcrypto.PubKey)(nil), &ethsecp256k1.PubKey{})
	registry.RegisterImplementations((*cryptotypes.PubKey)(nil), &ethsecp256k1.PubKey{})
}
