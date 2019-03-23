package analyzeyaml

import (
	. "github.com/lunfardo314/goq/abstract"
	"github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/quplayaml/quplayaml"
)

func AnalyzeFieldExpr(exprYAML *QuplaFieldExprYAML, module ModuleInterface, scope FuncDefInterface) (*qupla.QuplaFieldExpr, error) {
	var err error
	module.IncStat("numFieldExpr")
	ret := &qupla.QuplaFieldExpr{}
	ret.CondExpr, err = module.AnalyzeExpression(exprYAML.CondExpr, scope)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
