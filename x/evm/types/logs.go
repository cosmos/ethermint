package types

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// TransactionLogs define the logs generated from a transaction execution
// with a given hash. It it used for import/export data as transactions are not persisted
// on blockchain state after an upgrade.
type TransactionLogs struct {
	Hash ethcmn.Hash     `json:"hash"`
	Logs []*ethtypes.Log `json:"logs"`
}

// NewTransactionLog creates a new TransactionLog instance.
func NewTransactionLog(hash ethcmn.Hash, logs []*ethtypes.Log) TransactionLogs {
	return TransactionLogs{
		Hash: hash,
		Logs: logs,
	}
}

// MarshalLogs encodes an array of logs using amino
func MarshalLogs(logs []*ethtypes.Log) ([]byte, error) {
	return ModuleCdc.MarshalBinaryLengthPrefixed(logs)
}

// UnmarshalLogs decodes an amino-encoded byte array into an array of logs
func UnmarshalLogs(in []byte) ([]*ethtypes.Log, error) {
	logs := []*ethtypes.Log{}
	err := ModuleCdc.UnmarshalBinaryLengthPrefixed(in, &logs)
	return logs, err
}
