package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {

}

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	ModuleCdc = cdc.Seal()
}
