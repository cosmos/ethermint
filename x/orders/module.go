package orders

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/ethermint/x/orders/client/cli"
	"github.com/cosmos/ethermint/x/orders/types"
	grpc "github.com/gogo/protobuf/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/ethermint/ethereum/provider"
	"github.com/cosmos/ethermint/ethereum/registry"
	"github.com/cosmos/ethermint/eventdb"
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

// RegisterLegacyAminoCodec performs a no-op as the evm module doesn't support amino.
func (AppModuleBasic) RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {
}

// RegisterInterfaces registers interfaces and implementations of the orders module.
func (AppModuleBasic) RegisterInterfaces(interfaceRegistry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(interfaceRegistry)
}

// DefaultGenesis returns default genesis state as raw bytes for the orders
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

func (b AppModuleBasic) ValidateGenesis(cdc codec.JSONMarshaler, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genesisState types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genesisState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}

	return genesisState.Validate()
}

// RegisterRESTRoutes performs a no-op as the orders module doesn't expose REST
// endpoints
func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr *mux.Router) {
	return
}

// RegisterGRPCRoutes registers the gRPC Gateway routes for the orders module.
func (AppModuleBasic) RegisterGRPCRoutes(clientCtx client.Context, serveMux *runtime.ServeMux) {
	err := types.RegisterQueryHandlerClient(context.Background(), serveMux, types.NewQueryClient(clientCtx))
	if err != nil {
		panic("Failed to RegisterGRPCRoutes in orders module")
	}
}

// GetTxCmd returns the root tx command for the orders module.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

// GetQueryCmd returns no root query command for the orders module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

type AppModule struct {
	AppModuleBasic

	svcTags      metrics.Tags
	keeper       Keeper
	//cosmosClient loopback.CosmosClient
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
	//cosmosClient loopback.CosmosClient,
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
		//cosmosClient: cosmosClient,

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

func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.RouterKey, am.NewHandler())
}

func (am AppModule) QuerierRoute() string {
	return RouterKey
}


func (am AppModule) LegacyQuerierHandler(amino *codec.LegacyAmino) sdk.Querier {
	return nil
}

func (am AppModule) RegisterQueryService(server grpc.Server) {
	types.RegisterQueryServer(server, am.keeper)
}

func (am AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
	am.BeginBlocker(ctx)
}

func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState

	cdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(gs)
}




