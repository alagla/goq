package main

import (
	"github.com/lunfardo314/goq/abra"
	"github.com/lunfardo314/goq/analyzeyaml"
	. "github.com/lunfardo314/goq/cfg"
	"github.com/lunfardo314/goq/readyaml"
	"sort"
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

	Logf(0, "------ checking sizes")
	codeUnit.CalcSizes()
	printSizes(codeUnit)
}

type sizeInfo struct{ size, assumedSize int }

func printSizes(codeUnit *abra.CodeUnit) {
	blockMap := make(map[string]*sizeInfo)
	names := make([]string, 0, len(codeUnit.Code.Blocks))
	for _, b := range codeUnit.Code.Blocks {
		names = append(names, b.LookupName)
		blockMap[b.LookupName] = &sizeInfo{size: b.Size, assumedSize: b.AssumedSize}
	}
	sort.Strings(names)
	for _, n := range names {
		Logf(0, "%20s -> size = %d (%d)", n, blockMap[n].size, blockMap[n].assumedSize)
	}

}
