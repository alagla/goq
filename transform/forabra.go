package transform

import (
	. "github.com/lunfardo314/goq/cfg"
	. "github.com/lunfardo314/goq/qupla"
)

func PrepareModuleForAbra(module *QuplaModule, stats map[string]int) {
	Logf(3, "Preparing module for Abra..")

	// generate 1 output luts
	//for _, lutDef := range module.Luts{
	//
	//}
	// optimize slices, concats and merges. That will reduce number of branches
	for _, fun := range module.Functions {
		optimizeFunction4Abra(fun, stats)
	}

}

func optimizeFunction4Abra(def *Function, stats map[string]int) bool {
	var optSlices, optInlineSlices, optConcats, optMerges, optInlineCalls bool

	optSlices = optimizeSlices(def, stats)
	optInlineSlices = optimizeInlineSlices(def, stats)
	optConcats = optimizeConcats(def, stats)
	optMerges = optimizeMerges(def, stats)
	return optSlices || optInlineSlices || optConcats || optMerges || optInlineCalls
}
