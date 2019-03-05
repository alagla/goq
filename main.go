package main

import (
	"fmt"
	"github.com/lunfardo314/goq/cfg"
	dispatcher2 "github.com/lunfardo314/goq/dispatcher"
	"github.com/lunfardo314/goq/qupla"
	"github.com/lunfardo314/goq/quplayaml"
	"os"
)

const fname = "C:/Users/evaldas/Documents/proj/Java/github.com/qupla/src/main/resources/Qupla.yml"
const testout = "C:/Users/evaldas/Documents/proj/site_data/tmp/echotest.yml"

func main() {
	logf(0, "GOQ: Qupla analyzer and interpreter in pure Go ver. '%v'", cfg.Config.Version)
	logf(0, "Verbosity level: %v", cfg.Config.Verbosity)
	logf(0, "Reading Qupla module form file %v", fname)

	qupla.SetLog(nil, cfg.Config.Trace)

	moduleYAML, err := quplayaml.NewQuplaModuleFromYAML(fname)
	if err != nil {
		logf(0, "Error while parsing YAML file: %v", err)
		os.Exit(1)
	}
	// echo for testing
	logf(0, "Echo Qupla module to file %v", testout)
	_ = moduleYAML.WriteToFile(testout)

	module, succ := qupla.AnalyzeQuplaModule("single_module", moduleYAML, &qupla.ExpressionFactoryFromYAML{})
	module.PrintStats()

	dispatcher := dispatcher2.StartDispatcher()
	module.AttachToDispatcher(dispatcher)

	if !succ {
		logf(0, "Failed analyzing Qupla module")
	} else {
		module.Execute()
	}
	logf(0, "Ciao! I'll be back")
}

func logf(minVerbosity int, format string, args ...interface{}) {
	if cfg.Config.Verbosity < minVerbosity {
		return
	}
	fmt.Printf(format+"\n", args...)
}
