package qupla

import (
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
