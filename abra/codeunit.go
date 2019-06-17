package abra

import "fmt"

func (codeUnit *CodeUnit) FindLUTBlock(lookupName string) *Block {
	for _, b := range codeUnit.Code.Blocks {
		if b.BlockType == BLOCK_LUT && b.LookupName == lookupName {
			return b
		}
	}
	return nil
}

func (codeUnit *CodeUnit) AddLUTBlock(lut LUT, lookupName string) bool {
	if codeUnit.FindLUTBlock(lookupName) != nil {
		return false
	}
	codeUnit.Code.Blocks = append(codeUnit.Code.Blocks, lut.NewBlock(lookupName))
	return true
}

func (codeUnit *CodeUnit) FindBranchBlock(lookupName string) *Block {
	for _, b := range codeUnit.Code.Blocks {
		if b.BlockType == BLOCK_BRANCH && b.LookupName == lookupName {
			return b
		}
	}
	return nil
}

func (codeUnit *CodeUnit) AddBranchBlock(branch *Branch, lookupName string) bool {
	if codeUnit.FindBranchBlock(lookupName) != nil {
		return false
	}
	codeUnit.Code.Blocks = append(codeUnit.Code.Blocks, branch.NewBlock(lookupName))
	return true
}

func (codeUnit *CodeUnit) GetConcatBlockForSize(size int) *Block {
	lookupName := fmt.Sprintf("concat_branch_%d", size)
	ret := codeUnit.FindBranchBlock(lookupName)
	if ret != nil {
		return ret
	}
	ret = codeUnit.NewBranchBlock(lookupName, size)
	input := ret.Branch.AddInputSite(size)
	ret.Branch.AddOutputSite(NewMerge(input).NewSite(lookupName + "_out"))
	return ret
}

// returns or creates block which takes to output 0 trit of it input

func (codeUnit *CodeUnit) GetLsbSliceBlock() *Block {
	lookupName := "LSB_SLICE_BRANCH_BLOCK"
	ret := codeUnit.FindBranchBlock(lookupName)
	if ret != nil {
		return ret
	}
	ret = codeUnit.NewBranchBlock(lookupName, 1)
	input := ret.Branch.AddInputSite(1)
	output := NewMerge(input).NewSite(lookupName + "_merge_site")
	ret.Branch.AddOutputSite(output) //
	return ret
}
