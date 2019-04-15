package optimize

import (
	. "github.com/lunfardo314/goq/cfg"
	. "github.com/lunfardo314/goq/qupla"
	"sort"
)

func OptimizeModule(module *QuplaModule) {
	Logf(0, "Optimizing module..")
	Logf(1, "Inline call optimisation = %v", Config.OptimizeFunCallsInline)
	Logf(1, "Inline slice optimisation = %v", Config.OptimizeInlineSlices)
	Logf(1, "One time site optimisation = %v", Config.OptimizeOneTimeSites)
	Logf(1, "Concat optimisation = %v", Config.OptimizeConcats)

	stats := make(map[string]int)
	tmpKeys := make([]string, 0, len(module.Functions))

	for k := range module.Functions {
		tmpKeys = append(tmpKeys, k)
	}

	sort.Strings(tmpKeys)

	for _, funName := range tmpKeys {
		// optimize while there's something to optimize
		Logf(2, "Optimizing function '%v'", funName)
		for optimizeFunction(module.Functions[funName], stats) {
		}
	}
	LogStats(0, stats)

	//for _, exec := range module.Execs {
	//	OptimizeExecStmt(exec)
	//}
}
