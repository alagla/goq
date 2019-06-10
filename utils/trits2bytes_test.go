package utils

import (
	"testing"
)

var testStrings = []string{
	"kuku",
	"be or not to be",
	" !\"#$%&\\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~",
}

var testArray = []byte{1, 2, 0xF1, 0xFF, 0xEF, 100, 142, 22, 255, 0, 0, 0xFF, 0x01}

func TestTrits2Bytes(t *testing.T) {
	for i := range testStrings {
		trits := Bytes2Trits([]byte(testStrings[i]))
		str, err := Trits2Bytes(trits)
		if err != nil {
			t.Errorf("%v", err)
		} else if string(str) != testStrings[i] {
			t.Errorf("%s  !=  %s", testStrings[i], string(str))
		}
	}
	trits := Bytes2Trits(testArray)
	data, err := Trits2Bytes(trits)
	if err != nil {
		t.Errorf("%v", err)
	} else if !equal(data, testArray) {
		t.Errorf("%v  !=  %v", testArray, data)
	}
}

func equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
