package abra

import (
	"fmt"
)

const TRITCODE_VERSION = 0

func NewCodeUnit() *CodeUnit {
	return &CodeUnit{
		EntityAttachments: make([]*EntityAttachment, 0, 10),
		Code: &Code{
			TritcodeVersion: TRITCODE_VERSION,
			Blocks:          make([]*Block, 0, 100),
		},
	}
}

func (codeUnit *CodeUnit) AddNewBlock(block *Block) bool {
	for _, b := range codeUnit.Code.Blocks {
		if b.BlockType == block.BlockType && b.LookupName == block.LookupName {
			return false
		}
	}
	codeUnit.Code.Blocks = append(codeUnit.Code.Blocks, block)
	block.Index = len(codeUnit.Code.Blocks)
	return true
}

func (codeUnit *CodeUnit) AddNewLUTBlock(lookupName string, binaryEncodedLUT int64) *Block {
	ret := LUT(binaryEncodedLUT)
	block := ret.NewBlock(lookupName)
	if codeUnit.AddNewBlock(block) {
		return block
	}
	panic(fmt.Errorf("LUT block '%s' already exists", lookupName))
}

func (lut LUT) NewBlock(lookupName string) *Block {
	return &Block{
		BlockType:   BLOCK_LUT,
		LUT:         lut,
		LookupName:  lookupName,
		AssumedSize: 1,
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
	if len(reprString) != 27 {
		panic("wrong LUT reprString")
	}
	ret := codeUnit.FindLUTBlock(reprString)
	if ret != nil {
		return ret
	}
	lut := BinaryEncodedLUTFromString(reprString)
	ret = codeUnit.AddNewLUTBlock(reprString, lut)
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

func (codeUnit *CodeUnit) setZeroSizes() {
	for _, b := range codeUnit.Code.Blocks {
		switch b.BlockType {
		case BLOCK_LUT:
			b.Size = 0
		case BLOCK_EXTERNAL:
			b.Size = 0
		case BLOCK_BRANCH:
			b.Size = 0
			b.Branch.setZeroSize()
		}
	}
}

func (codeUnit *CodeUnit) CalcSizes() {
	codeUnit.setZeroSizes()

	notFinished := true
	var saveSize int
	for notFinished {
		notFinished = false
		for _, block := range codeUnit.Code.Blocks {
			saveSize = block.Size
			block.CalcSizes()
			notFinished = notFinished || saveSize != block.Size
		}
	}
}
