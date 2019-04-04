package cfg

type ConfigStruct struct {
	Version                string
	Verbosity              int
	OptimizeFunCallsInline bool
	OptimizeOneTimeSites   bool
	OptimizeInlineSlices   bool
	OptimizeConcats        bool
}

var Config = &ConfigStruct{
	Version:                "0.01 alpha",
	Verbosity:              2,
	OptimizeFunCallsInline: true,
	OptimizeOneTimeSites:   true,
	OptimizeInlineSlices:   true,
	OptimizeConcats:        true,
}
