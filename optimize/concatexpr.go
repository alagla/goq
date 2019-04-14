package optimize

import . "github.com/lunfardo314/goq/qupla"

// tree of concats optimized into one concat of many expressions

func optimizeConcatExpr(expr ExpressionInterface, numOptimized *int) ExpressionInterface {
	_, ok := expr.(*ConcatExpr)
	if !ok {
		return optimizeSubxpressions(expr, func(se ExpressionInterface) ExpressionInterface {
			return optimizeConcatExpr(se, numOptimized)
		})
	}
	subExpr := make([]ExpressionInterface, 0)
	for _, se := range expr.GetSubexpressions() {
		oe := optimizeConcatExpr(se, numOptimized)
		if ce, ok := oe.(*ConcatExpr); ok {
			*numOptimized++
			for _, e := range ce.GetSubexpressions() {
				subExpr = append(subExpr, e)
			}
		} else {
			subExpr = append(subExpr, oe)
		}
	}
	return NewConcatExpression("optimized", subExpr)
}
