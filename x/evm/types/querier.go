package types

import (
	"math/big"
)

type QueryResBalance struct {
	Balance *big.Int
}

func (q QueryResBalance) String() string {
	return q.Balance.String()
}

type QueryResBlockNumber struct {
	Number *big.Int
}

func (q QueryResBlockNumber) String() string {
	return q.Number.String()
}

type QueryResStorage struct {
	Value []byte `json:"value"`
}

func (q QueryResStorage) String() string {
	return string(q.Value)
}

type QueryResCode struct {
	Code []byte
}

func (q QueryResCode) String() string {
	return string(q.Code)
}
