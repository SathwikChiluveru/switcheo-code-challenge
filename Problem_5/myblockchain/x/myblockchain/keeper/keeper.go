package keeper

import (
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"myblockchain/x/myblockchain/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger
		storeKey     sdk.StoreKey

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	storeKey sdk.StoreKey, 

) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		storeKey:     storeKey,
		authority:    authority,
		logger:       logger,
	}
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// CRUD operations

func (k Keeper) CreateItem(ctx sdk.Context, id, name, details string) {
	store := ctx.KVStore(k.storeKey)
	item := types.NewItem(id, name, details)
	store.Set([]byte(id), k.cdc.MustMarshalBinaryBare(item))
}

func (k Keeper) GetItem(ctx sdk.Context, id string) types.Item {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(id)) {
		panic("Item not found")
	}
	bz := store.Get([]byte(id))
	var item types.Item
	k.cdc.MustUnmarshalBinaryBare(bz, &item)
	return item
}

func (k Keeper) GetAllItems(ctx sdk.Context) []types.Item {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	var items []types.Item
	for ; iterator.Valid(); iterator.Next() {
		var item types.Item
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &item)
		items = append(items, item)
	}

	return items
}

func (k Keeper) UpdateItem(ctx sdk.Context, id, name, details string) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(id)) {
		panic("Item not found")
	}
	item := types.NewItem(id, name, details)
	store.Set([]byte(id), k.cdc.MustMarshalBinaryBare(item))
}

func (k Keeper) DeleteItem(ctx sdk.Context, id string) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(id)) {
		panic("Item not found")
	}
	store.Delete([]byte(id))
}