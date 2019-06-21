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
	ret := branch.FindSite(lookupName)
	if ret != nil {
		return ret
	}
	// didn't find. Need to create one
	// first find or create the only lut for 1 trit constant
	lutRepr := Get1TritConstLutRepr(val)
	lutValConstBlock := codeUnit.GetLUTBlock(lutRepr)

	// now create site in the branch
	// it will always generate constant trit
	// the input for lut is 3 repeated 1-trit sites from lsb of the branches input
	any := branch.GetAnyTritInputSite(codeUnit)
	ret = NewKnot(lutValConstBlock, any, any, any).NewSite()
	branch.AddNewSite(ret, lookupName)
	return ret
}

func (branch *Branch) GetAnyTritInputSite(codeUnit *CodeUnit) *Site {
	lookupName := "any_input_site" // each branch will have the only site with this name
	ret := branch.FindSite(lookupName)
	if ret != nil {
		return ret
	}
	// must be the only LstSliceBlock in the code unit
	lstBlock := codeUnit.GetSlicingBranchBlock(branch.InputSites[0].Size, 0, 1)
	ret = NewKnot(lstBlock, branch.InputSites[0]).NewSite()
	branch.AddNewSite(ret, "")
	return ret
}

func (branch *Branch) GetTritConstSite(codeUnit *CodeUnit, val Trits) *Site {
	lookupName := TritsToString(val) + "_const_site"
	ret := branch.FindSite(lookupName)
	if ret != nil {
		return ret
	}
	inputs := make([]*Site, len(val))
	for i, trit := range val {
		inputs[i] = branch.Get1TritConstLutSite(codeUnit, trit)
	}

	concatBlock := codeUnit.GetConcatBlockForSize(len(val))
	ret = NewKnot(concatBlock, inputs...).NewSite()
	branch.AddNewSite(ret, lookupName)
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
	output := NewMerge(input).NewSite()
	output.SiteType = SITE_OUTPUT
	ret.Branch.AddNewSite(output, lookupName)
	return ret
}

// returns or creates block which takes to output least significant trit of it input
// ue to requirement to have exact size matches, there's one block per each

func (codeUnit *CodeUnit) GetSlicingBranchBlock(inputSize, offset, size int) *Block {
	if size == 0 {
		panic("zero sized slice not allowed")
	}
	lookupName := fmt.Sprintf("slicing_branch_%d_%d_%d", inputSize, offset, size)
	ret := codeUnit.FindBranchBlock(lookupName)
	if ret != nil {
		return ret
	}
	ret = codeUnit.AddNewBranchBlock(lookupName, size)
	if offset != 0 {
		ret.Branch.AddInputSite(offset)
	}
	theSlice := ret.Branch.AddInputSite(size)
	if offset+size < inputSize {
		ret.Branch.AddInputSite(inputSize - offset - size)
	}
	output := NewMerge(theSlice).NewSite()
	output.SiteType = SITE_OUTPUT
	ret.Branch.AddNewSite(output, lookupName)
	return ret
}

func (codeUnit *CodeUnit) GetNullifyLUTBlock(trueFalse bool) *Block {
	strRepr := GetNullifyLUTRepr(trueFalse)
	return codeUnit.GetLUTBlock(strRepr)
}

func (codeUnit *CodeUnit) GetNullifyBranchBlock(size int, trueFalse bool) *Block {
	lookupName := fmt.Sprintf("nullify_%v_arg_%d", trueFalse, size)
	ret := codeUnit.FindBranchBlock(lookupName)
	if ret != nil {
		return ret
	}
	ret = codeUnit.AddNewBranchBlock(lookupName, size)
	ret.Branch.AddInputSite(1) // condition
	for i := 0; i < size; i++ {
		ret.Branch.AddInputSite(1) // arg
	}
	nullifyLutBlock := codeUnit.GetNullifyLUTBlock(trueFalse)
	condInput := ret.Branch.InputSites[0]
	for i := 0; i < size; i++ {
		nullifyTritKnot :=
			NewKnot(nullifyLutBlock, condInput, ret.Branch.InputSites[i+1], condInput)
		nullifyTritSite := nullifyTritKnot.NewSite()
		nullifyTritSite.SiteType = SITE_OUTPUT
		ret.Branch.AddNewSite(nullifyTritSite, "")
	}
	return ret
}
