package rpc

import (
	"log"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const defaultWebsocketPort = 7545

// TODO: add logger
type server struct{}

func newServer(websocketAddr string) {
	s := new(server)
	// TODO: add codec to turn . into _
	ws := mux.NewRouter()
	ws.Handle("/", s)
	go func() {
		err := http.ListenAndServe(websocketAddr, ws)
		if err != nil {
			log.Println("http error:", err)
		}
	}()
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// TODO: write error
		log.Println("websocket upgrade failed; error:", err)
		return
	}

	s.readLoop(wsConn)
}

func (*server) readLoop(wsConn *websocket.Conn) {
	for {
		_, mb, err := wsConn.ReadMessage()
		if err != nil {
			// TODO: write error
			wsConn.Close()
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
func NewPublicPubSubAPI(cliCtx context.CLIContext, backend Backend, websocketAddr string) *PublicPubSubAPI {
	newServer(websocketAddr)

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
