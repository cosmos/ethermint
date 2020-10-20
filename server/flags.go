package server

// Tendermint full-node start flags
const (
	flagWithTendermint = "with-tendermint"
	flagAddress        = "address"
	flagTransport      = "transport"
	flagTraceStore     = "trace-store"
	flagCPUProfile     = "cpu-profile"
)

// GRPC-related flags.
const (
	flagGRPCEnable  = "grpc.enable"
	flagGRPCAddress = "grpc.address"
)

// Ethereum-related flags.
const (
	flagEthereumJSONRPCEnable    = "ethereum.enable-json-rpc"
	flagEthereumWebsocketEnable  = "ethereum.enable-ethereum-websocket"
	flagEthereumWebsocketAddress = "ethereum.websocket-address"
)
