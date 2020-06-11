package evm

import (
	"encoding/json"

	"github.com/gogo/protobuf/grpc"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/cosmos/ethermint/x/evm/client/cli"
	"github.com/cosmos/ethermint/x/evm/keeper"
	"github.com/cosmos/ethermint/x/evm/types"
)

var _ module.AppModuleBasic = AppModuleBasic{}
var _ module.AppModule = AppModule{}

// AppModuleBasic struct
type AppModuleBasic struct{}

// Name for app module basic
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterCodec registers types for module
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
}

// DefaultGenesis is json default structure
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis is the validation check of the Genesis
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONMarshaler, bz json.RawMessage) error {
	var genesisState types.GenesisState
	err := cdc.UnmarshalJSON(bz, &genesisState)
	if err != nil {
		return err
	}

	return genesisState.Validate()
}

// RegisterRESTRoutes registers the REST routes for the evm module.
func (AppModuleBasic) RegisterRESTRoutes(clientCtx client.Context, rtr *mux.Router) {}

// GetTxCmd returns the root tx command for the ibc module.
func (AppModuleBasic) GetTxCmd(clientCtx client.Context) *cobra.Command {
	return cli.GetTxCmd(clientCtx.Codec)
}

// GetQueryCmd Gets the root query command of this module
func (AppModuleBasic) GetQueryCmd(clientCtx client.Context) *cobra.Command {
	return cli.GetQueryCmd(clientCtx.Codec)
}

// RegisterInterfaceTypes registers module concrete types into protobuf Any.
func (AppModuleBasic) RegisterInterfaceTypes(registry cdctypes.InterfaceRegistry) {}

//____________________________________________________________________________

// AppModule implements an application module for the evm module.
type AppModule struct {
	AppModuleBasic
	keeper Keeper
	ak     types.AccountKeeper
}

// NewAppModule creates a new AppModule Object
func NewAppModule(k Keeper, ak types.AccountKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
		ak:             ak,
	}
}

// Name is module name
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterInvariants interface for registering invariants
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

// Route returns the sdk.Route for the evm module.
func (am AppModule) Route() *sdk.Route {
	return sdk.NewRoute(RouterKey, NewHandler(am.keeper))
}

// QuerierRoute sets up path for queries
func (am AppModule) QuerierRoute() string {
	return types.ModuleName
}

// NewQuerierHandler sets up new querier handler for module
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return keeper.NewQuerier(am.keeper)
}

// RegisterQueryService registers the gRPC query service for the evm module
func (am AppModule) RegisterQueryService(grpc.Server) {}

// BeginBlock function for module at start of each block
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	BeginBlock(am.keeper, ctx, req)
}

// EndBlock function for module at end of block
func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return EndBlock(am.keeper, ctx, req)
}

// InitGenesis instantiates the genesis state
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	return InitGenesis(ctx, am.keeper, genesisState)
}

// ExportGenesis exports the genesis state to be used by daemon
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper, am.ak)
	return cdc.MustMarshalJSON(gs)
}
