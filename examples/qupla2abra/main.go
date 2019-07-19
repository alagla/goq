package main

import (
	"github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/abra"
	cabra "github.com/lunfardo314/goq/abra/construct"
	gabra "github.com/lunfardo314/goq/abra/generate"
	rabra "github.com/lunfardo314/goq/abra/read"
	vabra "github.com/lunfardo314/goq/abra/validate"
	"github.com/lunfardo314/goq/analyzeyaml"
	. "github.com/lunfardo314/goq/cfg"
	"github.com/lunfardo314/goq/readyaml"
	"io/ioutil"
	"sort"
)

const (
	moduleName  = "QuplaTests"
	yamlPath    = "C:/Users/evaldas/Documents/proj/Go/src/github.com/lunfardo314/goq/examples/modules/"
	siteDataDir = "C:/Users/evaldas/Documents/proj/site_data/tritcode/"
)

func main() {
	Logf(0, "Loading Qupla module from %v")
	moduleYAML, err := readyaml.NewQuplaModuleFromYAML(yamlPath + moduleName + ".yml")
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
	tritcode = gabra.WriteCode(tritcode, codeUnit.Code)
	Logf(0, "Number of trits generated: %d", len(tritcode))

	trytecode := trinary.MustTritsToTrytes(tritcode)
	Logf(0, "Number of trytes generated: %d", len(trytecode))
	Logf(0, "First 200 of trytes: %s...", string(trytecode[:200]))

	packedtrits := trinary.TritsToBytes(tritcode)
	Logf(0, "Number of bytes in packed trits (5 trits per byte) generated: %d", len(packedtrits))

	fname := siteDataDir + moduleName + ".abra.trytes"
	Logf(0, "writing trytes to %s", fname)
	err = ioutil.WriteFile(fname, []byte(trytecode), 0644)
	if err != nil {
		panic(err)
	}

	fname = siteDataDir + moduleName + ".abra.packed"
	Logf(0, "writing packed trits to %s", fname)
	err = ioutil.WriteFile(fname, packedtrits, 0644)
	if err != nil {
		panic(err)
	}

	var tritsecho trinary.Trits
	tritsecho, err = abra.Trytes2Trits(trytecode)
	Logf(0, "reading back %d trits from trytes", len(tritsecho))

	Logf(0, "parsing back generated tritcode")
	var codeEcho *abra.CodeUnit
	codeEcho, err = rabra.ParseTritcode(tritsecho)
	if err != nil {
		panic(err)
	}
	Logf(0, "tritcode version is '%d'", codeEcho.Code.TritcodeVersion)
	Logf(0, "number of LUT blocks is '%d'", codeEcho.Code.NumLUTs)
	Logf(0, "number of branch blocks is '%d'", codeEcho.Code.NumBranches)
	Logf(0, "number of external blocks is '%d'", codeEcho.Code.NumExternalBlocks)

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
