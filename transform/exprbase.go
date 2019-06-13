package transform

import (
	. "github.com/lunfardo314/goq/qupla"
)

// calls optFun for each subexpression and replace

func transformSubexpressions(expr ExpressionInterface, optFun func(ExpressionInterface) ExpressionInterface) ExpressionInterface {
	subExpr := make([]ExpressionInterface, 0)
	var opt ExpressionInterface
	for _, se := range expr.GetSubexpressions() {
		opt = optFun(se)
		subExpr = append(subExpr, opt)
	}
	expr.SetSubexpressions(subExpr)
	return expr
}
