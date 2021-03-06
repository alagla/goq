package transform

import (
	. "github.com/lunfardo314/goq/qupla"
	"github.com/lunfardo314/goq/utils"
)

// tree of merges optimized into one merge of many expressions

func OptimizeMerges(def *Function, stats map[string]int) bool {
	before := utils.StatValue("numOptimizedMerges", stats)
	for _, site := range def.Sites {
		if site.NotUsed || site.IsState || site.IsParam || site.NumUses > 1 {
			continue
		}
		site.Assign = optimizeMergesInExpr(site.Assign, stats)
	}
	def.RetExpr = optimizeMergesInExpr(def.RetExpr, stats)
	return before != utils.StatValue("numOptimizedMerges", stats)
}

func optimizeMergesInExpr(expr ExpressionInterface, stats map[string]int) ExpressionInterface {
	_, ok := expr.(*MergeExpr)
	if !ok {
		return transformSubexpressions(expr, func(se ExpressionInterface) ExpressionInterface {
			return optimizeMergesInExpr(se, stats)
		})
	}
	subExpr := make([]ExpressionInterface, 0)
	for _, se := range expr.GetSubexpressions() {
		oe := optimizeMergesInExpr(se, stats)
		if ce, ok := oe.(*MergeExpr); ok {
			utils.IncStat("numOptimizedMerges", stats)
			for _, e := range ce.GetSubexpressions() {
				subExpr = append(subExpr, e)
			}
		} else {
			subExpr = append(subExpr, oe)
		}
	}
	return NewMergeExpression("optimized", subExpr)
}
