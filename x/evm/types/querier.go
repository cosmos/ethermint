package types

import (
	"math/big"

	ethcmn "github.com/ethereum/go-ethereum/common"
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
	Storage ethcmn.Hash
}

func (q QueryResStorage) String() string {
	return q.Storage.String()
}

type QueryResCode struct {
	Code []byte
}

func (q QueryResCode) String() string {
	return string(q.Code)
}
