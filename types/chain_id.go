package types

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	regexChainID     = `[a-z]*`
	regexSeparator   = `-{1}`
	regexEpoch       = `[1-9][0-9]*`
	ethermintChainID = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, regexChainID, regexSeparator, regexEpoch))
)

// IsValidChainID returns false if the given chain identifier is incorrectly formatted.
var IsValidChainID = ethermintChainID.MatchString

// ParseChainID parses a string chain identifier's epoch to an Ethereum-compatible
// chain-id in *big.Int format. The function returns an error if the chain-id has an invalid format
func ParseChainID(chainID string) (*big.Int, error) {
	chainID = strings.TrimSpace(chainID)

	matches := ethermintChainID.FindStringSubmatch(chainID)
	if matches == nil {
		return nil, sdkerrors.Wrap(ErrInvalidChainID, chainID)
	}

	// verify that the chain-id entered is a base 10 integer
	chainIDInt, ok := new(big.Int).SetString(matches[2], 10)
	if !ok {
		return nil, sdkerrors.Wrapf(ErrInvalidChainID, "epoch %s must be base-10 integer format", matches[2])
	}

	return chainIDInt, nil
}
