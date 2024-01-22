package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"myblockchain/x/myblockchain/types"
)

func (k msgServer) CreateItem(goCtx context.Context, msg *types.MsgCreateItem) (*types.MsgCreateItemResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgCreateItemResponse{}, nil
}
