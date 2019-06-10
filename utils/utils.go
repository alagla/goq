package utils

import (
	"fmt"
	. "github.com/iotaledger/iota.go/kerl"
	. "github.com/iotaledger/iota.go/trinary"
	"math/big"
	"time"
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

var (
	a0 = big.NewInt(-1)
	a1 = big.NewInt(0)
	a2 = big.NewInt(1)
)

func TritToBigInt(t int8) (*big.Int, error) {
	switch t {
	case -1:
		return a0, nil
	case 0:
		return a1, nil
	case 1:
		return a2, nil
	}
	return nil, fmt.Errorf("wrong trit value")
}

var big3 = big.NewInt(3)

func TritsToBigInt(t Trits) (*big.Int, error) {
	if err := ValidTrits(t); err != nil {
		return nil, err
	}
	ret := big.NewInt(0)
	var err error
	var trit *big.Int
	for i := len(t) - 1; i >= 0; i-- {
		if trit, err = TritToBigInt(t[i]); err != nil {
			return nil, err
		}
		ret.Mul(big3, ret)
		ret.Add(ret, trit)
	}
	return ret, nil
}

func MustTritsToBigInt(t Trits) *big.Int {
	bi, err := TritsToBigInt(t)
	if err != nil {
		panic(err)
	}
	return bi
}

type StringSet map[string]struct{}

func (s StringSet) Contains(el string) bool {
	_, exists := s[el]
	return exists
}

func (s StringSet) Append(el string) bool {
	_, exists := s[el]
	s[el] = struct{}{}
	return !exists
}

func (s StringSet) Delete(el string) bool {
	_, exists := s[el]
	delete(s, el)
	return exists
}

func (s StringSet) AppendAll(another StringSet) int {
	var ret int
	for el := range another {
		if s.Append(el) {
			ret++
		}
	}
	return ret
}

func (s StringSet) DeleteAll(another StringSet) int {
	var ret int
	for el := range another {
		if s.Delete(el) {
			ret++
		}
	}
	return ret
}

func (s StringSet) Join(d string) string {
	ret := ""
	first := true
	for str := range s {
		if !first {
			ret += d + str
		} else {
			ret = str
			first = false
		}
	}
	return ret
}

func (s StringSet) List() []string {
	ret := make([]string, 0, len(s))
	for str := range s {
		ret = append(ret, str)
	}
	return ret
}

func UnixMsNow() uint64 {
	return uint64(time.Now().UnixNano()) / uint64(time.Millisecond)
}

func ReprTrits(t Trits) string {
	bi, _ := TritsToBigInt(t)
	return fmt.Sprintf("%v, '%.40s'..", bi, TritsToString(t))
}

// 4 bits -> 3 trits
// enoded_value = dec_value + 7
var b2tencoding = []Trits{
	{-1, 1, -1}, // dec value = -7, endoded value = 0
	{0, 1, -1},  // dec value = -6, endoded value = 1
	{1, 1, -1},  // dec value = -5, endoded value = 2
	{-1, -1, 0}, // dec value = -4, endoded value = 3
	{0, -1, 0},  // dec value = -3, endoded value = 4
	{1, -1, 0},  // dec value = -2, endoded value = 5
	{-1, 0, 0},  // dec value = -1, endoded value = 6
	{0, 0, 0},   // dec value = 0, endoded value = 7
	{1, 0, 0},   // dec value = 1, endoded value = 8
	{-1, 1, 0},  // dec value = 2, endoded value = 9
	{0, 1, 0},   // dec value = 3, endoded value = 10
	{1, 1, 0},   // dec value = 4, endoded value = 11
	{-1, -1, 1}, // dec value = 5, endoded value = 12
	{0, -1, 1},  // dec value = 6, endoded value = 13
	{1, -1, 1},  // dec value = 7, endoded value = 14
	{-1, 0, 1},  // dec value = 8, endoded value = 15
}

func Bytes2Trits(data []byte, lengthMultiple int) Trits {
	var length int
	length = 6 * len(data)
	if lengthMultiple > 0 && length%lengthMultiple != 0 {
		length = length + (lengthMultiple - length%lengthMultiple)
	}
	ret := make(Trits, length)
	for i, b := range data {
		copy(ret[6*i:6*i+3], b2tencoding[b>>4])
		copy(ret[6*i+3:6*i+6], b2tencoding[b&0x0F])
	}
	return ret
}

func Trits2Bytes(trits Trits) ([]byte, error) {
	if len(trits)%6 != 0 {
		return nil, fmt.Errorf("length of trit vertor must be 6*n")
	}
	err := ValidTrits(trits)
	if err != nil {
		return nil, err
	}
	ret := make([]byte, len(trits)/6, len(trits)/6)
	var bt byte
	var decValue, encodedValue int8
	for i := 0; i < len(trits); i += 6 {
		bt = 0
		// 1st half byte
		decValue = trits[i+0] + 3*trits[i+1] + 9*trits[i+2]
		encodedValue = decValue + 7
		if encodedValue < 0 || encodedValue >= 16 {
			return nil, fmt.Errorf("wrong trit combination in half-byte encoding")
		}
		bt |= byte(encodedValue&0x0F) << 4
		// 2nd half byte
		decValue = trits[i+3] + 3*trits[i+4] + 9*trits[i+5]
		encodedValue = decValue + 7
		if encodedValue < 0 || encodedValue >= 16 {
			return nil, fmt.Errorf("wrong trit combination in half-byte encoding")
		}
		bt |= byte(encodedValue & 0x0F)
		ret[i/6] = bt
	}
	return ret, nil
}

func KerlHash243(trits Trits) (Trits, error) {
	k := NewKerl()
	if k == nil {
		return nil, fmt.Errorf("couldn't initialize Kerl instance")
	}
	var err error
	err = k.Absorb(trits)
	if err != nil {
		return nil, fmt.Errorf("absorb() failed: %s", err)
	}
	ts, err := k.Squeeze(243)
	if err != nil {
		return nil, fmt.Errorf("squeeze() failed: %v", err)
	}
	return ts, nil
}
