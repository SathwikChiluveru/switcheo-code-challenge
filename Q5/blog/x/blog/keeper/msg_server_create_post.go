package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"blog/x/blog/types"
)

func (k msgServer) CreatePost(goCtx context.Context, msg *types.MsgCreatePost) (*types.MsgCreatePostResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Breaking Change: Change the key used to store the post in the underlying key-value store
	storeKey := []byte("new_post_key/") + sdk.Uint64ToBigEndian(msg.Id).Bytes()

	var post = types.Post{
		Creator: msg.Creator,
		Title:   msg.Title,
		Body:    msg.Body,
	}
	id := k.AppendPost(ctx, post, storeKey) // Update the AppendPost function to accept the store key

	return &types.MsgCreatePostResponse{
		Id: id,
	}, nil
}
