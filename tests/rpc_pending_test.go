// This is a test utility for Ethermint's Web3 JSON-RPC services.
//
// To run these tests please first ensure you have the ethermintd running
// and have started the RPC service with `ethermintcli rest-server`.
//
// You can configure the desired HOST and MODE as well
package tests

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	addrA         = "0xc94770007dda54cF92009BFF0dE90c06F603a09f"
	addrAStoreKey = 0
)

var (
	MODE = os.Getenv("MODE")
	HOST = os.Getenv("HOST")

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

	if HOST == "" {
		HOST = "http://localhost:8545"
	}

	var err error
	from, err = getAddress()
	if err != nil {
		fmt.Printf("failed to get account: %s\n", err)
		os.Exit(1)
	}

	// Start all tests
	code := m.Run()
	os.Exit(code)
}

func getAddress() ([]byte, error) {
	rpcRes, err := callWithError("eth_accounts", []string{})
	if err != nil {
		return nil, err
	}

	var res []hexutil.Bytes
	err = json.Unmarshal(rpcRes.Result, &res)
	if err != nil {
		return nil, err
	}

	return res[0], nil
}

func createRequest(method string, params interface{}) Request {
	return Request{
		Version: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}
}

func call(t *testing.T, method string, params interface{}) *Response {
	req, err := json.Marshal(createRequest(method, params))
	require.NoError(t, err)

	var rpcRes *Response
	time.Sleep(1 * time.Second)
	/* #nosec */
	res, err := http.Post(HOST, "application/json", bytes.NewBuffer(req))
	require.NoError(t, err)

	decoder := json.NewDecoder(res.Body)
	rpcRes = new(Response)
	err = decoder.Decode(&rpcRes)
	require.NoError(t, err)

	err = res.Body.Close()
	require.NoError(t, err)
	require.Nil(t, rpcRes.Error)

	return rpcRes
}

func callWithError(method string, params interface{}) (*Response, error) {
	req, err := json.Marshal(createRequest(method, params))
	if err != nil {
		return nil, err
	}

	var rpcRes *Response
	time.Sleep(1 * time.Second)
	/* #nosec */
	res, err := http.Post(HOST, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(res.Body)
	rpcRes = new(Response)
	err = decoder.Decode(&rpcRes)
	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	if rpcRes.Error != nil {
		return nil, fmt.Errorf(rpcRes.Error.Message)
	}

	return rpcRes, nil
}

// turns a 0x prefixed hex string to a big.Int
func hexToBigInt(t *testing.T, in string) *big.Int {
	s := in[2:]
	b, err := hex.DecodeString(s)
	require.NoError(t, err)
	return big.NewInt(0).SetBytes(b)
}

// sendTestTransaction sends a dummy transaction
func sendTestTransaction(t *testing.T) hexutil.Bytes {
	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
	param[0]["to"] = "0x1122334455667788990011223344556677889900"
	param[0]["value"] = "0x1"
	rpcRes := call(t, "eth_sendTransaction", param)

	var hash hexutil.Bytes
	err := json.Unmarshal(rpcRes.Result, &hash)
	require.NoError(t, err)
	return hash
}

// deployTestContract deploys a contract that emits an event in the constructor
func deployTestContract(t *testing.T) (hexutil.Bytes, map[string]interface{}) {
	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
	param[0]["data"] = "0x6080604052348015600f57600080fd5b5060117f775a94827b8fd9b519d36cd827093c664f93347070a554f65e4a6f56cd73889860405160405180910390a2603580604b6000396000f3fe6080604052600080fdfea165627a7a723058206cab665f0f557620554bb45adf266708d2bd349b8a4314bdff205ee8440e3c240029"
	param[0]["gas"] = "0x200000"

	rpcRes := call(t, "eth_sendTransaction", param)

	var hash hexutil.Bytes
	err := json.Unmarshal(rpcRes.Result, &hash)
	require.NoError(t, err)

	receipt := waitForReceipt(t, hash)
	require.NotNil(t, receipt, "transaction failed")
	require.Equal(t, "0x1", receipt["status"].(string))

	return hash, receipt
}

func getTransactionReceipt(t *testing.T, hash hexutil.Bytes) map[string]interface{} {
	param := []string{hash.String()}
	rpcRes := call(t, "eth_getTransactionReceipt", param)

	receipt := make(map[string]interface{})
	err := json.Unmarshal(rpcRes.Result, &receipt)
	require.NoError(t, err)

	return receipt
}

func waitForReceipt(t *testing.T, hash hexutil.Bytes) map[string]interface{} {
	for i := 0; i < 12; i++ {
		receipt := getTransactionReceipt(t, hash)
		if receipt != nil {
			return receipt
		}

		time.Sleep(time.Second)
	}

	return nil
}

func TestEth_Pending_GetBalance(t *testing.T) {
	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
	param[0]["to"] = addrA
	param[0]["value"] = "0xA"
	param[0]["gasLimit"] = "0x5208"
	param[0]["gasPrice"] = "0x1"

	rpcRes := call(t, "eth_sendTransaction", param)

	var res hexutil.Big

	rpcRes = call(t, "eth_getBalance", []string{addrA, "pending"})
	err := res.UnmarshalJSON(rpcRes.Result)
	require.NoError(t, err)
	t.Logf("Got balance %s for %s\n", res.String(), addrA)

	if res.ToInt().Cmp(big.NewInt(10)) != 0 {
		t.Errorf("expected balance: %d, got: %s", 10, res.String())
	}

	rpcRes = call(t, "eth_getBalance", []string{addrA, "latest"})
	err = res.UnmarshalJSON(rpcRes.Result)
	require.NoError(t, err)
	t.Logf("Got balance %s for %s\n", res.String(), addrA)

	if res.ToInt().Cmp(big.NewInt(0)) != 0 {
		t.Errorf("expected balance: %d, got: %s", 0, res.String())
	}
}

func TestEth_Pending_GetTransactionCount(t *testing.T) {
	currentNonce := getNonce(t, "latest")
	t.Logf("Current nonce is %d", currentNonce)

	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
	param[0]["to"] = addrA
	param[0]["value"] = "0xA"
	param[0]["gasLimit"] = "0x5208"
	param[0]["gasPrice"] = "0x1"

	_ = call(t, "eth_sendTransaction", param)

	pendingNonce := getNonce(t, "pending")
	latestNonce := getNonce(t, "latest")
	t.Logf("Latest nonce is %d", latestNonce)
	require.Equal(t, currentNonce, latestNonce)
	t.Logf("Pending nonce is %d", pendingNonce)
	require.NotEqual(t, latestNonce, pendingNonce)
}

func TestEth_Pending_GetBlockTransactionCountByNumber(t *testing.T) {
	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
	param[0]["to"] = addrA
	param[0]["value"] = "0xA"
	param[0]["gasLimit"] = "0x5208"
	param[0]["gasPrice"] = "0x1"

	_ = call(t, "eth_sendTransaction", param)

	rpcRes := call(t, "eth_getBlockTransactionCountByNumber", []interface{}{"pending"})
	var pendingTxCount hexutil.Uint
	err := json.Unmarshal(rpcRes.Result, &pendingTxCount)
	require.NoError(t, err)
	t.Logf("Pending nonce is %d", pendingTxCount)

	rpcRes = call(t, "eth_getBlockTransactionCountByNumber", []interface{}{"latest"})
	var latestTxCount hexutil.Uint
	err = json.Unmarshal(rpcRes.Result, &latestTxCount)
	require.NoError(t, err)
	t.Logf("Latest nonce is %d", latestTxCount)

	require.NotEqual(t, pendingTxCount, latestTxCount)
}

func TestEth_Pending_GetBlockByNumber(t *testing.T) {
	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
	param[0]["to"] = addrA
	param[0]["value"] = "0xA"
	param[0]["gasLimit"] = "0x5208"
	param[0]["gasPrice"] = "0x1"

	_ = call(t, "eth_sendTransaction", param)

	rpcRes := call(t, "eth_getBlockByNumber", []interface{}{"pending", true})
	var pendingBlock map[string]interface{}
	err := json.Unmarshal(rpcRes.Result, &pendingTxCount)
	require.NoError(t, err)
	t.Logf("Pending nonce is %d", pendingTxCount)

}
