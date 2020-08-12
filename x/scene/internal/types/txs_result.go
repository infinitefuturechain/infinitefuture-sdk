package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TxIdResult struct {
	ID uint64         `json:"id" yaml:"id"`
	Tx sdk.TxResponse `json:"tx" yaml:"tx"`
}
type TxIdResults []TxIdResult

func (t TxIdResult) String() string {
	return fmt.Sprintf(`Tx:
	ID:           				%d
	Result:           			%s`,
		t.ID, t.Tx.String())
}
func (v TxIdResults) String() string {
	if len(v) == 0 {
		return "[]"
	}
	out := fmt.Sprintf("Txs :")
	for _, vot := range v {
		out += fmt.Sprintf("\n----\n  %s", vot.String())
	}
	return out
}
