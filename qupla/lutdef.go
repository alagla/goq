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

// lookup table adjusted for 3 inputs
func (lutDef *LutDef) LookupTable() [27]Trits {
	return lutDef.lookupTable
}

const (
	TRIT_MINUS1 = 0x0003
	TRIT_ZERO   = 0x0000
	TRIT_ONE    = 0x0001
	TRIT_NULL   = 0x0002
)

func binaryEncodeTrit(trit []int8) int64 {
	if len(trit) != 1 {
		panic("wrong param")
	}
	if trit == nil {
		return TRIT_NULL
	}
	switch trit[0] {
	case -1:
		return TRIT_MINUS1
	case 0:
		return TRIT_ZERO
	case 1:
		return TRIT_ONE
	}
	panic("wrong trit")
}

func charEncodeTrit(trit []int8) byte {
	if len(trit) != 1 {
		panic("wrong param")
	}
	if trit == nil {
		return '@'
	}
	switch trit[0] {
	case -1:
		return '-'
	case 0:
		return '0'
	case 1:
		return '1'
	}
	panic("wrong trit")
}

func (lutDef *LutDef) BinaryEncodedLUT() int64 {
	var ret int64
	var bet int64
	lt := lutDef.LookupTable()
	for i := 0; i < 27; i++ {
		bet = binaryEncodeTrit(lt[i])
		ret = ret << 2
		ret |= bet
	}
	return ret
}

func (lutDef *LutDef) GetTritcode() Trits {
	ret := IntToTrits(lutDef.BinaryEncodedLUT())
	ret = PadTrits(ret, 35)
	if len(ret) != 35 {
		panic("wrong LUT tritcode")
	}
	return ret
}

func (lutDef *LutDef) GetStringRepr() string {
	var ret [27]byte
	for i, t := range lutDef.LookupTable() {
		ret[i] = charEncodeTrit(t)
	}
	return string(ret[:])
}

func (lutDef *LutDef) GetBranch(numInputs int) {
	if numInputs != 1 && numInputs != 2 && numInputs != 3 {
		panic("wrong number of inputs")
	}

}
