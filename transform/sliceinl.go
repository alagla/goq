package transform

import (
	. "github.com/lunfardo314/goq/qupla"
	"github.com/lunfardo314/goq/utils"
)

// if inline slice does no slicing of the expression at all (takes whole vector)
// return that expression
// if it is slicing of Value expression (constant trit vector)
// optimize the value expression

func OptimizeInlineSlices(def *Function, stats map[string]int) bool {
	before := utils.StatValue("numOptimizedInlineSlices", stats)
	for _, site := range def.Sites {
		if site.NotUsed || site.IsState || site.IsParam || site.NumUses > 1 {
			continue
		}
		site.Assign = optimizeInlineSlicesInExpr(site.Assign, stats)
	}
	def.RetExpr = optimizeInlineSlicesInExpr(def.RetExpr, stats)
	return before != utils.StatValue("numOptimizedInlineSlices", stats)
}

func optimizeInlineSlicesInExpr(expr ExpressionInterface, stats map[string]int) ExpressionInterface {
	inlineSlice, ok := expr.(*SliceInline)
	if !ok {
		return transformSubexpressions(expr, func(se ExpressionInterface) ExpressionInterface {
			return optimizeInlineSlicesInExpr(se, stats)
		})
	}

	if inlineSlice.NoSlice {
		utils.IncStat("numOptimizedInlineSlices:eliminated", stats)
		return optimizeInlineSlicesInExpr(inlineSlice.GetSubExpr(0), stats)
	} else {
	}
	valueExpr, ok := inlineSlice.GetSubExpr(0).(*ValueExpr)
	if ok {
		utils.IncStat("numOptimizedInlineSlices:toValue", stats)
		return NewValueExpr(valueExpr.TritValue[inlineSlice.Offset:inlineSlice.SliceEnd])
	}
	return inlineSlice
}
