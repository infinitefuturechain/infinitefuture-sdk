package types

import (
	"fmt"
)

type TxId struct {
	ID     uint64 `json:"id" yaml:"id"`
	TxHash string `json:"tx_hash" yaml:"tx_hash"`
}
type TxIds []TxId

func (t TxId) String() string {
	return fmt.Sprintf(`Tx:
	ID:           				%d
	TxHash:           			%s`,
		t.ID, t.TxHash)
}
func (v TxIds) String() string {
	if len(v) == 0 {
		return "[]"
	}
	out := fmt.Sprintf("Txs :")
	for _, vot := range v {
		out += fmt.Sprintf("\n----\n  %s", vot.String())
	}
	return out
}
