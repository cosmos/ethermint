package coordinator

import (
	"bytes"
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	//log "github.com/xlab/suplog"

	"github.com/cosmos/ethermint/ethereum/registry"
	"github.com/cosmos/ethermint/metrics"
	"github.com/InjectiveLabs/zeroex-go"
)

func NewSpotOrderCoordinator(
	coordinatorFrom common.Address,
	contractSet registry.ContractsSet,
	ethSigner zeroex.Signer,
) SpotOrderCoordinator {
	return &ethCoordinator{
		coordinatorFrom: coordinatorFrom,
		contractSet:     contractSet,
		ethSigner:       ethSigner,

		svcTags: metrics.Tags{
			"module": "spot_coordinator",
		},
	}
}

type ethCoordinator struct {
	//abciClient      abci.Client
	coordinatorFrom common.Address
	contractSet     registry.ContractsSet
	ethSigner       zeroex.Signer

	svcTags metrics.Tags
}

func (e *ethCoordinator) ApproveTransaction(
	ctx context.Context,
	tx *zeroex.SignedTransaction,
	txOrigin common.Address,
	deadlineTimestamp time.Time,
) (*zeroex.SignedCoordinatorApproval, error) {
	metrics.ReportFuncCall(e.svcTags)
	doneFn := metrics.ReportFuncTiming(e.svcTags)
	defer doneFn()

	txData, err := tx.DecodeTransactionData()
	if err != nil {
		err = errors.Wrap(err, "failed to decode tx data")
		metrics.ReportFuncError(e.svcTags)
		return nil, err
	}

	var ordersToCheck []*zeroex.Order

	if txData.FunctionName == zeroex.BatchMatchOrdersWithMaximalFill {
		if txData.LeftOrders == nil {
			err = errors.New("incorrect parameters passed into batchMatchOrdersWithMaximalFill call: leftOrders not present")
			metrics.ReportFuncError(e.svcTags)
			return nil, err
		} else if len(txData.LeftOrders) != len(txData.RightOrders) {
			err = errors.New("incorrect parameters passed into batchMatchOrdersWithMaximalFill call: leftOrders different length than rightOrders")
			metrics.ReportFuncError(e.svcTags)
			return nil, err
		} else if len(txData.RightOrders) != len(txData.RightSignatures) {
			err = errors.New("incorrect parameters passed into batchMatchOrdersWithMaximalFill call: rightOrders different length than rightSignatures")
			metrics.ReportFuncError(e.svcTags)
			return nil, err
		} else if len(txData.RightSignatures) != len(txData.LeftSignatures) {
			err = errors.New("incorrect parameters passed into batchMatchOrdersWithMaximalFill call: leftSignatures different length than rightSignatures")
			metrics.ReportFuncError(e.svcTags)
			return nil, err
		}

		// NOTE: RightOrders contain the same limit order, so we only check the first one
		ordersToCheck = append(txData.LeftOrders, txData.RightOrders[0])
	} else if txData.LeftOrders != nil || txData.RightOrders != nil || txData.LeftSignatures != nil || txData.RightSignatures != nil {
		err = errors.New("only batchMatchOrdersWithMaximalFill is supported, batchMatchOrders is not supported")
		metrics.ReportFuncError(e.svcTags)
		return nil, err
	}

	if e.hasCancelledOrders(ctx, ordersToCheck) {
		err = errors.New("transaction contains soft-cancelled orders")
		metrics.ReportFuncError(e.svcTags)
		return nil, err
	} else if time.Unix(tx.ExpirationTimeSeconds.Int64(), 0).Before(deadlineTimestamp) {
		err = errors.New("transaction expiration time has passed already")
		metrics.ReportFuncError(e.svcTags)
		return nil, err
	}

	for _, order := range ordersToCheck {
		// if !bytes.Equal(order.FeeRecipientAddress.Bytes(), e.coordinatorFrom.Bytes()) {
		// 	err = errors.New("transaction contains orders with other fee recipients")
		// 	return nil, err
		// }

		if time.Unix(order.ExpirationTimeSeconds.Int64(), 0).Before(deadlineTimestamp) {
			err = errors.New("transaction contains expired orders")
			metrics.ReportFuncError(e.svcTags)
			return nil, err
		}

		if !bytes.Equal(order.ExchangeAddress.Bytes(), e.contractSet.ExchangeContract.Bytes()) {
			err = errors.New("transaction contains orders to different exchange")
			metrics.ReportFuncError(e.svcTags)
			return nil, err
		}
	}

	txHash, err := tx.ComputeTransactionHash()
	if err != nil {
		err = errors.Wrap(err, "failed to compute zeroex tx hash")
		metrics.ReportFuncError(e.svcTags)
		return nil, err
	}

	approval := &zeroex.CoordinatorApproval{
		TxOrigin:             txOrigin,
		TransactionHash:      txHash,
		TransactionSignature: tx.Signature,
		Domain: zeroex.EIP712Domain{
			VerifyingContract: e.contractSet.CoordinatorContract,
			ChainID:           tx.Domain.ChainID,
		},
	}

	signedApproval, err := zeroex.SignCoordinatorApproval(
		e.coordinatorFrom,
		e.ethSigner,
		approval,
	)
	if err != nil {
		err = errors.Wrap(err, "failed to sign approval")
		metrics.ReportFuncError(e.svcTags)
		return nil, err
	}

	return signedApproval, nil
}

func (e *ethCoordinator) hasCancelledOrders(ctx context.Context, orders []*zeroex.Order) bool {
	metrics.ReportFuncCall(e.svcTags)
	doneFn := metrics.ReportFuncTiming(e.svcTags)
	defer doneFn()

	orderHashes := make([]common.Hash, 0, len(orders))
	for _, order := range orders {
		orderHash, _ := order.ComputeOrderHash()
		orderHashes = append(orderHashes, orderHash)
	}

	//cancelledOrders, err := e.abciClient.FindSoftCancelledOrders(ctx, orderHashes)
	//if err != nil {
	//	if err == abci.ErrObjectNotFound {
	//		return false
	//	}
	//
	//	log.WithError(err).Warningln("failed to find soft-cancelled orders")
	//	metrics.ReportFuncError(e.svcTags)
	//	return false
	//}

	//return len(cancelledOrders) > 0
	return false
}
