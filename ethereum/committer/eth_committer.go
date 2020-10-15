package committer

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"strings"
	"time"

	"github.com/InjectiveLabs/zeroex-go"
	"github.com/InjectiveLabs/zeroex-go/wrappers"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	log "github.com/xlab/suplog"

	"github.com/cosmos/ethermint/ethereum/gasmeter"
	"github.com/cosmos/ethermint/ethereum/keystore"
	"github.com/cosmos/ethermint/ethereum/provider"
	"github.com/cosmos/ethermint/ethereum/registry"
	eth "github.com/cosmos/ethermint/ethereum/util"
	"github.com/cosmos/ethermint/metrics"
)

// NewEthCommitter returns an instance of EVMCommitter, which
// can be used to submit zeroex txns into Ethereum, Matic, and other EVM-compatible networks.
func NewEthCommitter(
	keystorePath string,
	fromAddress common.Address,
	fromPassphrase string,
	fromPrivateKey *string,
	evmProvider provider.EVMProvider,
	contractSet registry.ContractsSet,
) (EVMCommitter, error) {
	ks, err := keystore.New(keystorePath)
	if err != nil {
		return nil, err
	}

	var fromKey *ecdsa.PrivateKey
	if fromPrivateKey != nil && len(*fromPrivateKey) > 0 {
		fromKey, err = crypto.HexToECDSA(*fromPrivateKey)
		if err == nil {
			fromAddress = crypto.PubkeyToAddress(fromKey.PublicKey)
			fromPassphrase = ""
			log.WithField("from_address", fromAddress.String()).Infoln("using EVM committer with private key")
		} else {
			err = errors.Wrap(err, "private key provided, but failed to parse")
			return nil, err
		}
	} else {
		log.WithField("from_address", fromAddress.String()).Infoln("using EVM committer with key from keystore")
	}

	committer := &ethCommitter{
		svcTags: metrics.Tags{
			"module": "eth_comitter",
		},

		fromAddress:    fromAddress,
		fromPassphrase: fromPassphrase,
		fromKey:        fromKey,
		contractSet:    contractSet,
		evmProvider:    evmProvider,
		keystore:       ks,
		nonceCache:     eth.NewNonceCache(),
	}

	if committer.fromKey == nil {
		// if no private key provided, committer must be able to retrieve one from keystore
		if _, ok := ks.PrivateKey(fromAddress, fromPassphrase); !ok {
			err = errors.Errorf("failed to decode Ethereum wallet '%s' using provided passphrase", fromAddress.Hex())
			return nil, err
		}
	}


	committer.coordinatorContract, err = wrappers.NewCoordinator(contractSet.CoordinatorContract, evmProvider)
	if err != nil {
		err = errors.Wrap(err, "failed to init NewCoordinator client")
		log.WithError(err).Errorln("failed to init eth committer")
	}


	committer.nonceCache.Sync(fromAddress, func() (uint64, error) {
		nonce, err := evmProvider.PendingNonceAt(context.TODO(), fromAddress)
		return nonce, err
	})

	return committer, nil
}

type ethCommitter struct {
	fromAddress    common.Address
	fromPassphrase string
	fromKey        *ecdsa.PrivateKey

	contractSet        registry.ContractsSet
	coordinatorContract *wrappers.Coordinator
	gasStation          gasmeter.GasStation
	evmProvider         provider.EVMProvider
	keystore            keystore.EthKeyStore
	nonceCache          eth.NonceCache

	svcTags metrics.Tags
}

func (t *ethCommitter) ExchangeAddress() common.Address {
	return t.contractSet.ExchangeContract
}

func (t *ethCommitter) CoordinatorAddress() common.Address {
	return t.contractSet.CoordinatorContract
}

// protocolFee0xV3 is 150,000 * ordersFilled * gasPrice
// However that initial constant might vary upon deployment.
func protocolFee0xV3(gasPrice *big.Int, numOrders int) *big.Int {
	v := big.NewInt(0).Set(gasPrice)
	v.Mul(v, big.NewInt(150000))
	v.Mul(v, big.NewInt(int64(numOrders)))
	return v
}

func (e *ethCommitter) CommitZeroExTx(
	zeroExTx *zeroex.SignedTransaction,
	approvalSignature []byte,
) (txHash common.Hash, err error) {
	metrics.ReportFuncCall(e.svcTags)
	doneFn := metrics.ReportFuncTiming(e.svcTags)
	defer doneFn()

	var signerFn bind.SignerFn

	if e.fromKey != nil {
		signerFn = eth.SignerFnForPk(e.fromKey)
	} else {
		signerFn = e.keystore.SignerFn(e.fromAddress, e.fromPassphrase)
	}

	var ordersNum int
	txData, err := zeroExTx.DecodeTransactionData()
	if err == nil {
		if len(txData.LeftOrders) > 0 {
			// NB: len(left) == len(right)
			ordersNum = len(txData.LeftOrders) + len(txData.RightOrders)
		} else {
			ordersNum = len(txData.Orders)
		}
	}

	opts := &bind.TransactOpts{
		From:     e.fromAddress,
		Signer:   signerFn,
		GasPrice: zeroExTx.GasPrice,
		GasLimit: 6000000, // todo: no hardcoding
		Value:    protocolFee0xV3(zeroExTx.GasPrice, ordersNum),
	}

	resyncNonces := func(from common.Address) {
		e.nonceCache.Sync(from, func() (uint64, error) {
			nonce, err := e.evmProvider.PendingNonceAt(context.TODO(), from)
			if err != nil {
				log.WithError(err).Warningln("unable to acquire nonce")
			}

			return nonce, err
		})
	}

	if err := e.nonceCache.Serialize(e.fromAddress, func() (err error) {
		nonce := e.nonceCache.Incr(e.fromAddress)
		var resyncUsed bool

		for {
			opts.Nonce = big.NewInt(nonce)
			opts.Context, _ = context.WithTimeout(context.Background(), 20*time.Second)

			zeroExTxArg := wrappers.ZeroExTransaction{
				Salt:                  zeroExTx.Salt,
				ExpirationTimeSeconds: zeroExTx.ExpirationTimeSeconds,
				GasPrice:              zeroExTx.GasPrice,
				SignerAddress:         zeroExTx.SignerAddress,
				Data:                  zeroExTx.Data,
			}
			tx, err := e.coordinatorContract.ExecuteTransaction(opts,
				zeroExTxArg,
				zeroExTx.SignerAddress,
				zeroExTx.Signature,
				[][]byte{approvalSignature},
			)
			if err == nil {
				txHash = tx.Hash()
				return nil
			} else {
				log.WithError(err).Warningln("ExecuteTransaction failed with error")
			}

			switch {
			case strings.Contains(err.Error(), "invalid sender"):
				e.nonceCache.Decr(e.fromAddress)

				err := errors.New("failed to sign transaction")

				return err
			case strings.Contains(err.Error(), "nonce is too low"),
				strings.Contains(err.Error(), "nonce is too high"),
				strings.Contains(err.Error(), "the tx doesn't have the correct nonce"):

				if resyncUsed {
					log.Errorf("nonces synced, but still wrong nonce for %s: %d", e.fromAddress, nonce)
					err = errors.Wrapf(err, "nonce %d mismatch", nonce)

					return err
				}

				resyncNonces(e.fromAddress)

				resyncUsed = true
				// try again with new nonce
				nonce = e.nonceCache.Incr(e.fromAddress)
				opts.Nonce = big.NewInt(nonce)

				continue

			default:
				if strings.Contains(err.Error(), "known transaction") {
					// skip one nonce step, try to send again
					nonce := e.nonceCache.Incr(e.fromAddress)
					opts.Nonce = big.NewInt(nonce)
					continue
				}

				if strings.Contains(err.Error(), "VM Exception") {
					// a VM execution consumes gas and nonce is increasing
					return err
				}

				return err
			}
		}
	}); err != nil {
		metrics.ReportFuncError(e.svcTags)

		return common.Hash{}, err
	}

	return txHash, nil
}

func (e *ethCommitter) CommitFuturesTx(
	txData []byte,
) (txHash common.Hash, err error) {
	metrics.ReportFuncCall(e.svcTags)
	doneFn := metrics.ReportFuncTiming(e.svcTags)
	defer doneFn()

	var signerFn bind.SignerFn

	if e.fromKey != nil {
		signerFn = eth.SignerFnForPk(e.fromKey)
	} else {
		signerFn = e.keystore.SignerFn(e.fromAddress, e.fromPassphrase)
	}

	opts := &bind.TransactOpts{
		From:     e.fromAddress,
		Signer:   signerFn,
		GasPrice: big.NewInt(1000000000), // todo: no hardcoding
		GasLimit: 6000000,                // todo: no hardcoding
	}

	resyncNonces := func(from common.Address) {
		e.nonceCache.Sync(from, func() (uint64, error) {
			nonce, err := e.evmProvider.PendingNonceAt(context.TODO(), from)
			if err != nil {
				log.WithError(err).Warningln("unable to acquire nonce")
			}

			return nonce, err
		})
	}

	if err := e.nonceCache.Serialize(e.fromAddress, func() (err error) {
		nonce := e.nonceCache.Incr(e.fromAddress)
		var resyncUsed bool

		for {
			opts.Nonce = big.NewInt(nonce)
			opts.Context, _ = context.WithTimeout(context.Background(), 20*time.Second)

			tx := types.NewTransaction(
				opts.Nonce.Uint64(),
				e.contractSet.FuturesContract,
				nil, opts.GasLimit, opts.GasPrice, txData)
			signedTx, err := opts.Signer(types.HomesteadSigner{}, opts.From, tx)
			if err != nil {
				e.nonceCache.Decr(e.fromAddress)

				err := errors.Wrap(err, "failed to sign transaction")

				return err
			}

			txHash = signedTx.Hash()

			err = e.evmProvider.SendTransaction(opts.Context, signedTx)
			if err == nil {
				return nil
			} else {
				log.WithField("txHash", txHash.Hex()).WithError(err).Warningln("SendTransaction failed with error")
			}

			switch {
			case strings.Contains(err.Error(), "invalid sender"):
				e.nonceCache.Decr(e.fromAddress)

				err := errors.New("failed to sign transaction")

				return err
			case strings.Contains(err.Error(), "nonce is too low"),
				strings.Contains(err.Error(), "nonce is too high"),
				strings.Contains(err.Error(), "the tx doesn't have the correct nonce"):

				if resyncUsed {
					log.Errorf("nonces synced, but still wrong nonce for %s: %d", e.fromAddress, nonce)
					err = errors.Wrapf(err, "nonce %d mismatch", nonce)
					return err
				}

				resyncNonces(e.fromAddress)

				resyncUsed = true
				// try again with new nonce
				nonce = e.nonceCache.Incr(e.fromAddress)
				opts.Nonce = big.NewInt(nonce)

				continue

			default:
				if strings.Contains(err.Error(), "known transaction") {
					// skip one nonce step, try to send again
					nonce := e.nonceCache.Incr(e.fromAddress)
					opts.Nonce = big.NewInt(nonce)
					continue
				}

				if strings.Contains(err.Error(), "VM Exception") {
					// a VM execution consumes gas and nonce is increasing
					return err
				}

				e.nonceCache.Decr(e.fromAddress)

				return err
			}
		}
	}); err != nil {
		metrics.ReportFuncError(e.svcTags)

		return common.Hash{}, err
	}

	return txHash, nil
}
