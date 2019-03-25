package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/utils"
)

type LutDef struct {
	Name           string
	InputSize      int
	OutputSize     int
	LutLookupTable []Trits
}

//func (LutDef *LutDef) SetName(Name string) {
//	LutDef.Name = Name
//}

func (lutDef *LutDef) Size() int {
	return lutDef.OutputSize
}

func (lutDef *LutDef) Lookup(res, args Trits) bool {
	t := lutDef.LutLookupTable[utils.Trits3ToLutIdx(args)]
	if t == nil {
		return true
	}
	copy(res, t)
	return false
}
