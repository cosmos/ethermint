package rpc

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"

	emintcrypto "github.com/cosmos/ethermint/crypto"

	"github.com/spf13/viper"
)

const (
	flagUnlockKey = "unlock-key"
)

// Config contains configuration fields that determine the behavior of the RPC HTTP server.
// TODO: These may become irrelevant if HTTP config is handled by the SDK
type Config struct {
	// EnableRPC defines whether or not to enable the RPC server
	EnableRPC bool
	// RPCAddr defines the IP address to listen on
	RPCAddr string
	// RPCPort defines the port to listen on
	RPCPort int
	// RPCCORSDomains defines list of domains to enable CORS headers for (used by browsers)
	RPCCORSDomains []string
	// RPCVhosts defines list of domains to listen on (useful if Tendermint is addressable via DNS)
	RPCVHosts []string
}

func UnlockKeyFromNameAndPassphrase(accountName, passphrase string) (emintcrypto.PrivKeySecp256k1, error) {
	keystore, err := keyring.New(
		sdk.KeyringServiceName(),
		viper.GetString(flags.FlagKeyringBackend),
		viper.GetString(flags.FlagHome),
		os.Stdin,
		emintcrypto.EthSeckp256k1Option,
	)
	if err != nil {
		return nil, err
	}

	// With keyring key store, password is not required as it is pulled from the OS prompt
	// Exports private key from keyring using password
	armored, err := keystore.ExportPrivKeyArmor(accountName, passphrase)
	if err != nil {
		return nil, err
	}

	privKey, algo, err := crypto.UnarmorDecryptPrivKey(armored, passphrase)
	if err != nil {
		return nil, err
	}

	if algo != string(emintcrypto.EthSecp256k1Type) {
		panic(fmt.Sprintf("invalid Keyring algorithm for JSON-RPC key, expected '%s' got '%s'", string(emintcrypto.EthSecp256k1Type), algo))
	}

	ethSecp256k1, ok := privKey.(emintcrypto.PrivKeySecp256k1)
	if !ok {
		panic(fmt.Sprintf("invalid private key type: %T", privKey))
	}

	return ethSecp256k1, nil
}
