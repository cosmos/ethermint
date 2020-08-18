package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestParamsValidate(t *testing.T) {
	testCases := []struct {
		name     string
		params   Params
		expError bool
	}{
		{"default", DefaultParams(), false},
		{
			"valid",
			Params{
				HomesteadBlock:      sdk.OneInt(),
				DAOForkBlock:        sdk.OneInt(),
				EIP150Block:         sdk.OneInt(),
				EIP150Hash:          defaultEIP150Hash,
				EIP155Block:         sdk.OneInt(),
				EIP158Block:         sdk.OneInt(),
				ByzantiumBlock:      sdk.OneInt(),
				ConstantinopleBlock: sdk.OneInt(),
				PetersburgBlock:     sdk.OneInt(),
				IstanbulBlock:       sdk.OneInt(),
				MuirGlacierBlock:    sdk.OneInt(),
				YoloV1Block:         sdk.OneInt(),
				EWASMBlock:          sdk.OneInt(),
			},
			false,
		},
		{
			"empty",
			Params{},
			false,
		},
		{
			"invalid HomesteadBlock",
			Params{
				HomesteadBlock: sdk.NewInt(-1),
			},
			true,
		},
		{
			"invalid DAOForkBlock",
			Params{
				HomesteadBlock: sdk.OneInt(),
				DAOForkBlock:   sdk.NewInt(-1),
			},
			true,
		},
		{
			"invalid EIP150Block",
			Params{
				HomesteadBlock: sdk.OneInt(),
				DAOForkBlock:   sdk.OneInt(),
				EIP150Block:    sdk.NewInt(-1),
			},
			true,
		},
		{
			"invalid EIP155Block",
			Params{
				HomesteadBlock: sdk.OneInt(),
				DAOForkBlock:   sdk.OneInt(),
				EIP150Block:    sdk.OneInt(),
				EIP150Hash:     defaultEIP150Hash,
				EIP155Block:    sdk.NewInt(-1),
			},
			true,
		},
		{
			"invalid EIP158Block",
			Params{
				HomesteadBlock: sdk.OneInt(),
				DAOForkBlock:   sdk.OneInt(),
				EIP150Block:    sdk.OneInt(),
				EIP150Hash:     defaultEIP150Hash,
				EIP155Block:    sdk.OneInt(),
				EIP158Block:    sdk.NewInt(-1),
			},
			true,
		},
		{
			"invalid ByzantiumBlock",
			Params{
				HomesteadBlock: sdk.OneInt(),
				DAOForkBlock:   sdk.OneInt(),
				EIP150Block:    sdk.OneInt(),
				EIP150Hash:     defaultEIP150Hash,
				EIP155Block:    sdk.OneInt(),
				EIP158Block:    sdk.OneInt(),
				ByzantiumBlock: sdk.NewInt(-1),
			},
			true,
		},
		{
			"invalid ConstantinopleBlock",
			Params{
				HomesteadBlock:      sdk.OneInt(),
				DAOForkBlock:        sdk.OneInt(),
				EIP150Block:         sdk.OneInt(),
				EIP150Hash:          defaultEIP150Hash,
				EIP155Block:         sdk.OneInt(),
				EIP158Block:         sdk.OneInt(),
				ByzantiumBlock:      sdk.OneInt(),
				ConstantinopleBlock: sdk.NewInt(-1),
			},
			true,
		},
		{
			"invalid PetersburgBlock",
			Params{
				HomesteadBlock:      sdk.OneInt(),
				DAOForkBlock:        sdk.OneInt(),
				EIP150Block:         sdk.OneInt(),
				EIP150Hash:          defaultEIP150Hash,
				EIP155Block:         sdk.OneInt(),
				EIP158Block:         sdk.OneInt(),
				ByzantiumBlock:      sdk.OneInt(),
				ConstantinopleBlock: sdk.OneInt(),
				PetersburgBlock:     sdk.NewInt(-1),
			},
			true,
		},
		{
			"invalid IstanbulBlock",
			Params{
				HomesteadBlock:      sdk.OneInt(),
				DAOForkBlock:        sdk.OneInt(),
				EIP150Block:         sdk.OneInt(),
				EIP150Hash:          defaultEIP150Hash,
				EIP155Block:         sdk.OneInt(),
				EIP158Block:         sdk.OneInt(),
				ByzantiumBlock:      sdk.OneInt(),
				ConstantinopleBlock: sdk.OneInt(),
				IstanbulBlock:       sdk.NewInt(-1),
			},
			true,
		},
		{
			"invalid MuirGlacierBlock",
			Params{
				HomesteadBlock:      sdk.OneInt(),
				DAOForkBlock:        sdk.OneInt(),
				EIP150Block:         sdk.OneInt(),
				EIP150Hash:          defaultEIP150Hash,
				EIP155Block:         sdk.OneInt(),
				EIP158Block:         sdk.OneInt(),
				ByzantiumBlock:      sdk.OneInt(),
				ConstantinopleBlock: sdk.OneInt(),
				PetersburgBlock:     sdk.OneInt(),
				IstanbulBlock:       sdk.OneInt(),
				MuirGlacierBlock:    sdk.NewInt(-1),
			},
			true,
		},
		{
			"invalid YoloV1Block",
			Params{
				HomesteadBlock:      sdk.OneInt(),
				DAOForkBlock:        sdk.OneInt(),
				EIP150Block:         sdk.OneInt(),
				EIP150Hash:          defaultEIP150Hash,
				EIP155Block:         sdk.OneInt(),
				EIP158Block:         sdk.OneInt(),
				ByzantiumBlock:      sdk.OneInt(),
				ConstantinopleBlock: sdk.OneInt(),
				PetersburgBlock:     sdk.OneInt(),
				IstanbulBlock:       sdk.OneInt(),
				MuirGlacierBlock:    sdk.OneInt(),
				YoloV1Block:         sdk.NewInt(-1),
			},
			true,
		},
		{
			"invalid EWASMBlock",
			Params{
				HomesteadBlock:      sdk.OneInt(),
				DAOForkBlock:        sdk.OneInt(),
				EIP150Block:         sdk.OneInt(),
				EIP150Hash:          defaultEIP150Hash,
				EIP155Block:         sdk.OneInt(),
				EIP158Block:         sdk.OneInt(),
				ByzantiumBlock:      sdk.OneInt(),
				ConstantinopleBlock: sdk.OneInt(),
				PetersburgBlock:     sdk.OneInt(),
				IstanbulBlock:       sdk.OneInt(),
				MuirGlacierBlock:    sdk.OneInt(),
				YoloV1Block:         sdk.OneInt(),
				EWASMBlock:          sdk.NewInt(-1),
			},
			true,
		},
	}

	for _, tc := range testCases {
		err := tc.params.Validate()

		if tc.expError {
			require.Error(t, err, tc.name)
		} else {
			require.NoError(t, err, tc.name)
		}
	}
}

func TestValidate(t *testing.T) {
	require.Error(t, validateInt(""))
	require.Error(t, validateDAOForkSupport(""))
	require.Error(t, validateHash(1))
}
