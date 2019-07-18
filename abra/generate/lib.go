package generate

import (
	"bytes"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abra"
)

type Tritcode struct {
	Buf bytes.Buffer
}

func NewTritcode() *Tritcode {
	return &Tritcode{}
}

func (tcode *Tritcode) writeTrits(trits Trits) error {
	for _, v := range trits {
		_, err := tcode.Buf.Write([]byte{byte(v)})
		if err != nil {
			return err
		}
	}
	return nil
}

func (tcode *Tritcode) MustWriteTrits(trits Trits) {
	if err := tcode.writeTrits(trits); err != nil {
		panic(err)
	}
}

func (tcode *Tritcode) MustWritePosInt(n int) {
	if err := tcode.writeTrits(MustInt2PosIntAsPerSpec(n)); err != nil {
		panic(err)
	}
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

func (tcode *Tritcode) MustWriteCode(code *Code) int {
	switch {
	case code.NumLUTs > len(code.Blocks):
		panic("inconsistent Code 1")
	case code.NumLUTs+code.NumBranches > len(code.Blocks):
		panic("inconsistent Code 2")
	case code.NumLUTs+code.NumBranches+code.NumExternalBlocks != len(code.Blocks):
		panic("inconsistent Code 3")
	}

	tcode.MustWritePosInt(code.TritcodeVersion)

	tcode.MustWritePosInt(code.NumLUTs)
	// Expected to be sorted, first LUTs
	// write LUTs
	for i := 0; i < code.NumLUTs; i++ {
		if code.Blocks[i].BlockType != BLOCK_LUT {
			panic("inconsistent Code 4")
		}
		tcode.MustWriteLUTTrits(code.Blocks[i].LUT)
	}

	tcode.MustWritePosInt(code.NumBranches)
	// second Branches
	// write Branches
	for i := code.NumLUTs; i < code.NumLUTs+code.NumBranches; i++ {
		if code.Blocks[i].BlockType != BLOCK_BRANCH {
			panic("inconsistent Code 5")
		}
		tcode.MustWriteBranchTrits(code.Blocks[i].Branch)
	}

	tcode.MustWritePosInt(code.NumExternalBlocks)
	// second Branches
	// write Branches
	for i := code.NumLUTs + code.NumBranches; i < len(code.Blocks); i++ {
		if code.Blocks[i].BlockType != BLOCK_EXTERNAL {
			panic("inconsistent Code 6")
		}
		tcode.MustWriteExternalBlockTrits(code.Blocks[i].ExternalBlock)
	}
	// make len(result) % 3 == 0
	rem := tcode.Buf.Len() % 3
	for i := 0; i < rem; i++ {
		tcode.MustWriteTrits(Trits{0})
	}
	if tcode.Buf.Len()%3 != 0 {
		panic("tcode.Buf.Len() % 3 != 0")
	}
	return tcode.Buf.Len()
}

func (tcode *Tritcode) MustWriteLUTTrits(lut *LUT) {
	tlut := TritEncodeLUTBinary(lut.Binary)
	tcode.MustWritePosInt(len(tlut)) // must be 35
	tcode.MustWriteTrits(tlut)
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

func (tcode *Tritcode) MustWriteBranchTrits(branch *Branch) {
	switch {
	case branch.NumInputs > branch.SiteIndexCount:
		panic("something wrong with enumerating sites in branch 1")
	case branch.NumInputs+branch.NumBodySites > branch.SiteIndexCount:
		panic("something wrong with enumerating sites in branch 2")
	case branch.NumInputs+branch.NumBodySites+branch.NumOutputs > branch.SiteIndexCount:
		panic("something wrong with enumerating sites in branch 3")
	case branch.NumInputs+branch.NumBodySites+branch.NumOutputs+branch.NumStateSites != branch.SiteIndexCount:
		panic("something wrong with enumerating sites in branch 4")
	}
	// first writing branch to determine length in trits
	tbranch := NewTritcode()
	tbranch.MustWritePosInt(branch.NumInputs)

	for i := 0; i < branch.NumInputs; i++ {
		if branch.AllSites[i].SiteType != SITE_INPUT {
			panic("input site expected")
		}
		tbranch.MustWritePosInt(branch.AllSites[i].Size)
	}
	tbranch.MustWritePosInt(branch.NumBodySites)
	tbranch.MustWritePosInt(branch.NumOutputs)
	tbranch.MustWritePosInt(branch.NumStateSites)

	for i := branch.NumInputs; i < branch.NumInputs+branch.NumBodySites; i++ {
		if branch.AllSites[i].SiteType != SITE_BODY {
			panic("body site expected")
		}
		tbranch.MustWriteSiteDefinition(branch.AllSites[i])
	}
	for i := branch.NumInputs + branch.NumBodySites; i < branch.NumInputs+branch.NumBodySites+branch.NumOutputs; i++ {
		if branch.AllSites[i].SiteType != SITE_OUTPUT {
			panic("output site expected")
		}
		tbranch.MustWriteSiteDefinition(branch.AllSites[i])
	}

	for i := branch.NumInputs + branch.NumBodySites + branch.NumOutputs; i < branch.SiteIndexCount; i++ {
		if branch.AllSites[i].SiteType != SITE_STATE {
			panic("state site expected")
		}
		tbranch.MustWriteSiteDefinition(branch.AllSites[i])
	}

	// writing block
	tcode.MustWritePosInt(len(tbranch.Buf.Bytes()))
	tcode.MustWriteTrits(Bytes2Trits(tbranch.Buf.Bytes()))
}

func (tcode *Tritcode) MustWriteSiteDefinition(site *Site) {
	if site.SiteType == SITE_INPUT {
		tcode.MustWritePosInt(site.Size)
		return
	}
	if site.IsKnot {
		tcode.MustWriteTrits(Trits{-1})
		tcode.MustWritePosInt(len(site.Knot.Sites))
		for _, s := range site.Knot.Sites {
			tcode.MustWritePosInt(s.Index)
		}
		tcode.MustWritePosInt(site.Knot.Block.Index)
	} else {
		tcode.MustWriteTrits(Trits{1})
		tcode.MustWritePosInt(len(site.Merge.Sites))
		for _, s := range site.Merge.Sites {
			tcode.MustWritePosInt(s.Index)
		}
	}
}

func (tcode *Tritcode) MustWriteExternalBlockTrits(external *ExternalBlock) {
	panic("implement me")
}
