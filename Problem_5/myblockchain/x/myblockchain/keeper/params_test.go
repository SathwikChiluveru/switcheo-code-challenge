package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "myblockchain/testutil/keeper"
	"myblockchain/x/myblockchain/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.MyblockchainKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
