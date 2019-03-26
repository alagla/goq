package analyzeyaml

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/goq/readyaml"
	"strings"
)

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

	var inTrits, outTrits Trits
	var err error

	for _, entry := range defYAML.LutTable {
		sides = strings.Split(entry, "=")
		if len(sides) != 2 {
			return fmt.Errorf("wrong LUT entry: %v", entry)
		}
		sides[0] = strings.TrimSpace(sides[0])
		sides[1] = strings.TrimSpace(sides[1])

		if inTrits, err = quplaTritStringToTrits(sides[0]); err != nil || len(inTrits) != inputSize {
			return fmt.Errorf("wrong input trits in LUT entry: %v", entry)

		}
		if outTrits, err = quplaTritStringToTrits(sides[1]); err != nil || len(outTrits) != outputSize {
			return fmt.Errorf("wrong output trits in LUT entry: %v", entry)
		}
		idx := Trits3ToLutIdx(inTrits)
		if lutTable[idx] != nil {
			return fmt.Errorf("duplicated LUT entry '%v' in LUT '%v'", entry, name)
		}
		lutTable[idx] = outTrits
	}
	ret := NewLUTDef(name, inputSize, outputSize, lutTable)
	module.AddLutDef(name, ret)
	return nil
}

func quplaTritStringToTrits(str string) (Trits, error) {
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
