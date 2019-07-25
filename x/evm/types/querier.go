package types

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
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
	Value ethcmn.Hash `json:"value"`
}

func (q QueryResStorage) String() string {
	return q.Value.String()
}

type QueryResCode struct {
	Code []byte
}

func (q QueryResCode) String() string {
	return string(q.Code)
}
