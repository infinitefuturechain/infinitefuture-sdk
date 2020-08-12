package types

// Joe.He
import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/cosmos/cosmos-sdk/x/scene/utils"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	crkeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func broadcastTxToChain(cliCtx context.CLIContext, txBldr TxBuilder, sender sdk.AccAddress, msg sdk.Msg) (sdk.TxResponse, error) {
	validateErr := msg.ValidateBasic()
	if validateErr != nil {
		return sdk.TxResponse{}, validateErr
	}
	name, kb, err := getAccount(sender)
	if err != nil {
		return sdk.TxResponse{}, err
	}
	num, seq, err := NewAccountRetriever(cliCtx).GetAccountNumberSequence(sender)
	if err != nil {
		return sdk.TxResponse{}, err
	}
	txBldr = txBldr.WithKeybase(kb).WithAccountNumber(num).WithSequence(seq)
	txBytes, err := txBldr.BuildAndSign(name, utils.GetAutoSingerPassphrase(), []sdk.Msg{msg})
	if err != nil {
		return sdk.TxResponse{}, err
	}
	res, err := cliCtx.BroadcastTx(txBytes)
	if err != nil {
		return sdk.TxResponse{}, err
	}
	return res, nil
}

func BroadcastTxBySign(cliCtx context.CLIContext, txBldr TxBuilder, msg sdk.Msg) (sdk.TxResponse, error) {
	sender := msg.GetSigners()[0]
	res, err := broadcastTxToChain(cliCtx, txBldr, sender, msg)
	if err != nil {
		utils.PushSignAddr(sender)
		return sdk.TxResponse{}, err
	}
	if cliCtx.BroadcastMode == client.BroadcastBlock {
		utils.PushSignAddr(sender)
	} else {
		utils.CheckSignTxHash(cliCtx, sender, res.TxHash)
	}
	return res, nil
}

func getAccount(sender sdk.AccAddress) (name string, kb crkeys.Keybase, err error) {
	kb, err = utils.OpenKeyBase(sender)
	if err != nil {
		return "", kb, err
	}
	kbs, err := kb.List()
	if err != nil {
		return "", kb, err
	}
	for _, info := range kbs {
		if info.GetAddress().Equals(sender) {
			return info.GetName(), kb, nil
		}
	}
	return name, kb, sdk.ErrInvalidAddress("no signer:" + sender.String())
}

func BroadcastTx(cliCtx context.CLIContext, chainId, memo string, fees sdk.Coins, msg sdk.Msg, w http.ResponseWriter) {
	validateErr := msg.ValidateBasic()
	if validateErr != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, validateErr.Error())
		return
	}
	txBldr := NewTxBuilderFromRest(chainId, memo, fees).WithTxEncoder(GetTxEncoder(cliCtx.Codec))
	rep, err := BroadcastTxBySign(cliCtx, txBldr, msg)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	rest.PostProcessResponseBare(w, cliCtx, rep)
}
func GetTxEncoder(cdc *codec.Codec) (encoder sdk.TxEncoder) {
	encoder = sdk.GetConfig().GetTxEncoder()
	if encoder == nil {
		encoder = DefaultTxEncoder(cdc)
	}

	return encoder
}
