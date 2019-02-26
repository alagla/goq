package quplayaml

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
)

type QuplaConcatExpr struct {
	lhsExpr ExpressionInterface
	rhsExpr ExpressionInterface
}

func AnalyzeConcatExpr(exprYAML *QuplaConcatExprYAML, module ModuleInterface, scope FuncDefInterface) (*QuplaConcatExpr, error) {
	var err error
	module.IncStat("numConcat")

	ret := &QuplaConcatExpr{}
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

func (e *QuplaConcatExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.lhsExpr.Size() + e.rhsExpr.Size()
}

func (e *QuplaConcatExpr) Eval(proc ProcessorInterface, result Trits) bool {
	null := e.lhsExpr.Eval(proc, result)
	if null {
		return true
	}
	return e.rhsExpr.Eval(proc, result[e.lhsExpr.Size():])
}
