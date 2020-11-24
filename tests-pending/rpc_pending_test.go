// This is a test utility for Ethermint's Web3 JSON-RPC services.
//
// To run these tests please first ensure you have the ethermintd running
// and have started the RPC service with `ethermintcli rest-server`.
//
// You can configure the desired HOST and MODE as well
package pending

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	util "github.com/cosmos/ethermint/tests"
)

const (
	addrA         = "0xc94770007dda54cF92009BFF0dE90c06F603a09f"
	addrAStoreKey = 0
)

var (
	MODE = os.Getenv("MODE")
	from = []byte{}
)

type Request struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Response struct {
	Error  *RPCError       `json:"error"`
	ID     int             `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
}

func TestMain(m *testing.M) {
	if MODE != "pending" {
		_, _ = fmt.Fprintln(os.Stdout, "Skipping pending RPC test")
		return
	}

	var err error
	from, err = util.GetAddress()
	if err != nil {
		fmt.Printf("failed to get account: %s\n", err)
		os.Exit(1)
	}

	// Start all tests
	code := m.Run()
	os.Exit(code)
}

func TestEth_Pending_GetBalance(t *testing.T) {
	var res hexutil.Big
	rpcRes := util.Call(t, "eth_getBalance", []string{addrA, "latest"})
	err := res.UnmarshalJSON(rpcRes.Result)
	require.NoError(t, err)
	preTxLatestBalance := res.ToInt()

	rpcRes = util.Call(t, "eth_getBalance", []string{addrA, "pending"})
	err = res.UnmarshalJSON(rpcRes.Result)
	require.NoError(t, err)
	preTxPendingBalance := res.ToInt()

	t.Logf("Got pending balance %s for %s pre tx\n", preTxPendingBalance, addrA)
	t.Logf("Got latest balance %s for %s pre tx\n", preTxLatestBalance, addrA)

	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
	param[0]["to"] = addrA
	param[0]["value"] = "0xA"
	param[0]["gasLimit"] = "0x5208"
	param[0]["gasPrice"] = "0x1"

	_ = util.Call(t, "eth_sendTransaction", param)

	rpcRes = util.Call(t, "eth_getBalance", []string{addrA, "pending"})
	err = res.UnmarshalJSON(rpcRes.Result)
	require.NoError(t, err)
	postTxPendingBalance := res.ToInt()
	t.Logf("Got pending balance %s for %s post tx\n", postTxPendingBalance, addrA)

	require.Equal(t, preTxPendingBalance.Add(preTxPendingBalance, big.NewInt(10)), postTxPendingBalance)

	rpcRes = util.Call(t, "eth_getBalance", []string{addrA, "latest"})
	err = res.UnmarshalJSON(rpcRes.Result)
	require.NoError(t, err)
	postTxLatestBalance := res.ToInt()
	t.Logf("Got latest balance %s for %s post tx\n", postTxLatestBalance, addrA)

	require.Equal(t, preTxLatestBalance, postTxLatestBalance)
	// 1 if x > y; where x is postTxPendingBalance, y is preTxPendingBalance
	if postTxLatestBalance.Cmp(preTxPendingBalance) != 0 {
		t.Errorf("expected balance: %d, got: %s", 0, res.String())
	}
}

func TestEth_Pending_GetTransactionCount(t *testing.T) {
	currentNonce := util.GetNonce(t, "latest")
	t.Logf("Current nonce is %d", currentNonce)

	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
	param[0]["to"] = addrA
	param[0]["value"] = "0xA"
	param[0]["gasLimit"] = "0x5208"
	param[0]["gasPrice"] = "0x1"

	_ = util.Call(t, "eth_sendTransaction", param)

	pendingNonce := util.GetNonce(t, "pending")
	latestNonce := util.GetNonce(t, "latest")
	t.Logf("Latest nonce is %d", latestNonce)
	require.Equal(t, currentNonce, latestNonce)
	t.Logf("Pending nonce is %d", pendingNonce)
	require.NotEqual(t, latestNonce, pendingNonce)

	require.Greater(t, uint64(pendingNonce), uint64(latestNonce))
}

func TestEth_Pending_GetBlockTransactionCountByNumber(t *testing.T) {
	rpcRes := util.Call(t, "eth_getBlockTransactionCountByNumber", []interface{}{"pending"})
	var preTxPendingTxCount hexutil.Uint
	err := json.Unmarshal(rpcRes.Result, &preTxPendingTxCount)
	require.NoError(t, err)
	t.Logf("Pre tx pending nonce is %d", preTxPendingTxCount)

	rpcRes = util.Call(t, "eth_getBlockTransactionCountByNumber", []interface{}{"latest"})
	var preTxLatestTxCount hexutil.Uint
	err = json.Unmarshal(rpcRes.Result, &preTxLatestTxCount)
	require.NoError(t, err)
	t.Logf("Pre tx latest nonce is %d", preTxLatestTxCount)

	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
	param[0]["to"] = addrA
	param[0]["value"] = "0xA"
	param[0]["gasLimit"] = "0x5208"
	param[0]["gasPrice"] = "0x1"

	_ = util.Call(t, "eth_sendTransaction", param)

	rpcRes = util.Call(t, "eth_getBlockTransactionCountByNumber", []interface{}{"pending"})
	var postTxPendingTxCount hexutil.Uint
	err = json.Unmarshal(rpcRes.Result, &postTxPendingTxCount)
	require.NoError(t, err)
	t.Logf("Post tx pending nonce is %d", postTxPendingTxCount)

	rpcRes = util.Call(t, "eth_getBlockTransactionCountByNumber", []interface{}{"latest"})
	var postTxLatestTxCount hexutil.Uint
	err = json.Unmarshal(rpcRes.Result, &postTxLatestTxCount)
	require.NoError(t, err)
	t.Logf("Post tx latest nonce is %d", postTxLatestTxCount)

	require.Equal(t, uint64(preTxPendingTxCount)+uint64(1), uint64(postTxPendingTxCount))
	require.NotEqual(t, uint64(postTxPendingTxCount)-uint64(preTxPendingTxCount), uint64(postTxLatestTxCount)-uint64(preTxLatestTxCount))
}

func TestEth_Pending_GetBlockByNumber(t *testing.T) {
	rpcRes := util.Call(t, "eth_getBlockByNumber", []interface{}{"latest", true})
	var preTxLatestBlock map[string]interface{}
	err := json.Unmarshal(rpcRes.Result, &preTxLatestBlock)
	require.NoError(t, err)
	preTxLatestTxs := len(preTxLatestBlock["transactions"].([]interface{}))

	rpcRes = util.Call(t, "eth_getBlockByNumber", []interface{}{"pending", true})
	var preTxPendingBlock map[string]interface{}
	err = json.Unmarshal(rpcRes.Result, &preTxPendingBlock)
	require.NoError(t, err)
	preTxPendingTxs := len(preTxPendingBlock["transactions"].([]interface{}))

	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
	param[0]["to"] = addrA
	param[0]["value"] = "0xA"
	param[0]["gasLimit"] = "0x5208"
	param[0]["gasPrice"] = "0x1"

	_ = util.Call(t, "eth_sendTransaction", param)

	rpcRes = util.Call(t, "eth_getBlockByNumber", []interface{}{"pending", true})
	var postTxPendingBlock map[string]interface{}
	err = json.Unmarshal(rpcRes.Result, &postTxPendingBlock)
	require.NoError(t, err)
	postTxPendingTxs := len(postTxPendingBlock["transactions"].([]interface{}))
	require.Greater(t, postTxPendingTxs, preTxPendingTxs)

	rpcRes = util.Call(t, "eth_getBlockByNumber", []interface{}{"latest", true})
	var postTxLatestBlock map[string]interface{}
	err = json.Unmarshal(rpcRes.Result, &postTxLatestBlock)
	require.NoError(t, err)
	postTxLatestTxs := len(postTxLatestBlock["transactions"].([]interface{}))
	require.Equal(t, preTxLatestTxs, postTxLatestTxs)

	require.Greater(t, postTxPendingTxs, preTxPendingTxs)
}

func TestEth_Pending_GetTransactionByBlockNumberAndIndex(t *testing.T) {
	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
	param[0]["to"] = addrA
	param[0]["value"] = "0xA"
	param[0]["gasLimit"] = "0x5208"
	param[0]["gasPrice"] = "0x1"

	_ = util.Call(t, "eth_sendTransaction", param)

	rpcRes := util.Call(t, "eth_getTransactionByBlockNumberAndIndex", []interface{}{"pending", "0x1"})
	var pendingBlock map[string]interface{}
	err := json.Unmarshal(rpcRes.Result, &pendingBlock)
	require.NoError(t, err)

	require.Equal(t, pendingBlock["blockHash"], nil)
	require.Equal(t, pendingBlock["blockNumber"], nil)
	require.Equal(t, pendingBlock["transactionIndex"], nil)
	require.NotEmpty(t, pendingBlock["hash"])

	rpcRes = util.Call(t, "eth_getTransactionByBlockNumberAndIndex", []interface{}{"latest", "0x1"})
	var latestBlock map[string]interface{}
	err = json.Unmarshal(rpcRes.Result, &latestBlock)
	require.NoError(t, err)

	require.Empty(t, latestBlock)
}

func TestEth_Pending_GetTransactionByHash(t *testing.T) {
	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
	param[0]["to"] = addrA
	param[0]["value"] = "0xA"
	param[0]["gasLimit"] = "0x5208"
	param[0]["gasPrice"] = "0x1"

	txRes := util.Call(t, "eth_sendTransaction", param)
	var txHash common.Hash
	err := txHash.UnmarshalJSON(txRes.Result)
	require.NoError(t, err)

	rpcRes := util.Call(t, "eth_getTransactionByHash", []interface{}{txHash})
	var pendingBlock map[string]interface{}
	err = json.Unmarshal(rpcRes.Result, &pendingBlock)
	require.NoError(t, err)

	require.NotEmpty(t, pendingBlock)
	require.Equal(t, nil, pendingBlock["blockHash"])
	require.Equal(t, nil, pendingBlock["blockNumber"])
	require.Equal(t, nil, pendingBlock["transactionIndex"])
	require.NotEmpty(t, pendingBlock["hash"])
	require.NotEmpty(t, pendingBlock["value"], "0xa")
}

func TestEth_Pending_SendTransaction_PendingNonce(t *testing.T) {
	currNonce := util.GetNonce(t, "latest")
	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
	param[0]["to"] = addrA
	param[0]["value"] = "0xA"
	param[0]["gasLimit"] = "0x5208"
	param[0]["gasPrice"] = "0x1"

	// first transaction
	_ = util.Call(t, "eth_sendTransaction", param)
	pendingNonce1 := util.GetNonce(t, "pending")
	require.Greater(t, uint64(pendingNonce1), uint64(currNonce))

	// second transaction
	param[0]["to"] = "0x7f0f463c4d57b1bd3e3b79051e6c5ab703e803d9"
	_ = util.Call(t, "eth_sendTransaction", param)
	pendingNonce2 := util.GetNonce(t, "pending")
	require.Greater(t, uint64(pendingNonce2), uint64(currNonce))
	require.Greater(t, uint64(pendingNonce2), uint64(pendingNonce1))

	// third transaction
	param[0]["to"] = "0x7fb24493808b3f10527e3e0870afeb8a953052d2"
	_ = util.Call(t, "eth_sendTransaction", param)
	pendingNonce3 := util.GetNonce(t, "pending")
	require.Greater(t, uint64(pendingNonce3), uint64(currNonce))
	require.Greater(t, uint64(pendingNonce3), uint64(pendingNonce2))
}

// func TestEth_Call_Pending(t *testing.T) {
// 	param := make([]map[string]string, 1)
// 	param[0] = make(map[string]string)
// 	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
// 	param[0]["to"] = "0x0000000000000000000000000000000012341234"
// 	param[0]["value"] = "0xA"
// 	param[0]["gasLimit"] = "0x5208"
// 	param[0]["gasPrice"] = "0x1"

// 	rpcRes := util.Call(t, "eth_sendTransaction", param)

// 	var hash hexutil.Bytes
// 	err := json.Unmarshal(rpcRes.Result, &hash)
// 	require.NoError(t, err)

// 	param = make([]map[string]string, 1)
// 	param[0] = make(map[string]string)
// 	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
// 	param[0]["to"] = "0x0000000000000000000000000000000012341234"
// 	param[0]["value"] = "0xA"
// 	param[0]["gasLimit"] = "0x5208"
// 	param[0]["gasPrice"] = "0x1"

// 	rpcRes = util.Call(t, "eth_call", []interface{}{param[0], "pending"})
// 	err = json.Unmarshal(rpcRes.Result, &hash)
// 	require.NoError(t, err)
// }
