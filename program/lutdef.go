package program

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	"strings"
)

type QuplaLutDef struct {
	LutTable []string `yaml:"lutTable"`
	//----
	name           string
	analyzed       bool
	inputSize      int
	outputSize     int
	lutLookupTable []Trits
}

var pow3 = []int{1, 3, 9, 27}

func (lutDef *QuplaLutDef) SetName(name string) {
	lutDef.name = name
}

func (lutDef *QuplaLutDef) Analyze(module *QuplaModule) (*QuplaLutDef, error) {
	if lutDef.analyzed {
		return lutDef, nil
	}
	lutDef.analyzed = true
	module.IncStat("numLUTDef")

	if len(lutDef.LutTable) == 0 {
		return nil, fmt.Errorf("No LUT entries found")
	}
	if len(lutDef.LutTable) > 27 {
		return nil, fmt.Errorf("lut table can't have more than 27 entries")
	}
	inputs := make([]Trits, 0, len(lutDef.LutTable))
	outputs := make([]Trits, 0, len(lutDef.LutTable))
	lutDef.inputSize = 0
	lutDef.outputSize = 0
	for _, entry := range lutDef.LutTable {
		sides := strings.Split(entry, "=")
		if len(sides) != 2 {
			return nil, fmt.Errorf("wrong LUT entry: %v", entry)
		}
		sides[0] = strings.TrimSpace(sides[0])
		sides[1] = strings.TrimSpace(sides[1])

		if lutDef.inputSize == 0 {
			lutDef.inputSize = len(sides[0])
			lutDef.outputSize = len(sides[1])
			if lutDef.inputSize < 1 || lutDef.inputSize > 3 || lutDef.outputSize < 1 {
				return nil, fmt.Errorf("wrong input or output size")
			}
		}
		if len(sides[0]) != lutDef.inputSize {
			return nil, fmt.Errorf("input len expected to be %v", lutDef.inputSize)
		}
		if len(sides[1]) != lutDef.outputSize {
			return nil, fmt.Errorf("ouput len expected to be %v", lutDef.outputSize)
		}
		inTrits, err := quplaTritStringToTrits(sides[0])
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, inTrits)
		outTrits, err := quplaTritStringToTrits(sides[1])
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, outTrits)
	}
	// index it to the final table
	lutDef.lutLookupTable = make([]Trits, pow3[lutDef.inputSize])
	for i, inp := range inputs {
		idx := tritsToIdx(inp)
		if lutDef.lutLookupTable[idx] != nil {
			return nil, fmt.Errorf("duplicated input in LUT table")
		}
		lutDef.lutLookupTable[idx] = outputs[i]
	}
	return lutDef, nil
}

func (lutDef *QuplaLutDef) Size() int64 {
	return int64(lutDef.outputSize)
}

func (lutDef *QuplaLutDef) lookupWithCheck(t Trits) (Trits, error) {
	if len(t) != lutDef.inputSize {
		return nil, fmt.Errorf("wrong input size")
	}
	return lutDef.lookup(t), nil
}

func (lutDef *QuplaLutDef) lookup(t Trits) Trits {
	return lutDef.lutLookupTable[tritsToIdx(t)]
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
