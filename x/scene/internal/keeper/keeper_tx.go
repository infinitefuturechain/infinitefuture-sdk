package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/scene/internal/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

const (
	AllTxScene = "all"
)

func (k Keeper) AddSenderTxScene(ctx sdk.Context, sender sdk.AccAddress) (err sdk.Error) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(fmt.Sprintf("%X", tmhash.Sum(ctx.TxBytes())))
	v := sender.String()
	store.Set(KeyTxScene(v, k.getNewTxSceneID(store, v)), bz)
	return
}
func (k Keeper) AddTxScene(ctx sdk.Context, module string, signers []sdk.AccAddress) (err sdk.Error) {
	if ctx.BlockHeight() < 10 {
		return
	}
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(fmt.Sprintf("%X", tmhash.Sum(ctx.TxBytes())))

	sceneId := k.getNewTxSceneID(store, module)
	store.Set(KeyTxScene(module, sceneId), bz)

	sceneId = k.getNewTxSceneID(store, AllTxScene)
	store.Set(KeyTxScene(AllTxScene, sceneId), bz)

	for _, v := range signers {
		signer := v.String()
		store.Set(KeyTxScene(signer, k.getNewTxSceneID(store, signer)), bz)
	}

	return
}
func (k Keeper) getTxScene(store sdk.KVStore, module string, sceneId uint64) (txHash string) {
	bz := store.Get(KeyTxScene(module, sceneId))
	if len(bz) > 0 {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &txHash)
	}
	return
}

func (k Keeper) GetTxSceneTxs(ctx sdk.Context, module string, numLatest uint64, limit int, sort string) []*types.TxId {
	store := ctx.KVStore(k.storeKey)
	list := make([]*types.TxId, 0, limit)
	maxSceneId := k.getLastTxSceneID(store, module) - 1
	if "asc" == sort {
		if numLatest == 0 {
			numLatest = 1
		}
		for txID := numLatest; txID <= maxSceneId; txID++ {
			list = append(list, &types.TxId{ID: txID, TxHash: k.getTxScene(store, module, txID)})
			if len(list) >= limit {
				break
			}
		}
		return list
	}
	if numLatest == 0 {
		numLatest = maxSceneId
	}
	for txID := numLatest; txID > 0; txID-- {
		list = append(list, &types.TxId{ID: txID, TxHash: k.getTxScene(store, module, txID)})
		if len(list) >= limit {
			break
		}
	}
	return list
}

func (k Keeper) getNewTxSceneID(store sdk.KVStore, sceneKey string) (sceneId uint64) {
	sceneId = k.getLastTxSceneID(store, sceneKey)
	store.Set(KeyNextTxSceneID(sceneKey), k.cdc.MustMarshalBinaryLengthPrefixed(sceneId+1))
	return
}

func (k Keeper) getLastTxSceneID(store sdk.KVStore, sceneKey string) (sceneId uint64) {
	bz := store.Get(KeyNextTxSceneID(sceneKey))
	sceneId = 1
	if len(bz) > 0 {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &sceneId)
	}
	return
}
