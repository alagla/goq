package abra

func (code *Code) FindLUT(lookupName string) *LUT {
	for _, b := range code.Blocks {
		if b.BlockType == BLOCK_LUT && b.LookupName == lookupName {
			return b.LUT
		}
	}
	return nil
}

func (code *Code) AddLUT(lut *LUT, lookupName string) bool {
	if code.FindLUT(lookupName) != nil {
		return false
	}
	code.Blocks = append(code.Blocks, lut.NewBlock(lookupName))
	return true
}

func (code *Code) FindBranch(lookupName string) *Branch {
	for _, b := range code.Blocks {
		if b.BlockType == BLOCK_BRANCH && b.LookupName == lookupName {
			return b.Branch
		}
	}
	return nil
}

func (code *Code) AddBranch(branch *Branch, lookupName string) bool {
	if code.FindBranch(lookupName) != nil {
		return false
	}
	code.Blocks = append(code.Blocks, branch.NewBlock(lookupName))
	return true
}
