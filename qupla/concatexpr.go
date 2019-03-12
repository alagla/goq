package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/quplayaml/quplayaml"
)

type QuplaConcatExpr struct {
	QuplaExprBase
}

func AnalyzeConcatExpr(exprYAML *QuplaConcatExprYAML, module ModuleInterface, scope FuncDefInterface) (*QuplaConcatExpr, error) {
	module.IncStat("numConcat")

	ret := &QuplaConcatExpr{
		QuplaExprBase: NewQuplaExprBase(exprYAML.Source),
	}
	if lhsExpr, err := module.AnalyzeExpression(exprYAML.Lhs, scope); err != nil {
		return nil, err
	} else {
		ret.AppendSubExpr(lhsExpr)
	}
	if rhsExpr, err := module.AnalyzeExpression(exprYAML.Rhs, scope); err != nil {
		return nil, err
	} else {
		ret.AppendSubExpr(rhsExpr)
	}
	return ret, nil
}

func (e *QuplaConcatExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.subexpr[0].Size() + e.subexpr[1].Size()
}

func (e *QuplaConcatExpr) Eval(proc ProcessorInterface, result Trits) bool {
	null := proc.Eval(e.subexpr[0], result)
	if null {
		return true
	}
	return proc.Eval(e.subexpr[1], result[e.subexpr[0].Size():])
}
