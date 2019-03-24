package analyzeyaml

import (
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/quplayaml/quplayaml"
)

func AnalyzeConcatExpr(exprYAML *QuplaConcatExprYAML, module *QuplaModule, scope *Function) (*ConcatExpr, error) {
	module.IncStat("numConcat")

	ret := &ConcatExpr{
		ExpressionBase: NewExpressionBase(exprYAML.Source),
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
	return ret, nil
}
