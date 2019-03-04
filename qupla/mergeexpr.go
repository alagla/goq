package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/quplayaml"
)

type QuplaMergeExpr struct {
	QuplaExprBase
}

func AnalyzeMergeExpr(exprYAML *QuplaMergeExprYAML, module ModuleInterface, scope FuncDefInterface) (*QuplaMergeExpr, error) {
	module.IncStat("numMergeExpr")

	ret := &QuplaMergeExpr{
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
	if ret.subexpr[0].Size() != ret.subexpr[1].Size() {
		return nil, fmt.Errorf("operand sizes must be equal in merge expression, scope %v", scope.GetName())
	}
	return ret, nil
}

func (e *QuplaMergeExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.subexpr[0].Size()
}

func (e *QuplaMergeExpr) Eval(proc ProcessorInterface, result Trits) bool {
	null := proc.Eval(e.subexpr[0], result)
	if null {
		return proc.Eval(e.subexpr[1], result)
	}
	return false
}
