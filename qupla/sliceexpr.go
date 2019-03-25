package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
)

type SliceExpr struct {
	ExpressionBase
	LocalVarIdx int
	VarScope    *Function
	offset      int
	size        int
	sliceEnd    int
}

func NewQuplaSliceExpr(src string, offset, size int) *SliceExpr {
	return &SliceExpr{
		ExpressionBase: NewExpressionBase(src),
		offset:         offset,
		size:           size,
		sliceEnd:       offset + size,
	}
}

func (e *SliceExpr) Size() int {
	if e == nil {
		return 0
	}
	return e.size
}

func (e *SliceExpr) Eval(frame *EvalFrame, result Trits) bool {
	restmp, null := frame.EvalVar(e.LocalVarIdx)
	if null {
		return true
	}
	copy(result, restmp[e.offset:e.sliceEnd])
	return false
}
