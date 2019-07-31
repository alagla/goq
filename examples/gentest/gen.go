package main

import (
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/abra"
	"github.com/lunfardo314/goq/abra/construct"
	"github.com/lunfardo314/goq/abra/generate"
	"github.com/lunfardo314/goq/utils"
)

func constructObjectsAndTests(codeUnit *abra.CodeUnit) []*generate.AbraTest {
	ret := make([]*generate.AbraTest, 0, 100)

	ret = createNullifyBlocks(codeUnit, ret)
	ret = createConcatBlocks(codeUnit, ret)
	ret = createSliceBlocks(codeUnit, ret)
	ret = createConstBlocks(codeUnit, ret)

	ret = createLUTTests(codeUnit, ret)
	return ret
}

func createNullifyBlocks(codeUnit *abra.CodeUnit, ret []*generate.AbraTest) []*generate.AbraTest {
	var b *abra.Block
	// nullify true
	b = construct.GetNullifyBranchBlock(codeUnit, 1, true)
	ret = addTest(ret, Trits{-1, -1}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{-1, 0}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{-1, 1}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{0, -1}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{0, 0}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{0, 1}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{1, -1}, Trits{-1}, false, b.LookupName)
	ret = addTest(ret, Trits{1, 0}, Trits{0}, false, b.LookupName)
	ret = addTest(ret, Trits{1, 1}, Trits{1}, false, b.LookupName)

	// nullify false
	b = construct.GetNullifyBranchBlock(codeUnit, 1, false)
	ret = addTest(ret, Trits{-1, -1}, Trits{-1}, false, b.LookupName)
	ret = addTest(ret, Trits{-1, 0}, Trits{0}, false, b.LookupName)
	ret = addTest(ret, Trits{-1, 1}, Trits{1}, false, b.LookupName)
	ret = addTest(ret, Trits{0, -1}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{0, 0}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{0, 1}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{1, -1}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{1, 0}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{1, 1}, nil, false, b.LookupName)

	// nullify_true 3
	b = construct.GetNullifyBranchBlock(codeUnit, 3, true)
	ret = addTest(ret, Trits{-1, -1, 0, 0}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{-1, 0, 1}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{-1, 1, 1, 1}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{0, -1, 1, 0}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{0, 0, -1, 1}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{0, 1, 1, -1}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{1, -1, 0, 1}, Trits{-1, 0, 1}, false, b.LookupName)
	ret = addTest(ret, Trits{1, 1, 1, 1}, Trits{1, 1, 1}, false, b.LookupName)
	ret = addTest(ret, Trits{1, 0, -1, -1}, Trits{0, -1, -1}, false, b.LookupName)

	// nullify_false 3
	b = construct.GetNullifyBranchBlock(codeUnit, 3, false)
	ret = addTest(ret, Trits{-1, -1, 1, 0}, Trits{-1, 1, 0}, false, b.LookupName)
	ret = addTest(ret, Trits{-1, 0, 1, 1}, Trits{0, 1, 1}, false, b.LookupName)
	ret = addTest(ret, Trits{-1, 1, 0, -1}, Trits{1, 0, -1}, false, b.LookupName)
	ret = addTest(ret, Trits{0, -1, 0, 1}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{0, 0, 1, 1}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{0, 1, -1, 0}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{1, -1, 0, -1}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{1, 0, 0, 0}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{1, 1, 1, 1}, nil, false, b.LookupName)

	// nullify_true 9
	b = construct.GetNullifyBranchBlock(codeUnit, 9, true)
	ret = addTest(ret, Trits{-1, -1, 0, 0, 1, 1, 1, 0, 0, 0}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{-1, 0, 1, 0, 1, 1, 1, 0, 0, 0}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{-1, 1, 1, 1, 0, 1, 1, 1, 0, 1}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{0, -1, 1, 0, 1, 1, 1, 1, 0, -1}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{0, 0, -1, 1, 0, 1, 1, 1, 0, 0}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{0, 1, 1, -1, 1, 1, 0, -1, 0, 0}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{1, -1, 0, 1, 1, 1, 1, 1, 1, 1}, Trits{-1, 0, 1, 1, 1, 1, 1, 1, 1}, false, b.LookupName)
	ret = addTest(ret, Trits{1, 1, 1, 1, 1, 0, 1, 1, -1, 1}, Trits{1, 1, 1, 1, 0, 1, 1, -1, 1}, false, b.LookupName)
	ret = addTest(ret, Trits{1, 0, 0, -1, 1, 1, -1, 0, 1, -1}, Trits{0, 0, -1, 1, 1, -1, 0, 1, -1}, false, b.LookupName)

	// nullify_false 9
	b = construct.GetNullifyBranchBlock(codeUnit, 9, false)
	ret = addTest(ret, Trits{-1, -1, 0, 1, 1, 1, 1, 1, 1, 1}, Trits{-1, 0, 1, 1, 1, 1, 1, 1, 1}, false, b.LookupName)
	ret = addTest(ret, Trits{-1, 1, 1, 1, 1, 0, 1, 1, -1, 1}, Trits{1, 1, 1, 1, 0, 1, 1, -1, 1}, false, b.LookupName)
	ret = addTest(ret, Trits{-1, 0, 0, -1, 1, 1, -1, 0, 1, -1}, Trits{0, 0, -1, 1, 1, -1, 0, 1, -1}, false, b.LookupName)
	ret = addTest(ret, Trits{0, -1, 1, 0, 1, 1, 1, 1, 0, -1}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{0, 0, -1, 1, 0, 1, 1, 1, 0, 0}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{0, 1, 1, -1, 1, 1, 0, -1, 0, 0}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{1, -1, 0, 0, 1, 1, 1, 0, 0, 0}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{1, 0, 1, 0, 1, 1, 1, 0, 0, 0}, nil, false, b.LookupName)
	ret = addTest(ret, Trits{1, 1, 1, 1, 0, 1, 1, 1, 0, 1}, nil, false, b.LookupName)

	return ret
}

func createConcatBlocks(codeUnit *abra.CodeUnit, ret []*generate.AbraTest) []*generate.AbraTest {
	var b *abra.Block
	b = construct.GetConcatBlockForSize(codeUnit, 1)
	ret = addTest(ret, Trits{-1}, Trits{-1}, false, b.LookupName)
	ret = addTest(ret, Trits{0}, Trits{0}, false, b.LookupName)
	ret = addTest(ret, Trits{1}, Trits{1}, false, b.LookupName)

	b = construct.GetConcatBlockForSize(codeUnit, 3)
	ret = addTest(ret, Trits{-1, 0, 0}, Trits{-1, 0, 0}, false, b.LookupName)
	ret = addTest(ret, Trits{0, -1, 1}, Trits{0, -1, 1}, false, b.LookupName)
	ret = addTest(ret, Trits{1, 1, 0}, Trits{1, 1, 0}, false, b.LookupName)

	b = construct.GetConcatBlockForSize(codeUnit, 9)
	ret = addTest(ret, Trits{-1, 0, 0, -1, 0, 0, 1, 0, 0}, Trits{-1, 0, 0, -1, 0, 0, 1, 0, 0}, false, b.LookupName)
	ret = addTest(ret, Trits{0, -1, 1, -1, 1, 0, -1, 0, 0}, Trits{0, -1, 1, -1, 1, 0, -1, 0, 0}, false, b.LookupName)
	ret = addTest(ret, Trits{1, 1, 0, -1, -1, 0, -1, 0, 1}, Trits{1, 1, 0, -1, -1, 0, -1, 0, 1}, false, b.LookupName)
	return ret
}

func createSliceBlocks(codeUnit *abra.CodeUnit, ret []*generate.AbraTest) []*generate.AbraTest {
	var b *abra.Block
	b = construct.GetSliceBranchBlock(codeUnit, 1, 0, 1)
	ret = addTest(ret, Trits{-1}, Trits{-1}, false, b.LookupName)
	ret = addTest(ret, Trits{0}, Trits{0}, false, b.LookupName)
	ret = addTest(ret, Trits{1}, Trits{1}, false, b.LookupName)

	b = construct.GetSliceBranchBlock(codeUnit, 3, 0, 1)
	ret = addTest(ret, Trits{-1, 0, 0}, Trits{-1}, false, b.LookupName)
	ret = addTest(ret, Trits{0, 1, 1}, Trits{0}, false, b.LookupName)
	ret = addTest(ret, Trits{1, -1, 0}, Trits{1}, false, b.LookupName)

	b = construct.GetSliceBranchBlock(codeUnit, 3, 0, 2)
	ret = addTest(ret, Trits{-1, 0, 0}, Trits{-1, 0}, false, b.LookupName)
	ret = addTest(ret, Trits{0, 1, 1}, Trits{0, 1}, false, b.LookupName)
	ret = addTest(ret, Trits{1, -1, 0}, Trits{1, -1}, false, b.LookupName)

	b = construct.GetSliceBranchBlock(codeUnit, 3, 1, 2)
	ret = addTest(ret, Trits{-1, 0, 0}, Trits{0, 0}, false, b.LookupName)
	ret = addTest(ret, Trits{0, 1, 1}, Trits{1, 1}, false, b.LookupName)
	ret = addTest(ret, Trits{1, -1, 0}, Trits{-1, 0}, false, b.LookupName)

	b = construct.GetSliceBranchBlock(codeUnit, 3, 0, 3)
	ret = addTest(ret, Trits{-1, 0, 0}, Trits{-1, 0, 0}, false, b.LookupName)
	ret = addTest(ret, Trits{0, 1, 1}, Trits{0, 1, 1}, false, b.LookupName)
	ret = addTest(ret, Trits{1, -1, 0}, Trits{1, -1, 0}, false, b.LookupName)

	return ret
}

func createConstBlocks(codeUnit *abra.CodeUnit, ret []*generate.AbraTest) []*generate.AbraTest {
	var b *abra.Block

	vals := [][]int8{
		{1, 1, 1}, {1, 1, 1}, {1}, {0}, {-1}, {1, 0, 0, -1, 0, -1}, {1, 0, -1, -1, 0, -1, 1, 1}, {1, 0, 1, -1, 0, 0, 1},
	}
	for _, v := range vals {
		b = construct.GetConstTritVectorBlock(codeUnit, v)
		ret = addTest(ret, Trits{-1}, v, false, b.LookupName)
		ret = addTest(ret, Trits{0}, v, false, b.LookupName)
		ret = addTest(ret, Trits{1}, v, false, b.LookupName)
	}
	return ret
}

func addTest(ret []*generate.AbraTest, input, expected Trits, isFloat bool, blockName string) []*generate.AbraTest {
	return append(ret, newTest(input, expected, isFloat, blockName))
}

func newTest(input, expected Trits, isFloat bool, blockName string) *generate.AbraTest {
	exp := ""
	if expected == nil {
		exp = "@"
	} else {
		exp = utils.TritsToString(expected)
	}
	return &generate.AbraTest{
		BlockIndex: -1, // will be assigned later
		Input:      utils.TritsToString(input),
		Expected:   exp,
		IsFloat:    isFloat,
		Comment:    blockName,
	}
}

func createLUTTests(codeUnit *abra.CodeUnit, ret []*generate.AbraTest) []*generate.AbraTest {

	// gen tests for luts
	for _, b := range codeUnit.Code.Blocks {
		if b.BlockType != abra.BLOCK_LUT {
			continue
		}
		strRepr := abra.StringFromBinaryEncodedLUT(b.LUT.Binary)
		for i, tripl := range utils.GetTriplets() {
			var exp Trits
			c := ([]byte(strRepr))[i]
			if c == '@' {
				exp = nil
			} else {
				exp = []int8{charToTrit(c)}
			}
			ret = addTest(ret, tripl, exp, false, b.LookupName)
		}
	}
	return ret
}

func charToTrit(c byte) int8 {
	switch c {
	case '-':
		return -1
	case '0':
		return 0
	case '1':
		return 1
	}
	panic("wrong trit")
}
