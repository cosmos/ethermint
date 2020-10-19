package server

import (
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/server/config"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
)

const (
	defaultMinGasPrices = ""
)

// Config defines the server's top level configuration
type Config struct {
	config.BaseConfig `mapstructure:",squash"`

	// Telemetry defines the application telemetry configuration
	Telemetry telemetry.Config       `mapstructure:"telemetry"`
	API       APIConfig              `mapstructure:"api"`
	GRPC      config.GRPCConfig      `mapstructure:"grpc"`
	StateSync config.StateSyncConfig `mapstructure:"state-sync"`
}

// APIConfig defines the API listener configuration.
type APIConfig struct {
	// Enable defines if the API server should be enabled.
	Enable bool `mapstructure:"enable"`

	// Swagger defines if swagger documentation should automatically be registered.
	Swagger bool `mapstructure:"swagger"`

	// EnableUnsafeCORS defines if CORS should be enabled (unsafe - use it at your own risk)
	EnableUnsafeCORS bool `mapstructure:"enabled-unsafe-cors"`

	// Address defines the API server to listen on
	Address string `mapstructure:"address"`

	// MaxOpenConnections defines the number of maximum open connections
	MaxOpenConnections uint `mapstructure:"max-open-connections"`

	// RPCReadTimeout defines the Tendermint RPC read timeout (in seconds)
	RPCReadTimeout uint `mapstructure:"rpc-read-timeout"`

	// RPCWriteTimeout defines the Tendermint RPC write timeout (in seconds)
	RPCWriteTimeout uint `mapstructure:"rpc-write-timeout"`

	// RPCMaxBodyBytes defines the Tendermint maximum response body (in bytes)
	RPCMaxBodyBytes uint `mapstructure:"rpc-max-body-bytes"`

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
		BaseConfig: config.BaseConfig{
			MinGasPrices:      defaultMinGasPrices,
			InterBlockCache:   true,
			Pruning:           storetypes.PruningOptionDefault,
			PruningKeepRecent: "0",
			PruningKeepEvery:  "0",
			PruningInterval:   "0",
			MinRetainBlocks:   0,
			IndexEvents:       make([]string, 0),
		},
		Telemetry: telemetry.Config{
			Enabled:      false,
			GlobalLabels: [][]string{},
		},
		API: APIConfig{
			Enable:             false,
			Swagger:            false,
			EnableUnsafeCORS:   false,
			Address:            "tcp://0.0.0.0:1317",
			MaxOpenConnections: 1000,
			RPCReadTimeout:     10,
			RPCWriteTimeout:    0,
			RPCMaxBodyBytes:    1000000,
			EnableJSONRPC:      false,
			EnableWebsocket:    false,
			WebsocketAddress:   "tcp://0.0.0.0:8546",
		},
		GRPC: config.GRPCConfig{
			Enable:  true,
			Address: config.DefaultGRPCAddress,
		},
		StateSync: config.StateSyncConfig{
			SnapshotInterval:   0,
			SnapshotKeepRecent: 2,
		},
	}
}

// GetConfig returns a fully parsed Config object.
func GetConfig(v *viper.Viper) Config {
	globalLabelsRaw := v.Get("telemetry.global-labels").([]interface{})
	globalLabels := make([][]string, 0, len(globalLabelsRaw))
	for _, glr := range globalLabelsRaw {
		labelsRaw := glr.([]interface{})
		if len(labelsRaw) == 2 {
			globalLabels = append(globalLabels, []string{labelsRaw[0].(string), labelsRaw[1].(string)})
		}
	}

	return Config{
		BaseConfig: config.BaseConfig{
			MinGasPrices:      v.GetString("minimum-gas-prices"),
			InterBlockCache:   v.GetBool("inter-block-cache"),
			Pruning:           v.GetString("pruning"),
			PruningKeepRecent: v.GetString("pruning-keep-recent"),
			PruningKeepEvery:  v.GetString("pruning-keep-every"),
			PruningInterval:   v.GetString("pruning-interval"),
			HaltHeight:        v.GetUint64("halt-height"),
			HaltTime:          v.GetUint64("halt-time"),
			IndexEvents:       v.GetStringSlice("index-events"),
			MinRetainBlocks:   v.GetUint64("min-retain-blocks"),
		},
		Telemetry: telemetry.Config{
			ServiceName:             v.GetString("telemetry.service-name"),
			Enabled:                 v.GetBool("telemetry.enabled"),
			EnableHostname:          v.GetBool("telemetry.enable-hostname"),
			EnableHostnameLabel:     v.GetBool("telemetry.enable-hostname-label"),
			EnableServiceLabel:      v.GetBool("telemetry.enable-service-label"),
			PrometheusRetentionTime: v.GetInt64("telemetry.prometheus-retention-time"),
			GlobalLabels:            globalLabels,
		},
		API: APIConfig{
			Enable:             v.GetBool("api.enable"),
			Swagger:            v.GetBool("api.swagger"),
			Address:            v.GetString("api.address"),
			MaxOpenConnections: v.GetUint("api.max-open-connections"),
			RPCReadTimeout:     v.GetUint("api.rpc-read-timeout"),
			RPCWriteTimeout:    v.GetUint("api.rpc-write-timeout"),
			RPCMaxBodyBytes:    v.GetUint("api.rpc-max-body-bytes"),
			EnableUnsafeCORS:   v.GetBool("api.enabled-unsafe-cors"),
			EnableJSONRPC:      v.GetBool("api.enable-json-rpc"),
			EnableWebsocket:    v.GetBool("api.enable-websocket"),
			WebsocketAddress:   v.GetString("api.websocket-address"),
		},
		GRPC: config.GRPCConfig{
			Enable:  v.GetBool("grpc.enable"),
			Address: v.GetString("grpc.address"),
		},
		StateSync: config.StateSyncConfig{
			SnapshotInterval:   v.GetUint64("state-sync.snapshot-interval"),
			SnapshotKeepRecent: v.GetUint32("state-sync.snapshot-keep-recent"),
		},
	}
}
