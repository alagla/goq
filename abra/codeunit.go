package abra

func (codeUnit *CodeUnit) FindLUTBlock(lookupName string) *Block {
	for _, b := range codeUnit.Code.Blocks {
		if b.BlockType == BLOCK_LUT && b.LookupName == lookupName {
			return b
		}
	}
	return nil
}

func (codeUnit *CodeUnit) FindBranchBlock(lookupName string) *Block {
	for _, b := range codeUnit.Code.Blocks {
		if b.BlockType == BLOCK_BRANCH && b.LookupName == lookupName {
			return b
		}
	}
	return nil
}
