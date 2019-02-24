package program

import "github.com/iotaledger/iota.go/trinary"

type QuplaFieldExpr struct {
	FieldName       string                  `yaml:"fieldName"`
	CondExprWrapper *QuplaExpressionWrapper `yaml:"condExpr"`
	//---
	condExpr ExpressionInterface
}

func (e *QuplaFieldExpr) Analyze(module *QuplaModule, scope *QuplaFuncDef) (ExpressionInterface, error) {
	var err error
	module.IncStat("numFieldExpr")
	e.condExpr, err = e.CondExprWrapper.Analyze(module, scope)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (e *QuplaFieldExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.condExpr.Size()
}

func (e *QuplaFieldExpr) Eval(_ *CallFrame, _ trinary.Trits) bool {
	return true
}
