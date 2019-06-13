package transform

import (
	. "github.com/lunfardo314/goq/cfg"
	. "github.com/lunfardo314/goq/qupla"
)

func optimizeFunction(def *Function, stats map[string]int) bool {
	var optSlices, optInlineSlices, optConcats, optMerges, optInlineCalls bool

	if Config.OptimizeOneTimeSites {
		optSlices = optimizeSlices(def, stats)
	}
	if Config.OptimizeInlineSlices {
		optInlineSlices = optimizeInlineSlices(def, stats)
	}
	if Config.OptimizeConcats {
		optConcats = optimizeConcats(def, stats)
	}
	if Config.OptimizeMerges {
		optMerges = optimizeMerges(def, stats)
	}
	if Config.OptimizeFunCallsInline {
		optInlineCalls = optimizeFunctionByInlining(def, stats)
	}
	return optSlices || optInlineSlices || optConcats || optMerges || optInlineCalls
}

func IncStat(key string, stats map[string]int) {
	_, ok := stats[key]
	if !ok {
		stats[key] = 0
	}
	stats[key]++
}

func StatValue(key string, stats map[string]int) int {
	v, ok := stats[key]
	if !ok {
		return 0
	}
	return v
}

func isExpandableInline(funExpr *FunctionExpr, ctx *Function) bool {
	if ctx.WasInline(funExpr.FuncDef.Name) {
		return false // recursive expansion inline is not allowed
	}
	for _, site := range funExpr.FuncDef.Sites {
		if site.IsState {
			return false
		}
		if !site.IsParam && !site.NotUsed {
			return false
		}
	}
	return true
}

func optimizeFunctionByInlining(def *Function, stats map[string]int) bool {
	before := StatValue("numFuncCallInlined", stats)
	for _, site := range def.Sites {
		if site.NotUsed || site.IsState || site.IsParam || site.NumUses > 1 {
			continue
		}
		site.Assign = expandInlineExpr(site.Assign, def, stats)
	}
	def.RetExpr = expandInlineExpr(def.RetExpr, def, stats)
	return before != StatValue("numFuncCallInlined", stats)

}

func expandInlineExpr(expr ExpressionInterface, ctx *Function, stats map[string]int) ExpressionInterface {
	var ret ExpressionInterface
	if funcExpr, ok := expr.(*FunctionExpr); ok && isExpandableInline(funcExpr, ctx) {
		ret = expandInlineFuncCall(funcExpr)
		ctx.AppendInline(funcExpr.FuncDef.Name)
		IncStat("numFuncCallInlined", stats)
		return ret
	}
	ret = expr.Copy()
	transformSubexpressions(ret, func(se ExpressionInterface) ExpressionInterface {
		return expandInlineExpr(se, ctx, stats)
	})
	return ret
}

func expandInlineFuncCall(funExpr *FunctionExpr) ExpressionInterface {
	return inlineCopy(funExpr.FuncDef.RetExpr, funExpr)
}

func inlineCopy(expr ExpressionInterface, scope *FunctionExpr) ExpressionInterface {
	if sliceExpr, ok := expr.(*SliceExpr); ok {
		if !sliceExpr.Site().IsParam {
			panic("can't expand inline non-param slice")
		}
		return NewSliceInline(sliceExpr, scope.GetSubExpr(sliceExpr.Site().Idx))
	}
	ret := expr.Copy()
	transformSubexpressions(ret, func(se ExpressionInterface) ExpressionInterface {
		return inlineCopy(se, scope)
	})
	return ret
}

//func InlineExpression(expr ExpressionInterface, def *Function) ExpressionInterface {
//	switch e := expr.(type) {
//	case *SliceExpr:
//		panic("can't inline slice expression")
//	case *FunctionExpr:
//		return ExpandInlineFunCall(e, def)
//	}
//}
//
//func ExpandInlineFunCall(funExpr *FunctionExpr, def *Function) ExpressionInterface {
//	if !funExpr.FuncDef.ZeroInternalSites() || funExpr.FuncDef == def {
//		// inline only if there's no internal sites
//		// don't do recursive inlining
//		return funExpr
//	}
//	return funExpr.FuncDef.RetExpr.InlineCopy(funExpr)
//}
