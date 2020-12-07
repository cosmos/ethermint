package types

// Constants to match up protocol versions and messages
const (
	eth63 = 63
	eth64 = 64
	eth65 = 65
)

// ProtocolVersions are the supported versions of the eth protocol (first is primary).
var ProtocolVersions = []uint{eth65, eth64, eth63}
