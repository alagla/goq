package analyzeyaml

import (
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/goq/readyaml"
)

func AnalyzeFieldExpr(exprYAML *QuplaFieldExprYAML, module *QuplaModule, scope *Function) (*QuplaFieldExpr, error) {
	module.IncStat("numFieldExpr")
	ret := &QuplaFieldExpr{}
	condExpr, err := AnalyzeExpression(exprYAML.CondExpr, module, scope)
	if err != nil {
		return nil, err
	}
	ret.AppendSubExpr(condExpr)
	return ret, nil
}
