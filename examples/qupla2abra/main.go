package main

import (
	"github.com/lunfardo314/goq/abra"
	"github.com/lunfardo314/goq/analyzeyaml"
	. "github.com/lunfardo314/goq/cfg"
	"github.com/lunfardo314/goq/readyaml"
)

const yamlToLoad = "C:/Users/evaldas/Documents/proj/Go/src/github.com/lunfardo314/goq/examples/modules/QuplaTests.yml"

func main() {
	Logf(0, "Loading Qupla module from %v", yamlToLoad)
	moduleYAML, err := readyaml.NewQuplaModuleFromYAML(yamlToLoad)
	if err != nil {
		Logf(0, "Error while parsing YAML file: %v", err)
		moduleYAML = nil
		return
	}
	// analyze loaded module and produce interpretable IR
	module, succ := analyzeyaml.AnalyzeQuplaModule("Qupla Module", moduleYAML)
	if !succ {
		Logf(0, "Failed to lead module: %v", err)
		return
	}
	module.PrintStats()

	Logf(0, "------------------------")
	Logf(0, "generating Abra code")

	codeUnit := abra.NewCodeUnit()
	module.GetAbra(codeUnit)
}
