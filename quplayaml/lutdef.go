package quplayaml

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	"strings"
)

type QuplaLutDef struct {
	name           string
	inputSize      int
	outputSize     int
	lutLookupTable []Trits
}

var pow3 = []int{1, 3, 9, 27}

func (lutDef *QuplaLutDef) SetName(name string) {
	lutDef.name = name
}

func AnalyzeLutDef(name string, defYAML *QuplaLutDefYAML, module ModuleInterface) error {
	ret := &QuplaLutDef{
		name: name,
	}
	module.IncStat("numLUTDef")

	if len(defYAML.LutTable) == 0 {
		return fmt.Errorf("no LUT entries found")
	}
	if len(defYAML.LutTable) > 27 {
		return fmt.Errorf("lut table can't have more than 27 entries")
	}
	inputs := make([]Trits, 0, len(defYAML.LutTable))
	outputs := make([]Trits, 0, len(defYAML.LutTable))
	ret.inputSize = 0
	ret.outputSize = 0
	for _, entry := range defYAML.LutTable {
		sides := strings.Split(entry, "=")
		if len(sides) != 2 {
			return fmt.Errorf("wrong LUT entry: %v", entry)
		}
		sides[0] = strings.TrimSpace(sides[0])
		sides[1] = strings.TrimSpace(sides[1])

		if ret.inputSize == 0 {
			ret.inputSize = len(sides[0])
			ret.outputSize = len(sides[1])
			if ret.inputSize < 1 || ret.inputSize > 3 || ret.outputSize < 1 {
				return fmt.Errorf("wrong input or output size")
			}
		}
		if len(sides[0]) != ret.inputSize {
			return fmt.Errorf("input len expected to be %v", ret.inputSize)
		}
		if len(sides[1]) != ret.outputSize {
			return fmt.Errorf("ouput len expected to be %v", ret.outputSize)
		}
		inTrits, err := quplaTritStringToTrits(sides[0])
		if err != nil {
			return err
		}
		inputs = append(inputs, inTrits)
		outTrits, err := quplaTritStringToTrits(sides[1])
		if err != nil {
			return err
		}
		outputs = append(outputs, outTrits)
	}
	// index it to the final table
	ret.lutLookupTable = make([]Trits, pow3[ret.inputSize])
	for i, inp := range inputs {
		idx := tritsToIdx(inp)
		if ret.lutLookupTable[idx] != nil {
			return fmt.Errorf("duplicated input in LUT table")
		}
		ret.lutLookupTable[idx] = outputs[i]
	}
	module.AddLutDef(name, ret)
	return nil
}

func (lutDef *QuplaLutDef) Size() int64 {
	return int64(lutDef.outputSize)
}

func (lutDef *QuplaLutDef) Lookup(res, args Trits) bool {
	t := lutDef.lutLookupTable[tritsToIdx(args)]
	if t == nil {
		return true
	}
	copy(res, t)
	return false
}

func tritsToIdx(trits Trits) int64 {
	return TritsToInt(trits) + int64(pow3[len(trits)]/2)
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
