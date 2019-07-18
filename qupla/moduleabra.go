package qupla

import (
	"github.com/lunfardo314/goq/abra"
	cabra "github.com/lunfardo314/goq/abra/construct"
	vabra "github.com/lunfardo314/goq/abra/validate"
	. "github.com/lunfardo314/goq/cfg"
)

func (module *QuplaModule) GetAbra(codeUnit *abra.CodeUnit) {
	// TODO environments etc
	Logf(2, "---- generating LUT blocks")
	count := 0
	for _, lut := range module.Luts {
		strRepr := lut.GetStringRepr()
		if cabra.FindLUTBlock(codeUnit, strRepr) != nil {
			continue
		}
		cabra.MustAddNewLUTBlock(codeUnit, strRepr, lut.Name)
		count++
	}

	Logf(2, "---- generating branch blocks")
	for _, fun := range module.Functions {
		fun.GetAbraBranchBlock(codeUnit)
	}

	vabra.SortAndEnumerateBlocks(codeUnit)
	vabra.SortAndEnumerateSites(codeUnit)

	Logf(0, "total %d LUTs, %d branches, %d external blocks",
		codeUnit.Code.NumLUTs, codeUnit.Code.NumBranches, codeUnit.Code.NumExternalBlocks)

	for _, block := range codeUnit.Code.Blocks {
		switch block.BlockType {
		case abra.BLOCK_LUT:
			Logf(2, "  LUT  %4d %20s -> '%s'", block.Index, block.LUT.Name, block.LookupName)
		case abra.BLOCK_BRANCH:
			st := vabra.GetStats(block.Branch)
			Logf(2, "  BRCH %4d %40s -> in: %2d, out: %2d, body: %2d, state: %2d, knots: %2d, merges: %2d inSizes: %v=%d",
				block.Index, block.LookupName, st.NumInputs, st.NumOutputs, st.NumBodySites, st.NumStateSites, st.NumKnots, st.NumMerges, st.InputSizes, st.InputSize)
		case abra.BLOCK_EXTERNAL:
			panic("implement me")
		}
	}
}
