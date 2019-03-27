package analyzeyaml

import (
	"fmt"
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/goq/readyaml"
)

func AnalyzeLutExpr(exprYAML *QuplaLutExprYAML, module *QuplaModule, scope *Function) (*LutExpr, error) {
	var ae ExpressionInterface
	var err error
	ret := &LutExpr{
		ExpressionBase: NewExpressionBase(exprYAML.Source),
	}
	ret.LutDef, err = module.FindLUTDef(exprYAML.Name)
	if err != nil {
		return nil, err
	}
	module.IncStat("numLUTExpr")

	for _, a := range exprYAML.Args {
		ae, err = AnalyzeExpression(a, module, scope)
		if err != nil {
			return nil, err
		}
		if err = RequireSize(ae, 1); err != nil {
			return nil, fmt.Errorf("LUT expression with '%v': %v", ret.LutDef.Name, err)
		}
		ret.AppendSubExpr(ae)
	}
	if ret.LutDef.InputSize != ret.NumSubExpr() {
		return nil, fmt.Errorf("idx arg doesnt't match input dimension of the LUT %v", ret.LutDef.Name)
	}
	return ret, nil
}
