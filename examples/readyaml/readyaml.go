package main

import (
	"fmt"
	. "github.com/lunfardo314/goq/readyaml"
	"os"
	"time"
)

const fname = "./examples/readyaml/QuplaTests.yml"

func main() {
	fmt.Printf("Reading from file '%v'\n", fname)

	start := time.Now()
	module, err := NewQuplaModuleFromYAML(fname)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Module '%v' has been loaded succesfully in %v\n", module.Name, time.Since(start))
	fmt.Printf("    Global types: %v\n", len(module.Types))
	fmt.Printf("    LUTs: %v\n", len(module.Luts))
	fmt.Printf("    Functions: %v\n", len(module.Functions))
	fmt.Printf("    Executable statements: %v\n", len(module.Execs))
}
