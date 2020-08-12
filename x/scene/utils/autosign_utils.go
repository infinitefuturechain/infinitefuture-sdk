package utils

// Joe.He
import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/go-bip39"

	keys2 "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys"

	"github.com/cosmos/cosmos-sdk/client/context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	mnemonicEntropySize = 256
)

var (
	passphrase   = "123456789"
	autoSignAddr []string
	ch           = make(chan sdk.AccAddress, 0)
)

func GetAutoSingerPassphrase() string {
	return passphrase
}
func CheckSignTxHash(cliCtx context.CLIContext, addr sdk.AccAddress, hashHexStr string) {
	go CheckTxHash(ch, cliCtx, addr, hashHexStr)
}

func IsSignAddr(addr string) bool {
	if len(autoSignAddr) == 0 {
		return false
	}
	for _, v := range autoSignAddr {
		if v == addr {
			return true
		}
	}
	return false
}
func PopSignAddr() sdk.AccAddress {
	if len(ch) == 0 {
		return nil
	}
	return <-ch
}
func PushSignAddr(addr sdk.AccAddress) {
	if addr == nil {
		return
	}
	ch <- addr
}
func SetSignAddr(addr string) {
	if len(addr) == 0 {
		return
	}
	autoSignAddr = strings.Split(addr, ",")
	ch = make(chan sdk.AccAddress, len(autoSignAddr))
	for _, v := range autoSignAddr {
		addr, err := sdk.AccAddressFromBech32(v)
		if err != nil {
			panic(sdk.ErrInvalidAddress(v))
		}
		if !Exist(addr) {
			panic(sdk.ErrInternal("not a auto singer：" + v))
		}
		ch <- addr
	}
}
func GenSignAccount(amount int, cliCtx context.CLIContext, rootDir string) error {
	node, err := cliCtx.GetNode()
	if err != nil {
		return err
	}
	block := int64(1)
	v, err := node.Block(&block)
	if err != nil {
		return err
	}
	chainId := v.Block.ChainID
	if !strings.HasSuffix(rootDir, "/") {
		rootDir += "/"
	}
	entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
	if err != nil {
		panic(err)
	}
	mnemonic, err := bip39.NewMnemonic(entropySeed[:])
	if err != nil {
		panic(err)
	}
	lcd := ""
	mnemonics := mnemonic + "\n\n"
	kb := keys.NewInMemory()
	for i := 0; i < amount; i++ {
		info, err := kb.CreateAccount(fmt.Sprintf("user%d", i), mnemonic, "", passphrase, uint32(i), 0)
		if err != nil {
			panic(err)
		}
		addSign(info.GetAddress(), rootDir, mnemonic, passphrase, uint32(i))
		if len(lcd) == 0 {
			lcd = info.GetAddress().String()
		} else {
			lcd = lcd + "," + info.GetAddress().String()
		}
		mnemonics += fmt.Sprintf("%s:%d\n", info.GetAddress().String(), i)
	}

	subPath := strconv.FormatInt(time.Now().Unix(), 10)
	os.Mkdir(rootDir+subPath, os.ModePerm)

	lcd = version.ClientName + " rest-server --home " + rootDir + " --chain-id=" + chainId + "  --trust-node --laddr tcp://0.0.0.0:1317 --sign-addr " + lcd
	fmt.Println("lcd：", rootDir+subPath+"/lcd.sh")
	ioutil.WriteFile(rootDir+subPath+"/lcd.sh", []byte(lcd), os.ModePerm)

	fmt.Println("mnemonic：", rootDir+subPath+"/mnemonic.txt")
	ioutil.WriteFile(rootDir+subPath+"/mnemonic.txt", []byte(mnemonics), os.ModePerm)

	return nil
}

func addSign(addr sdk.AccAddress, rootDir, mnemonic, encryptPassword string, account uint32) {
	kb, _ := keys2.NewKeyBaseFromDir(rootDir + SingerPath + addr.String())
	_, err := kb.CreateAccount(addr.String(), mnemonic, "", encryptPassword, account, 0)
	if err != nil {
		panic(err)
	}
}
