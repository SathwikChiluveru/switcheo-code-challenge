package keeper

import (
	"myblockchain/x/myblockchain/types"
)

var _ types.QueryServer = Keeper{}
