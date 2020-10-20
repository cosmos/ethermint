package config

import (
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/server/config"
)

const (
	// DefaultEthereumWebsocketAddress is the default address the Ethereum websocket server binds to.
	DefaultEthereumWebsocketAddress = "tcp://0.0.0.0:8546"
)

// Config defines the server's top level configuration
type Config struct {
	*config.Config
	Ethereum EthereumConfig `mapstructure:"ethereum"`
}

// EthereumConfig defines the Ethereum API listener configuration.
type EthereumConfig struct {
	// EnableJSONRPC defines if the JSON-RPC server should be enabled.
	EnableJSONRPC bool `mapstructure:"enable-json-rpc"`

	// EnableWebsocket defines if the Ethereum websocker server should be enabled.
	EnableWebsocket bool `mapstructure:"enable-ethereum-websocket"`

	// Address defines the Websocket server address to listen on
	WebsocketAddress string `mapstructure:"websocket-address"`
}

// DefaultConfig returns server's default configuration.
func DefaultConfig() *Config {
	return &Config{
		Config: config.DefaultConfig(),
		Ethereum: EthereumConfig{
			EnableJSONRPC:    false,
			EnableWebsocket:  false,
			WebsocketAddress: DefaultEthereumWebsocketAddress,
		},
	}
}

// GetConfig returns a fully parsed Config object.
func GetConfig(v *viper.Viper) Config {
	sdkConfig := config.GetConfig(v)
	return Config{
		Config: &sdkConfig,
		Ethereum: EthereumConfig{
			EnableJSONRPC:    v.GetBool("ethereum.enable-json-rpc"),
			EnableWebsocket:  v.GetBool("ethereum.enable-websocket"),
			WebsocketAddress: v.GetString("ethereum.websocket-address"),
		},
	}
}
