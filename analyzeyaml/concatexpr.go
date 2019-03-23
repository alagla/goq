package analyzeyaml

import (
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/quplayaml/quplayaml"
)

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
