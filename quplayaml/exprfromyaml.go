package quplayaml

import (
	"fmt"
	. "github.com/lunfardo314/goq/abstract"
)

type ExpressionFactoryFromYAML struct{}

func (ef *ExpressionFactoryFromYAML) AnalyzeExpression(
	dataYAML interface{}, module ModuleInterface, scope FuncDefInterface) (ExpressionInterface, error) {
	switch data := dataYAML.(type) {
	case *QuplaConstNumberYAML:
		return AnalyzeConstNumber(data, module, scope)
	case *QuplaConstTypeNameYAML:
		return AnalyzeConstTypeName(data, module, scope)
	case *QuplaConstTermYAML:
		return AnalyzeConstTerm(data, module, scope)
	case *QuplaConstExprYAML:
		return AnalyzeConstExpr(data, module, scope)
	case *QuplaCondExprYAML:
		return AnalyzeCondExpr(data, module, scope)
	case *QuplaLutExprYAML:
		return AnalyzeLutExpr(data, module, scope)
	case *QuplaSliceExprYAML:
		return AnalyzeSliceExpr(data, module, scope)
	case *QuplaValueExprYAML:
		return AnalyzeValueExpr(data, module, scope)
	case *QuplaFuncExprYAML:
		return AnalyzeFuncExpr(data, module, scope)
	case *QuplaFieldExprYAML:
		return AnalyzeFieldExpr(data, module, scope)
	case *QuplaConcatExprYAML:
		return AnalyzeConcatExpr(data, module, scope)
	case *QuplaMergeExprYAML:
		return AnalyzeMergeExpr(data, module, scope)
	case *QuplaTypeExprYAML:
		return AnalyzeTypeExpr(data, module, scope)
	case *QuplaExpressionYAML:
		r, err := data.Unwrap()
		if err != nil {
			return nil, err
		}
		if r == nil {
			return &QuplaNullExpr{}, nil
		}
		return ef.AnalyzeExpression(r, module, scope)
	}
	return nil, fmt.Errorf("wrong QuplaYAML object type")
}
