package scene

import (
	"github.com/cosmos/cosmos-sdk/x/scene/internal/keeper"
	"github.com/cosmos/cosmos-sdk/x/scene/internal/types"
)

const (
	ModuleName        = types.ModuleName
	DefaultParamspace = types.DefaultParamspace
	StoreKey          = types.StoreKey
	QuerierRoute      = types.QuerierRoute
	QueryTxScenes     = types.QueryTxScenes
	RouterKey         = types.RouterKey
)

type (
	Keeper = keeper.Keeper
)

var (
	ModuleCdc     = types.ModuleCdc
	NewKeeper     = keeper.NewKeeper
	RegisterCodec = types.RegisterCodec
)
