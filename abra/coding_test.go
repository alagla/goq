package abra

import (
	"testing"
)

type testData struct {
	input    int
	expected string
}

var testval1 = []testData{
	{0, "0"},
	{1, "10"},
	{20, "--1-10"},
}

func TestMustInt2PosIntAsPerSpec(t *testing.T) {
	for _, v := range testval1 {
		rs := MustInt2PosIntAsPerSpec(v.input)
		if len(rs) == 0 || rs[len(rs)-1] != 0 {
			t.Errorf("wrong resutl '%v'", rs)
		}
		r := TritsToString(rs)
		if r != v.expected && len(r)+1 == len(v.expected) {
			t.Errorf("MustInt2PosIntAsPerSpec(%d) -> %v != expectd %s\n", v.input, r, v.expected)
		}
	}
}

var testval2 = []testData{
	{0, "0"},
	{1, "101"},
	{-1, "10-"},
	{20, "--10-1-1"},
}

func TestMustEncodeIndex(t *testing.T) {
	for _, v := range testval2 {
		r := TritsToString(MustEncodeIndex(v.input))
		if r != v.expected && len(r)+1 == len(v.expected) {
			t.Errorf("MustInt2PosIntAsPerSpec(%d) -> %v != expectd %s\n", v.input, r, v.expected)
		}
	}
}

var lut_test = []string{
	"---------------------------",
	"-------------01111111111111",
	"----01111----01111----01111",
	"--0-00000-00000001000001011",
	"--1--1--1--1--1--1--1--1--1",
	"-00000001-00000001-00000001",
	"-1--1--1--1--1--1--1--1--1-",
	"-1-1-1-1-1-1-1-1-1-1-1-1-1-",
	"-11-11-11-11-11-11-11-11-11",
	"-111-111--111-111--111-111-",
	"-@-@@@-@-@@@@@@@@@-@-@@@-@1",
	"-@-@@@-@1-@-@@@-@1-@-@@@-@1",
	"-@1@@@1@--@1@@@1@--@1@@@1@-",
	"-@1@@@1@-@@@@@@@@@1@-@@@-@1",
	"-@1@@@1@1-@1@@@1@1-@1@@@1@1",
	"-@1@@@1@1@@@@@@@@@1@1@@@1@1",
	"-@@0@@1@@-@@0@@1@@-@@0@@1@@",
	"000000000000000000000000000",
	"01-1-0-011-0-0101--0101-1-0",
	"011-01--0011-01--0011-01--0",
	"1---1---11---1---11---1---1",
	"1--1--1--1--1--1--1--1--1--",
	"1-0-0101-1-0-0101-1-0-0101-",
	"1-11-11-11-11-11-11-11-11-1",
	"10-000-0110-000-0110-000-01",
	"10-10-10-10-10-10-10-10-10-",
	"11-11-11-11-11-11-11-11-11-",
	"111111111111111111111111111",
	"1@-1@-1@-1@-1@-1@-1@-1@-1@-",
	"1@-@@@-@-1@-@@@-@-1@-@@@-@-",
	"1@-@@@-@-@@@@@@@@@-@-@@@-@-",
	"1@-@@@-@11@-@@@-@11@-@@@-@1",
	"1@-@@@-@1@@@@@@@@@-@1@@@1@-",
	"1@1@@@1@-1@1@@@1@-1@1@@@1@-",
	"1@1@@@1@1@@@@@@@@@1@1@@@1@-",
	"@@-@@0@@1@@-@@0@@1@@-@@0@@1",
}

func TestLUTEncode(t *testing.T) {
	for _, s := range lut_test {
		bin := BinaryEncodedLUTFromString(s)
		echo := StringFromBinaryEncodedLUT(bin)
		if s != echo {
			t.Errorf("s != echo for '%s' != '%s'", s, echo)
		}
	}
}
