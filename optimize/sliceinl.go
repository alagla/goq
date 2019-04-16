package optimize

import (
	. "github.com/lunfardo314/goq/qupla"
)

// if inline slice does no slicing of the expression at all (takes whole vector)
// return that expression
// if it is slicing of Value expression (constant trit vector)
// optimize the value expression

func optimizeInlineSlices(def *Function, stats map[string]int) bool {
	//if strings.HasPrefix(def.Name, "fixSign"){
	//	fmt.Println("kuku")
	//}
	before := StatValue("numOptimizedInlineSlices", stats)
	for _, site := range def.LocalVars {
		if site.NotUsed || site.IsState || site.IsParam || site.NumUses > 1 {
			continue
		}
		site.Assign = optimizeInlineSlicesInExpr(site.Assign, stats)
	}
	def.RetExpr = optimizeInlineSlicesInExpr(def.RetExpr, stats)
	return before != StatValue("numOptimizedInlineSlices", stats)
}

func optimizeInlineSlicesInExpr(expr ExpressionInterface, stats map[string]int) ExpressionInterface {
	inlineSlice, ok := expr.(*SliceInline)
	if !ok {
		return transformSubexpressions(expr, func(se ExpressionInterface) ExpressionInterface {
			return optimizeInlineSlicesInExpr(se, stats)
		})
	}

	if inlineSlice.NoSlice {
		IncStat("numOptimizedInlineSlices:eliminated", stats)
		return inlineSlice.GetSubExpr(0)
	}
	valueExpr, ok := inlineSlice.GetSubExpr(0).(*ValueExpr)
	if ok {
		IncStat("numOptimizedInlineSlices:toValue", stats)
		return NewValueExpr(valueExpr.TritValue[inlineSlice.Offset:inlineSlice.SliceEnd])
	}
	return inlineSlice
}
