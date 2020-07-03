package rpc

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	// "github.com/cosmos/cosmos-sdk/client/context"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"

	context "github.com/cosmos/cosmos-sdk/client/context"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// SubscriptionResponseJSON for json subscription responses
type SubscriptionResponseJSON struct {
	Jsonrpc string  `json:"jsonrpc"`
	Result  rpc.ID  `json:"result"`
	ID      float64 `json:"id"`
}

type SubscriptionNotification struct {
	Jsonrpc string              `json:"jsonrpc"`
	Method  string              `json:"method"`
	Params  *SubscriptionResult `json:"params"`
}

type SubscriptionResult struct {
	Subscription rpc.ID           `json:"subscription"`
	Result       *ethtypes.Header `json:"result"`
}

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
	//wsConn  *websocket.Conn
	api *pubSubAPI
}

func newWebsocketsServer(cliCtx context.CLIContext, wsAddr string) *websocketsServer {
	return &websocketsServer{
		rpcAddr: viper.GetString("laddr"),
		wsAddr:  wsAddr,
		api:     newPubSubAPI(cliCtx),
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
	//s.wsConn.Close()
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
			err = wsConn.WriteJSON(res)
			if err != nil {
				log.Println("websocket failed write message", "error", err)
			}
			continue
		}

		// TODO: check if method == eth_subscribe or eth_unsubscribe
		method := msg["method"]
		if method.(string) == "eth_subscribe" {
			id, err := s.api.subscribe(wsConn, "")
			if err != nil {
				log.Println("failed to subscribe; error", err)
				continue
			}

			res := &SubscriptionResponseJSON{
				Jsonrpc: "2.0",
				ID:      1,
				Result:  id,
			}

			err = wsConn.WriteJSON(res)
			if err != nil {
				log.Println("failed to write json response", err)
				continue
			}
		}

		// otherwise, call the usual rpc server to respond
		tcpConn, err := net.Dial("tcp", "localhost:8545")
		if err != nil {
			log.Println("cannot connect to tcp:localhost:8545", err)
			continue
		}

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
		req.Write(tcpConn)

		respBytes, err := ioutil.ReadAll(tcpConn)
		if err != nil {
			log.Println("error reading response body; error", err)
			return
		}

		respbuf := &bytes.Buffer{}
		respbuf.Write(respBytes)
		resp, err := http.ReadResponse(bufio.NewReader(respbuf), req)
		if err != nil {
			log.Println("could not read response; error", err)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("could not read body from response; error", err)
			continue
		}

		var wsSend interface{}
		err = json.Unmarshal(body, &wsSend)
		if err != nil {
			log.Println("error json unmarshal rpc response; error", err)
			continue
		}

		err = wsConn.WriteJSON(wsSend)
		if err != nil {
			log.Println("error writing json response; error", err)
			continue
		}
	}
}

type wsSubscription struct {
	sub  *Subscription
	conn *websocket.Conn
}

// pubSubAPI is the eth_ prefixed set of APIs in the Web3 JSON-RPC spec
type pubSubAPI struct {
	cliCtx    context.CLIContext
	events    *EventSystem
	filtersMu sync.Mutex
	filters   map[rpc.ID]*wsSubscription
}

// newPubSubAPI creates an instance of the ethereum PubSub API.
func newPubSubAPI(cliCtx context.CLIContext) *pubSubAPI {
	return &pubSubAPI{
		cliCtx:  cliCtx,
		events:  NewEventSystem(cliCtx.Client),
		filters: make(map[rpc.ID]*wsSubscription),
	}
}

func (api *pubSubAPI) subscribe(conn *websocket.Conn, method string) (rpc.ID, error) {
	// TODO: switch method

	sub, _, err := api.events.SubscribeNewHeads()
	if err != nil {
		return rpc.ID(0), fmt.Errorf("error creating block filter: %s", err.Error())
	}

	api.filtersMu.Lock()
	api.filters[sub.ID()] = &wsSubscription{
		sub:  sub,
		conn: conn,
	}
	api.filtersMu.Unlock()

	go func(headersCh <-chan coretypes.ResultEvent, errCh <-chan error) {
		for {
			select {
			case ev := <-headersCh:
				data, _ := ev.Data.(tmtypes.EventDataNewBlockHeader)
				header := EthHeaderFromTendermint(data.Header)

				api.filtersMu.Lock()
				if f, found := api.filters[sub.ID()]; found {
					// write to ws conn
					res := &SubscriptionNotification{
						Jsonrpc: "2.0",
						Method:  "eth_subscription",
						Params: &SubscriptionResult{
							Subscription: sub.ID(),
							Result:       header,
						},
					}

					err = f.conn.WriteJSON(res)
					if err != nil {
						log.Println("error writing header")
					}
				}
				api.filtersMu.Unlock()
			case <-errCh:
				api.filtersMu.Lock()
				delete(api.filters, sub.ID())
				api.filtersMu.Unlock()
				return
			}
		}
	}(sub.eventCh, sub.Err())

	return sub.ID(), nil
}

func (api *pubSubAPI) unsubscribe(id rpc.ID) {
	api.filtersMu.Lock()
	delete(api.filters, id)
	api.filtersMu.Unlock()
}
