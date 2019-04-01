package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
)

type MergeExpr struct {
	ExpressionBase
}

func (e *MergeExpr) Size() int {
	if e == nil {
		return 0
	}
	return e.subExpr[0].Size()
}

func (e *MergeExpr) Eval(frame *EvalFrame, result Trits) bool {
	if e.subExpr[0].Eval(frame, result) {
		return e.subExpr[1].Eval(frame, result)
	}
	return false
}
