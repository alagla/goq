package transform

import (
	. "github.com/lunfardo314/goq/qupla"
	"github.com/lunfardo314/goq/utils"
)

// All not params and not state vars slices which are used only once
// within function will be replaced by SliceInline
// This will eliminate unnecessary Eval call and unnecessary caching

func OptimizeSlices(def *Function, stats map[string]int) bool {
	before := utils.StatValue("numOptimizedSlices", stats)
	for _, site := range def.Sites {
		if site.NotUsed || site.IsState || site.IsParam || site.NumUses > 1 {
			continue
		}
		site.Assign = optimizeSlicesInExpr(site.Assign, stats)
	}
	def.RetExpr = optimizeSlicesInExpr(def.RetExpr, stats)
	return before != utils.StatValue("numOptimizedSlices", stats)
}

func optimizeSlicesInExpr(expr ExpressionInterface, stats map[string]int) ExpressionInterface {
	sliceExpr, ok := expr.(*SliceExpr)
	if !ok {
		return transformSubexpressions(expr, func(se ExpressionInterface) ExpressionInterface {
			return optimizeSlicesInExpr(se, stats)
		})
	}
	site := sliceExpr.Site()
	if site.IsState || site.IsParam || site.NumUses > 1 {
		return expr
	}
	// slice expressions optimize along chain of assignments
	opt := optimizeSlicesInExpr(site.Assign, stats)
	site.NotUsed = true
	utils.IncStat("numOptimizedSlices", stats)
	return NewSliceInline(sliceExpr, opt)
}
