package keeper

import (
	"fmt"
)

var (
	TxSceneKeyIDPrefix = "tsi"
	TxSceneKeyPrefix   = "tsk"
)

func KeyNextTxSceneID(sceneKey string) []byte {
	return []byte(fmt.Sprintf("%s:%s", TxSceneKeyIDPrefix, sceneKey))
}
func KeyTxScene(sceneKey string, id uint64) []byte {
	return []byte(fmt.Sprintf("%s:%s:%d", TxSceneKeyPrefix, sceneKey, id))
}
func KeyTxScenePrefix() []byte {
	return []byte(fmt.Sprintf("%s:", TxSceneKeyPrefix))
}
