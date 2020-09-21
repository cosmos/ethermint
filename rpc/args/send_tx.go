package args

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// SendTxArgs represents the arguments to submit a new transaction into the transaction pool.
// Duplicate struct definition since geth struct is in internal package
// Ref: https://github.com/ethereum/go-ethereum/blob/release/1.9/internal/ethapi/api.go#L1346
type SendTxArgs struct {
	From     common.Address  `json:"from"`
	To       *common.Address `json:"to"`
	Gas      *hexutil.Uint64 `json:"gas"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
	Value    *hexutil.Big    `json:"value"`
	Nonce    *hexutil.Uint64 `json:"nonce"`
	// We accept "data" and "input" for backwards-compatibility reasons. "input" is the
	// newer name and should be preferred by clients.
	Data  *hexutil.Bytes `json:"data"`
	Input *hexutil.Bytes `json:"input"`
}

func (args SendTxArgs) String() string {
	if args.Data != nil {
		return fmt.Sprintf("to=0x%x from=0x%x nonce=%d gasPrice=%d gasLimit=%d value=%d data=0x%x",
			*args.To, args.From, args.Nonce, args.GasPrice, args.Gas, args.Value, args.Data)
	}

	return fmt.Sprintf("to=0x%x from=0x%x nonce=%d gasPrice=%d gasLimit=%d value=%d input=0x%x",
		*args.To, args.From, args.Nonce, args.GasPrice, args.Gas, args.Value, args.Input)
}