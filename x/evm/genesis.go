package evm

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ethermint "github.com/cosmos/ethermint/types"
	"github.com/cosmos/ethermint/x/evm/keeper"
	"github.com/cosmos/ethermint/x/evm/types"

	ethcmn "github.com/ethereum/go-ethereum/common"

	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis initializes genesis state based on exported genesis
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) []abci.ValidatorUpdate {
	for _, account := range data.Accounts {
		// FIXME: this will override bank InitGenesis balance!

		address := ethcmn.HexToAddress(account.Address)
		balance := new(big.Int).SetBytes(account.Balance)
		k.SetBalance(ctx, address, balance)
		k.SetCode(ctx, address, account.Code)

		for _, storage := range account.Storage {
			k.SetState(ctx, address, ethcmn.HexToHash(storage.Key), ethcmn.HexToHash(storage.Value))
		}
	}

	var err error
	for _, txLog := range data.TxsLogs {
		err = k.SetLogs(ctx, ethcmn.HexToHash(txLog.Hash), txLog.EthLogs())
		if err != nil {
			panic(err)
		}
	}

	k.SetChainConfig(ctx, *data.ChainConfig)
	k.SetParams(ctx, data.Params)

	// set state objects and code to store
	_, err = k.Commit(ctx, false)
	if err != nil {
		panic(err)
	}

	// set storage to store
	// NOTE: don't delete empty object to prevent import-export simulation failure
	err = k.Finalise(ctx, false)
	if err != nil {
		panic(err)
	}

	return []abci.ValidatorUpdate{}
}

// ExportGenesis exports genesis state of the EVM module
func ExportGenesis(ctx sdk.Context, k keeper.Keeper, ak types.AccountKeeper) *types.GenesisState {
	// nolint: prealloc
	var ethGenAccounts []types.GenesisAccount
	accounts := ak.GetAllAccounts(ctx)

	for _, account := range accounts {

		ethAccount, ok := account.(*ethermint.EthAccount)
		if !ok {
			continue
		}

		addr := ethAccount.EthAddress()

		storage, err := k.GetAccountStorage(ctx, addr)
		if err != nil {
			panic(err)
		}

		genAccount := types.GenesisAccount{
			Address: addr.String(),
			Balance: k.GetBalance(ctx, addr).Bytes(),
			Code:    k.GetCode(ctx, addr),
			Storage: storage,
		}

		ethGenAccounts = append(ethGenAccounts, genAccount)
	}

	config, _ := k.GetChainConfig(ctx)

	return &types.GenesisState{
		Accounts:    ethGenAccounts,
		TxsLogs:     k.GetAllTxLogs(ctx),
		ChainConfig: &config,
		Params:      k.GetParams(ctx),
	}
}
