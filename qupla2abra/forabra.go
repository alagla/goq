package qupla2abra

import (
	. "github.com/lunfardo314/goq/qupla"
	"github.com/lunfardo314/goq/transform"
)

func optimizeFunction4Abra(def *Function, stats map[string]int) bool {
	var optSlices, optInlineSlices, optConcats, optMerges, optInlineCalls bool

	optSlices = transform.OptimizeSlices(def, stats)
	optInlineSlices = transform.OptimizeInlineSlices(def, stats)
	optConcats = transform.OptimizeConcats(def, stats)
	optMerges = transform.OptimizeMerges(def, stats)
	return optSlices || optInlineSlices || optConcats || optMerges || optInlineCalls
}
