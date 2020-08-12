package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/scene/internal/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the mint store
type Keeper struct {
	cdc      *codec.Codec
	storeKey sdk.StoreKey
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: key,
	}
}
func (k Keeper) GetCodec() *codec.Codec {
	return k.cdc
}

//______________________________________________________________________

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
