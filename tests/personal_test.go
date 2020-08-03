package tests

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func TestPersonal_ListAccounts(t *testing.T) {
	rpcRes := call(t, "personal_listAccounts", []string{})

	var res []hexutil.Bytes
	err := json.Unmarshal(rpcRes.Result, &res)
	require.NoError(t, err)
	require.Equal(t, 1, len(res))
}