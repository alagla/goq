package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
)

type QuplaLutExpr struct {
	QuplaExprBase
	ArgExpr []ExpressionInterface
	LutDef  *QuplaLutDef
}

func (e *QuplaLutExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.LutDef.Size()
}

func (e *QuplaLutExpr) Eval(proc ProcessorInterface, result Trits) bool {
	var buf [3]int8 // no more than 3 inputs
	for i, a := range e.ArgExpr {
		if proc.Eval(a, buf[i:i+1]) {
			return true
		}
	}
	lutArg := buf[:e.LutDef.InputSize]
	null := e.LutDef.Lookup(result, lutArg)
	return null
}
