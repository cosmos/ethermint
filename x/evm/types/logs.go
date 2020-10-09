package types

import (
	"bytes"
	"errors"
	"fmt"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// NewTransactionLogs creates a new NewTransactionLogs instance.
func NewTransactionLogs(hash ethcmn.Hash, ethlogs []*ethtypes.Log) TransactionLogs {
	return TransactionLogs{
		Hash: hash,
		// TODO: logs
	}
}

// Validate performs a basic validation of a GenesisAccount fields.
func (tx TransactionLogs) Validate() error {
	if bytes.Equal(tx.Hash.Bytes(), ethcmn.Hash{}.Bytes()) {
		return fmt.Errorf("hash cannot be the empty %s", tx.Hash.String())
	}

	for i, log := range tx.Logs {
		if err := ValidateLog(log); err != nil {
			return fmt.Errorf("invalid log %d: %w", i, err)
		}
		if !bytes.Equal(log.TxHash.Bytes(), tx.Hash.Bytes()) {
			return fmt.Errorf("log tx hash mismatch (%s ≠ %s)", log.TxHash.String(), tx.Hash.String())
		}
	}
	return nil
}

// Validate performs a basic validation of an ethereum Log fields.
func (log *Log) Validate() error {
	if bytes.Equal(ethcmn.HexToAddress(log.Address).Bytes(), ethcmn.Address{}.Bytes()) {
		return fmt.Errorf("log address cannot be empty %s", log.Address.String())
	}
	if bytes.Equal(log.BlockHash.Bytes(), ethcmn.Hash{}.Bytes()) {
		return fmt.Errorf("block hash cannot be the empty %s", log.BlockHash.String())
	}
	if log.BlockNumber == 0 {
		return errors.New("block number cannot be zero")
	}
	if bytes.Equal(log.TxHash.Bytes(), ethcmn.Hash{}.Bytes()) {
		return fmt.Errorf("tx hash cannot be the empty %s", log.TxHash.String())
	}
	return nil
}
