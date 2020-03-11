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
		Logs:    []*ethtypes.Log{},
		Ret:     ret,
	}

	enc, err := EncodeResultData(data)
	if err != nil {
		t.Fatal(err)
	}

	res, err := DecodeResultData(enc)
	if err != nil {
		t.Fatal(err)
	}

	require.NoError(t, err)
	require.Equal(t, addr, res.Address)
	require.Equal(t, bloom, res.Bloom)
	require.Equal(t, data.Logs, res.Logs)
	require.Equal(t, ret, res.Ret)
}
