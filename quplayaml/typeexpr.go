package quplayaml

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
)

// ----- ?????? do we need it?
type QuplaTypeExpr struct {
	expr   ExpressionInterface
	size   int64
	fields map[string]ExpressionInterface
}

func AnalyzeTypeExpr(exprYAML *QuplaTypeExprYAML, module ModuleInterface, scope FuncDefInterface) (*QuplaTypeExpr, error) {
	ret := &QuplaTypeExpr{
		fields: make(map[string]ExpressionInterface),
	}
	var err error
	module.IncStat("numTypeExpr")

	if ret.expr, err = module.AnalyzeExpression(exprYAML.TypeExpr, scope); err != nil {
		return nil, err
	}

	if ret.size, err = GetConstValue(ret.expr); err != nil {
		return nil, err
	}

	var fe ExpressionInterface
	for name, expr := range exprYAML.Fields {
		if fe, err = module.AnalyzeExpression(expr, scope); err != nil {
			return nil, err
		}
		ret.fields[name] = fe
	}
	if err = ret.CheckSize(); err != nil {
		return nil, err
	}
	return ret, nil
}

func (e *QuplaTypeExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.size
}

func (e *QuplaTypeExpr) CheckSize() error {
	var sumFld int64

	for _, f := range e.fields {
		sumFld += f.Size()
	}
	if sumFld != e.size {
		return fmt.Errorf("sum of field sizes != type end")
	}
	return nil
}

func (e *QuplaTypeExpr) Eval(_ *CallFrame, _ Trits) bool {
	return true
}
