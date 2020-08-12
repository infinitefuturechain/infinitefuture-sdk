package scene

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"strconv"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryTxScenes:
			return queryTxScenes(ctx, path[1], path[2], path[3], path[4], keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown scenes query endpoint")
		}
	}
}
func queryTxScenes(ctx sdk.Context, sceneKey, numLatestStr, limitStr, sort string, keeper Keeper) ([]byte, sdk.Error) {
	numLatest, err := strconv.ParseUint(numLatestStr, 0, 64)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}

	txs := keeper.GetTxSceneTxs(ctx, sceneKey, numLatest, limit, sort)
	bz, err := codec.MarshalJSONIndent(keeper.GetCodec(), txs)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
