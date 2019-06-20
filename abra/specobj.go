package abra

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
)

// finds or creates site which returns constant 'val'

func (branch *Branch) Get1TritConstLutSite(codeUnit *CodeUnit, val int8) *Site {
	// first try to find if there's constant lut site for val
	// if not, create one
	lookupName := fmt.Sprintf("1trit_const_site_%s", TritName(val))
	ret := branch.FindBodySite(lookupName)
	if ret != nil {
		return ret
	}
	// didn't find. Need to create one
	// first find or create the only lut for 1 trit constant
	lutRepr := Get1TritConstLutRepr(val)
	lutValConstBlock := codeUnit.FindLUTBlock(lutRepr)
	if lutValConstBlock == nil {
		lut := BinaryEncodedLUTFromString(lutRepr)
		lutValConstBlock = codeUnit.NewLUTBlock(lutRepr, lut)
	}
	// now create site in the branch
	// it will always generate constant trit
	// the input for lut is 3 repeated 1-trit sites from lsb of the branches input
	any := branch.GetAnyTritInputSite(codeUnit)
	ret = branch.AddKnotSiteForInputs(lutValConstBlock, lookupName, any, any, any)
	return ret
}

func (branch *Branch) GetAnyTritInputSite(codeUnit *CodeUnit) *Site {
	lookupName := "any_input_site" // each branch will have the only site with this name
	ret := branch.FindBodySite(lookupName)
	if ret != nil {
		return ret
	}
	// must be the only LstSliceBlock in the code unit
	lstBlock := codeUnit.GetLstSliceBlock(branch.InputSites[0].Size) // get Lst branch for the size of first input
	ret = NewKnot(lstBlock, branch.InputSites[0]).NewSite(lookupName + "_knot")
	branch.AddBodySite(ret)
	return ret
}

func (branch *Branch) GetTritConstSite(codeUnit *CodeUnit, val Trits) *Site {
	lookupName := TritsToString(val) + "_const_site"
	ret := branch.FindBodySite(lookupName)
	if ret != nil {
		return ret
	}
	inputs := make([]*Site, len(val))
	for i, trit := range val {
		inputs[i] = branch.Get1TritConstLutSite(codeUnit, trit)
	}

	concatBlock := codeUnit.GetConcatBlockForSize(len(val))
	ret = branch.AddKnotSiteForInputs(concatBlock, lookupName, inputs...)
	return ret
}

func (codeUnit *CodeUnit) GetConcatBlockForSize(size int) *Block {
	lookupName := fmt.Sprintf("concat_branch_%d", size)
	ret := codeUnit.FindBranchBlock(lookupName)
	if ret != nil {
		return ret
	}
	ret = codeUnit.AddNewBranchBlock(lookupName, size)
	input := ret.Branch.AddInputSite(size)
	ret.Branch.AddOutputSite(NewMerge(input).NewSite(lookupName + "_out"))
	return ret
}

// returns or creates block which takes to output least significant trit of it input
// ue to requirement to have exact size matches, there's one block per each

func (codeUnit *CodeUnit) GetLstSliceBlock(inputSize int) *Block {
	lookupName := fmt.Sprintf("LST_SLICE_BRANCH_BLOCK_%d", inputSize)
	ret := codeUnit.FindBranchBlock(lookupName)
	if ret != nil {
		return ret
	}
	ret = codeUnit.AddNewBranchBlock(lookupName, 1)
	inputLst := ret.Branch.AddInputSite(1) // two inputs
	ret.Branch.AddInputSite(inputSize - 1) //
	output := NewMerge(inputLst).NewSite("")
	ret.Branch.AddOutputSite(output)
	return ret
}

func (codeUnit *CodeUnit) GetSlicingBranch(offset, size int) *Block {
	return nil // TODO
}
