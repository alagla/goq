package analyzeyaml

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/goq/readyaml"
	"strings"
)

// lookup table entry is adjusted for 3 inputs regardless of the input size
// this is needed for abra

func AnalyzeLutDef(name string, defYAML *QuplaLutDefYAML, module *QuplaModule) error {
	module.IncStat("numLUTDef")

	if len(defYAML.LutTable) == 0 {
		return fmt.Errorf("no LUT entries found")
	}
	if len(defYAML.LutTable) > 27 {
		return fmt.Errorf("lut table can't have more than 27 entries")
	}
	lutTable := make([]Trits, 27)

	sides := strings.Split(defYAML.LutTable[0], "=")
	sides[0] = strings.TrimSpace(sides[0])
	sides[1] = strings.TrimSpace(sides[1])
	inputSize := len(sides[0])
	if inputSize != 1 && inputSize != 2 && inputSize != 3 {
		return fmt.Errorf("lut input size can be 1, 2 or 3 only")
	}
	outputSize := len(sides[1])

	var outTrits Trits
	var entries3 []Trits
	var err error

	for _, entry := range defYAML.LutTable {
		sides = strings.Split(entry, "=")
		if len(sides) != 2 {
			return fmt.Errorf("wrong LUT entry: %v", entry)
		}
		sides[0] = strings.TrimSpace(sides[0])
		sides[1] = strings.TrimSpace(sides[1])

		if outTrits, err = quplaTritStringToTrits(sides[1], outputSize); err != nil {
			return fmt.Errorf("wrong output trits in LUT entry: %v", entry)
		}

		if entries3, err = genAll3TritEntries(sides[0], inputSize); err != nil {
			return fmt.Errorf("wrong input trits in LUT entry: %v", entry)

		}
		for _, entry := range entries3 {
			idx := Trits3ToLutIdx(entry)
			if lutTable[idx] != nil {
				return fmt.Errorf("duplicated LUT entry '%v' in LUT '%v'", entry, name)
			}
			lutTable[idx] = outTrits
		}
	}
	ret := NewLUTDef(name, inputSize, outputSize, lutTable)
	module.AddLutDef(name, ret)
	return nil
}

func genAll3TritEntries(str string, expectedSize int) ([]Trits, error) {
	trits, err := quplaTritStringToTrits(str, expectedSize)
	if err != nil {
		return nil, err
	}
	switch len(trits) {
	case 3:
		return []Trits{trits}, nil
	case 2:
		ret := make([]Trits, 0, 3)
		for t2 := int8(-1); t2 <= 1; t2++ {
			e := Trits{trits[0], trits[1], t2}
			ret = append(ret, e)
		}
		return ret, nil
	case 1:
		ret := make([]Trits, 0, 9)
		for t1 := int8(-1); t1 <= 1; t1++ {
			for t2 := int8(-1); t2 <= 1; t2++ {
				e := Trits{trits[0], t1, t2}
				ret = append(ret, e)
			}
		}
		return ret, nil
	}
	panic(fmt.Errorf("inconsistency in genAll3TritEntries"))
}

func quplaTritStringToTrits(str string, expectedSize int) (Trits, error) {
	if len(str) != expectedSize {
		return nil, fmt.Errorf("unexpected size of trit string")
	}
	ret := make([]int8, len(str))
	var idx int8
	for i, s := range str {
		idx = int8(strings.Index("-01", string(s)))
		if idx < 0 {
			return nil, fmt.Errorf("wrong character in trit string %v", str)
		}
		ret[i] = idx - 1
	}
	return ret, nil
}
