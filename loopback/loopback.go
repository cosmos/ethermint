package loopback

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	log "github.com/xlab/suplog"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/cosmos/ethermint/metrics"
)

type CosmosClient interface {
	CanSignTransactions() bool
	FromAddress() sdk.AccAddress
	SyncBroadcastMsg(msg sdk.Msg) (*sdk.TxResponse, error)
	QueueBroadcastMsg(msg sdk.Msg) error
	CliCtx() context.CLIContext
	Close()
}

func NewCosmosClient(cdc *codec.Codec, canSign bool, passphrase string) CosmosClient {
	cc := &cosmosClient{
		cdc:     cdc,
		canSign: canSign,

		svcTags: metrics.Tags{
			"svc": "loopback",
		},

		syncMux: new(sync.Mutex),
	}

	if canSign {
		cc.provider = &passphraseProvider{
			pass: passphrase,
		}

		cc.txBldr = auth.NewTxBuilderFromCLI(cc.provider).WithTxEncoder(utils.GetTxEncoder(cdc))
		cc.passphrase = passphrase
		cc.msgC = make(chan sdk.Msg, 512)
		cc.doneC = make(chan bool, 1)

		time.Sleep(5 * time.Second)

		go cc.runBatchBroadcast()
	}

	return cc
}

type passphraseProvider struct {
	pass string
}

func (p *passphraseProvider) Read(buf []byte) (n int, err error) {
	if len(buf) < len(p.pass)+1 {
		panic(fmt.Sprintf("buf is too short: %d", len(buf)))
	}

	log.Infoln("reading passphrase into buf for keyring unlock")
	defer func() {
		log.Infoln("done reading passphrase")
	}()

	return bytes.NewBufferString(p.pass + "\n").Read(buf)
}

func (c *cosmosClient) newCliCtx() context.CLIContext {
	c.cliCtxOnce.Do(func() {
		c.cliCtx = context.NewCLIContextWithInput(c.provider).WithCodec(c.cdc)
		c.cliCtx.SkipConfirm = true
		c.cliCtx.TrustNode = true
		c.fromAcc = c.cliCtx.GetFromAddress()
	})

	return c.cliCtx
}

type cosmosClient struct {
	svcTags metrics.Tags
	params  *Params

	cdc      *codec.Codec
	txBldr   auth.TxBuilder
	provider io.Reader

	cliCtx     context.CLIContext
	cliCtxOnce sync.Once

	doneC   chan bool
	msgC    chan sdk.Msg
	syncMux *sync.Mutex

	fromAcc         sdk.AccAddress
	fromAccSequence uint64
	fromAccNumber   *uint64
	passphrase      string

	closed  int64
	canSign bool
}

func (c *cosmosClient) CanSignTransactions() bool {
	return c.canSign
}

func (c *cosmosClient) CliCtx() context.CLIContext {
	c.newCliCtx()
	return c.cliCtx
}

func (c *cosmosClient) FromAddress() sdk.AccAddress {
	if !c.canSign {
		return sdk.AccAddress{}
	}

	c.newCliCtx()
	return c.fromAcc
}

var (
	ErrQueueClosed    = errors.New("queue is closed")
	ErrEnqueueTimeout = errors.New("enqueue timeout")
	ErrReadOnly       = errors.New("client is in read-only mode")
)

func (c *cosmosClient) SyncBroadcastMsg(msg sdk.Msg) (*sdk.TxResponse, error) {
	c.syncMux.Lock()
	defer c.syncMux.Unlock()

	cliCtx := c.newCliCtx()

	txBldr, err := c.prepareTxBuilder(cliCtx, c.txBldr)
	if err != nil {
		return nil, err
	}

	fromName := cliCtx.GetFromName()
	txBytes, err := txBldr.BuildAndSign(fromName, c.passphrase, []sdk.Msg{msg})
	if err != nil {
		err = errors.Wrap(err, "buildAndSign failed")
		return nil, err
	}

	// broadcast to a Tendermint node
	resp, err := cliCtx.BroadcastTxCommit(txBytes)
	if err != nil {
		c.rollbackSequence()
		err = errors.Wrap(err, "broadcastTxCommit failed")
		return nil, err
	}

	return &resp, nil
}

func (c *cosmosClient) QueueBroadcastMsg(msg sdk.Msg) error {
	if !c.canSign {
		return ErrReadOnly
	} else if atomic.LoadInt64(&c.closed) == 1 {
		return ErrQueueClosed
	}

	t := time.NewTimer(10 * time.Second)
	select {
	case <-t.C:
		return ErrEnqueueTimeout
	case c.msgC <- msg:
		t.Stop()
	}

	return nil
}

func (c *cosmosClient) Close() {
	if !c.canSign {
		return
	}

	if atomic.CompareAndSwapInt64(&c.closed, 0, 1) {
		close(c.msgC)
	}

	<-c.doneC
}

const (
	msgCommitBatchSizeLimit = 512
	msgCommitBatchTimeLimit = 500 * time.Millisecond
)

func (c *cosmosClient) runBatchBroadcast() {
	expirationTimer := time.NewTimer(msgCommitBatchTimeLimit)
	msgBatch := make([]sdk.Msg, 0, msgCommitBatchSizeLimit)

	resetBatch := func() {
		msgBatch = msgBatch[:0]

		expirationTimer.Reset(msgCommitBatchTimeLimit)
	}

	submitBatch := func() {
		c.syncMux.Lock()
		defer c.syncMux.Unlock()

		metrics.ReportClosureFuncCall("submitBatch", c.svcTags)
		doneFn := metrics.ReportClosureFuncTiming("submitBatch", c.svcTags)
		defer doneFn()

		ts := time.Now()

		cliCtx := c.newCliCtx()
		if err := c.buildSignBroadcast(cliCtx, msgBatch); err != nil {
			log.WithField("size", len(msgBatch)).WithError(err).Errorln("failed to commit msg batch")
			metrics.ReportClosureFuncError("submitBatch", c.svcTags)
			return
		}

		metrics.SidechainSubmitBatchSize(len(msgBatch), c.svcTags)
		log.WithField("size", len(msgBatch)).Infof("commit msg batch took %s", time.Since(ts))
	}

	for {
		select {
		case msg, ok := <-c.msgC:
			if !ok {
				// exit required
				if len(msgBatch) > 0 {
					submitBatch()
				}

				close(c.doneC)
				return
			}

			log.WithField("type", msg.Type()).Infoln("batching msg")

			msgBatch = append(msgBatch, msg)

			if len(msgBatch) >= msgCommitBatchSizeLimit {
				submitBatch()
				resetBatch()
			}
		case <-expirationTimer.C:
			if len(msgBatch) > 0 {
				submitBatch()
			}

			resetBatch()
		}
	}
}

func (c *cosmosClient) buildSignBroadcast(cliCtx context.CLIContext, msgs []sdk.Msg) error {
	txBldr, err := c.prepareTxBuilder(cliCtx, c.txBldr)
	if err != nil {
		return err
	}

	fromName := cliCtx.GetFromName()
	txBytes, err := txBldr.BuildAndSign(fromName, c.passphrase, msgs)
	if err != nil {
		err = errors.Wrap(err, "buildAndSign failed")
		return err
	}

	// broadcast to a Tendermint node
	if _, err := cliCtx.BroadcastTxCommit(txBytes); err != nil {
		c.rollbackSequence()
		err = errors.Wrap(err, "broadcastTxCommit failed")
		return err
	}

	return nil
}

func (c *cosmosClient) rollbackSequence() {
	if c.fromAccNumber != nil {
		c.fromAccSequence--
	}
}

// PrepareTxBuilder populates a TxBuilder in preparation for the build of a Tx.
func (c *cosmosClient) prepareTxBuilder(
	cliCtx context.CLIContext,
	txBldr authtypes.TxBuilder,
) (authtypes.TxBuilder, error) {
	if c.fromAccNumber != nil {
		c.fromAccSequence++
		txBldr = txBldr.WithSequence(c.fromAccSequence)

		if txBldr.AccountNumber() != *c.fromAccNumber {
			txBldr = txBldr.WithAccountNumber(*c.fromAccNumber)
		}

		return txBldr, nil
	}

	from := cliCtx.GetFromAddress()
	accGetter := authtypes.NewAccountRetriever(cliCtx)

	if err := accGetter.EnsureExists(from); err != nil {
		return txBldr, err
	}

	num, seq, err := accGetter.GetAccountNumberSequence(from)
	if err != nil {
		return txBldr, err
	}

	c.fromAccNumber = &num
	c.fromAccSequence = seq
	txBldr = txBldr.WithAccountNumber(num)
	txBldr = txBldr.WithSequence(seq)

	return txBldr, nil
}
