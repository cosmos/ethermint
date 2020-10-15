package orders

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/ethermint/ethereum/provider"
	"github.com/cosmos/ethermint/ethereum/registry"
	"github.com/cosmos/ethermint/eventdb"
	"github.com/cosmos/ethermint/loopback"
	"github.com/cosmos/ethermint/x/orders/internal/types"
	"github.com/cosmos/ethermint/metrics"
)

// type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// app module Basics object
type AppModuleBasic struct{}

func (AppModuleBasic) Name() string {
	return ModuleName
}

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// Validation check of the Genesis
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	// Once json successfully marshalled, passes along to genesis.go
	return ValidateGenesis(data)
}

// Register rest routes
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	return
}

// Get the root query command of this module
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return nil
}

// Get the root tx command of this module
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return nil
}

type AppModule struct {
	AppModuleBasic

	svcTags      metrics.Tags
	keeper       Keeper
	cosmosClient loopback.CosmosClient
	isExportOnly bool

	ethOrderEventDB           eventdb.OrderEventDB
	ethFuturesPositionEventDB eventdb.FuturesPositionEventDB

	ethProvider          func() provider.EVMProvider
	ethProviderStreaming func() provider.EVMProvider
	ethContracts         registry.ContractDiscoverer

	evmSyncStatus        types.EvmSyncStatus
	futuresEvmSyncStatus types.EvmSyncStatus
}

// NewAppModule creates a new AppModule Object
func NewAppModule(
	keeper Keeper,
	isExportOnly bool,
	cosmosClient loopback.CosmosClient,
	ethOrderEventDB eventdb.OrderEventDB,
	ethFuturesPositionEventDB eventdb.FuturesPositionEventDB,
	ethProvider func() provider.EVMProvider,
	ethProviderStreaming func() provider.EVMProvider,
	ethContracts registry.ContractDiscoverer,
) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},

		svcTags: metrics.Tags{
			"svc": "orders_m",
		},
		keeper:       keeper,
		isExportOnly: isExportOnly,
		cosmosClient: cosmosClient,

		ethOrderEventDB:           ethOrderEventDB,
		ethFuturesPositionEventDB: ethFuturesPositionEventDB,

		ethProvider:          ethProvider,
		ethProviderStreaming: ethProviderStreaming,
		ethContracts:         ethContracts,
	}
}

func (AppModule) Name() string {
	return ModuleName
}

func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

func (am AppModule) Route() string {
	return RouterKey
}

func (am AppModule) NewHandler() sdk.Handler {
	return NewOrderMsgHandler(
		am.keeper,
		am.isExportOnly,
		am.ethOrderEventDB,
		am.ethFuturesPositionEventDB,
		am.ethProvider,
		am.ethContracts,
	)
}

func (am AppModule) QuerierRoute() string {
	return ModuleName
}

func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
	// am.BeginBlocker(ctx)
}

func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}
