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
}
