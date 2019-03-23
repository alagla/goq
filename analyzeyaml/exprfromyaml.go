package analyzeyaml

import (
	. "fmt"
	. "github.com/lunfardo314/goq/abstract"
	"github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/quplayaml/quplayaml"
)

type ExpressionFactoryFromYAML struct{}

func (ef *ExpressionFactoryFromYAML) AnalyzeExpression(
	dataYAML interface{}, module ModuleInterface, scope FuncDefInterface) (ExpressionInterface, error) {
	switch data := dataYAML.(type) {
	case *QuplaConstNumberYAML:
		return AnalyzeConstNumber(data, module, scope)
	case *QuplaConstTypeNameYAML:
		return qupla.AnalyzeConstTypeName(data, module, scope)
	case *QuplaConstTermYAML:
		return AnalyzeConstTerm(data, module, scope)
	case *QuplaConstExprYAML:
		return AnalyzeConstExpr(data, module, scope)
	case *QuplaCondExprYAML:
		return AnalyzeCondExpr(data, module, scope)
	case *QuplaLutExprYAML:
		return qupla.AnalyzeLutExpr(data, module, scope)
	case *QuplaSliceExprYAML:
		return qupla.AnalyzeSliceExpr(data, module, scope)
	case *QuplaValueExprYAML:
		return qupla.AnalyzeValueExpr(data, module, scope)
	case *QuplaSizeofExprYAML:
		return qupla.AnalyzeSizeofExpr(data, module, scope)
	case *QuplaFuncExprYAML:
		return qupla.AnalyzeFuncExpr(data, module, scope)
	case *QuplaFieldExprYAML:
		return AnalyzeFieldExpr(data, module, scope)
	case *QuplaConcatExprYAML:
		return AnalyzeConcatExpr(data, module, scope)
	case *QuplaMergeExprYAML:
		return AnalyzeMergeExpr(data, module, scope)
	case *QuplaTypeExprYAML:
		return qupla.AnalyzeTypeExpr(data, module, scope)
	case *QuplaNullExprYAML:
		return qupla.AnalyzeNullExpr(data, module, scope)
	case *QuplaExpressionYAML:
		r, err := data.Unwrap()
		if err != nil {
			return nil, err
		}
		if r == nil {
			return &qupla.QuplaNullExpr{}, nil
		}
		return ef.AnalyzeExpression(r, module, scope)
	}
	return nil, Errorf("wrong QuplaYAML object type")
}
