package generate

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abra"
)

type Tritcode Trits

func NewTritcode() Tritcode {
	return make(Trits, 0, 4096)
}

func WriteTrits(tcode Tritcode, trits Trits) Tritcode {
	return append(tcode, trits...)
}

func WritePosInt(tcode Tritcode, n int) Tritcode {
	return WriteTrits(tcode, MustInt2PosIntAsPerSpec(n))
}

//code:
//[ tritcode version (positive integer [0])
//, number of lookup table blocks (positive integer)
//, 35-trit lookup tables (27 nullable trits in bct)...
//, number of branch blocks (positive integer)
//, branch block definitions ...
//, number of external blocks (positive integer)
//, external block definitions...
//]

// always returns trits with len(trits) % 3 == 0

func WriteCode(tcode Tritcode, code *Code) Tritcode {
	switch {
	case code.NumLUTs > len(code.Blocks):
		panic("inconsistent Code 1")
	case code.NumLUTs+code.NumBranches > len(code.Blocks):
		panic("inconsistent Code 2")
	case code.NumLUTs+code.NumBranches+code.NumExternalBlocks != len(code.Blocks):
		panic("inconsistent Code 3")
	}

	tcode = WritePosInt(tcode, code.TritcodeVersion)

	tcode = WritePosInt(tcode, code.NumLUTs)
	// Expected to be sorted, first LUTs
	// write LUTs
	for i := 0; i < code.NumLUTs; i++ {
		if code.Blocks[i].BlockType != BLOCK_LUT {
			panic("inconsistent Code 4")
		}
		tcode = WriteLUTTrits(tcode, code.Blocks[i].LUT)
	}

	tcode = WritePosInt(tcode, code.NumBranches)
	// second Branches
	// write Branches
	for i := code.NumLUTs; i < code.NumLUTs+code.NumBranches; i++ {
		if code.Blocks[i].BlockType != BLOCK_BRANCH {
			panic("inconsistent Code 5")
		}
		tcode = WriteBranchTrits(tcode, code.Blocks[i].Branch)
	}

	tcode = WritePosInt(tcode, code.NumExternalBlocks)
	// second Branches
	// write Branches
	for i := code.NumLUTs + code.NumBranches; i < len(code.Blocks); i++ {
		if code.Blocks[i].BlockType != BLOCK_EXTERNAL {
			panic("inconsistent Code 6")
		}
		tcode = tcode.WriteExternalBlockTrits(code.Blocks[i].ExternalBlock)
	}
	// make len(result) % 3 == 0
	rem := len(tcode) % 3
	for ; rem%3 != 0; rem = (rem + 1) % 3 {
		tcode = WriteTrits(tcode, Trits{0})
	}
	if len(tcode)%3 != 0 {
		panic("tcode.Buf.Len() % 3 != 0")
	}
	return tcode
}

func WriteLUTTrits(tcode Tritcode, lut *LUT) Tritcode {
	tlut := TritEncodeLUTBinary(lut.Binary)
	tcode = WritePosInt(tcode, len(tlut)) // must be 35
	return WriteTrits(tcode, tlut)
}

//branch:
//[ number of inputs (positive integer)
//, input lengths (positive integers)...
//, number of body sites (positive integer)
//, number of output sites (positive integer)
//, number of memory latch sites (positive integer)
//, body site definitions...
//, output site definitions...
//, memory latch site definitions...
//]

func WriteBranchTrits(tcode Tritcode, branch *Branch) Tritcode {
	switch {
	case branch.NumInputs > branch.NumSites:
		panic("something wrong with enumerating sites in branch 1")
	case branch.NumInputs+branch.NumBodySites > branch.NumSites:
		panic("something wrong with enumerating sites in branch 2")
	case branch.NumInputs+branch.NumBodySites+branch.NumOutputs > branch.NumSites:
		panic("something wrong with enumerating sites in branch 3")
	case branch.NumInputs+branch.NumBodySites+branch.NumOutputs+branch.NumStateSites != branch.NumSites:
		panic("something wrong with enumerating sites in branch 4")
	}
	// first writing branch to determine length in trits
	tbranch := NewTritcode()
	tbranch = WritePosInt(tbranch, branch.NumInputs)

	for i := 0; i < branch.NumInputs; i++ {
		if branch.AllSites[i].SiteType != SITE_INPUT {
			panic("input site expected")
		}
		tbranch = WritePosInt(tbranch, branch.AllSites[i].Size)
	}
	tbranch = WritePosInt(tbranch, branch.NumBodySites)
	tbranch = WritePosInt(tbranch, branch.NumOutputs)
	tbranch = WritePosInt(tbranch, branch.NumStateSites)

	for i := branch.NumInputs; i < branch.NumInputs+branch.NumBodySites; i++ {
		if branch.AllSites[i].SiteType != SITE_BODY {
			panic("body site expected")
		}
		tbranch = tbranch.WriteSiteDefinition(branch.AllSites[i])
	}
	for i := branch.NumInputs + branch.NumBodySites; i < branch.NumInputs+branch.NumBodySites+branch.NumOutputs; i++ {
		if branch.AllSites[i].SiteType != SITE_OUTPUT {
			panic("output site expected")
		}
		tbranch = tbranch.WriteSiteDefinition(branch.AllSites[i])
	}

	for i := branch.NumInputs + branch.NumBodySites + branch.NumOutputs; i < branch.NumSites; i++ {
		if branch.AllSites[i].SiteType != SITE_STATE {
			panic("state site expected")
		}
		tbranch = tbranch.WriteSiteDefinition(branch.AllSites[i])
	}

	// writing block
	tcode = WritePosInt(tcode, len(tbranch))
	tcode = WriteTrits(tcode, tbranch)
	return tcode
}

func (tcode Tritcode) WriteSiteDefinition(site *Site) Tritcode {
	if site.SiteType == SITE_INPUT {
		return WritePosInt(tcode, site.Size)
	}
	if site.IsKnot {
		tcode = WriteTrits(tcode, Trits{-1})
		tcode = WritePosInt(tcode, len(site.Knot.Sites))
		for _, s := range site.Knot.Sites {
			tcode = WritePosInt(tcode, s.Index)
		}
		tcode = WritePosInt(tcode, site.Knot.Block.Index)
	} else {
		tcode = WriteTrits(tcode, Trits{1})
		tcode = WritePosInt(tcode, len(site.Merge.Sites))
		for _, s := range site.Merge.Sites {
			tcode = WritePosInt(tcode, s.Index)
		}
	}
	return tcode
}

func (tcode Tritcode) WriteExternalBlockTrits(external *ExternalBlock) Tritcode {
	panic("implement me")
}
