package types

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"gopkg.in/yaml.v2"
)

// ChainConfig defines the Ethereum ChainConfig parameters
type ChainConfig struct {
	HomesteadBlock sdk.Int `json:"homestead_block"` // Homestead switch block (nil = no fork, 0 = already homestead)

	DAOForkBlock   sdk.Int `json:"dao_fork_block"`   // TheDAO hard-fork switch block (nil = no fork)
	DAOForkSupport bool    `json:"dao_fork_support"` // Whether the nodes supports or opposes the DAO hard-fork

	// EIP150 implements the Gas price changes (https://github.com/ethereum/EIPs/issues/150)
	EIP150Block sdk.Int `json:"eip150_block"` // EIP150 HF block (nil = no fork)
	EIP150Hash  string  `json:"eip150_hash"`  // EIP150 HF hash (needed for header only clients as only gas pricing changed)

	EIP155Block sdk.Int `json:"eip155_block"` // EIP155 HF block
	EIP158Block sdk.Int `json:"eip158_block"` // EIP158 HF block

	ByzantiumBlock      sdk.Int `json:"byzantium_block"`      // Byzantium switch block (nil = no fork, 0 = already on byzantium)
	ConstantinopleBlock sdk.Int `json:"constantinople_block"` // Constantinople switch block (nil = no fork, 0 = already activated)
	PetersburgBlock     sdk.Int `json:"petersburg_block"`     // Petersburg switch block (nil = same as Constantinople)
	IstanbulBlock       sdk.Int `json:"istanbul_block"`       // Istanbul switch block (nil = no fork, 0 = already on istanbul)
	MuirGlacierBlock    sdk.Int `json:"muir_glacier_block"`   // Eip-2384 (bomb delay) switch block (nil = no fork, 0 = already activated)

	YoloV1Block sdk.Int `json:"yoloV1_block"` // YOLO v1: https://github.com/ethereum/EIPs/pull/2657 (Ephemeral testnet)
	EWASMBlock  sdk.Int `json:"ewasm_block"`  // EWASM switch block (nil = no fork, 0 = already activated)
}

// EthereumConfig returns an Ethereum ChainConfig for EVM state transitions
func (cc ChainConfig) EthereumConfig(chainID *big.Int) *params.ChainConfig {
	return &params.ChainConfig{
		ChainID:             chainID,
		HomesteadBlock:      getBlockValue(cc.HomesteadBlock),
		DAOForkBlock:        getBlockValue(cc.DAOForkBlock),
		DAOForkSupport:      cc.DAOForkSupport,
		EIP150Block:         getBlockValue(cc.EIP150Block),
		EIP150Hash:          common.HexToHash(cc.EIP150Hash),
		EIP155Block:         getBlockValue(cc.EIP155Block),
		EIP158Block:         getBlockValue(cc.EIP158Block),
		ByzantiumBlock:      getBlockValue(cc.ByzantiumBlock),
		ConstantinopleBlock: getBlockValue(cc.ConstantinopleBlock),
		PetersburgBlock:     getBlockValue(cc.PetersburgBlock),
		IstanbulBlock:       getBlockValue(cc.IstanbulBlock),
		MuirGlacierBlock:    getBlockValue(cc.MuirGlacierBlock),
		YoloV1Block:         getBlockValue(cc.YoloV1Block),
		EWASMBlock:          getBlockValue(cc.EWASMBlock),
	}
}

// String implements the fmt.Stringer interface
func (cc ChainConfig) String() string {
	out, _ := yaml.Marshal(cc)
	return string(out)
}

// DefaultChainConfig returns default evm parameters
func DefaultChainConfig() ChainConfig {
	return ChainConfig{
		HomesteadBlock:      sdk.ZeroInt(),
		DAOForkBlock:        sdk.ZeroInt(),
		DAOForkSupport:      true,
		EIP150Block:         sdk.ZeroInt(),
		EIP150Hash:          "0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0",
		EIP155Block:         sdk.ZeroInt(),
		EIP158Block:         sdk.ZeroInt(),
		ByzantiumBlock:      sdk.ZeroInt(),
		ConstantinopleBlock: sdk.ZeroInt(),
		PetersburgBlock:     sdk.ZeroInt(),
		IstanbulBlock:       sdk.NewInt(-1),
		MuirGlacierBlock:    sdk.NewInt(-1),
		YoloV1Block:         sdk.NewInt(-1),
		EWASMBlock:          sdk.NewInt(-1),
	}
}

func getBlockValue(block sdk.Int) *big.Int {
	if block.BigInt() == nil || block.IsNegative() {
		return nil
	}

	return block.BigInt()
}

func (cc ChainConfig) Validate() error {
	if err := validateHash(cc.EIP150Hash); err != nil {
		return err
	}

	return nil
}


func validateHash(hex string) error {
	bz := common.FromHex(hex)
	lenHex := len(bz)
	if lenHex > 0 && lenHex != common.HashLength {
		return fmt.Errorf("invalid hash length, expected %d, got %d", common.HashLength, lenHex)
	}

	return nil
}