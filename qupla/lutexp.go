package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/abra"
)

type LutExpr struct {
	ExpressionBase
	LutDef *LutDef
}

func (e *LutExpr) Size() int {
	if e == nil {
		return 0
	}
	return e.LutDef.Size()
}

func (e *LutExpr) Copy() ExpressionInterface {
	return &LutExpr{
		ExpressionBase: e.copyBase(),
		LutDef:         e.LutDef,
	}
}

func (e *LutExpr) Eval(frame *EvalFrame, result Trits) bool {
	var buf [3]int8 // no more than 3 inputs
	for i, a := range e.subExpr {
		if a.Eval(frame, buf[i:i+1]) {
			return true
		}
	}
	lutArg := buf[:e.LutDef.InputSize]
	null := e.LutDef.Lookup(result, lutArg)
	return null
}

func (e *LutExpr) GetAbraSite(branch *abra.Branch, codeUnit *abra.CodeUnit, lookupName string) *abra.Site {
	lut := codeUnit.FindLUTBlock(e.LutDef.GetStringRepr())
	if lut == nil {
		panic(fmt.Errorf("can't find lut block '%s'", e.LutDef.GetStringRepr()))
	}
	in0 := e.GetSubExpr(0).GetAbraSite(branch, codeUnit, "")
	in1 := in0
	if e.LutDef.InputSize > 1 {
		in1 = e.GetSubExpr(1).GetAbraSite(branch, codeUnit, "")
	}
	in2 := in1
	if e.LutDef.InputSize > 2 {
		in2 = e.GetSubExpr(2).GetAbraSite(branch, codeUnit, "")
	}
	ret := abra.NewKnot(lut, in0, in1, in2).NewSite()
	ret.SetLookupName(lookupName)
	return branch.GenOrUpdateSite(ret)
}
