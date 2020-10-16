package types

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/cosmos/ethermint/crypto/ethsecp256k1"
)

func TestValidateGenesisAccount(t *testing.T) {
	testCases := []struct {
		name           string
		genesisAccount GenesisAccount
		expPass        bool
	}{
		{
			"valid genesis account",
			GenesisAccount{
				Address: ethcmn.BytesToAddress([]byte{1, 2, 3, 4, 5}).String(),
				Balance: big.NewInt(1).Bytes(),
				Code:    []byte{1, 2, 3},
				Storage: Storage{
					NewState(ethcmn.BytesToHash([]byte{1, 2, 3}), ethcmn.BytesToHash([]byte{1, 2, 3})),
				},
			},
			true,
		},
		{
			"empty account address bytes",
			GenesisAccount{
				Address: ethcmn.Address{}.String(),
				Balance: big.NewInt(1).Bytes(),
			},
			false,
		},
		{
			"nil account balance",
			GenesisAccount{
				Address: ethcmn.BytesToAddress([]byte{1, 2, 3, 4, 5}).String(),
				Balance: nil,
			},
			false,
		},
		{
			"nil account balance",
			GenesisAccount{
				Address: ethcmn.BytesToAddress([]byte{1, 2, 3, 4, 5}).String(),
				Balance: big.NewInt(-1),
			},
			false,
		},
		{
			"empty code bytes",
			GenesisAccount{
				Address: ethcmn.BytesToAddress([]byte{1, 2, 3, 4, 5}).String(),
				Balance: big.NewInt(1),
				Code:    []byte{},
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
	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	addr := ethcrypto.PubkeyToAddress(priv.ToECDSA().PublicKey)

	testCases := []struct {
		name     string
		genState *GenesisState
		expPass  bool
	}{
		{
			name:     "default",
			genState: DefaultGenesisState(),
			expPass:  true,
		},
		{
			name: "valid genesis",
			genState: &GenesisState{
				Accounts: []GenesisAccount{
					{
						Address: ethcmn.BytesToAddress([]byte{1, 2, 3, 4, 5}).String(),
						Balance: big.NewInt(1),
						Code:    []byte{1, 2, 3},
						Storage: Storage{
							{Key: ethcmn.BytesToHash([]byte{1, 2, 3}).String()},
						},
					},
				},
				TxsLogs: []TransactionLogs{
					{
						Hash: ethcmn.BytesToHash([]byte("tx_hash")).String(),
						Logs: []*ethtypes.Log{
							{
								Address:     addr,
								Topics:      []ethcmn.Hash{ethcmn.BytesToHash([]byte("topic"))},
								Data:        []byte("data"),
								BlockNumber: 1,
								TxHash:      ethcmn.BytesToHash([]byte("tx_hash")),
								TxIndex:     1,
								BlockHash:   ethcmn.BytesToHash([]byte("block_hash")),
								Index:       1,
								Removed:     false,
							},
						},
					},
				},
				ChainConfig: DefaultChainConfig(),
				Params:      DefaultParams(),
			},
			expPass: true,
		},
		{
			name:     "empty genesis",
			genState: &GenesisState{},
			expPass:  false,
		},
		{
			name: "invalid genesis",
			genState: &GenesisState{
				Accounts: []GenesisAccount{
					{
						Address: ethcmn.Address{}.String(),
					},
				},
			},
			expPass: false,
		},
		{
			name: "duplicated genesis account",
			genState: &GenesisState{
				Accounts: []GenesisAccount{
					{
						Address: ethcmn.BytesToAddress([]byte{1, 2, 3, 4, 5}).String(),
						Balance: big.NewInt(1),
						Code:    []byte{1, 2, 3},
						Storage: Storage{
							NewState(ethcmn.BytesToHash([]byte{1, 2, 3}), ethcmn.BytesToHash([]byte{1, 2, 3})),
						},
					},
					{
						Address: ethcmn.BytesToAddress([]byte{1, 2, 3, 4, 5}).String(),
						Balance: big.NewInt(1),
						Code:    []byte{1, 2, 3},
						Storage: Storage{
							NewState(ethcmn.BytesToHash([]byte{1, 2, 3}), ethcmn.BytesToHash([]byte{1, 2, 3})),
						},
					},
				},
			},
			expPass: false,
		},
		{
			name: "duplicated tx log",
			genState: &GenesisState{
				Accounts: []GenesisAccount{
					{
						Address: ethcmn.BytesToAddress([]byte{1, 2, 3, 4, 5}),
						Balance: big.NewInt(1),
						Code:    []byte{1, 2, 3},
						Storage: Storage{
							{Key: ethcmn.BytesToHash([]byte{1, 2, 3}).String()},
						},
					},
				},
				TxsLogs: []TransactionLogs{
					{
						Hash: ethcmn.BytesToHash([]byte("tx_hash")).String(),
						Logs: []*ethtypes.Log{
							{
								Address:     addr,
								Topics:      []ethcmn.Hash{ethcmn.BytesToHash([]byte("topic"))},
								Data:        []byte("data"),
								BlockNumber: 1,
								TxHash:      ethcmn.BytesToHash([]byte("tx_hash")),
								TxIndex:     1,
								BlockHash:   ethcmn.BytesToHash([]byte("block_hash")),
								Index:       1,
								Removed:     false,
							},
						},
					},
					{
						Hash: ethcmn.BytesToHash([]byte("tx_hash")).String(),
						Logs: []*ethtypes.Log{
							{
								Address:     addr,
								Topics:      []ethcmn.Hash{ethcmn.BytesToHash([]byte("topic"))},
								Data:        []byte("data"),
								BlockNumber: 1,
								TxHash:      ethcmn.BytesToHash([]byte("tx_hash")),
								TxIndex:     1,
								BlockHash:   ethcmn.BytesToHash([]byte("block_hash")),
								Index:       1,
								Removed:     false,
							},
						},
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid tx log",
			genState: &GenesisState{
				Accounts: []GenesisAccount{
					{
						Address: ethcmn.BytesToAddress([]byte{1, 2, 3, 4, 5}).String(),
						Balance: big.NewInt(1),
						Code:    []byte{1, 2, 3},
						Storage: Storage{
							{Key: ethcmn.BytesToHash([]byte{1, 2, 3}).String()},
						},
					},
				},
				TxsLogs: []TransactionLogs{NewTransactionLogs(ethcmn.Hash{}, nil)},
			},
			expPass: false,
		},
		{
			name: "invalid params",
			genState: &GenesisState{
				ChainConfig: DefaultChainConfig(),
				Params:      Params{},
			},
			expPass: false,
		},
		{
			name: "invalid chain config",
			genState: &GenesisState{
				ChainConfig: &ChainConfig{},
				Params:      DefaultParams(),
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
