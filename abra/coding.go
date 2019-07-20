package abra

import (
	"fmt"
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
	if val != -1 && val != 0 && val != 1 {
		panic("wrong trit value")
	}
	return strings.Repeat(TritName(val), 27)
}

// https://github.com/iotaledger/omega-docs/blob/master/qbc/abra/Spec.md#lut-definition
//
// LUT definition
// The lookup table is encoded as 27 nullable trits, which fits in a 35 -trit number as 27 binary-coded trits.
// A lookup table which returns 0 for any input would look, in binary, like 3F_FF_FF_FF_FF.
//
// Since this value only covers for any non-null possible inputs, we start encoding by starting
// at all negatives (first input as lowest-endian), ---, and continuing to increment: 0--, 1--, -0-, 00-, 10-, ..., 111.
//
// Thus, the most-significant pair of bits (binary-coded trits) corresponds to 111, and the least
// significant pair of bits corresponds to ---.
//
//This final value is treated as a binary number, and encoded within a 35-trit vector.

func BinaryEncodedLUTFromString(strRepr string) uint64 {
	var ret uint64
	var bet uint64
	bytes := []byte(strRepr)
	for i := 0; i < 27; i++ {
		bet = uint64(TritFromByteRepr(bytes[i]))
		ret = ret << 2
		ret |= bet
	}
	return ret
}

func StringFromBinaryEncodedLUT(binary uint64) string {
	ret := make([]byte, 27)
	for i := 26; i >= 0; i-- {
		switch binary & 0x3 {
		case TRIT_NULL:
			ret[i] = '@'
		case TRIT_ONE:
			ret[i] = '1'
		case TRIT_ZERO:
			ret[i] = '0'
		case TRIT_MINUS1:
			ret[i] = '-'
		}
		binary >>= 2
	}
	return string(ret)
}

// https://github.com/iotaledger/omega-docs/blob/master/qbc/abra/Spec.md#lut-definition
// Encoding
// Positive integers (as listed above) are encoded as binary.1/-, little endian, terminated with 0.
//

func MustInt2PosIntAsPerSpec(n int) Trits {
	if n < 0 {
		panic("MustPosIntToPosIntAsPerSpec: non negative argument is expected")
	}
	var buf [64]int8
	// int can be max 64 bit long, normally 32. So msb != 1
	c := 0
	for ; n != 0; n >>= 1 {
		if n&0x1 == 0 {
			buf[c] = -1
		} else {
			buf[c] = 1
		}
		c++
	}
	// max c can be 63
	ret := make([]int8, c+1)
	copy(ret, buf[:c]) // always 0 in the end
	return ret
}

// Site indices may be positive or negative, so the minimum number of trits to encode the site is
// given first (positive integer), followed by the site value. 0 indicates both 0 trits and the value 0.
// 101 encodes 1, 10- encodes minus 1, 1101-- encodes minus 11.

func MustEncodeIndex(n int) Trits {
	if n == 0 {
		return []int8{0}
	}
	nt := IntToTrits(int64(n))
	lenpref := len(nt)
	tlen := MustInt2PosIntAsPerSpec(lenpref)
	ret := make([]int8, len(tlen)+len(nt))
	copy(ret, tlen)
	copy(ret[len(tlen):], nt)
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

func MustEncodeTritsToBytes(trits Trits) []byte {
	return []byte(MustTritsToTrytes(trits))
}

func TritEncodeLUTBinary(lutBin uint64) Trits {
	ret := IntToTrits(int64(lutBin))
	ret = PadTrits(ret, 35)
	if len(ret) != 35 {
		panic("wrong LUT tritcode")
	}
	return ret
}

func Bytes2Trits(bytes []byte) Trits {
	ret := make(Trits, len(bytes))
	for i, t := range bytes {
		ret[i] = int8(t)
	}
	return ret
}

var tryteAlpabet = map[byte][3]int8{
	'9': {0, 0, 0},    //    0
	'A': {1, 0, 0},    //    1
	'B': {-1, 1, 0},   //	    2
	'C': {0, 1, 0},    //	    3
	'D': {1, 1, 0},    //	    4
	'E': {-1, -1, 1},  //	    5
	'F': {0, -1, 1},   //	    6
	'G': {1, -1, 1},   //	    7
	'H': {-1, 0, 1},   //	    8
	'I': {0, 0, 1},    //	    9
	'J': {1, 0, 1},    //	   10
	'K': {-1, 1, 1},   //	   11
	'L': {0, 1, 1},    //	   12
	'M': {1, 1, 1},    //	   13
	'N': {-1, -1, -1}, //	  -13
	'O': {0, -1, -1},  //	  -12
	'P': {1, -1, -1},  //	  -11
	'Q': {-1, 0, -1},  //	  -10
	'R': {0, 0, -1},   //	   -9
	'S': {1, 0, -1},   //	   -8
	'T': {-1, 1, -1},  //	   -7
	'U': {0, 1, -1},   //	   -6
	'V': {1, 1, -1},   //	   -5
	'W': {-1, -1, 0},  //	   -4
	'X': {0, -1, 0},   //	   -3
	'Y': {1, -1, 0},   //	   -2
	'Z': {-1, 0, 0},   //    -1
}

func Trytes2Trits(trytes string) (Trits, error) {
	ret := make(Trits, len(trytes)*3)
	c := 0
	for i, t := range []byte(trytes) {
		trits, ok := tryteAlpabet[t]
		if !ok {
			return nil, fmt.Errorf("wrong tryte character at pos %d", i)
		}
		ret[c] = trits[0]
		c++
		ret[c] = trits[1]
		c++
		ret[c] = trits[2]
		c++
	}
	return ret, nil
}
