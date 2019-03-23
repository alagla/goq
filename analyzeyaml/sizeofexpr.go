package analyzeyaml

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/quplayaml/quplayaml"
	"strconv"
)

func AnalyzeSizeofExpr(exprYAML *QuplaSizeofExprYAML, module *QuplaModule) (*QuplaSizeofExpr, error) {
	module.IncStat("numSizeofExpr")

	if exprYAML.Trits == "" {
		return nil, fmt.Errorf("invalid trit string in SizeofExpr '%v'", exprYAML.Trits)
	}
	t := make([]int8, len(exprYAML.Trits))
	for i := range exprYAML.Trits {
		switch exprYAML.Trits[i] {
		case '-':
			t[i] = -1
		case '0':
			t[i] = 0
		case '1':
			t[i] = 1
		default:
			return nil, fmt.Errorf("invalid trit string '%v'", exprYAML.Trits)
		}
	}
	var orig int
	var err error
	if exprYAML.Value == "-" {
		orig = -1
	} else {
		orig, err = strconv.Atoi(exprYAML.Value)
		if err != nil {
			return nil, fmt.Errorf("wrong 'value' field in ValueExpr")
		}
	}

	value := int64(orig)
	var tritValue Trits
	if tritValue, err = NewTrits(t); err != nil {
		return nil, err
	}

	if value != TritsToInt(tritValue) {
		return nil, fmt.Errorf("wrong 'value' ('%v') or 'trits' ('%v') field in value expr",
			exprYAML.Value, exprYAML.Trits)
	}
	return NewQuplaSizeofExpr(value, tritValue), nil
}
