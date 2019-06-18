package abra

func (external *ExternalBlock) NewBlock() *Block {
	return &Block{
		BlockType:     BLOCK_EXTERNAL,
		ExternalBlock: external,
	}
}
