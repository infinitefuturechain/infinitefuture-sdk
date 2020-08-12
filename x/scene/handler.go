package scene

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "scene" type messages.
func NewHandler(_ Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		
		switch msg := msg.(type) {
		default:
			errMsg := fmt.Sprintf("unrecognized org message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}
