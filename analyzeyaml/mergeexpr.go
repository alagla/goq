package analyzeyaml

import (
	. "fmt"
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/quplayaml/quplayaml"
)

func AnalyzeMergeExpr(exprYAML *QuplaMergeExprYAML, module *QuplaModule, scope *QuplaFuncDef) (*QuplaMergeExpr, error) {
	module.IncStat("numMergeExpr")

	ret := &QuplaMergeExpr{
		QuplaExprBase: NewQuplaExprBase(exprYAML.Source),
	}

	if lhsExpr, err := AnalyzeExpression(exprYAML.Lhs, module, scope); err != nil {
		return nil, err
	} else {
		ret.AppendSubExpr(lhsExpr)
	}
	if rhsExpr, err := AnalyzeExpression(exprYAML.Rhs, module, scope); err != nil {
		return nil, err
	} else {
		ret.AppendSubExpr(rhsExpr)
	}
	if ret.GetSubExpr(0).Size() != ret.GetSubExpr(1).Size() {
		return nil, Errorf("operand sizes must be equal in merge expression, scope %v", scope.Name)
	}
	return ret, nil
}
