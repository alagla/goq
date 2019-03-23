package analyzeyaml

import (
	"fmt"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/quplayaml/quplayaml"
)

func AnalyzeLutExpr(exprYAML *QuplaLutExprYAML, module *QuplaModule, scope *QuplaFuncDef) (*QuplaLutExpr, error) {
	var ae ExpressionInterface
	var err error
	ret := &QuplaLutExpr{
		QuplaExprBase: NewQuplaExprBase(exprYAML.Source),
	}
	ret.LutDef, err = module.FindLUTDef(exprYAML.Name)
	if err != nil {
		return nil, err
	}
	module.IncStat("numLUTExpr")

	ret.ArgExpr = make([]ExpressionInterface, 0, len(exprYAML.Args))
	for _, a := range exprYAML.Args {
		ae, err = AnalyzeExpression(a, module, scope)
		if err != nil {
			return nil, err
		}
		if err = RequireSize(ae, 1); err != nil {
			return nil, fmt.Errorf("LUT expression with '%v': %v", ret.LutDef.Name, err)
		}
		ret.ArgExpr = append(ret.ArgExpr, ae)
	}
	if ret.LutDef.InputSize != len(ret.ArgExpr) {
		return nil, fmt.Errorf("idx arg doesnt't match input dimension of the LUT %v", ret.LutDef.Name)
	}
	return ret, nil
}
