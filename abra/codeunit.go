package abra

func (codeUnit *CodeUnit) NewLUTBlock(lookupName string, binaryEncodedLUT int64) *Block {
	ret := LUT(binaryEncodedLUT)
	block := ret.NewBlock(lookupName)
	if codeUnit.addBlock(block) {
		return block
	}
	return nil
}

func (lut LUT) NewBlock(lookupName string) *Block {
	return &Block{
		BlockType:  BLOCK_LUT,
		LUT:        lut,
		LookupName: lookupName,
	}
}

func (codeUnit *CodeUnit) FindLUTBlock(lookupName string) *Block {
	for _, b := range codeUnit.Code.Blocks {
		if b.BlockType == BLOCK_LUT && b.LookupName == lookupName {
			return b
		}
	}
	return nil
}

func (codeUnit *CodeUnit) GetLUTBlock(reprString string) *Block {
	ret := codeUnit.FindLUTBlock(reprString)
	if ret != nil {
		return ret
	}
	lut := BinaryEncodedLUTFromString(reprString)
	ret = codeUnit.NewLUTBlock(reprString, lut)
	return ret
}

func (codeUnit *CodeUnit) FindBranchBlock(lookupName string) *Block {
	for _, b := range codeUnit.Code.Blocks {
		if b.BlockType == BLOCK_BRANCH && b.LookupName == lookupName {
			return b
		}
	}
	return nil
}
