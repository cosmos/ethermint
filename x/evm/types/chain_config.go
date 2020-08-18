package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

// GenerateChainConfig returns an Ethereum chainconfig for EVM state transitions
func GenerateChainConfig(chainID *big.Int, parameters Params) *params.ChainConfig {
	return &params.ChainConfig{
		ChainID:             chainID,
		HomesteadBlock:      parameters.HomesteadBlock.BigInt(),
		DAOForkBlock:        parameters.DAOForkBlock.BigInt(),
		DAOForkSupport:      parameters.DAOForkSupport,
		EIP150Block:         parameters.EIP150Block.BigInt(),
		EIP150Hash:          common.HexToHash(parameters.EIP150Hash),
		EIP155Block:         parameters.EIP155Block.BigInt(),
		EIP158Block:         parameters.EIP158Block.BigInt(),
		ByzantiumBlock:      parameters.ByzantiumBlock.BigInt(),
		ConstantinopleBlock: parameters.ConstantinopleBlock.BigInt(),
		PetersburgBlock:     parameters.PetersburgBlock.BigInt(),
		IstanbulBlock:       parameters.IstanbulBlock.BigInt(),
		MuirGlacierBlock:    parameters.MuirGlacierBlock.BigInt(),
		YoloV1Block:         parameters.YoloV1Block.BigInt(),
		EWASMBlock:          nil,
	}
}
