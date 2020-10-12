package types

import (
	"bytes"
	"errors"
	"fmt"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// NewTransactionLogs creates a new NewTransactionLogs instance.
func NewTransactionLogs(hash ethcmn.Hash, logs []*Log) TransactionLogs {
	return TransactionLogs{
		Hash: hash.String(),
		Logs: logs,
	}
}

// NewTransactionLogsFromEth creates a new NewTransactionLogs instance using []*ethtypes.Log.
func NewTransactionLogsFromEth(hash ethcmn.Hash, ethlogs []*ethtypes.Log) TransactionLogs {
	return TransactionLogs{
		Hash: hash.String(),
		// TODO: logs
	}
}

// Validate performs a basic validation of a GenesisAccount fields.
func (tx TransactionLogs) Validate() error {
	if bytes.Equal(ethcmn.Hex2Bytes(tx.Hash), ethcmn.Hash{}.Bytes()) {
		return fmt.Errorf("hash cannot be the empty %s", tx.Hash)
	}

	for i, log := range tx.Logs {
		if err := log.Validate(); err != nil {
			return fmt.Errorf("invalid log %d: %w", i, err)
		}
		if log.TxHash != tx.Hash {
			return fmt.Errorf("log tx hash mismatch (%s ≠ %s)", log.TxHash, tx.Hash)
		}
	}
	return nil
}

// Validate performs a basic validation of an ethereum Log fields.
func (log *Log) Validate() error {
	if bytes.Equal(ethcmn.HexToAddress(log.Address).Bytes(), ethcmn.Address{}.Bytes()) {
		return fmt.Errorf("log address cannot be empty %s", log.Address)
	}
	if bytes.Equal(log.BlockHash.Bytes(), ethcmn.Hash{}.Bytes()) {
		return fmt.Errorf("block hash cannot be the empty %s", log.BlockHash)
	}
	if log.BlockNumber == 0 {
		return errors.New("block number cannot be zero")
	}
	if bytes.Equal(log.TxHash.Bytes(), ethcmn.Hash{}.Bytes()) {
		return fmt.Errorf("tx hash cannot be the empty %s", log.TxHash)
	}
	return nil
}
