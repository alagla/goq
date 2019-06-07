package cfg

type ConfigStruct struct {
	Version                string
	Verbosity              int
	OptimizeFunCallsInline bool
	OptimizeOneTimeSites   bool
	OptimizeInlineSlices   bool
	OptimizeConcats        bool
	OptimizeMerges         bool
}

var Config = &ConfigStruct{
	Version:                "0.01 alpha",
	Verbosity:              2,
	OptimizeFunCallsInline: false,
	OptimizeOneTimeSites:   false,
	OptimizeInlineSlices:   false,
	OptimizeConcats:        false,
	OptimizeMerges:         false,
}
