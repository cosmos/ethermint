package loopback

import (
	"os"
	"path/filepath"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/tendermint/tendermint/libs/log"
	tmliteProxy "github.com/tendermint/tendermint/lite/proxy"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

type Params struct {
	Home      string
	ChainID   string
	NodeURI   string
	FromName  string
	Broadcast BroadcaseMode
}

type BroadcaseMode string

const (
	// BroadcastBlock defines a tx broadcasting mode where the client waits for
	// the tx to be committed in a block.
	BroadcastBlock BroadcaseMode = "block"
	// BroadcastSync defines a tx broadcasting mode where the client waits for
	// a CheckTx execution response only.
	BroadcastSync BroadcaseMode = "sync"
	// BroadcastAsync defines a tx broadcasting mode where the client returns
	// immediately.
	BroadcastAsync BroadcaseMode = "async"
)

func checkParams(params *Params) *Params {
	if params == nil {
		params = &Params{}
	}
	if len(params.Home) == 0 {
		params.Home = os.ExpandEnv("$HOME/.injectived")
	}
	if len(params.ChainID) == 0 {
		params.ChainID = "localnet"
	}
	if len(params.NodeURI) == 0 {
		params.NodeURI = "tcp://localhost:26657"
	}
	if len(params.Broadcast) == 0 {
		params.Broadcast = BroadcastSync
	}
	return params
}

func NewCosmosClientWithParams(cdc *codec.Codec, passphrase string, params *Params) CosmosClient {
	params = checkParams(params)
	cc := &cosmosClient{
		cdc:        cdc,
		params:     params,
		passphrase: passphrase,

		msgC:  make(chan sdk.Msg, 512),
		doneC: make(chan bool, 1),
	}
	cc.cliCtxOnce.Do(func() {
		cc.cliCtx = cc.newContextFromParams()
		cc.fromAcc = cc.cliCtx.GetFromAddress()
	})

	txEncoder := utils.GetTxEncoder(cdc)
	txBuilder := auth.NewTxBuilder(txEncoder, 0, 0, 0, 0, false, params.ChainID, "", nil, nil)
	cc.txBldr = txBuilder

	time.Sleep(5 * time.Second)
	go cc.runBatchBroadcast()

	return cc
}

const tmVerifierCacheSize = 10

func (c *cosmosClient) newContextFromParams() context.CLIContext {
	keybase := keys.New("keys", filepath.Join(c.params.Home, "keys"))
	nodeClient, err := rpchttp.New(c.params.NodeURI, "/websocket")
	if err != nil {
		panic(err)
	}
	verifier, err := tmliteProxy.NewVerifier(
		c.params.ChainID, filepath.Join(c.params.Home, ".lite_verifier"),
		nodeClient, log.NewNopLogger(), tmVerifierCacheSize,
	)

	from, _, err := getFromFields(keybase, c.params.FromName)
	if err != nil {
		panic(err)
	}

	ctx := context.CLIContext{
		Codec:         c.cdc,
		Client:        nodeClient,
		Keybase:       keybase,
		Output:        os.Stderr,
		OutputFormat:  "json",
		NodeURI:       c.params.NodeURI,
		FromName:      c.params.FromName,
		FromAddress:   from,
		TrustNode:     true,
		BroadcastMode: string(c.params.Broadcast),
		Verifier:      verifier,
		SkipConfirm:   true,
	}

	return ctx
}

func getFromFields(keybase keys.Keybase, from string) (sdk.AccAddress, string, error) {
	if from == "" {
		return nil, "", nil
	}

	var info keys.Info
	if addr, err := sdk.AccAddressFromBech32(from); err == nil {
		info, err = keybase.GetByAddress(addr)
		if err != nil {
			return nil, "", err
		}
	} else {
		info, err = keybase.Get(from)
		if err != nil {
			return nil, "", err
		}
	}

	return info.GetAddress(), info.GetName(), nil
}
