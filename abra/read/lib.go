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

	err := ParseCode(tReader, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func ParseCode(tReader *tritReader, codeUnit *abra.CodeUnit) error {
	var err error
	var ver int
	ver, err = ParsePosInt(tReader)
	if err != nil {
		return err
	}
	if ver != codeUnit.Code.TritcodeVersion {
		return fmt.Errorf("expected tritcode version %d, got %d", codeUnit.Code.TritcodeVersion, ver)
	}
	codeUnit.Code.NumLUTs, err = ParsePosInt(tReader)
	if err != nil {
		return err
	}
	for i := 0; i < codeUnit.Code.NumLUTs; i++ {
		err = ParseLUTBlock(tReader, codeUnit)
		if err != nil {
			return err
		}
	}
	// TODO branches and externals
	return nil
}

func ParseLUTBlock(tReader *tritReader, codeUnit *abra.CodeUnit) error {
	n, err := ParsePosInt(tReader)
	if err != nil {
		return err
	}
	if n != 35 {
		return fmt.Errorf("expected PosInt == 35 at position %d", tReader.curPos)
	}
	var trits Trits
	trits, err = ReadNTrits(tReader, 35)
	if err != nil {
		return err
	}
	strRepr := abra.StringFromBinaryEncodedLUT(uint64(TritsToInt(trits)))
	_, err = construct.AddNewLUTBlock(codeUnit, strRepr, "")
	if err != nil {
		return err
	}
	return nil
}

func ParsePosInt(tReader *tritReader) (int, error) {
	buf := make(Trits, 0, 31)
	exit := false
	for !exit {
		if len(buf) >= 31 {
			return -1, fmt.Errorf("ParsePosInt: wrong PosInt: longer that 31 bit at position %d", tReader.curPos-1)
		}
		t, eof := tReader.readTrit()
		switch {
		case eof:
			return -1, fmt.Errorf("ParsePosInt: unexpected EOF at position %d", tReader.curPos)
		case t == 0:
			exit = true
		case t == -1:
			buf = append(buf, 0)
		case t == 1:
			buf = append(buf, 1)
		default:
			return -1, fmt.Errorf("ParsePosInt: wrong trit at position %d", tReader.curPos-1)
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

func ReadNTrits(tReader *tritReader, n int) (Trits, error) {
	ret := make(Trits, n)
	var eof bool
	for i := 0; i < n; i++ {
		ret[i], eof = tReader.readTrit()
		if eof {
			return nil, fmt.Errorf("unexpected EOF at pos %d", tReader.curPos)
		}
	}
	return ret, nil
}
