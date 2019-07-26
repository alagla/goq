package abra

import (
	"fmt"
	. "github.com/lunfardo314/goq/cfg"
)

func PrintBlocks(codeUnit *CodeUnit, printSite func(int) bool, level int) {
	Logf(level, "------------------- Listing all blocks")
	Logf(level, "total %d LUTs, %d branches, %d external blocks",
		codeUnit.Code.NumLUTs, codeUnit.Code.NumBranches, codeUnit.Code.NumExternalBlocks)

	for _, block := range codeUnit.Code.Blocks {
		switch block.BlockType {
		case BLOCK_LUT:
			Logf(2, "  LUT  %4d %20s -> '%s' inSize = %d outSize = %d",
				block.Index, block.LUT.Name, block.LookupName, block.SizeIn, block.SizeOut)
		case BLOCK_BRANCH:
			st := GetStats(block.Branch)
			Logf(2, "  BRCH %4d %40s-> n_in: %2d, n_body: %2d, n_out: %2d, n_state: %2d, inSizes: %v=%d allSizes = %v outSiteIdx = %v",
				block.Index, block.LookupName, st.NumInputs, st.NumBodySites,
				st.NumOutputs, st.NumStateSites, st.InputSizes, st.InputSize, st.AllSizes, st.OutSites)
			if nil != printSite && printSite(block.Index) {
				PrintSites(codeUnit, block.Index, level)
			}
		case BLOCK_EXTERNAL:
			panic("implement me")
		}
	}
}

func PrintSites(codeUnit *CodeUnit, blockIndex int, level int) {
	if codeUnit.Code.Blocks[blockIndex].BlockType != BLOCK_BRANCH {
		return
	}
	Logf(level, "+++++++++++ Block #%d sites:", blockIndex)
	for _, s := range codeUnit.Code.Blocks[blockIndex].Branch.AllSites {
		if s.SiteType == SITE_INPUT {
			Logf(level, "      #%d  'input'.  Size = %d", s.Index, s.Size)
			continue
		}
		tn := ""
		switch s.SiteType {
		case SITE_BODY:
			tn = "body"
		case SITE_OUTPUT:
			tn = "output"
		case SITE_STATE:
			tn = "state"
		}
		instr := ""
		if s.IsKnot {
			instr = fmt.Sprintf("knot branch = #%d inp = %v inp_sizes = %v",
				s.Knot.Block.Index, getIndices(s.Knot.Sites), getSizes(s.Knot.Sites))
		} else {
			instr = fmt.Sprintf("merge inp = %v inp_sizes = %v", getIndices(s.Merge.Sites), getSizes(s.Merge.Sites))
		}
		Logf(level, "      #%d  '%s'. Size = %d. %s", s.Index, tn, s.Size, instr)
	}
}

func getIndices(sites []*Site) []int {
	ret := make([]int, 0, 10)
	for _, s := range sites {
		ret = append(ret, s.Index)
	}
	return ret
}

func getSizes(sites []*Site) []int {
	ret := make([]int, 0, 10)
	for _, s := range sites {
		ret = append(ret, s.Size)
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
	OutSites      []int
	AllSizes      []int
}

func GetStats(branch *Branch) *BranchStats {
	ret := &BranchStats{
		InputSizes: make([]int, 0, 5),
		OutSites:   make([]int, 0, 5),
		AllSizes:   make([]int, 0, 5),
	}
	for _, s := range branch.AllSites {
		ret.AllSizes = append(ret.AllSizes, s.Size)
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
			ret.OutSites = append(ret.OutSites, s.Index)
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
