package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
)

type MergeExpr struct {
	ExpressionBase
}

func (e *MergeExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.subexpr[0].Size()
}

func (e *MergeExpr) Eval(frame *EvalFrame, result Trits) bool {
	if e.subexpr[0].Eval(frame, result) {
		return e.subexpr[1].Eval(frame, result)
	}
	return false
}
