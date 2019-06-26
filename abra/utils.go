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

func GetNullifyLUTRepr(trueFalse bool) string {
	if trueFalse {
		return "@@-@@0@@1@@-@@0@@1@@-@@0@@1"
	} else {
		return "-@@0@@1@@-@@0@@1@@-@@0@@1@@"
	}
}
