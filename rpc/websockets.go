package rpc

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var wsPort = 7545

type server struct{}

func newServer() {
	s := new(server)
	// TODO: add codec to turn . into _
	ws := mux.NewRouter()
	ws.Handle("/", s)
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", wsPort), ws)
		if err != nil {
			log.Println("http error:", err)
		}
	}()
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var upg = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	ws, err := upg.Upgrade(w, r, nil)
	if err != nil {
		log.Println("websocket upgrade failed; error:", err)
		return
	}

	for {
		_, mb, err := ws.ReadMessage()
		if err != nil {
			log.Println("failed to read message; error", err)
			return
		}

		log.Println("got websockets message!", mb)
	}
}

// PublicPubSubAPI is the eth_ prefixed set of APIs in the Web3 JSON-RPC spec
type PublicPubSubAPI struct {
	cliCtx  context.CLIContext
	backend Backend
}

// NewPublicPubSubAPI creates an instance of the public ETH Web3 PubSub API.
func NewPublicPubSubAPI(cliCtx context.CLIContext, backend Backend) *PublicPubSubAPI {
	newServer()

	return &PublicPubSubAPI{
		cliCtx:  cliCtx,
		backend: backend,
	}
}

func (p *PublicPubSubAPI) Subscribe() (rpc.ID, error) {
	return rpc.ID(0), nil
}

func (p *PublicPubSubAPI) Unsubscribe(id rpc.ID) {

}
