package main

import (
	"bufio"
	"fmt"
	"github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/abra"
	"github.com/lunfardo314/goq/abra/construct"
	"github.com/lunfardo314/goq/abra/generate"
	"github.com/lunfardo314/goq/abra/validate"
	"io/ioutil"
	"os"
)

const (
	moduleName  = "GeneratedTest"
	siteDataDir = "C:/Users/evaldas/Documents/proj/site_data/tritcode/"
	assumeSizes = true
)

func main() {
	fmt.Println("Generating test code unit...")
	codeUnit := construct.NewCodeUnit()
	tests := constructObjectsAndTests(codeUnit)

	fmt.Println("Calculating sizes...")
	validate.CalcAllSizes(codeUnit)

	fmt.Println("Sorting and enumerating blocks...")
	validate.SortAndEnumerateBlocks(codeUnit)

	fmt.Println("Sorting and enumerating sites...")
	validate.SortAndEnumerateSites(codeUnit)

	fmt.Printf("Validating code unit (assume sizes = %v)...\n", assumeSizes)

	errs := validate.Validate(codeUnit, assumeSizes)
	if len(errs) != 0 {
		fmt.Println("Validation errors")
		for _, err := range errs {
			fmt.Printf("%v\n", err)
		}
		fmt.Println("Can't continue. Exit..")
		os.Exit(0)
	} else {
		fmt.Println("Validation OK")
	}
	fmt.Println("Generating tritcode")
	tritcode := generate.NewTritcode()
	tritcode = generate.WriteCode(tritcode, codeUnit.Code)
	fmt.Printf("Number of trits generated: %d\n", len(tritcode))

	fmt.Println("Converting to trytes")
	trytecode := trinary.MustTritsToTrytes(tritcode)
	fmt.Printf("Number of trytes generated: %d\n", len(trytecode))

	fname := siteDataDir + moduleName + ".abra.trytes"
	fmt.Printf("Saving trytes to %s\n", fname)
	err := ioutil.WriteFile(fname, []byte(trytecode), 0644)
	if err != nil {
		panic(err)
	}

	fname = siteDataDir + moduleName + ".abra.txt"
	fmt.Printf("Writing Abra code uni in human readable form to %s\n", fname)
	err = generate.SaveReadable(codeUnit, fname)
	if err != nil {
		panic(err)
	}

	fname = siteDataDir + moduleName + ".abra.test"

	assignBlockIndicesToTests(codeUnit, tests)

	fmt.Printf("Writing %d Abra tests to %s\n", len(tests), fname)
	err = writeTestsToFile(tests, fname)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Ciao..")
}

func assignBlockIndicesToTests(codeUnit *abra.CodeUnit, tests []*generate.AbraTest) {
	for _, t := range tests {
		block := construct.FindBlock(codeUnit, t.Comment)
		if block == nil {
			panic("can't find " + t.Comment)
		}
		t.BlockIndex = block.Index
	}
}

func writeTestsToFile(tests []*generate.AbraTest, fname string) error {
	fout, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer fout.Close()
	w := bufio.NewWriter(fout)
	for i, t := range tests {
		err = generate.WriteAbraTest(w, t, i)
		if err != nil {
			return err
		}
	}
	w.Flush()
	return nil

}
