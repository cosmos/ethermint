package types

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	ethcmn "github.com/ethereum/go-ethereum/common"
)

func TestValidateGenesisAccount(t *testing.T) {
	testCases := []struct {
		name           string
		genesisAccount GenesisAccount
		expPass        bool
	}{
		{
			"empty account address bytes",
			GenesisAccount{
				Address: ethcmn.Address{},
				Balance: big.NewInt(1),
			},
			false,
		},
		{
			"nil account balance",
			GenesisAccount{
				Address: ethcmn.BytesToAddress([]byte{1, 2, 3, 4, 5}),
				Balance: nil,
			},
			false,
		},
		{
			"empty code bytes",
			GenesisAccount{
				Address: ethcmn.BytesToAddress([]byte{1, 2, 3, 4, 5}),
				Balance: big.NewInt(1),
				Code:    []byte{},
			},
			false,
		},
		{
			"empty storage key bytes",
			GenesisAccount{
				Address: ethcmn.BytesToAddress([]byte{1, 2, 3, 4, 5}),
				Balance: big.NewInt(1),
				Code:    []byte{1, 2, 3},
				Storage: []GenesisStorage{
					{Key: ethcmn.Hash{}},
				},
			},
			false,
		},
		{
			"duplicated storage key",
			GenesisAccount{
				Address: ethcmn.BytesToAddress([]byte{1, 2, 3, 4, 5}),
				Balance: big.NewInt(1),
				Code:    []byte{1, 2, 3},
				Storage: []GenesisStorage{
					{Key: ethcmn.Hash{1, 2}},
					{Key: ethcmn.Hash{1, 2}},
				},
			},
			false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		err := tc.genesisAccount.Validate()
		if tc.expPass {
			require.NoError(t, err, tc.name)
		} else {
			require.Error(t, err, tc.name)
		}
	}
}

func TestValidateGenesis(t *testing.T) {

	testCases := []struct {
		name     string
		genState GenesisState
		expPass  bool
	}{
		{
			name:     "default",
			genState: DefaultGenesisState(),
			expPass:  true,
		},
		{
			name: "empty account address bytes",
			genState: GenesisState{
				Accounts: []GenesisAccount{
					{
						Address: ethcmn.Address{},
						Balance: big.NewInt(1),
					},
				},
			},
			expPass: false,
		},
		{
			name: "nil account balance",
			genState: GenesisState{
				Accounts: []GenesisAccount{
					{
						Address: ethcmn.BytesToAddress([]byte{1, 2, 3, 4, 5}),
						Balance: nil,
					},
				},
			},
			expPass: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		err := tc.genState.Validate()
		if tc.expPass {
			require.NoError(t, err, tc.name)
		} else {
			require.Error(t, err, tc.name)
		}
	}
}
