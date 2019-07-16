package validate

import (
	"fmt"
	. "github.com/lunfardo314/goq/abra"
)

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

func Validate(codeUnit *CodeUnit) []error {
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
			if block.Branch.Size == 0 || GetBranchInputSize(block.Branch) == 0 {
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
						if GetBranchInputSize(s.Knot.Block.Branch) != GetKnotInputSize(s.Knot) {
							ret = append(ret, fmt.Errorf("sum of sizes of the inputs (%d) != branch size (%d) in knot '%s' of branch '%s'",
								GetKnotInputSize(s.Knot), GetBranchInputSize(s.Knot.Block.Branch), s.LookupName, block.LookupName))
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

type BranchStats struct {
	NumSites      int
	NumInputs     int
	NumBodySites  int
	NumStateSites int
	NumOutputs    int
	NumKnots      int
	NumMerges     int
	InputSizes    []int
	InputSize     int
}

func GetStats(branch *Branch) *BranchStats {
	ret := &BranchStats{
		InputSizes: make([]int, 0, 5),
	}
	for _, s := range branch.AllSites {
		switch s.SiteType {
		case SITE_INPUT:
			ret.NumInputs++
			ret.InputSizes = append(ret.InputSizes, s.Size)
		case SITE_BODY:
			ret.NumBodySites++
		case SITE_STATE:
			ret.NumStateSites++
		case SITE_OUTPUT:
			ret.NumOutputs++
		}
		if s.IsKnot {
			ret.NumKnots++
		} else {
			ret.NumMerges++
		}
	}
	ret.NumSites = len(branch.AllSites)
	for _, s := range ret.InputSizes {
		ret.InputSize += s
	}
	return ret
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

func AssertValid(branch *Branch) {
	// assert
	for _, s := range branch.AllSites {
		if s.SiteType == SITE_OUTPUT && s.Merge == nil && s.Knot == nil {
			panic("assert: wrong output")
		}
	}
	for _, s := range branch.AllSites {
		AssertValidSite(s)
	}
}

func GetBlockInputSize(block *Block) int {
	switch block.BlockType {
	case BLOCK_LUT:
		return 3
	case BLOCK_BRANCH:
		return GetBranchInputSize(block.Branch)
	}
	panic("implement me")
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

func AssertValidSite(site *Site) {
	switch site.SiteType {
	case SITE_INPUT:
		if site.Knot != nil {
			panic("invalid site 1")
		}
		if site.Merge != nil {
			panic("invalid site 2")
		}
	case SITE_BODY:
		if (site.Merge == nil) == (site.Knot == nil) {
			panic("invalid site 3")
		}
		if site.IsKnot {
			if len(site.Knot.Sites) == 0 {
				panic("invalid site 4")
			}
		} else {
			if len(site.Merge.Sites) == 0 {
				panic("invalid site 5")
			}
		}
	case SITE_OUTPUT:
		if (site.Merge == nil) == (site.Knot == nil) {
			panic("invalid site 6")
		}
		if site.IsKnot {
			if len(site.Knot.Sites) == 0 {
				panic("invalid site 7")
			}
		} else {
			if len(site.Merge.Sites) == 0 {
				panic("invalid site 8")
			}
		}
	case SITE_STATE:
		if (site.Merge == nil) == (site.Knot == nil) {
			panic("invalid site 9")
		}
		if site.IsKnot {
			if len(site.Knot.Sites) == 0 {
				panic("invalid site 10")
			}
		} else {
			if len(site.Merge.Sites) == 0 {
				panic("invalid site 11")
			}
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
