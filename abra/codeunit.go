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

func (codeUnit *CodeUnit) Validate() []error {
	ret := make([]error, 0, 10)
	for _, block := range codeUnit.Code.Blocks {
		if block.AssumedSize != block.Size || block.Size == 0 {
			ret = append(ret, fmt.Errorf("AssumedSize (%d) != Size (%d) in block '%s'",
				block.AssumedSize, block.Size, block.LookupName))
		}
		switch block.BlockType {
		case BLOCK_LUT:
			if block.Size != 1 {
				ret = append(ret, fmt.Errorf("LUT size != 1 in '%s'", block.LookupName))
			}
		case BLOCK_BRANCH:
			if block.Branch.Size == 0 || block.Branch.GetInputSize() == 0 {
				ret = append(ret, fmt.Errorf("wrong branch size %d in '%s'", block.Branch.Size, block.LookupName))
			}
			for _, s := range block.Branch.AllSites {
				if s.Size == 0 || s.Size != s.AssumedSize {
					ret = append(ret, fmt.Errorf("site.AssumedSize (%d) != site.Size (%d) in site '%s' of block '%s'",
						s.AssumedSize, s.Size, s.LookupName, block.LookupName))
				}
				if s.SiteType != SITE_INPUT && s.Knot == nil && s.Merge == nil {
					ret = append(ret, fmt.Errorf("inconsistent site '%s' in branch '%s'", s.LookupName, block.LookupName))
					continue
				}
				if s.SiteType != SITE_INPUT && s.IsKnot {
					if s.Knot.Block.BlockType == BLOCK_BRANCH {
						if s.Knot.Block.Branch.GetInputSize() != s.Knot.GetInputSize() {
							ret = append(ret, fmt.Errorf("sum of sizes of the inputs (%d) != branch size (%d) in knot '%s' of branch '%s'",
								s.Knot.GetInputSize(), s.Knot.Block.Branch.GetInputSize(), s.LookupName, block.LookupName))
						}
					}
				}
			}
		case BLOCK_EXTERNAL:
			panic("implement me")
		}
	}
	return ret
}
