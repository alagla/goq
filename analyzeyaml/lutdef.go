package analyzeyaml

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/qupla"
	"github.com/lunfardo314/goq/utils"
	. "github.com/lunfardo314/quplayaml/quplayaml"
	"strings"
)

var pow3 = []int{1, 3, 9, 27}

func AnalyzeLutDef(name string, defYAML *QuplaLutDefYAML, module *QuplaModule) error {
	ret := &QuplaLutDef{
		Name: name,
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
	ret.InputSize = 0
	ret.OutputSize = 0
	for _, entry := range defYAML.LutTable {
		sides := strings.Split(entry, "=")
		if len(sides) != 2 {
			return fmt.Errorf("wrong LUT entry: %v", entry)
		}
		sides[0] = strings.TrimSpace(sides[0])
		sides[1] = strings.TrimSpace(sides[1])

		if ret.InputSize == 0 {
			ret.InputSize = len(sides[0])
			ret.OutputSize = len(sides[1])
			if ret.InputSize < 1 || ret.InputSize > 3 || ret.OutputSize < 1 {
				return fmt.Errorf("wrong input or output size")
			}
		}
		if len(sides[0]) != ret.InputSize {
			return fmt.Errorf("input len expected to be %v", ret.InputSize)
		}
		if len(sides[1]) != ret.OutputSize {
			return fmt.Errorf("ouput len expected to be %v", ret.OutputSize)
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
	ret.LutLookupTable = make([]Trits, pow3[ret.InputSize])
	for i, inp := range inputs {
		idx := utils.Trits3ToLutIdx(inp)
		if ret.LutLookupTable[idx] != nil {
			return fmt.Errorf("duplicated input in LUT table")
		}
		ret.LutLookupTable[idx] = outputs[i]
	}
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
