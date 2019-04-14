package optimize

import (
	. "github.com/lunfardo314/goq/cfg"
	. "github.com/lunfardo314/goq/qupla"
)

func optimizeFunction(def *Function) bool {
	var numOptimizedSlices, numOptimizedInlineSlices, numOptimizedConcat int

	if Config.OptimizeOneTimeSites {
		def.RetExpr = optimizeSlices(def, def.RetExpr, &numOptimizedSlices)
	}
	if Config.OptimizeInlineSlices {
		def.RetExpr = optimizeInlineSlices(def.RetExpr, &numOptimizedInlineSlices)
	}
	if Config.OptimizeConcats {
		def.RetExpr = optimizeConcatExpr(def.RetExpr, &numOptimizedConcat)
	}
	return numOptimizedSlices+numOptimizedInlineSlices+numOptimizedConcat > 0
}

func InlineExpression(expr ExpressionInterface, def *Function) ExpressionInterface {
	switch e := expr.(type) {
	case *SliceExpr:
		panic("can't inline slice expression")
	case *FunctionExpr:
		return InlineFunctionCall(e, def)
	}
}

func InlineFunctionCall(funExpr *FunctionExpr, def *Function) ExpressionInterface {
	if !funExpr.FuncDef.ZeroInternalSites() || funExpr.FuncDef == def {
		// inline only if there's no internal sites
		// don't do recursive inlining
		return funExpr
	}
	return funExpr.FuncDef.RetExpr.InlineCopy(funExpr)
}
