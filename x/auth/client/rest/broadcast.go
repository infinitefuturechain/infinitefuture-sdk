package rest

import (
	"fmt"
	"github.com/buger/jsonparser"
	utils2 "github.com/cosmos/cosmos-sdk/x/scene/utils"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

const (
	RestDelegation = "delegation"
)

// BroadcastReq defines a tx broadcasting request.
type BroadcastReq struct {
	Tx   types.StdTx `json:"tx" yaml:"tx"`
	Mode string      `json:"mode" yaml:"mode"`
	// Joe.He
	ChainID string `json:"chain_id" yaml:"chain_id"`
}

// BroadcastTxRequest implements a tx broadcasting handler that is responsible
// for broadcasting a valid and signed tx to a full node. The tx can be
// broadcasted via a sync|async|block mechanism.
func BroadcastTxRequest(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BroadcastReq

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		// Joe.He
		if "true" == r.URL.Query().Get(RestDelegation) {
			sender := utils2.PopSignAddr()
			if sender == nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, "no signer")
				return
			}
			body, err := jsonparser.Set(body, []byte(fmt.Sprintf("\"%s\"", sender.String())), "tx", "msg", "[0]", "value", "sender")
			if err != nil {
				utils2.PushSignAddr(sender)
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			err = cliCtx.Codec.UnmarshalJSON(body, &req)
			if err != nil {
				utils2.PushSignAddr(sender)
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			cliCtx = cliCtx.WithBroadcastMode(req.Mode)
			types.BroadcastTx(cliCtx, req.ChainID, req.Tx.Memo, req.Tx.Fee.Amount, req.Tx.Msgs[0], w)
			return
		}
		err = cliCtx.Codec.UnmarshalJSON(body, &req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		txBytes, err := cliCtx.Codec.MarshalBinaryLengthPrefixed(req.Tx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithBroadcastMode(req.Mode)

		res, err := cliCtx.BroadcastTx(txBytes)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponseBare(w, cliCtx, res)
	}
}
