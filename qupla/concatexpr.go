package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/quplayaml"
)

type QuplaConcatExpr struct {
	QuplaExprBase
	lhsExpr ExpressionInterface
	rhsExpr ExpressionInterface
}

func AnalyzeConcatExpr(exprYAML *QuplaConcatExprYAML, module ModuleInterface, scope FuncDefInterface) (*QuplaConcatExpr, error) {
	var err error
	module.IncStat("numConcat")

	ret := &QuplaConcatExpr{
		QuplaExprBase: NewQuplaExprBase(exprYAML.Source),
	}
	if ret.lhsExpr, err = module.AnalyzeExpression(exprYAML.Lhs, scope); err != nil {
		return nil, err
	}
	if ret.rhsExpr, err = module.AnalyzeExpression(exprYAML.Rhs, scope); err != nil {
		return nil, err
	}
	if ret.rhsExpr.Size() == 0 || ret.lhsExpr.Size() == 0 {
		return nil, fmt.Errorf("size of concat opeation can't be 0: scope '%v'", scope.GetName())
	}
	return ret, nil
}

func (e *QuplaConcatExpr) HasState() bool {
	return e.lhsExpr.HasState() || e.rhsExpr.HasState()
}

func (e *QuplaConcatExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.lhsExpr.Size() + e.rhsExpr.Size()
}

func (e *QuplaConcatExpr) Eval(proc ProcessorInterface, result Trits) bool {
	null := proc.Eval(e.lhsExpr, result)
	if null {
		return true
	}
	return proc.Eval(e.rhsExpr, result[e.lhsExpr.Size():])
}
