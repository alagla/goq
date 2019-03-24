package analyzeyaml

import (
	"fmt"
	"github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/quplayaml/quplayaml"
)

func AnalyzeFuncExpr(exprYAML *QuplaFuncExprYAML, module *QuplaModule, scope *QuplaFuncDef) (*QuplaFuncExpr, error) {
	funcDef := module.FindFuncDef(exprYAML.Name)
	if funcDef == nil {
		return nil, fmt.Errorf("can't find function '%v'", exprYAML.Name)
	}
	ret := NewQuplaFuncExpr(exprYAML.Source, funcDef)
	module.IncStat("numFuncExpr")

	var tmpSubexpr = make([]abstract.ExpressionInterface, 0, len(exprYAML.Args))
	for _, arg := range exprYAML.Args {
		if fe, err := AnalyzeExpression(arg, module, scope); err != nil {
			return nil, err
		} else {
			ret.AppendSubExpr(fe)
			tmpSubexpr = append(tmpSubexpr, fe)
		}
	}
	if err := funcDef.CheckArgSizes(tmpSubexpr); err != nil {
		return nil, err
	}
	return ret, nil
}