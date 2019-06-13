package qupla2abra

import (
	. "github.com/lunfardo314/goq/cfg"
	. "github.com/lunfardo314/goq/qupla"
	"github.com/lunfardo314/goq/transform"
	"github.com/lunfardo314/goq/utils"
)

func PrepareModuleForAbra(module *QuplaModule, stats map[string]int) {
	Logf(3, "Preparing module for Abra..")

	// generate 1 output luts
	// for each >1 size output lut are created several new luts for each output trit position
	// with the name lutDef.MustGetProjectionName(outpos)

	utils.SetStat("numLuts", stats, len(module.Luts))
	utils.SetStat("numLutsSize>1", stats, 0)
	utils.SetStat("numLutProjectionsGenerated", stats, 0)

	for _, lutDef := range module.Luts {
		if lutDef.Size() > 1 {
			utils.IncStat("numLutsSize>1", stats)
			for outpos := 0; outpos < lutDef.Size(); outpos++ {
				newLutDef := lutDef.MakeAdjustedProjection(outpos)
				module.AddLutDef(newLutDef.Name, newLutDef)
				utils.IncStat("numLutProjectionsGenerated", stats)
			}
		}
	}

	// optimize slices, concats and merges. That will reduce number of branches
	for _, fun := range module.Functions {
		optimizeFunction4Abra(fun, stats)
	}

}

func optimizeFunction4Abra(def *Function, stats map[string]int) bool {
	var optSlices, optInlineSlices, optConcats, optMerges, optInlineCalls bool

	optSlices = transform.OptimizeSlices(def, stats)
	optInlineSlices = transform.OptimizeInlineSlices(def, stats)
	optConcats = transform.OptimizeConcats(def, stats)
	optMerges = transform.OptimizeMerges(def, stats)
	return optSlices || optInlineSlices || optConcats || optMerges || optInlineCalls
}
