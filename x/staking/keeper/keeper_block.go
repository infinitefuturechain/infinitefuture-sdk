package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	types2 "github.com/tendermint/tendermint/types"
)

// Joe.He
func KeyValBlock(valAddr sdk.ValAddress, id uint64) []byte {
	return []byte(fmt.Sprintf("vblock:%s:%d", valAddr.String(), id))
}
func KeyValBlocksID(valAddr sdk.ValAddress) []byte {
	return []byte(fmt.Sprintf("nvid:%s", valAddr.String()))
}

func (k Keeper) AddStakingBlock(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	consAddr := sdk.ConsAddress(ctx.BlockHeader().ProposerAddress)
	if consAddr == nil {
		return
	}
	valAddr := k.ValidatorByConsAddr(ctx, consAddr).GetOperator()

	id := k.getNewValBlockID(store, valAddr)
	store.Set(KeyValBlock(valAddr, id), k.cdc.MustMarshalBinaryLengthPrefixed(ctx.BlockHeight()))
}

func (k Keeper) GetStakingBlocks(ctx sdk.Context, valAddr sdk.ValAddress, startId uint64, limit int) (blocks types.ValidatorBlocks) {
	store := ctx.KVStore(k.storeKey)

	maxID := k.getLastValBlockID(store, valAddr)
	if startId == 0 {
		startId = maxID
	}
	blocks = types.ValidatorBlocks{Total: maxID}
	blockList := make([]types2.BlockMeta, 0)
	blockStore := server.ServerNode.BlockStore()
	for id := startId; id > 0; id-- {
		bz := store.Get(KeyValBlock(valAddr, id))
		if len(bz) == 0 {
			continue
		}
		var block int64
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &block)
		blockList = append(blockList, *blockStore.LoadBlockMeta(block))
		if len(blockList) >= limit {
			break
		}
	}
	blocks.Blocks = blockList
	return
}
func (k Keeper) getLastValBlockID(store sdk.KVStore, valAddr sdk.ValAddress) (id uint64) {
	bz := store.Get(KeyValBlocksID(valAddr))
	if len(bz) == 0 {
		id = 1
	} else {
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &id)
	}
	return
}
func (k Keeper) getNewValBlockID(store sdk.KVStore, valAddr sdk.ValAddress) (id uint64) {
	id = k.getLastValBlockID(store, valAddr)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(id + 1)
	store.Set(KeyValBlocksID(valAddr), bz)
	return
}
