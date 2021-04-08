package types

import (
	"testing"

	"github.com/gogo/protobuf/proto"

	"github.com/stretchr/testify/require"
)

func TestMarshalMsgEthereumTxResponse(t *testing.T) {
	testhash := []byte{77}
	msg := &MsgEthereumTxResponse{
		Bloom: []byte{1, 2, 3},
		TxLogs: TransactionLogs{
			Logs: []*Log{
				{
					Address:     string(testhash),
					Topics:      []string{string(testhash)},
					Data:        []byte("data"),
					BlockNumber: 1,
					TxHash:      testhash,
					TxIndex:     1,
					BlockHash:   testhash,
					Index:       1,
					Removed:     false,
				},
			},
		},
		Ret:             []byte{7},
		ContractAddress: []byte{9, 9},
	}

	b, err := proto.Marshal(msg)
	require.NoError(t, err)

	res := &MsgEthereumTxResponse{}
	err = proto.Unmarshal(b, res)
	require.NoError(t, err)
	require.Equal(t, msg, res)
}
