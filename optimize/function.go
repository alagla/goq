package optimize

import (
	. "github.com/lunfardo314/goq/cfg"
	. "github.com/lunfardo314/goq/qupla"
)

func optimizeFunction(def *Function) {
	if Config.OptimizeOneTimeSites {
		before := def.ZeroInternalSites()
		def.RetExpr = optimizeSlices(def, def.RetExpr)

		_, _, _, numVars, numUnusedVars := def.NumSites()
		Logf(5, "Optimized %v sites out of %v in '%v'", numUnusedVars, numVars, def.Name)
		after := def.ZeroInternalSites()
		if !before && after {
			Logf(5, "'%v' became inlineable", def.Name)
		}
	}
	if Config.OptimizeInlineSlices {
		def.RetExpr = optimizeInlineSlices(def.RetExpr)
	}
	if Config.OptimizeConcats {
		def.RetExpr = optimizeConcatExpr(def.RetExpr)
	}
}
