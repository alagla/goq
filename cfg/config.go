package cfg

type ConfigStruct struct {
	Version              string
	Trace                bool
	ExecuteStatefulExecs bool
	Verbosity            int
	ExecTests            bool
	ExecEvals            bool
}

var Config = &ConfigStruct{
	Version:              "0.01 alpha",
	Trace:                false,
	ExecuteStatefulExecs: false,
	Verbosity:            1,
	ExecTests:            true,
	ExecEvals:            true,
}
