package validate

import (
	"fmt"
	"github.com/lunfardo314/goq/abra"
)

func CompareCodeUnits(codeUnit1, codeUnit2 *abra.CodeUnit) []error {
	ret := make([]error, 0, 10)
	if codeUnit1.Code.TritcodeVersion != codeUnit2.Code.TritcodeVersion {
		ret = append(ret, fmt.Errorf("codeUnit1.Code.TritcodeVersion != codeUnit2.Code.TritcodeVersion"))
	}
	if codeUnit1.Code.NumLUTs != codeUnit2.Code.NumLUTs {
		ret = append(ret, fmt.Errorf("codeUnit1.Code.NumLUTs != codeUnit2.Code.NumLUTs"))
	}
	if codeUnit1.Code.NumBranches != codeUnit2.Code.NumBranches {
		ret = append(ret, fmt.Errorf("codeUnit1.Code.NumBranches != codeUnit2.Code.NumBranches"))
	}
	if codeUnit1.Code.NumExternalBlocks != codeUnit2.Code.NumExternalBlocks {
		ret = append(ret, fmt.Errorf("codeUnit1.Code.NumExternalBlocks != codeUnit2.Code.NumExternalBlocks"))
	}
	if len(codeUnit1.Code.Blocks) != len(codeUnit2.Code.Blocks) {
		ret = append(ret, fmt.Errorf("len(codeUnit1.Code.Blocks) != len(codeUnit2.Code.Blocks)"))
		return ret //doesn't make sense to continue
	}

	for i := range codeUnit1.Code.Blocks {
		errs := CompareBlocks(codeUnit1.Code.Blocks[i], codeUnit2.Code.Blocks[i], i)
		ret = append(ret, errs...)
	}
	return ret
}

func CompareBlocks(block1, block2 *abra.Block, assumedIndex int) []error {
	ret := make([]error, 0, 10)
	if block1.Index != block2.Index {
		ret = append(ret, fmt.Errorf("block #%d: block1.Index != block2.Index", assumedIndex))
	}
	if block1.SizeOut != block2.SizeOut {
		ret = append(ret, fmt.Errorf("block #%d: block1.Size != block2.Size", assumedIndex))
	}
	if block1.BlockType != block2.BlockType {
		ret = append(ret, fmt.Errorf("block #%d: block1.BlockType != block2.BlockType", assumedIndex))
		return ret
	}
	switch block1.BlockType {
	case abra.BLOCK_LUT:
		if block1.LUT.Binary != block2.LUT.Binary {
			ret = append(ret, fmt.Errorf("LUT block #%d: block1.LUT.Binary != block2.LUT.Binary", assumedIndex))
		}
	case abra.BLOCK_BRANCH:
		if block1.Branch.NumInputs != block2.Branch.NumInputs {
			ret = append(ret, fmt.Errorf("branch block #%d: block1.Branch.NumInputs != block2.Branch.NumInputs", assumedIndex))
		}
		if block1.Branch.NumBodySites != block2.Branch.NumBodySites {
			ret = append(ret, fmt.Errorf("branch block #%d: block1.Branch.NumBodySites != block2.Branch.NumBodySites", assumedIndex))
		}
		if block1.Branch.NumOutputs != block2.Branch.NumOutputs {
			ret = append(ret, fmt.Errorf("branch block #%d: block1.Branch.NumOutputs != block2.Branch.NumOutputs", assumedIndex))
		}
		if block1.Branch.NumStateSites != block2.Branch.NumStateSites {
			ret = append(ret, fmt.Errorf("branch block #%d: block1.Branch.NumStateSites != block2.Branch.NumStateSites", assumedIndex))
		}
		if block1.Branch.NumSites != block2.Branch.NumSites {
			ret = append(ret, fmt.Errorf("branch block #%d: block1.Branch.NumStateSites != block2.Branch.NumStateSites", assumedIndex))
		}
		if len(block1.Branch.AllSites) != len(block2.Branch.AllSites) {
			ret = append(ret, fmt.Errorf("branch block #%d: len(block1.Branch.AllSites) != len(block2.Branch.AllSites)", assumedIndex))
		}

		for i := range block1.Branch.AllSites {
			err := CompareSites(block1.Branch.AllSites[i], block2.Branch.AllSites[i])
			if err != nil {
				ret = append(ret, fmt.Errorf("branch block #%d, site #%d: %v", assumedIndex, i, err))
			}
		}
	case abra.BLOCK_EXTERNAL:
		panic("implement me")
	default:
		ret = append(ret, fmt.Errorf("block #%d: wrong BlockType", assumedIndex))
		return ret
	}
	return ret
}

func CompareSites(site1, site2 *abra.Site) error {
	if site1.SiteType != site2.SiteType {
		return fmt.Errorf("site1.SiteType != site2.SiteType")
	}
	if site1.Index != site2.Index {
		return fmt.Errorf("site1.Index != site2.Index")
	}
	if site1.Size != site2.Size {
		return fmt.Errorf("site1.Size != site2.Size")
	}
	if site1.SiteType == abra.SITE_INPUT {
		return nil
	}
	if site1.IsKnot != site2.IsKnot {
		return fmt.Errorf("site1.IsKnot != site2.IsKnot")
	}
	if site1.IsKnot {
		if len(site1.Knot.Sites) != len(site2.Knot.Sites) {
			return fmt.Errorf("len(site1.Knot.Sites) != len(site2.Knot.Sites")
		}
		for i := range site1.Knot.Sites {
			if site1.Knot.Sites[i] == nil || site2.Knot.Sites[i] == nil || site1.Knot.Sites[i].Index != site2.Knot.Sites[i].Index {
				return fmt.Errorf("site1.Knot.Sites[i] == nil || site2.Knot.Sites[i] == nil || site1.Knot.Sites[i] != site2.Knot.Sites[i] at i = %d", i)
			}
		}
		if site1.Knot.Block == nil || site2.Knot.Block == nil || site1.Knot.Block.Index != site2.Knot.Block.Index {
			return fmt.Errorf("site1.Knot.Block == nil || site2.Knot.Block == nil ||site1.Knot.Block.Index != site2.Knot.Block.Index")
		}
	} else {
		if len(site1.Merge.Sites) != len(site2.Merge.Sites) {
			return fmt.Errorf("len(site1.Merge.Sites) != len(site2.Merge.Sites)")
		}
		for i := range site1.Merge.Sites {
			if site1.Merge.Sites[i] == nil || site2.Merge.Sites[i] == nil || site1.Merge.Sites[i].Index != site2.Merge.Sites[i].Index {
				return fmt.Errorf("site1.Merge.Sites[i] == nil || site2.Merge.Sites[i] == nil || site1.Merge.Sites[i] != site2.Merge.Sites[i] at i = %d", i)
			}
		}
	}
	return nil
}
