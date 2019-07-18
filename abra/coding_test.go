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
