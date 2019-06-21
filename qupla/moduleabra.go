package qupla

import "github.com/lunfardo314/goq/abra"

func (module *QuplaModule) GetAbra(codeUnit *abra.CodeUnit) {
	// TODO environments etc
	for _, lut := range module.Luts {
		strRepr := lut.GetStringRepr()
		codeUnit.NewLUTBlock(strRepr, abra.BinaryEncodedLUTFromString(strRepr))
	}

	for _, fun := range module.Functions {
		fun.GetAbraBranchBlock(codeUnit)
	}
}
