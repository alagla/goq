package analyzeyaml

import (
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/goq/readyaml"
)

func AnalyzeConcatExpr(exprYAML *QuplaConcatExprYAML, module *QuplaModule, scope *Function) (*ConcatExpr, error) {
	module.IncStat("numConcat")
	lhsExpr, err := AnalyzeExpression(exprYAML.Lhs, module, scope)
	if err != nil {
		return nil, err
	}
	rhsExpr, err := AnalyzeExpression(exprYAML.Rhs, module, scope)
	if err != nil {
		return nil, err
	}
	return NewConcatExpression(exprYAML.Source, []ExpressionInterface{lhsExpr, rhsExpr}), nil
}
