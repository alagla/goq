package analyzeyaml

import (
	. "fmt"
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/goq/readyaml"
)

func AnalyzeExpression(dataYAML interface{}, module *QuplaModule, scope *Function) (ExpressionInterface, error) {
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
		return AnalyzeValueExpr(data, module)
	case *QuplaSizeofExprYAML:
		return AnalyzeSizeofExpr(data, module)
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
	case *QuplaNullExprYAML:
		return AnalyzeNullExpr(module)
	case *QuplaExpressionYAML:
		r, err := data.Unwrap()
		if err != nil {
			return nil, err
		}
		if r == nil {
			return &NullExpr{}, nil
		}
		return AnalyzeExpression(r, module, scope)
	}
	return nil, Errorf("wrong QuplaYAML object type")
}
