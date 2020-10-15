package committer

import (
	"github.com/ethereum/go-ethereum/common"

	zeroex "github.com/InjectiveLabs/zeroex-go"
)

// EVMCommitter defines an interface for submitting zeroex and futures transactions
// into Ethereum, Matic, and other EVM-compatible networks.
type EVMCommitter interface {
	CoordinatorAddress() common.Address
	ExchangeAddress() common.Address
	CommitZeroExTx(tx *zeroex.SignedTransaction, approvalSignature []byte) (txHash common.Hash, err error)
	CommitFuturesTx(txData []byte) (txHash common.Hash, err error)
}
