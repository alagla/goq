package abra

import (
	. "github.com/iotaledger/iota.go/trinary"
	"strings"
)

func TritsToString(trits Trits) string {
	if trits == nil {
		return "<nil>"
	}
	b := make([]byte, len(trits), len(trits))
	for i := range trits {
		switch trits[i] {
		case -1:
			b[i] = '-'
		case 0:
			b[i] = '0'
		case 1:
			b[i] = '1'
		default:
			b[i] = '?'
		}
	}
	return string(b)
}

func TritName(trit int8) string {
	switch trit {
	case -1:
		return "-1"
	case 0:
		return "0"
	case 1:
		return "1"
	}
	return "undef"
}

func Get1TritConstLutRepr(val int8) string {
	return strings.Repeat(TritName(val), 27)
}

func BinaryEncodedLUTFromString(strRepr string) int64 {
	var ret int64
	var bet int64
	bytes := []byte(strRepr)
	for i := 0; i < 27; i++ {
		bet = binaryEncodeTrit([]int8{int8(bytes[i])})
		ret = ret << 2
		ret |= bet
	}
	return ret
}

const (
	TRIT_MINUS1 = 0x0002
	TRIT_ZERO   = 0x0000
	TRIT_ONE    = 0x0001
	TRIT_NULL   = 0x0003
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
