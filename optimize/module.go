package optimize

import (
	. "github.com/lunfardo314/goq/cfg"
	. "github.com/lunfardo314/goq/qupla"
)

func OptimizeModule(module *QuplaModule) {
	Logf(0, "Optimizing module..")
	Logf(1, "Inline call optimisation = %v", Config.OptimizeFunCallsInline)
	Logf(1, "Inline slice optimisation = %v", Config.OptimizeInlineSlices)
	Logf(1, "One time site optimisation = %v", Config.OptimizeOneTimeSites)
	Logf(1, "Concat optimisation = %v", Config.OptimizeConcats)

	for _, fun := range module.Functions {
		optimizeFunction(fun)
	}
	for _, exec := range module.Execs {
		OptimizeExecStmt(exec)
	}
}
