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
