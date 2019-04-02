package cfg

type ConfigStruct struct {
	Version        string
	Verbosity      int
	OptimizeInline bool
}

var Config = &ConfigStruct{
	Version:        "0.01 alpha",
	Verbosity:      2,
	OptimizeInline: true,
}
