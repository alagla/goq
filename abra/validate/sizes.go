package validate

import . "github.com/lunfardo314/goq/abra"

func setZeroSizes(codeUnit *CodeUnit) {
	for _, b := range codeUnit.Code.Blocks {
		switch b.BlockType {
		case BLOCK_LUT:
			b.Size = 0
		case BLOCK_EXTERNAL:
			b.Size = 0
		case BLOCK_BRANCH:
			b.Size = 0
			setZeroSizesInBranch(b.Branch)
		}
	}
}

func CalcAllSizes(codeUnit *CodeUnit) {
	setZeroSizes(codeUnit)

	notFinished := true
	var saveSize int
	for notFinished {
		notFinished = false
		for _, block := range codeUnit.Code.Blocks {
			saveSize = block.Size
			CalcSizesInBlock(block)
			notFinished = notFinished || saveSize != block.Size
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

func setZeroSizesInBranch(branch *Branch) {
	for _, s := range branch.AllSites {
		if s.SiteType != SITE_INPUT {
			s.Size = 0
		}
	}
	branch.Size = 0
}

func CalcSizesInBlock(block *Block) {
	switch block.BlockType {
	case BLOCK_LUT:
		block.Size = 1
		return
	case BLOCK_EXTERNAL:
		panic("implement me")
	case BLOCK_BRANCH:
		CalcSizesInBranch(block.Branch)
		block.Size = block.Branch.Size
		return
	}
	panic("wrong block type")
}

func CalcSizesInBranch(branch *Branch) {
	for _, s := range branch.AllSites {
		CalcSiteSize(s)
	}
	branch.Size = 0
	for _, s := range branch.AllSites {
		if s.SiteType == SITE_OUTPUT {
			if s.Size == 0 {
				branch.Size = 0
				return
			}
			branch.Size += s.Size
		}
	}
}

func CalcSiteSize(site *Site) {
	if site.SiteType == SITE_INPUT {
		return
	}
	if site.IsKnot {
		site.Size = CalcKnotSize(site.Knot)
	} else {
		site.Size = CalcMergeSize(site.Merge)
	}
	site.CalculatedSize = true
}

func CalcKnotSize(knot *Knot) int {
	return knot.Block.Size
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
