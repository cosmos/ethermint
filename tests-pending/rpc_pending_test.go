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

	zeroString = "0x0"
	from       = []byte{}
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
	preTxBalance := res.ToInt()
	t.Logf("Got balance %s for %s before sending tx\n", preTxBalance, addrA)

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
	t.Logf("Got pending balance %s for %s post tx\n", res.String(), addrA)

	if res.ToInt().Cmp(preTxBalance.Add(preTxBalance, big.NewInt(10))) != 0 {
		t.Errorf("expected balance: %d, got: %s", 10, res.String())
	}

	rpcRes = util.Call(t, "eth_getBalance", []string{addrA, "latest"})
	err = res.UnmarshalJSON(rpcRes.Result)
	require.NoError(t, err)
	t.Logf("Got latest balance %s for %s post tx\n", res.String(), addrA)

	if res.ToInt().Cmp(preTxBalance) != 0 {
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
	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
	param[0]["to"] = addrA
	param[0]["value"] = "0xA"
	param[0]["gasLimit"] = "0x5208"
	param[0]["gasPrice"] = "0x1"

	_ = util.Call(t, "eth_sendTransaction", param)

	rpcRes := util.Call(t, "eth_getBlockTransactionCountByNumber", []interface{}{"pending"})
	var pendingTxCount hexutil.Uint
	err := json.Unmarshal(rpcRes.Result, &pendingTxCount)
	require.NoError(t, err)
	t.Logf("Pending nonce is %d", pendingTxCount)

	rpcRes = util.Call(t, "eth_getBlockTransactionCountByNumber", []interface{}{"latest"})
	var latestTxCount hexutil.Uint
	err = json.Unmarshal(rpcRes.Result, &latestTxCount)
	require.NoError(t, err)
	t.Logf("Latest nonce is %d", latestTxCount)

	require.NotEqual(t, uint64(pendingTxCount), uint64(latestTxCount))
	require.Greater(t, uint64(pendingTxCount), uint64(latestTxCount))
}

func TestEth_Pending_GetBlockByNumber(t *testing.T) {
	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
	param[0]["to"] = addrA
	param[0]["value"] = "0xA"
	param[0]["gasLimit"] = "0x5208"
	param[0]["gasPrice"] = "0x1"

	_ = util.Call(t, "eth_sendTransaction", param)

	rpcRes := util.Call(t, "eth_getBlockByNumber", []interface{}{"pending", true})
	var pendingBlock map[string]interface{}
	err := json.Unmarshal(rpcRes.Result, &pendingBlock)
	require.NoError(t, err)
	pendingBlockArr := len(pendingBlock["transactions"].([]interface{}))
	// in case there are other tx's inside the pending queue from prev test
	require.NotEqual(t, pendingBlockArr, 0)

	rpcRes = util.Call(t, "eth_getBlockByNumber", []interface{}{"latest", true})
	var latestBlock map[string]interface{}
	err = json.Unmarshal(rpcRes.Result, &latestBlock)
	require.NoError(t, err)
	latestBlockArr := len(latestBlock["transactions"].([]interface{}))
	require.Equal(t, latestBlockArr, 0)

	require.Greater(t, pendingBlockArr, latestBlockArr)
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

	rpcRes = util.Call(t, "eth_getBlockByNumber", []interface{}{"latest", true})
	var latestBlock map[string]interface{}
	err = json.Unmarshal(rpcRes.Result, &latestBlock)
	fmt.Println("latestBlock: ", latestBlock)
	require.NoError(t, err)

	require.NotEmpty(t, latestBlock["timestamp"])
	require.NotEmpty(t, latestBlock["gasUsed"])
	require.Equal(t, latestBlock["logsBloom"], "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
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
	require.Equal(t, pendingBlock["blockHash"], nil)
	require.Equal(t, pendingBlock["blockNumber"], nil)
	require.Equal(t, pendingBlock["transactionIndex"], nil)
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

	t.Logf("currNonce: %d", currNonce)

	// first transaction
	_ = util.Call(t, "eth_sendTransaction", param)
	pendingNonce1 := util.GetNonce(t, "pending")
	require.Greater(t, pendingNonce1, currNonce)

	// second transaction
	param[0]["to"] = "0x7f0f463c4d57b1bd3e3b79051e6c5ab703e803d9"
	_ = util.Call(t, "eth_sendTransaction", param)
	pendingNonce2 := util.GetNonce(t, "pending")
	require.Greater(t, pendingNonce2, currNonce)
	require.Equal(t, pendingNonce1+hexutil.Uint64(1), pendingNonce2)

	// third transaction
	param[0]["to"] = "0x7fb24493808b3f10527e3e0870afeb8a953052d2"
	_ = util.Call(t, "eth_sendTransaction", param)
	pendingNonce3 := util.GetNonce(t, "pending")
	require.Greater(t, pendingNonce3, currNonce)
	require.Equal(t, pendingNonce1+hexutil.Uint64(2), pendingNonce3)
}
