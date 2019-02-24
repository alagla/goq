package program

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
)

type QuplaValueExpr struct {
	Trits  string `yaml:"trits"`
	Trytes string `yaml:"trytes"`
	//-----
	TritValue Trits
}

func (e *QuplaValueExpr) Analyze(module *QuplaModule, scope *QuplaFuncDef) (ExpressionInterface, error) {
	module.IncStat("numValueExpr")

	if e.Trits == "" {
		return nil, fmt.Errorf("invalid trit string '%v'", e.Trits)
	}
	t := make([]int8, len(e.Trits))
	for i := range e.Trits {
		switch e.Trits[i] {
		case '-':
			t[i] = -1
		case '0':
			t[i] = 0
		case '1':
			t[i] = 1
		default:
			return nil, fmt.Errorf("invalid trit string '%v'", e.Trits)
		}
	}
	var err error
	if e.TritValue, err = NewTrits(t); err != nil {
		return nil, err
	}
	return e, nil
}

func (e *QuplaValueExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return int64(len(e.TritValue))
}

func (e *QuplaValueExpr) Eval(_ Trits) bool {
	return true
}
