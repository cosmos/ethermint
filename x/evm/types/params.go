package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// DefaultParamspace for params keeper
	DefaultParamspace = ModuleName

	homesteadBlock      = "HomesteadBlock"
	daoForkBlock        = "DAOForkBlock"
	eip150Block         = "EIP150Block"
	eip150Hash          = "EIP150Hash"
	eip155Block         = "EIP155Block"
	eip158Block         = "EIP158Block"
	byzantiumBlock      = "ByzantiumBlock"
	constantinopleBlock = "ConstantinopleBlock"
	petersburgBlock     = "PetersburgBlock"
	istanbulBlock       = "IstanbulBlock"
	muirGlacierBlock    = "MuirGlacierBlock"
	yoloV1Block         = "YoloV1Block"
	eWASMBlock          = "EWASMBlock"
)

// Parameter keys
var (
	ParamStoreKeyHomesteadBlock = []byte(homesteadBlock)

	ParamStoreKeyDAOForkBlock   = []byte(daoForkBlock)
	ParamStoreKeyDAOForkSupport = []byte("DAOForkSupport")

	ParamStoreKeyEIP150Block = []byte(eip150Block)
	ParamStoreKeyEIP150Hash  = []byte(eip150Hash)

	ParamStoreKeyEIP155Block = []byte(eip155Block)
	ParamStoreKeyEIP158Block = []byte(eip158Block)

	ParamStoreKeyByzantiumBlock      = []byte(byzantiumBlock)
	ParamStoreKeyConstantinopleBlock = []byte(constantinopleBlock)
	ParamStoreKeyPetersburgBlock     = []byte(petersburgBlock)
	ParamStoreKeyIstanbulBlock       = []byte(istanbulBlock)
	ParamStoreKeyMuirGlacierBlock    = []byte(muirGlacierBlock)

	ParamStoreKeyYoloV1Block = []byte(yoloV1Block)
	ParamStoreKeyEWASMBlock  = []byte(eWASMBlock)
)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// Params defines the Ethereum ChainConfig as SDK parameters
type Params struct {
	HomesteadBlock sdk.Int `json:"homesteadBlock,omitempty"` // Homestead switch block (nil = no fork, 0 = already homestead)

	DAOForkBlock   sdk.Int `json:"daoForkBlock,omitempty"`   // TheDAO hard-fork switch block (nil = no fork)
	DAOForkSupport bool    `json:"daoForkSupport,omitempty"` // Whether the nodes supports or opposes the DAO hard-fork

	// EIP150 implements the Gas price changes (https://github.com/ethereum/EIPs/issues/150)
	EIP150Block sdk.Int `json:"eip150Block,omitempty"` // EIP150 HF block (nil = no fork)
	EIP150Hash  string  `json:"eip150Hash,omitempty"`  // EIP150 HF hash (needed for header only clients as only gas pricing changed)

	EIP155Block sdk.Int `json:"eip155Block,omitempty"` // EIP155 HF block
	EIP158Block sdk.Int `json:"eip158Block,omitempty"` // EIP158 HF block

	ByzantiumBlock      sdk.Int `json:"byzantiumBlock,omitempty"`      // Byzantium switch block (nil = no fork, 0 = already on byzantium)
	ConstantinopleBlock sdk.Int `json:"constantinopleBlock,omitempty"` // Constantinople switch block (nil = no fork, 0 = already activated)
	PetersburgBlock     sdk.Int `json:"petersburgBlock,omitempty"`     // Petersburg switch block (nil = same as Constantinople)
	IstanbulBlock       sdk.Int `json:"istanbulBlock,omitempty"`       // Istanbul switch block (nil = no fork, 0 = already on istanbul)
	MuirGlacierBlock    sdk.Int `json:"muirGlacierBlock,omitempty"`    // Eip-2384 (bomb delay) switch block (nil = no fork, 0 = already activated)

	YoloV1Block sdk.Int `json:"yoloV1Block,omitempty"` // YOLO v1: https://github.com/ethereum/EIPs/pull/2657 (Ephemeral testnet)
	EWASMBlock  sdk.Int `json:"ewasmBlock,omitempty"`  // EWASM switch block (nil = no fork, 0 = already activated)
}

// DefaultParams returns default evm parameters
func DefaultParams() Params {
	return Params{
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
		IstanbulBlock:       sdk.Int{},
		MuirGlacierBlock:    sdk.Int{},
		YoloV1Block:         sdk.Int{},
		EWASMBlock:          sdk.Int{},
	}
}

func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyHomesteadBlock, &p.HomesteadBlock, validateInt),
		paramtypes.NewParamSetPair(ParamStoreKeyDAOForkBlock, &p.DAOForkBlock, validateInt),
		paramtypes.NewParamSetPair(ParamStoreKeyDAOForkSupport, &p.DAOForkSupport, validateDAOForkSupport),
		paramtypes.NewParamSetPair(ParamStoreKeyEIP150Block, &p.EIP150Block, validateInt),
		paramtypes.NewParamSetPair(ParamStoreKeyEIP150Hash, &p.EIP150Hash, validateHash),
		paramtypes.NewParamSetPair(ParamStoreKeyEIP155Block, &p.EIP155Block, validateInt),
		paramtypes.NewParamSetPair(ParamStoreKeyEIP158Block, &p.EIP158Block, validateInt),
		paramtypes.NewParamSetPair(ParamStoreKeyByzantiumBlock, &p.ByzantiumBlock, validateInt),
		paramtypes.NewParamSetPair(ParamStoreKeyConstantinopleBlock, &p.ConstantinopleBlock, validateInt),
		paramtypes.NewParamSetPair(ParamStoreKeyPetersburgBlock, &p.PetersburgBlock, validateInt),
		paramtypes.NewParamSetPair(ParamStoreKeyIstanbulBlock, &p.IstanbulBlock, validateInt),
		paramtypes.NewParamSetPair(ParamStoreKeyMuirGlacierBlock, &p.MuirGlacierBlock, validateInt),
		paramtypes.NewParamSetPair(ParamStoreKeyYoloV1Block, &p.YoloV1Block, validateInt),
		paramtypes.NewParamSetPair(ParamStoreKeyEWASMBlock, &p.EWASMBlock, validateInt),
	}
}

// ValidateBasic performs basic validation on evm parameters.
func (p Params) ValidateBasic() error {
	if err := validateInt(p.HomesteadBlock); err != nil {
		return fmt.Errorf("%s: %w", homesteadBlock, err)
	}
	if err := validateInt(p.DAOForkBlock); err != nil {
		return fmt.Errorf("%s: %w", daoForkBlock, err)
	}
	if err := validateDAOForkSupport(p.DAOForkSupport); err != nil {
		return err
	}
	if err := validateInt(p.EIP150Block); err != nil {
		return fmt.Errorf("%s: %w", eip150Block, err)
	}
	if err := validateHash(p.EIP150Hash); err != nil {
		return err
	}
	if err := validateInt(p.EIP155Block); err != nil {
		return fmt.Errorf("%s: %w", eip155Block, err)
	}
	if err := validateInt(p.EIP158Block); err != nil {
		return fmt.Errorf("%s: %w", eip158Block, err)
	}
	if err := validateInt(p.ByzantiumBlock); err != nil {
		return fmt.Errorf("%s: %w", byzantiumBlock, err)
	}
	if err := validateInt(p.ConstantinopleBlock); err != nil {
		return fmt.Errorf("%s: %w", constantinopleBlock, err)
	}
	if err := validateInt(p.PetersburgBlock); err != nil {
		return fmt.Errorf("%s: %w", petersburgBlock, err)
	}
	if err := validateInt(p.IstanbulBlock); err != nil {
		return fmt.Errorf("%s: %w", istanbulBlock, err)
	}
	if err := validateInt(p.MuirGlacierBlock); err != nil {
		return fmt.Errorf("%s: %w", muirGlacierBlock, err)
	}
	if err := validateInt(p.YoloV1Block); err != nil {
		return fmt.Errorf("%s: %w", yoloV1Block, err)
	}
	if err := validateInt(p.EWASMBlock); err != nil {
		return fmt.Errorf("%s: %w", eWASMBlock, err)
	}
	return nil
}

func validateInt(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.BigInt() == nil {
		return nil
	}

	if v.IsNegative() {
		return fmt.Errorf("parameter value cannot be negative: %s", v)
	}

	return nil
}

func validateDAOForkSupport(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateHash(i interface{}) error {
	hex, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	bz := common.FromHex(hex)
	if len(bz) != common.HashLength {
		return fmt.Errorf("invalid hash length, expected 32, got %d", len(bz))
	}

	return nil
}
