package abragen

import (
	"github.com/lunfardo314/goq/abra"
	"github.com/lunfardo314/goq/qupla"
)

func GenAbraCodeUnit(module *qupla.QuplaModule) *abra.CodeUnit {
	ret := abra.NewCodeUnit()

	for _, lut := range module.Luts {
		ret.NewLUT(lut.GetStringRepr(), lut.BinaryEncodedLUT())
	}
	for _, fun := range module.Functions {
		ret.NewBranch(fun.Name, fun.Size())
	}

	var fun *qupla.Function
	for _, b := range ret.Code.Blocks {
		if b.BlockType != abra.BLOCK_BRANCH {
			continue
		}
		fun = module.FindFuncDef(b.LookupName)
		GenAbraBranch(fun, b.Branch, ret)
	}
	return ret
}
