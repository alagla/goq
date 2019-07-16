package abra

import (
	. "github.com/iotaledger/iota.go/trinary"
	"strings"
)

const (
	TRIT_MINUS1 = 0x0002
	TRIT_ZERO   = 0x0000
	TRIT_ONE    = 0x0001
	TRIT_NULL   = 0x0003
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
			b[i] = '@'
		}
	}
	return string(b)
}

func TritName(trit int8) string {
	switch trit {
	case -1:
		return "-"
	case 0:
		return "0"
	case 1:
		return "1"
	}
	return "?"
}

func Get1TritConstLutRepr(val int8) string {
	return strings.Repeat(TritName(val), 27)
}

func BinaryEncodedLUTFromString(strRepr string) int64 {
	var ret int64
	var bet int64
	bytes := []byte(strRepr)
	for i := 0; i < 27; i++ {
		bet = int64(TritFromByteRepr(bytes[i]))
		ret = ret << 2
		ret |= bet
	}
	return ret
}

func TritFromByteRepr(c byte) int8 {
	switch c {
	case '-':
		return TRIT_MINUS1
	case '0':
		return TRIT_ZERO
	case '1':
		return TRIT_ONE
	case '@':
		return TRIT_NULL
	}
	panic("wrong trit repr")
}

// true = 1, false = -

// nullify true/false (return second arg when first is 1 (true) / - (false)
// true  --> "@@-@@0@@1@@-@@0@@1@@-@@0@@1"
// false --> "-@@0@@1@@-@@0@@1@@-@@0@@1@@"

// --- = @ / -  (-13)
// 0-- = @ / @  (-12)
// 1-- = - / @  (-11)
// -0- = @ / 0 (-10)
// 00- = @ / @  (-9)
// 10- = 0 / @  (-8)
// -1- = @ / 1 (-7)
// 01- = @ / @  (-6)
// 11- = 1 / @  (-5)
// --0 = @ / - (-4)
// 0-0 = @ / @  (-3)
// 1-0 = - / @  (-2)
// -00 = @ / 0 (-1)
// 000 = @ / @  (0)
// 100 = 0 / @  (1)
// -10 = @ / 1 (2)
// 010 = @ / @  (3)
// 110 = 1 / @  (4)
// --1 = @ / - (5)
// 0-1 = @ / @  (6)
// 1-1 = - / @  (7)
// -01 = @ / 0 (8)
// 001 = @ / @  (9)
// 101 = 0 / @  (10)
// -11 = @ / 1 (11)
// 011 = @ / @  (12)
// 111 = 1 / @  (13)

func GetNullifyLUTRepr(trueFalse bool) string {
	if trueFalse {
		return "@@-@@0@@1@@-@@0@@1@@-@@0@@1"
	} else {
		return "-@@0@@1@@-@@0@@1@@-@@0@@1@@"
	}
}
