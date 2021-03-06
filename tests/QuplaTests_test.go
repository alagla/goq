package tests

import (
	"github.com/lunfardo314/goq/analyzeyaml"
	"github.com/lunfardo314/goq/cfg"
	"github.com/lunfardo314/goq/readyaml"
	"github.com/lunfardo314/goq/transform"
	"testing"
)

func Test_QuplaTests(t *testing.T) {
	moduleTest("../examples/modules/QuplaTests.yml", false, t)
}

func Test_QuplaTestsChain(t *testing.T) {
	moduleTest("../examples/modules/QuplaTests.yml", true, t)
}

func Test_Fibonacci(t *testing.T) {
	moduleTest("../examples/modules/Fibonacci.yml", false, t)
}

func Test_FibonacciChain(t *testing.T) {
	moduleTest("../examples/modules/Fibonacci.yml", true, t)
}

func Test_Examples(t *testing.T) {
	moduleTest("../examples/modules/Examples.yml", false, t)
}

func Test_ExamplesChain(t *testing.T) {
	moduleTest("../examples/modules/Examples.yml", true, t)
}

func Test_Curl(t *testing.T) {
	moduleTest("../examples/modules/Curl.yml", false, t)
}

func Test_CurlChain(t *testing.T) {
	moduleTest("../examples/modules/Curl.yml", true, t)
}

func moduleTest(fname string, chain bool, t *testing.T) {

	cfg.Logf(0, "---------------------------\nTesting QuplaYAML module %v. Chain mode = %v", fname, chain)
	if !check0environments(t) {
		return
	}

	moduleYAML, err := readyaml.NewQuplaModuleFromYAML(fname)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	module, succ := analyzeyaml.AnalyzeQuplaModule(fname, moduleYAML)

	if succ {
		stats := make(map[string]int)

		transform.OptimizeModule(module, stats)
		cfg.LogStats(1, stats)

		succ = module.AttachToSupervisor(sv)
	} else {
		t.Errorf("Failed to load module from '%v'", fname)
		return
	}
	if err := module.RunExecs(sv, -1, -1, chain); err != nil {
		t.Errorf("Failed to run '%v': %v", fname, err)
	}
	if err := sv.ClearEnvironments(); err != nil {
		t.Errorf("ClearEnvironments: %v", err)
	}
}
