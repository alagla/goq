package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
)

type SliceExpr struct {
	ExpressionBase
	LocalVarIdx int64
	VarScope    *Function
	offset      int64
	size        int64
}

func NewQuplaSliceExpr(src string, offset, size int64) *SliceExpr {
	return &SliceExpr{
		ExpressionBase: NewExpressionBase(src),
		offset:         offset,
		size:           size,
	}
}

func (e *SliceExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.size
}

func (e *SliceExpr) Eval(proc ProcessorInterface, result Trits) bool {
	restmp, null := proc.EvalVar(e.LocalVarIdx)
	if null {
		return true
	}
	copy(result, restmp[e.offset:e.offset+e.size])
	return false
}
