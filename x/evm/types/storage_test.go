package types

import (
	"testing"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestStorageValidate(t *testing.T) {
	testCases := []struct {
		name    string
		storage Storage
		expPass bool
	}{
		{
			"valid storage",
			Storage{
				NewState(ethcmn.BytesToHash([]byte{1, 2, 3}), ethcmn.BytesToHash([]byte{1, 2, 3})),
			},
			true,
		},
		{
			"empty storage key bytes",
			Storage{
				{Key: ethcmn.Hash{}},
			},
			false,
		},
		{
			"duplicated storage key",
			Storage{
				{Key: ethcmn.BytesToHash([]byte{1, 2, 3})},
				{Key: ethcmn.BytesToHash([]byte{1, 2, 3})},
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		err := tc.storage.Validate()
		if tc.expPass {
			require.NoError(t, err, tc.name)
		} else {
			require.Error(t, err, tc.name)
		}
	}
}
