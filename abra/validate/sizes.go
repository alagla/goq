package validate

import (
	"fmt"
	. "github.com/lunfardo314/goq/abra"
)

func setZeroSizes(codeUnit *CodeUnit) {
	for _, b := range codeUnit.Code.Blocks {
		switch b.BlockType {
		case BLOCK_LUT, BLOCK_EXTERNAL:
			b.SizeOut = 0
			b.SizeIn = 0
		case BLOCK_BRANCH:
			b.SizeOut = 0
			b.SizeIn = 0
			for _, s := range b.Branch.AllSites {
				if s.SiteType != SITE_INPUT {
					s.Size = 0
				}
			}
		}
	}
}

func CalcAllSizes(codeUnit *CodeUnit) {
	setZeroSizes(codeUnit)
	changedBlock := true
	for changedBlock {
		changedBlock = false
		for _, block := range codeUnit.Code.Blocks {
			changed := CalcSizesInBlock(block)
			changedBlock = changedBlock || changed
		}
	}
}

func GetBranchInputSize(branch *Branch) int {
	ret := 0
	for _, s := range branch.AllSites {
		if s.SiteType == SITE_INPUT {
			ret += s.Size
		}
	}
	return ret
}

func CalcSizesInBlock(block *Block) bool {
	//cfg.Logf(0, "calculating size block #%d", block.Index)
	switch block.BlockType {
	case BLOCK_LUT:
		if block.SizeOut != 1 || block.SizeIn != 3 {
			block.SizeOut = 1
			block.SizeIn = 3
			return true
		}
		return false
	case BLOCK_EXTERNAL:
		panic("implement me")
	case BLOCK_BRANCH:
		return CalcSizesInBranchBlock(block)
	}
	panic("wrong block type")
}

func CalcSizesInBranchBlock(block *Block) bool {
	if block.Index == 36 {
		fmt.Println("kuku")
	}
	if block.BlockType != BLOCK_BRANCH {
		panic("inconsistency")
	}
	anything_changed := false
	saveIn := block.SizeIn
	saveOut := block.SizeOut

	branch := block.Branch
	block.SizeIn = GetBranchInputSize(branch)

	anySiteChanged := true
	for anySiteChanged {
		anySiteChanged = false
		for _, s := range branch.AllSites {
			changed := CalcSiteSize(s)
			anySiteChanged = anySiteChanged || changed
			anything_changed = anything_changed || anySiteChanged
		}
	}

	block.SizeOut = 0
	for _, s := range branch.AllSites {
		if s.SiteType == SITE_OUTPUT {
			if s.Size == 0 {
				block.SizeOut = 0
				break
			}
			block.SizeOut += s.Size
		}
	}
	return anything_changed || saveIn != block.SizeIn || saveOut != block.SizeOut
}

func CalcSiteSize(site *Site) bool {
	saveSize := site.Size
	if site.SiteType != SITE_INPUT {
		if site.IsKnot {
			site.Size = CalcKnotSize(site.Knot)
		} else {
			site.Size = CalcMergeSize(site.Merge)
		}
	}
	site.CalculatedSize = true
	return saveSize != site.Size
}

func CalcKnotSize(knot *Knot) int {
	return knot.Block.SizeOut
}

func GetKnotInputSize(knot *Knot) int {
	ret := 0
	for _, s := range knot.Sites {
		if s.Size == 0 {
			return 0
		}
		ret += s.Size
	}
	return ret
}

func CalcMergeSize(merge *Merge) int {
	ret := 0
	for _, s := range merge.Sites {
		if s.Size != 0 {
			ret = s.Size
		}
	}
	return ret
}
