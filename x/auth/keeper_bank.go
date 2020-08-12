// Joe.He
package auth

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"strings"
)

var (
	minId              = fmt.Sprintf("%020d", 0)
	maxId              = fmt.Sprintf("9%020d", 9)
	holderBytes        = make([]byte, 1, 1)
	keyHolderPrefix    = "holder"
	keyHolderKeyPrefix = "hkey"
)

func init() {
	holderBytes[0] = uint8(0)
}

func KeyHolder(address string, coin sdk.Coin) []byte {
	return []byte(fmt.Sprintf("%s:%s:%020s:%s", keyHolderPrefix, coin.Denom, coin.Amount.String(), address))
}
func KeyHolderDenomStr(denom, key string) []byte {
	return []byte(fmt.Sprintf("%s:%s:%s:", keyHolderPrefix, denom, key))
}
func KeyHolderFromKey(address, denom, amount string) []byte {
	return []byte(fmt.Sprintf("%s:%s:%020s:%s", keyHolderPrefix, denom, amount, address))
}
func KeyHolderKey(address, denom string) []byte {
	return []byte(fmt.Sprintf("%s:%s:%s", keyHolderKeyPrefix, address, denom))
}
func KeyHolderKeyPrefix(address string) []byte {
	return []byte(fmt.Sprintf("%s:%s:", keyHolderKeyPrefix, address))
}

func (k AccountKeeper) SetHolder(ctx sdk.Context, account exported.Account, newCoins sdk.Coins, genesis bool) {
	store := ctx.KVStore(k.key)
	addr := account.GetAddress().String()

	if genesis {
		for _, coin := range newCoins {
			store.Set(KeyHolder(addr, coin), holderBytes)
			store.Set(KeyHolderKey(addr, coin.Denom), []byte(coin.Amount.String()))
		}
		return
	}

	//for _, coin := range account.GetCoins() {
	//	store.Delete(KeyHolder(addr, coin))
	//}
	list := k.getHolders(store, addr)
	for _, v := range list {
		store.Delete(v)
	}

	for _, coin := range newCoins {
		store.Set(KeyHolderKey(addr, coin.Denom), []byte(coin.Amount.String()))
		store.Set(KeyHolder(addr, coin), holderBytes)
	}
}
func (k AccountKeeper) getHolders(store sdk.KVStore, addr string) [][]byte {
	iterator := sdk.KVStorePrefixIterator(store, KeyHolderKeyPrefix(addr))
	defer iterator.Close()
	list := make([][]byte, 0)
	for ; iterator.Valid(); iterator.Next() {
		keys := strings.Split(string(iterator.Key()), ":")
		list = append(list, KeyHolderFromKey(addr, keys[2], string(iterator.Value())))
	}
	return list
}
func (k AccountKeeper) GetHolders(ctx sdk.Context, denom string, startId string, limit int) []string {
	store := ctx.KVStore(k.key)

	startKey, endKey := k.getIteratorHolderRange(denom, startId)
	iterator := store.ReverseIterator(startKey, endKey)
	defer iterator.Close()
	list := make([]string, 0, limit)
	for ; iterator.Valid(); iterator.Next() {
		bz := iterator.Value()
		if len(bz) == 0 {
			continue
		}
		k := strings.Split(string(iterator.Key()), ":")
		list = append(list, k[2]+":"+k[3])
		if len(list) >= limit {
			break
		}
	}
	return list
}
func (k AccountKeeper) getIteratorHolderRange(denom string, startTxId string) (startId, endId []byte) {
	if len(startTxId) == 0 || "all" == startTxId {
		return KeyHolderDenomStr(denom, minId), KeyHolderDenomStr(denom, maxId)
	}
	return KeyHolderDenomStr(denom, minId), KeyHolderDenomStr(denom, startTxId)
}
