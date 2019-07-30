package main

import (
	"bufio"
	"fmt"
	"github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/abra"
	"github.com/lunfardo314/goq/abra/construct"
	"github.com/lunfardo314/goq/abra/generate"
	"github.com/lunfardo314/goq/abra/validate"
	"github.com/lunfardo314/goq/utils"
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
	constructObjects(codeUnit)

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
	tests := createTests(codeUnit)
	fmt.Printf("Generates %d tests\n", len(tests))
	fname = siteDataDir + moduleName + ".abra.test"

	fmt.Printf("Writing Abra tests to %s\n", fname)
	err = writeTestsToFile(tests, fname)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Ciao..")
}

func constructObjects(codeUnit *abra.CodeUnit) {
	construct.GetNullifyBranchBlock(codeUnit, 1, true)
	construct.GetNullifyBranchBlock(codeUnit, 3, true)
	construct.GetNullifyBranchBlock(codeUnit, 9, true)

	construct.GetNullifyBranchBlock(codeUnit, 1, false)
	construct.GetNullifyBranchBlock(codeUnit, 3, false)
	construct.GetNullifyBranchBlock(codeUnit, 9, false)

	construct.GetConcatBlockForSize(codeUnit, 1)
	construct.GetConcatBlockForSize(codeUnit, 3)
	construct.GetConcatBlockForSize(codeUnit, 9)

	construct.GetSliceBranchBlock(codeUnit, 1, 0, 1)
	construct.GetSliceBranchBlock(codeUnit, 3, 0, 1)
	construct.GetSliceBranchBlock(codeUnit, 3, 0, 2)
	construct.GetSliceBranchBlock(codeUnit, 3, 1, 2)
	construct.GetSliceBranchBlock(codeUnit, 3, 0, 3)

	construct.GetConstTritVectorBlock(codeUnit, trinary.Trits{0})
	construct.GetConstTritVectorBlock(codeUnit, trinary.Trits{1})
	construct.GetConstTritVectorBlock(codeUnit, trinary.Trits{-1})

	construct.GetConstTritVectorBlock(codeUnit, trinary.Trits{-1, -1, -1})
	construct.GetConstTritVectorBlock(codeUnit, trinary.Trits{-1, 0, 0})
	construct.GetConstTritVectorBlock(codeUnit, trinary.Trits{1, 0, -1})
}

func createTests(codeUnit *abra.CodeUnit) []*generate.AbraTest {
	ret := make([]*generate.AbraTest, 0, len(codeUnit.Code.Blocks))

	// gen tests for luts
	for _, b := range codeUnit.Code.Blocks {
		if b.BlockType != abra.BLOCK_LUT {
			continue
		}
		strRepr := abra.StringFromBinaryEncodedLUT(b.LUT.Binary)
		for i, tripl := range utils.GetTriplets() {
			t := generate.AbraTest{
				BlockIndex: b.Index,
				Input:      utils.TritsToString(tripl),
				Expected:   string([]byte{[]byte(strRepr)[i]}),
				IsFloat:    false,
				Comment:    b.LookupName,
			}
			ret = append(ret, &t)
		}
	}

	return ret
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
