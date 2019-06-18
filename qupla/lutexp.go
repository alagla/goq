package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/abra"
)

type LutExpr struct {
	ExpressionBase
	LutDef *LutDef
}

func (e *LutExpr) GenAbraSite(branch *abra.Branch, codeUnit *abra.CodeUnit) *abra.Site {
	panic("implement me")
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
