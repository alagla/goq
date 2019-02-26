package quplayaml

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
)

type QuplaMergeExpr struct {
	lhsExpr ExpressionInterface
	rhsExpr ExpressionInterface
}

func AnalyzeMergeExpr(exprYAML *QuplaMergeExprYAML, module ModuleInterface, scope FuncDefInterface) (*QuplaMergeExpr, error) {
	var err error
	module.IncStat("numMergeExpr")

	ret := &QuplaMergeExpr{}
	ret.lhsExpr, err = module.AnalyzeExpression(exprYAML.Lhs, scope)
	if err != nil {
		return nil, err
	}
	if IsNullExpr(ret.lhsExpr) {
		return nil, fmt.Errorf("constant null in merge expression, scope %v", scope.GetName())
	}
	ret.rhsExpr, err = module.AnalyzeExpression(exprYAML.Rhs, scope)
	if err != nil {
		return nil, err
	}
	if IsNullExpr(ret.rhsExpr) {
		return nil, fmt.Errorf("constant null in merge expression, scope %v", scope.GetName())
	}
	if ret.lhsExpr.Size() != ret.rhsExpr.Size() {
		return nil, fmt.Errorf("operand sizes must be equal in merge expression, scope %v", scope.GetName())
	}
	return ret, nil
}

func (e *QuplaMergeExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.lhsExpr.Size()
}

func (e *QuplaMergeExpr) Eval(proc ProcessorInterface, result Trits) bool {
	null := e.lhsExpr.Eval(proc, result)
	if null {
		return e.rhsExpr.Eval(proc, result)
	}
	return false
}
