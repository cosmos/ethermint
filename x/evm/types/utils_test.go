package types

import (
	"testing"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func TestEvmDataEncoding(t *testing.T) {
	addr := ethcmn.HexToAddress("0x12345")
	bloom := ethtypes.BytesToBloom([]byte{0x1, 0x3})
	ret := []byte{0x5, 0x8}

	data := &ResultData{
		Address: addr,
		Bloom:   bloom,
		Logs: []*ethtypes.Log{{
			Data:        []byte{1, 2, 3, 4},
			BlockNumber: 17,
		}},
		Ret: ret,
	}

	enc, err := EncodeResultData(data)
	require.NoError(t, err)

	res, err := DecodeResultData(enc)
	require.NoError(t, err)
	require.Equal(t, addr, res.Address)
	require.Equal(t, bloom, res.Bloom)
	require.Equal(t, data.Logs, res.Logs)
	require.Equal(t, ret, res.Ret)
}
