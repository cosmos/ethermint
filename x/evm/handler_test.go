package evm

import (
	"math/big"
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ethermint/crypto"
	"github.com/cosmos/ethermint/x/evm/types"
	eminttypes 	"github.com/cosmos/ethermint/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"
)

func TestHandler_Logs(t *testing.T) {
	// create logger, codec and root multi-store
	cdc := newTestCodec()

	// The ParamsKeeper handles parameter storage for the application
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	paramsKeeper := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	// Set specific supspaces
	authSubspace := paramsKeeper.Subspace(auth.DefaultParamspace)
	ak := auth.NewAccountKeeper(cdc, accKey, authSubspace, eminttypes.ProtoBaseAccount)
	ek := NewKeeper(ak, storageKey, codeKey, blockKey, cdc)

	gasLimit := uint64(100000)
	gasPrice := big.NewInt(1000000)
	address := ethcmn.BytesToAddress([]byte{0})

	priv1, _ := crypto.GenerateKey()

	tx := types.NewEthereumTxMsg(1, &address, big.NewInt(0), gasLimit, gasPrice, []byte{})
	tx.Sign(big.NewInt(1), priv1.ToECDSA())

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	// mount stores
	keys := []*sdk.KVStoreKey{accKey, storageKey, codeKey, blockKey}
	for _, key := range keys {
		cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, nil)
	}

	err := cms.LoadLatestVersion()
	if err != nil {
		t.Fatal(err)
	}
	
	ms := cms.CacheMultiStore()
	ctx := sdk.NewContext(ms, abci.Header{}, false, logger)
	ctx = ctx.WithBlockHeight(1).WithChainID("1")

	result := handleETHTxMsg(ctx, ek, tx)
	t.Log(result)
}