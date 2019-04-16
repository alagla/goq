package optimize

import (
	. "github.com/lunfardo314/goq/cfg"
	. "github.com/lunfardo314/goq/qupla"
)

func optimizeFunction(def *Function, stats map[string]int) bool {
	var optSlices, optInlineSlices, optConcats, optMerges bool

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
	return optSlices || optInlineSlices || optConcats || optMerges
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
