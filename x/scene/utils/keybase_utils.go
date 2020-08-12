package utils

// Joe.He
import (
	"encoding/hex"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/cosmos/cosmos-sdk/crypto/keys"

	ckeys "github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	LcdHome    = ""
	SingerPath = "singer/"
)

func OpenKeyBase(addr sdk.AccAddress) (keys.Keybase, error) {
	path := LcdHome
	if strings.HasSuffix(path, "/") {
		path = path + SingerPath + addr.String()
	} else {
		path = path + "/" + SingerPath + addr.String()
	}
	kb, err := ckeys.NewKeyBaseFromDir(path)
	if err != nil {
		return nil, err
	}
	return kb, nil
}

func Exist(addr sdk.AccAddress) bool {
	kb, err := OpenKeyBase(addr)
	if err != nil {
		return false
	}
	kbs, err := kb.List()
	if err != nil {
		return false
	}
	for _, info := range kbs {
		if info.GetAddress().Equals(addr) {
			return true
		}
	}
	return false
}

func CheckTxHash(txCh chan sdk.AccAddress, cliCtx context.CLIContext, addr sdk.AccAddress, hashHexStr string) {
	hash, err := hex.DecodeString(hashHexStr)
	if err != nil {
		return
	}
	node, _ := cliCtx.GetNode()
	if err != nil {
		time.Sleep(3 * time.Second)
		txCh <- addr
		return
	}
	resultStatus, err := node.Status()
	if err != nil {
		time.Sleep(3 * time.Second)
		txCh <- addr
		return
	}
	latestBlockHeight := resultStatus.SyncInfo.LatestBlockHeight
	for {
		time.Sleep(1 * time.Second)
		_, txErr := node.Tx(hash, true)
		if txErr != nil {
			resultStatus, err := node.Status()
			if err != nil {
				time.Sleep(3 * time.Second)
				txCh <- addr
				break
			}
			if resultStatus.SyncInfo.LatestBlockHeight >= latestBlockHeight+10 {
				txCh <- addr
				break
			}
			continue
		}
		txCh <- addr
		break
	}
}
