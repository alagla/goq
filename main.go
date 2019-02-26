package main

import (
	"github.com/lunfardo314/goq/quplayaml"
	"os"
)

const fname = "C:/Users/evaldas/Documents/proj/Java/github.com/qupla/src/main/resources/Qupla.yml"
const testout = "C:/Users/evaldas/Documents/proj/site_data/tmp/echotest.yml"

func main() {
	quplayaml.SetLog(nil, true)

	moduleYAML, err := quplayaml.NewQuplaModuleFromYAML(fname)
	if err != nil {
		errorf("err: %v", err)
		os.Exit(1)
	}
	// echo for testing
	_ = moduleYAML.WriteToFile(testout)

	module, succ := quplayaml.AnalyzeQuplaModule(moduleYAML, &quplayaml.ExpressionFactoryFromYAML{})
	module.PrintStats()

	if !succ {
		errorf("Failed analyzing Qupla module")
	} else {
		module.Execute()
	}
	infof("Ciao!")
}
