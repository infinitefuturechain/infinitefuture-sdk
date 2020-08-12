package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	authutils "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/scene/internal/types"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	sceneKey  = "sceneKey"
	numLatest = "numLatest"
	limit     = "limit"
	sort      = "sort"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	r.HandleFunc(fmt.Sprintf("/%s/%s/{%s}/{%s}/{%s}/{%s}", queryRoute, types.QueryTxScenes,
		sceneKey, numLatest, limit, sort), queryTxSceneHandlerFn(queryRoute, cliCtx)).Methods("GET")
}

func queryTxSceneHandlerFn(queryRoute string, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s/%s/%s/%s", queryRoute, types.QueryTxScenes,
			vars[sceneKey], vars[numLatest], vars[limit], vars[sort]), nil)

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if len(res) == 0 {
			rest.PostProcessResponse(w, cliCtx, make([]types.TxIdResult, 0, 0))
			return
		}
		var ls types.TxIds
		cliCtx.Codec.MustUnmarshalJSON(res, &ls)
		list := make([]types.TxIdResult, 0, len(ls))
		cliCtx.TrustNode = true
		for _, txId := range ls {
			output, err := authutils.QueryTx(cliCtx, txId.TxHash)
			if err != nil {
				continue
			}
			output.Logs = nil
			v := types.TxIdResult{ID: txId.ID, Tx: output}
			list = append(list, v)
		}
		rest.PostProcessResponse(w, cliCtx, list)
	}
}
