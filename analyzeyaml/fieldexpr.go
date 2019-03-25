package analyzeyaml

import (
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/goq/readyaml"
)

func AnalyzeFieldExpr(exprYAML *QuplaFieldExprYAML, module *QuplaModule, scope *Function) (*QuplaFieldExpr, error) {
	var err error
	module.IncStat("numFieldExpr")
	ret := &QuplaFieldExpr{}
	ret.CondExpr, err = AnalyzeExpression(exprYAML.CondExpr, module, scope)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
