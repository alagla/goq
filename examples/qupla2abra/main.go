package main

import (
	"github.com/lunfardo314/goq/abra"
	cabra "github.com/lunfardo314/goq/abra/construct"
	gabra "github.com/lunfardo314/goq/abra/generate"
	vabra "github.com/lunfardo314/goq/abra/validate"
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

	codeUnit := cabra.NewCodeUnit()
	module.GetAbra(codeUnit)

	Logf(0, "------ validating entire code unit")
	vabra.CalcAllSizes(codeUnit)
	errs := vabra.Validate(codeUnit)
	if len(errs) == 0 {
		Logf(0, "code unit validate OK")
		printSizes(codeUnit)
	} else {
		Logf(0, "Validation errors in code unit")
		for _, err := range errs {
			Logf(0, "    ->  %v", err)
		}
	}
	Logf(0, "Generating Abra tritcode")

	tritcode := gabra.NewTritcode()
	numtrits := tritcode.MustWriteCode(codeUnit.Code)
	Logf(0, "Number of trits generated: %d", numtrits)

	trytecode := abra.MustBTrits2Bytes(tritcode.Buf.Bytes())
	Logf(0, "Number of trytes generated: %d", len(trytecode))

	Logf(0, "First 200 of trytes: %s...", string(trytecode[:200]))
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
		Logf(3, "%20s -> size = %d (%d)", n, blockMap[n].size, blockMap[n].assumedSize)
	}

}
