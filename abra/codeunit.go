package abra

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
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

func (codeUnit *CodeUnit) NewEntityAttachment(codeHash Hash) *EntityAttachment {
	return &EntityAttachment{
		CodeHash:    codeHash,
		Attachments: make([]*Attachment, 0, 10),
	}
}

func (codeUnit *CodeUnit) NewAttachment(branch *Branch) *Attachment {
	return &Attachment{
		Branch:                 branch,
		MaximumRecursionDepth:  10, // tmp
		InputEnvironments:      make([]*InputEnvironmentData, 0, 5),
		OutputEnvironments:     make([]*OutputEnvironmentData, 0, 5),
		InputEnvironmentsDict:  make(map[Hash]*InputEnvironmentData),
		OutputEnvironmentsDict: make(map[Hash]*OutputEnvironmentData),
	}
}

func (att *Attachment) Join(envHash Hash, limit int) *InputEnvironmentData {
	ret := &InputEnvironmentData{
		EnvironmentHash: envHash,
		Limit:           limit,
	}
	att.InputEnvironments = append(att.InputEnvironments, ret)
	return ret
}

func (att *Attachment) Affect(envHash Hash, delay int) *OutputEnvironmentData {
	ret := &OutputEnvironmentData{
		EnvironmentHash: envHash,
		Delay:           delay,
	}
	att.OutputEnvironments = append(att.OutputEnvironments, ret)
	return ret
}

func (codeUnit *CodeUnit) CheckSizes() map[string]interface{} {
	ret := make(map[string]interface{})
	var res interface{}
	for i, b := range codeUnit.Code.Blocks {
		if b.LookupName == "qupla_function_abs_108" {
			fmt.Printf("kuku\n")
		}
		sz, err := b.GetSize()
		if err == nil {
			res = sz
		} else {
			res = err
		}
		ret[fmt.Sprintf("%s.%d", b.LookupName, i)] = res
	}
	return ret
}
