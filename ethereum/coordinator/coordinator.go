package coordinator

import (
	"context"
	"time"

	"github.com/InjectiveLabs/zeroex-go"
	"github.com/ethereum/go-ethereum/common"
)

type SpotOrderCoordinator interface {
	ApproveTransaction(
		ctx context.Context,
		tx *zeroex.SignedTransaction,
		txOrigin common.Address,
		deadlineTimestamp time.Time,
	) (approval *zeroex.SignedCoordinatorApproval, err error)
}
