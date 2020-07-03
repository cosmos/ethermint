package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	// "github.com/cosmos/cosmos-sdk/client/context"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// "github.com/ethereum/go-ethereum/rpc"
)

// ErrorResponseJSON json for error responses
type ErrorResponseJSON struct {
	Jsonrpc string            `json:"jsonrpc"`
	Error   *ErrorMessageJSON `json:"error"`
	ID      *big.Int          `json:"id"`
}

// ErrorMessageJSON json for error messages
type ErrorMessageJSON struct {
	Code    *big.Int `json:"code"`
	Message string   `json:"message"`
}

// TODO: add logger
type websocketsServer struct {
	rpcAddr string // listen address of rest-server
	wsAddr  string // listen address of ws server
	wsConn  *websocket.Conn
}

func newWebsocketsServer(wsAddr string) *websocketsServer {
	return &websocketsServer{
		rpcAddr: viper.GetString("laddr"),
		wsAddr:  wsAddr,
	}
}

func (s *websocketsServer) start() {
	ws := mux.NewRouter()
	ws.Handle("/", s)

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%s", s.wsAddr), ws)
		if err != nil {
			log.Println("http error:", err)
		}
	}()
}

func (s *websocketsServer) stop() {
	s.wsConn.Close()
}

func (s *websocketsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("websocket upgrade failed; error:", err)
		return
	}

	s.wsConn = wsConn
	s.readLoop(wsConn)
}

func (s *websocketsServer) readLoop(wsConn *websocket.Conn) {
	for {
		_, mb, err := wsConn.ReadMessage()
		if err != nil {
			// TODO: write error
			wsConn.Close()
			log.Println("failed to read message; error", err)
			return
		}

		log.Println("got websockets message!", mb)

		// determine if request is for subscribe method type
		var msg map[string]interface{}
		err = json.Unmarshal(mb, &msg)
		if err != nil {
			log.Println("websocket failed to unmarshal request message; error", err)
			res := &ErrorResponseJSON{
				Jsonrpc: "2.0",
				Error: &ErrorMessageJSON{
					Code:    big.NewInt(-32600),
					Message: "Invalid request",
				},
				ID: nil,
			}
			err = s.wsConn.WriteJSON(res)
			if err != nil {
				log.Println("[rpc] websocket failed write message", "error", err)
			}
			continue
		}

		// TODO: check if method == eth_subscribe or eth_unsubscribe
		//method := msg["method"]

		// otherwise, call the usual rpc server to respond
		client := &http.Client{}
		buf := &bytes.Buffer{}
		_, err = buf.Write(mb)
		if err != nil {
			log.Println("failed to write message to buffer; error", err)
			return
		}

		req, err := http.NewRequest("POST", s.rpcAddr, buf)
		if err != nil {
			log.Println("failed request to rpc service; error", err)
			return
		}

		req.Header.Set("Content-Type", "application/json;")

		res, err := client.Do(req)
		if err != nil {
			log.Println("websocket error calling rpc; error", err)
			return
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println("error reading response body; error", err)
			return
		}

		err = res.Body.Close()
		if err != nil {
			log.Println("error closing response body; error", err)
			return
		}

		var wsSend interface{}
		err = json.Unmarshal(body, &wsSend)
		if err != nil {
			log.Println("error unmarshal rpc response; error", err)
			return
		}

		err = s.wsConn.WriteJSON(wsSend)
		if err != nil {
			log.Println("error writing json response; error", err)
			return
		}
	}
}

// // PublicPubSubAPI is the eth_ prefixed set of APIs in the Web3 JSON-RPC spec
// type PublicPubSubAPI struct {
// 	cliCtx  context.CLIContext
// 	backend Backend
// }

// // NewPublicPubSubAPI creates an instance of the public ETH Web3 PubSub API.
// func NewPublicPubSubAPI(cliCtx context.CLIContext, backend Backend, websocketAddr string) *PublicPubSubAPI {
// 	newServer(websocketAddr)

// 	return &PublicPubSubAPI{
// 		cliCtx:  cliCtx,
// 		backend: backend,
// 	}
// }

// func (p *PublicPubSubAPI) Subscribe() (rpc.ID, error) {
// 	return rpc.ID(0), nil
// }

// func (p *PublicPubSubAPI) Unsubscribe(id rpc.ID) {

// }
