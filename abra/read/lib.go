package read

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/abra"
	"github.com/lunfardo314/goq/abra/construct"
)

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

type tritReader struct {
	trits  Trits
	curPos int
}

// return trit, false or -100, true (eof)
func (tr *tritReader) readTrit() (int8, bool) {
	if tr.curPos >= len(tr.trits) {
		return -100, true
	}
	ret := tr.trits[tr.curPos]
	tr.curPos++
	return ret, false
}

func ParseTritcode(trits Trits) (*abra.CodeUnit, error) {
	ret := construct.NewCodeUnit()
	tReader := &tritReader{trits: trits}

	var err error
	var ver int
	ver, err = ReadPosInt(tReader)
	if err != nil {
		return nil, err
	}
	if ver != ret.Code.TritcodeVersion {
		return nil, fmt.Errorf("expected tritcode version %d, got %d", ret.Code.TritcodeVersion, ver)
	}
	return ret, nil
}

func ReadPosInt(tReader *tritReader) (int, error) {
	buf := make(Trits, 0, 31)
	exit := false
	for !exit {
		if len(buf) >= 31 {
			return -1, fmt.Errorf("ReadPosInt: wrong PosInt: longer that 31 bit at position %d", tReader.curPos-1)
		}
		t, eof := tReader.readTrit()
		switch {
		case eof:
			return -1, fmt.Errorf("ReadPosInt: unexpected EOF at position %d", tReader.curPos)
		case t == 0:
			exit = true
		case t == -1:
			buf = append(buf, 0)
		case t == 1:
			buf = append(buf, 1)
		default:
			return -1, fmt.Errorf("ReadPosInt: wrong trit at position %d", tReader.curPos-1)
		}
	}
	ret := 0
	for i := len(buf) - 1; i >= 0; i-- {
		ret <<= 1
		if buf[i] == 1 {
			ret |= 0x1
		}
	}
	return ret, nil
}
