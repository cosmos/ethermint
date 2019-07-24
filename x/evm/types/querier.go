package types

import "math/big"

type QueryResBlockNumber struct {
	Number *big.Int
}

func (q QueryResBlockNumber) String() string {
	return q.Number.String()
}
