package types

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseChainID(t *testing.T) {
	testCases := []struct {
		name     string
		chainID  string
		expError bool
		expInt   *big.Int
	}{
		{
			"valid chain-id, single digit", "ethermint-1", false, big.NewInt(1),
		},
		{
			"valid chain-id, multiple digits", "aragonchain-256", false, big.NewInt(256),
		},
		{
			"invalid chain-id, double dash", "aragon-chain-1", true, nil,
		},
		{
			"invalid chain-id, uppercases", "ETHERMINT-1", true, nil,
		},
		{
			"invalid chain-id, mixed cases", "Ethermint-1", true, nil,
		},
		{
			"invalid chain-id, special chars", "$&*#!-1", true, nil,
		},
		{
			"invalid epoch, cannot start with 0", "ethermint-001", true, nil,
		},
		{
			"invalid epoch, cannot invalid base", "ethermint-0x212", true, nil,
		},
		{
			"invalid epoch, non-integer", "ethermint-ethermint", true, nil,
		},
		{
			"blank chain ID", " ", true, nil,
		},
		{
			"empty chain ID", "", true, nil,
		},
	}

	for _, tc := range testCases {
		intChainID, err := ParseChainID(tc.chainID)
		if tc.expError {
			require.Error(t, err, tc.name)
			require.Nil(t, intChainID)
		} else {
			require.NoError(t, err, tc.name)
			require.Equal(t, tc.expInt, intChainID, tc.name)
		}
	}
}
