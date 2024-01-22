// x/myblockchain/module.go
package myblockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"myblockchain/x/myblockchain/keeper"
	"myblockchain/x/myblockchain/types"
)

var (
	_ module.AppModuleBasic      = (*AppModule)(nil)
	_ module.AppModuleSimulation = (*AppModule)(nil)
	_ module.HasGenesis          = (*AppModule)(nil)
	_ module.HasInvariants       = (*AppModule)(nil)
	_ module.HasConsensusVersion = (*AppModule)(nil)

	_ appmodule.AppModule       = (*AppModule)(nil)
	_ appmodule.HasBeginBlocker = (*AppModule)(nil)
	_ appmodule.HasEndBlocker   = (*AppModule)(nil)
)

type AppModuleBasic struct {
	cdc codec.BinaryCodec
	AppModuleBasic
	Keeper
}

func NewAppModuleBasic(cdc codec.BinaryCodec) AppModuleBasic {
	return AppModuleBasic{cdc: cdc}
}

func (AppModuleBasic) Name() string {
	return types.ModuleName
}

func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {}

func (a AppModuleBasic) RegisterInterfaces(reg cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(reg)
}

func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesis())
}

func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return genState.Validate()
}

func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

type AppModule struct {
	AppModuleBasic
	keeper        keeper.Keeper
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
}

func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(cdc),
		keeper:         keeper,
		accountKeeper:  accountKeeper,
		bankKeeper:     bankKeeper,
	}
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}

func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) {
	var genState types.GenesisState
	cdc.MustUnmarshalJSON(gs, &genState)

	InitGenesis(ctx, am.keeper, genState)
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genState := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(genState)
}

func (AppModule) ConsensusVersion() uint64 { return 1 }

func (am AppModule) BeginBlock(_ context.Context) error {
	return nil
}

func (am AppModule) EndBlock(_ context.Context) error {
	return nil
}

func (am AppModule) IsOnePerModuleType() {}

func (am AppModule) IsAppModule() {}

func (am AppModule) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rtr.HandleFunc("/api/resource", am.CreateResourceHandler).Methods("POST")
	rtr.HandleFunc("/api/resources", am.ListResourcesHandler).Methods("GET")
	rtr.HandleFunc("/api/resource/{id}", am.GetResourceHandler).Methods("GET")
	rtr.HandleFunc("/api/resource/{id}", am.UpdateResourceHandler).Methods("PUT")
	rtr.HandleFunc("/api/resource/{id}", am.DeleteResourceHandler).Methods("DELETE")
}

func (am AppModule) GetTxCmd(*codec.Codec) *cobra.Command {
	return nil
}

func (am AppModule) GetQueryCmd(*codec.Codec) *cobra.Command {
	return nil
}

func init() {
	appmodule.Register(
		&modulev1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	StoreService store.KVStoreService
	Cdc          codec.Codec
	Config       *modulev1.Module
	Logger       log.Logger

	AccountKeeper types.AccountKeeper
	BankKeeper    types.BankKeeper
}

type ModuleOutputs struct {
	depinject.Out

	MyblockchainKeeper keeper.Keeper
	Module             appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}
	k := keeper.NewKeeper(
		in.Cdc,
		in.StoreService,
		in.Logger,
		authority.String(),
	)
	m := NewAppModule(
		in.Cdc,
		k,
		in.AccountKeeper,
		in.BankKeeper,
	)

	return ModuleOutputs{MyblockchainKeeper: k, Module: m}
}

// Handlers

func (am AppModule) CreateResourceHandler(w http.ResponseWriter, r *http.Request) {
	var resource types.Resource
	err := json.NewDecoder(r.Body).Decode(&resource)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err = am.keeper.CreateResource(ctx, resource.ID, resource.Name, resource.Details)
	if err != nil {
		http.Error(w, "Error creating resource", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (am AppModule) ListResourcesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	resources := am.keeper.GetAllResources(ctx)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resources)
}

func (am AppModule) GetResourceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resourceID := vars["id"]

	ctx := context.Background()
	resource := am.keeper.GetResource(ctx, resourceID)
	if resource == nil {
		http.Error(w, "Resource not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resource)
}

func (am AppModule) UpdateResourceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resourceID := vars["id"]

	var updatedResource types.Resource
	err := json.NewDecoder(r.Body).Decode(&updatedResource)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err = am.keeper.UpdateResource(ctx, resourceID, updatedResource.Name, updatedResource.Details)
	if err != nil {
		http.Error(w, "Error updating resource", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (am AppModule) DeleteResourceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resourceID := vars["id"]

	ctx := context.Background()
	err := am.keeper.DeleteResource(ctx, resourceID)
	if err != nil {
		http.Error(w, "Error deleting resource", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
