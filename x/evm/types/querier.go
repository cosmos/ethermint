package types

import "math/big"

type QueryResBlockNumber struct {
	number *big.Int
}

func (q QueryResBlockNumber) String() string {
	return q.number.String()
}
