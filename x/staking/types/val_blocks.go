package types

import (
	"fmt"
	"github.com/tendermint/tendermint/types"
)

// Joe.He
type ValidatorBlocks struct {
	Total  uint64            `json:"total" yaml:"total"`
	Blocks []types.BlockMeta `json:"blocks" yaml:"blocks"`
}

func (p ValidatorBlocks) String() string {
	return fmt.Sprintf(`ValidatorBlocks:
  Total:    %d
  Blocks:    %v`, p.Total, p.Blocks)
}
