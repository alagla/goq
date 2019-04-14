package optimize

import . "github.com/lunfardo314/goq/qupla"

// if inline slice does no slicing of the expression at all (takes whole vector)
// return that expression
// if it is slicing of Value expression (constant trit vector)
// optimize the value expression

func optimizeInlineSlices(expr ExpressionInterface, numOptimized *int) ExpressionInterface {
	inlineSlice, ok := expr.(*SliceInline)
	if !ok {
		return optimizeSubxpressions(expr, func(se ExpressionInterface) ExpressionInterface {
			return optimizeInlineSlices(se, numOptimized)
		})
	}
	*numOptimized++

	if inlineSlice.NoSlice {
		return inlineSlice.Expr
	}
	valueExpr, ok := inlineSlice.Expr.(*ValueExpr)
	if !ok {
		return inlineSlice
	}
	return NewValueExpr(valueExpr.TritValue[inlineSlice.Offset:inlineSlice.SliceEnd])
}
