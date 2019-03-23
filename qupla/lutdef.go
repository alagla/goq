package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/utils"
)

type QuplaLutDef struct {
	Name           string
	InputSize      int
	OutputSize     int
	LutLookupTable []Trits
}

//func (LutDef *QuplaLutDef) SetName(Name string) {
//	LutDef.Name = Name
//}

func (lutDef *QuplaLutDef) Size() int64 {
	return int64(lutDef.OutputSize)
}

func (lutDef *QuplaLutDef) Lookup(res, args Trits) bool {
	t := lutDef.LutLookupTable[utils.Trits3ToLutIdx(args)]
	if t == nil {
		return true
	}
	copy(res, t)
	return false
}
