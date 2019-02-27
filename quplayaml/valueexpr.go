package quplayaml

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
)

type QuplaValueExpr struct {
	TritValue Trits
}

func AnalyzeValueExpr(exprYAML *QuplaValueExprYAML, module ModuleInterface, _ FuncDefInterface) (*QuplaValueExpr, error) {
	module.IncStat("numValueExpr")

	if exprYAML.Trits == "" {
		return nil, fmt.Errorf("invalid trit string '%v'", exprYAML.Trits)
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
	ret := &QuplaValueExpr{}
	var err error
	if ret.TritValue, err = NewTrits(t); err != nil {
		return nil, err
	}
	return ret, nil
}

func (e *QuplaValueExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return int64(len(e.TritValue))
}

func (e *QuplaValueExpr) Eval(_ ProcessorInterface, result Trits) bool {
	if e.TritValue == nil {
		return true
	}
	copy(result, e.TritValue)
	return false
}
