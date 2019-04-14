package optimize

import . "github.com/lunfardo314/goq/qupla"

// expression 'expr' must be within context of the function 'def'
// All non param and non state slices which are used only once
// within function will be replaced by SliceInline
// This will eliminate unnecessary call Eval call and unnecessary caching

func optimizeSlices(def *Function, expr ExpressionInterface) ExpressionInterface {
	sliceExpr, ok := expr.(*SliceExpr)
	if !ok {
		return optimizeSubxpressions(expr, func(se ExpressionInterface) ExpressionInterface {
			return optimizeSlices(def, se)
		})
	}
	site := sliceExpr.Site()
	if site.IsState || site.IsParam || site.NumUses > 1 {
		return expr
	}
	// slice expressions optimize along chain of assignments
	opt := optimizeSlices(def, def.LocalVars[site.Idx].Assign)
	site.NotUsed = true
	return NewSliceInline(sliceExpr, opt)
}
