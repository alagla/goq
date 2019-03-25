package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
)

type ConcatExpr struct {
	ExpressionBase
}

func (e *ConcatExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.subexpr[0].Size() + e.subexpr[1].Size()
}

func (e *ConcatExpr) Eval(frame *EvalFrame, result Trits) bool {
	null := e.subexpr[0].Eval(frame, result)
	if null {
		return true
	}
	return e.subexpr[1].Eval(frame, result[e.subexpr[0].Size():])
}
