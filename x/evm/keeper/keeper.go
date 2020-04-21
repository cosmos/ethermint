package keeper

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/ethermint/x/evm/types"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"math/big"
)

// Keeper wraps the CommitStateDB, allowing us to pass in SDK context while adhering
// to the StateDB interface.
type Keeper struct {
	// Amino codec
	cdc *codec.Codec
	// Store key required to update the block bloom filter mappings needed for the
	// Web3 API
	blockKey      sdk.StoreKey
	CommitStateDB *types.CommitStateDB
	TxCount       int
	Bloom         *big.Int
}

// NewKeeper generates new evm module keeper
func NewKeeper(
	cdc *codec.Codec, blockKey, codeKey, storeKey sdk.StoreKey,
	ak types.AccountKeeper,
) Keeper {
	return Keeper{
		cdc:           cdc,
		blockKey:      blockKey,
		CommitStateDB: types.NewCommitStateDB(sdk.Context{}, codeKey, storeKey, ak),
		TxCount:       0,
		Bloom:         big.NewInt(0),
	}
}

// ----------------------------------------------------------------------------
// Block hash mapping functions
// May be removed when using only as module (only required by rpc api)
// ----------------------------------------------------------------------------

// SetBlockHashMapping sets the mapping from block consensus hash to block height
func (k *Keeper) SetBlockHashMapping(ctx sdk.Context, hash []byte, height int64) {
	store := ctx.KVStore(k.blockKey)
	if !bytes.Equal(hash, []byte{}) {
		bz := sdk.Uint64ToBigEndian(uint64(height))
		store.Set(hash, bz)
	}
}

// GetBlockHashMapping gets block height from block consensus hash
func (k *Keeper) GetBlockHashMapping(ctx sdk.Context, hash []byte) (height int64) {
	store := ctx.KVStore(k.blockKey)
	bz := store.Get(hash)
	if bytes.Equal(bz, []byte{}) {
		panic(fmt.Errorf("block with hash %s not found", ethcmn.BytesToHash(hash)))
	}

	height = int64(binary.BigEndian.Uint64(bz))
	return height
}

// ----------------------------------------------------------------------------
// Block bloom bits mapping functions
// May be removed when using only as module (only required by rpc api)
// ----------------------------------------------------------------------------

// SetBlockBloomMapping sets the mapping from block height to bloom bits
func (k *Keeper) SetBlockBloomMapping(ctx sdk.Context, bloom ethtypes.Bloom, height int64) error {
	store := ctx.KVStore(k.blockKey)
	bz := sdk.Uint64ToBigEndian(uint64(height))
	if len(bz) == 0 {
		return fmt.Errorf("block with bloombits %v not found", bloom)
	}

	store.Set(types.BloomKey(bz), bloom.Bytes())
	return nil
}

// GetBlockBloomMapping gets bloombits from block height
func (k *Keeper) GetBlockBloomMapping(ctx sdk.Context, height int64) (ethtypes.Bloom, error) {
	store := ctx.KVStore(k.blockKey)
	bz := sdk.Uint64ToBigEndian(uint64(height))
	bloom := store.Get(types.BloomKey(bz))
	if len(bloom) == 0 {
		return ethtypes.BytesToBloom([]byte{}), fmt.Errorf("block with height %d not found", height)
	}

	return ethtypes.BytesToBloom(bloom), nil
}

// SetTransactionLogs sets the transaction's logs in the KVStore
func (k *Keeper) SetTransactionLogs(ctx sdk.Context, hash []byte, logs []*ethtypes.Log) {
	store := ctx.KVStore(k.blockKey)
	if logs == nil || len(logs) == 0 {
		panic("cannot set empty or nil logs")
	}
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(logs)
	store.Set(types.LogsKey(hash), bz)
}

// GetTransactionLogs gets the logs for a transaction from the KVStore
func (k *Keeper) GetTransactionLogs(ctx sdk.Context, hash []byte) ([]*ethtypes.Log, error) {
	store := ctx.KVStore(k.blockKey)
	bz := store.Get(types.LogsKey(hash))
	if len(bz) == 0 {
		return nil, errors.New("cannot get transaction logs")
	}

	var logs []*ethtypes.Log
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &logs)
	return logs, nil
}
