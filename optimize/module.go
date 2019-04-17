package optimize

import (
	. "github.com/lunfardo314/goq/cfg"
	. "github.com/lunfardo314/goq/qupla"
	"sort"
)

func OptimizeModule(module *QuplaModule, stats map[string]int) {
	Logf(1, "Optimizing module..")
	Logf(1, "Inline call optimisation = %v", Config.OptimizeFunCallsInline)
	Logf(1, "Inline slice optimisation = %v", Config.OptimizeInlineSlices)
	Logf(1, "One time site optimisation = %v", Config.OptimizeOneTimeSites)
	Logf(1, "Concat optimisation = %v", Config.OptimizeConcats)

	tmpKeys := make([]string, 0, len(module.Functions))

	for k := range module.Functions {
		tmpKeys = append(tmpKeys, k)
	}
	sort.Strings(tmpKeys)

	for _, funName := range tmpKeys {
		// optimize while there's something to optimize
		Logf(3, "Optimizing function '%v'", funName)
		fstats := make(map[string]int)
		for optimizeFunction(module.Functions[funName], fstats) {
		}
		if len(fstats) != 0 {
			LogStats(3, fstats)
		}
		AddStats(stats, fstats)
	}
}
