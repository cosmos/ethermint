package types

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"
	tmlog "github.com/tendermint/tendermint/libs/log"
	tmdb "github.com/tendermint/tm-db"

	sdkcodec "github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"

	ethcmn "github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/ethermint/codec"
	"github.com/cosmos/ethermint/crypto"
	"github.com/cosmos/ethermint/types"
)

type JournalTestSuite struct {
	suite.Suite

	address ethcmn.Address
	journal *journal
	ctx     sdk.Context
	stateDB *CommitStateDB
}

func newTestCodec() *sdkcodec.Codec {
	cdc := sdkcodec.New()

	RegisterCodec(cdc)
	types.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	crypto.RegisterCodec(cdc)
	sdkcodec.RegisterCrypto(cdc)

	return cdc
}

func (suite *JournalTestSuite) SetupTest() {
	privkey, err := crypto.GenerateKey()
	suite.Require().NoError(err)

	suite.address = ethcmn.BytesToAddress(privkey.PubKey().Address().Bytes())
	suite.journal = newJournal()

	db := tmdb.NewDB("state", tmdb.GoLevelDBBackend, "temp")
	defer func() {
		os.RemoveAll("temp")
	}()

	cms := store.NewCommitMultiStore(db)

	// The ParamsKeeper handles parameter storage for the application
	bankKey := sdk.NewKVStoreKey(bank.StoreKey)
	authKey := sdk.NewKVStoreKey(auth.StoreKey)
	storeKey := sdk.NewKVStoreKey(StoreKey)

	// mount stores
	keys := []*sdk.KVStoreKey{authKey, bankKey, storeKey}
	for _, key := range keys {
		cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, nil)
	}

	cms.SetPruning(store.PruneNothing)

	// load latest version (root)
	err = cms.LoadLatestVersion()
	suite.Require().NoError(err)

	cdc := newTestCodec()
	appCodec := codec.NewAppCodec(cdc)
	authclient.Codec = appCodec

	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	paramsKeeper := params.NewKeeper(appCodec, keyParams, tkeyParams)
	// Set specific supspaces
	authSubspace := paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := paramsKeeper.Subspace(bank.DefaultParamspace)
	ak := auth.NewAccountKeeper(appCodec, authKey, authSubspace, types.ProtoAccount)
	bk := bank.NewBaseKeeper(appCodec, bankKey, ak, bankSubspace, nil)

	ms := cms.CacheMultiStore()
	suite.ctx = sdk.NewContext(ms, abci.Header{}, false, tmlog.NewNopLogger())

	suite.stateDB = NewCommitStateDB(suite.ctx, storeKey, ak, bk)

	// acc := types.EthAccount{
	// 	BaseAccount: auth.NewBaseAccount(sdk.AccAddress(suite.address.Bytes()), nil, 0, 0),
	// 	CodeHash:    ethcrypto.Keccak256(nil),
	// }

	// ak.SetAccount(suite.ctx, acc)
	bk.SetBalance(suite.ctx, sdk.AccAddress(suite.address.Bytes()), sdk.NewCoin(types.DenomDefault, sdk.NewInt(100)))

}

func TestJournalTestSuite(t *testing.T) {
	suite.Run(t, new(JournalTestSuite))
}

func (suite *JournalTestSuite) TestJournal_append_revert() {
	testCases := []struct {
		name  string
		entry journalEntry
	}{
		{
			"createObjectChange",
			createObjectChange{
				account: &suite.address,
			},
		},
		{
			"resetObjectChange",
			resetObjectChange{
				prev: &stateObject{
					address: suite.address,
					balance: sdk.OneInt(),
				},
			},
		},
		{
			"suicideChange",
			suicideChange{
				account:     &suite.address,
				prev:        false,
				prevBalance: sdk.OneInt(),
			},
		},
		{
			"balanceChange",
			balanceChange{
				account: &suite.address,
				prev:    sdk.OneInt(),
			},
		},
		{
			"nonceChange",
			nonceChange{
				account: &suite.address,
				prev:    1,
			},
		},
		{
			"storageChange",
			storageChange{
				account:   &suite.address,
				key:       ethcmn.BytesToHash([]byte("key")),
				prevValue: ethcmn.BytesToHash([]byte("value")),
			},
		},
		{
			"codeChange",
			codeChange{
				account:  &suite.address,
				prevCode: []byte("code"),
				prevHash: []byte("hash"),
			},
		},
		{
			"touchChange",
			touchChange{
				account: &suite.address,
			},
		},
		{
			"refundChange",
			refundChange{
				prev: 1,
			},
		},
		{
			"addPreimageChange",
			addPreimageChange{
				hash: ethcmn.BytesToHash([]byte("hash")),
			},
		},
	}
	var dirtyCount int
	for i, tc := range testCases {
		suite.journal.append(tc.entry)
		suite.Require().Equal(suite.journal.length(), i+1, tc.name)
		if tc.entry.dirtied() != nil {
			dirtyCount++
			suite.Require().Equal(dirtyCount, suite.journal.dirties[suite.address], tc.name)
		}
	}

	// for i, tc := range testCases {
	// suite.journal.revert(suite.stateDB, len(testCases)-1-i)
	// suite.Require().Equal(suite.journal.length(), len(testCases)-1-i, tc.name)
	// if tc.entry.dirtied() != nil {
	// 	dirtyCount--
	// 	suite.Require().Equal(dirtyCount, suite.journal.dirties[suite.address], tc.name)
	// }
	// }

	// verify the dirty entry
	// count, ok := suite.journal.dirties[suite.address]
	// suite.Require().False(ok)
	// suite.Require().Zero(count)
}

func (suite *JournalTestSuite) TestJournal_dirty() {
	// dirty entry hasn't been set
	count, ok := suite.journal.dirties[suite.address]
	suite.Require().False(ok)
	suite.Require().Zero(count)

	// update dirty count
	suite.journal.dirty(suite.address)
	suite.Require().Equal(1, suite.journal.dirties[suite.address])
}
