package types

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"
)

var (
	regexChainID     = `[a-z]*`
	regexSeparator   = `-{1}`
	regexEpoch       = `[1-9][0-9]*`
	ethermintChainID = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, regexChainID, regexSeparator, regexEpoch))
)

// ValidChainID returns false if the given chain identifier is incorrectly formatted.
var ValidChainID = ethermintChainID.MatchString

// ParseChainID parses a string chain identifier to an Ethereum-compatible
// chain-id in *big.Int format.
func ParseChainID(chainID string) (*big.Int, error) {
	chainID = strings.TrimSpace(chainID)

	matches := ethermintChainID.FindStringSubmatch(chainID)
	if matches == nil {
		return nil, fmt.Errorf("invalid chain-id: %s,", chainID)
	}

	// verify that the chain-id entered is a base 10 integer
	chainIDInt, ok := new(big.Int).SetString(matches[2], 10)
	if !ok {
		return nil, fmt.Errorf("invalid chain-id epoch: %s, must be base-10 integer format", matches[2])
	}

	return chainIDInt, nil
}
