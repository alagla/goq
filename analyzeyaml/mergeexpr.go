package analyzeyaml

import (
	. "fmt"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/quplayaml/quplayaml"
)

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
	if ret.GetSubExpr(0).Size() != ret.GetSubExpr(1).Size() {
		return nil, Errorf("operand sizes must be equal in merge expression, scope %v", scope.GetName())
	}
	return ret, nil
}
