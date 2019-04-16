package tests

import (
	"github.com/lunfardo314/goq/analyzeyaml"
	"github.com/lunfardo314/goq/cfg"
	"github.com/lunfardo314/goq/optimize"
	"github.com/lunfardo314/goq/qupla"
	"github.com/lunfardo314/goq/readyaml"
	"sort"
	"testing"
)

func Test_Opt0QuplaTests(t *testing.T) {
	moduleOptimize_0_Test("../examples/modules/QuplaTests.yml", t)
}

func Test_Opt1QuplaTests(t *testing.T) {
	moduleOptimize_1_Test("../examples/modules/QuplaTests.yml", t)
	moduleOptimize_1_Test("../examples/modules/Examples.yml", t)
	moduleOptimize_1_Test("../examples/modules/Curl.yml", t)
	moduleOptimize_1_Test("../examples/modules/Fibonacci.yml", t)
}

func moduleOptimize_0_Test(fname string, t *testing.T) {

	cfg.Logf(0, "---------------------------\nTesting QuplaYAML module %v optimization", fname)
	moduleYAML, err := readyaml.NewQuplaModuleFromYAML(fname)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	var module *qupla.QuplaModule
	var succ bool

	name := fname + "-opt1"
	if module, succ = analyzeyaml.AnalyzeQuplaModule(name, moduleYAML); !succ {
		t.Errorf("Failed to analyze module '%v'", name)
		return
	}
	stats := make(map[string]int)
	optimize.OptimizeModule(module, stats)
	cfg.Logf(0, "Optimization stats for %v (1st optimization)", name)
	cfg.LogStats(0, stats)

	stats = make(map[string]int)
	optimize.OptimizeModule(module, stats)
	if len(stats) != 0 {
		t.Errorf("Some part left unoptimized in '%v'", name)
		return
	}
}

func moduleOptimize_1_Test(fname string, t *testing.T) {

	cfg.Logf(0, "------------------\nTesting QuplaYAML module %v optimization: testing if two optimisation has same result", fname)
	moduleYAML, err := readyaml.NewQuplaModuleFromYAML(fname)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	var module1, module2 *qupla.QuplaModule
	var succ bool

	name := fname + "-opt1"
	if module1, succ = analyzeyaml.AnalyzeQuplaModule(name, moduleYAML); !succ {
		t.Errorf("Failed to analyze module '%v'", name)
		return
	}
	stats1 := make(map[string]int)
	optimize.OptimizeModule(module1, stats1)
	cfg.Logf(1, "Optimization stats for %v (1st optimization)", name)
	cfg.LogStats(1, stats1)

	name = fname + "-opt2"
	if module2, succ = analyzeyaml.AnalyzeQuplaModule(name, moduleYAML); !succ {
		t.Errorf("Failed to analyze module '%v'", name)
		return
	}
	stats2 := make(map[string]int)
	optimize.OptimizeModule(module2, stats2)
	cfg.Logf(1, "Optimization stats for %v (2nd optimization)", name)
	cfg.LogStats(1, stats2)

	funNames := make([]string, 0)
	for n := range module1.Functions {
		funNames = append(funNames, n)
	}
	sort.Strings(funNames)

	for _, funName := range funNames {
		compareStats(module1, module2, funName, t)
	}
}

func compareStats(module1, module2 *qupla.QuplaModule, funname string, t *testing.T) {
	fun1 := module1.Functions[funname]
	fun2 := module2.Functions[funname]
	stats1 := fun1.Stats()
	stats2 := fun2.Stats()
	for k1, v1 := range stats1 {
		v2, ok := stats2[k1]
		if !ok || v1 != v2 {
			//if strings.HasPrefix(funname, "fixSign"){
			//	fmt.Println("kuku")
			//}
			t.Errorf("Two optimizations of function '%v' differ in '%v': %v != %v",
				funname, k1, v1, v2)
		}
	}
}
