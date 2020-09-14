package types

import (
	"fmt"
	"math/big"
)

// ParseChainID parses a string chain identifier to an Ethereum-compatible
// chain-id in *big.Int format.
func ParseChainID(chainID string) (*big.Int, error) {
	// verify that the chain-id entered is a base 10 integer
	chainIDInt, ok := new(big.Int).SetString(chainID, 10)
	if !ok {
		return nil, fmt.Errorf("invalid chainID: %s, must be base-10 integer format", chainID)
	}

	return chainIDInt, nil
}
