package analyzeyaml

import (
	"fmt"
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/goq/readyaml"
)

func AnalyzeFuncExpr(exprYAML *QuplaFuncExprYAML, module *QuplaModule, scope *Function) (*FunctionExpr, error) {
	funcDef := module.FindFuncDef(exprYAML.Name)
	if funcDef == nil {
		return nil, fmt.Errorf("can't find function '%v'", exprYAML.Name)
	}
	ret := NewFunctionExpr(exprYAML.Source, funcDef)
	module.IncStat("numFuncExpr")

	var tmpSubexpr = make([]ExpressionInterface, 0, len(exprYAML.Args))
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
