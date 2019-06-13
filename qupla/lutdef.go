package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
)

type LutDef struct {
	Name        string
	InputSize   int
	outputSize  int
	lookupTable [27]Trits
}

func NewLUTDef(name string, inSize int, outSize int, lookupTable []Trits) *LutDef {
	ret := &LutDef{
		Name:       name,
		InputSize:  inSize,
		outputSize: outSize,
	}
	copy(ret.lookupTable[:], lookupTable)
	return ret
}

func (lutDef *LutDef) Size() int {
	return lutDef.outputSize
}

func (lutDef *LutDef) Lookup(res, args Trits) bool {
	idx := Trits3ToLutIdx(args)
	t := lutDef.lookupTable[idx]
	if t == nil {
		return true
	}
	copy(res, t)
	return false
}

func Trits3ToLutIdx(trits Trits) int {
	idx := int8(0)
	switch len(trits) {
	case 1:
		idx = 13 + trits[0]
	case 2:
		idx = 13 + trits[0] + trits[1]*3
	default:
		idx = 13 + trits[0] + trits[1]*3 + trits[2]*9
	}
	return int(idx)
}

func (lutDef *LutDef) MustGetProjectionName(outpos int) string {
	if outpos < 0 || outpos >= lutDef.Size() {
		panic(fmt.Errorf("wrong arg num"))
	}
	return lutDef.Name + fmt.Sprintf("_proj_arg_%d", outpos)
}

//  it is assumed, that lookup table is adjusted for 3 inputs regardless real input size
//
// creates new lut def out of the old one, with same args gives outpos position of the result
// with output size = 1
// arg = 0,1,2

func (lutDef *LutDef) MakeAdjustedProjection(outpos int) *LutDef {
	if outpos < 0 || outpos >= lutDef.InputSize {
		panic(fmt.Errorf("wrong lut argument index"))
	}
	ret := LutDef{
		Name:       lutDef.MustGetProjectionName(outpos),
		InputSize:  lutDef.InputSize,
		outputSize: 1,
	}
	// assumed, that lookupTable is adjusted for all three outputs
	for i := range lutDef.lookupTable {
		ret.lookupTable[i] = lutDef.lookupTable[i][outpos : outpos+1]
	}
	return &ret
}
