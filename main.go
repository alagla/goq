package main

import (
	"github.com/lunfardo314/goq/qupla"
	"github.com/lunfardo314/goq/quplayaml"
	"os"
)

const fname = "C:/Users/evaldas/Documents/proj/Java/github.com/qupla/src/main/resources/Qupla.yml"
const testout = "C:/Users/evaldas/Documents/proj/site_data/tmp/echotest.yml"

const trace = true

func main() {
	qupla.SetLog(nil, trace)

	moduleYAML, err := quplayaml.NewQuplaModuleFromYAML(fname)
	if err != nil {
		errorf("err: %v", err)
		os.Exit(1)
	}
	// echo for testing
	_ = moduleYAML.WriteToFile(testout)

	module, succ := qupla.AnalyzeQuplaModule(moduleYAML, &qupla.ExpressionFactoryFromYAML{})
	module.PrintStats()

	if !succ {
		errorf("Failed analyzing Qupla module")
	} else {
		module.Execute(true)
	}
	infof("Ciao!")
}
