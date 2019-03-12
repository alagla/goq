package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/quplayaml/quplayaml"
	"strconv"
)

type QuplaSizeofExpr struct {
	QuplaExprBase
	Value     int64
	TritValue Trits
}

func AnalyzeSizeofExpr(exprYAML *QuplaSizeofExprYAML, module ModuleInterface, _ FuncDefInterface) (*QuplaSizeofExpr, error) {
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

	ret := &QuplaSizeofExpr{}
	ret.Value = int64(orig)
	if ret.TritValue, err = NewTrits(t); err != nil {
		return nil, err
	}

	if ret.Value != TritsToInt(ret.TritValue) {
		return nil, fmt.Errorf("wrong 'value' ('%v') or 'trits' ('%v') field in value expr",
			exprYAML.Value, exprYAML.Trits)
	}
	return ret, nil
}

func (e *QuplaSizeofExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return int64(len(e.TritValue))
}

func (e *QuplaSizeofExpr) Eval(_ ProcessorInterface, result Trits) bool {
	if e.TritValue == nil {
		return true
	}
	copy(result, e.TritValue)
	return false
}
