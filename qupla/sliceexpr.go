package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
)

type SliceExpr struct {
	ExpressionBase
	//LocalVarIdx int
	//VarScope    *Function
	vi       *QuplaSite
	offset   int
	size     int
	sliceEnd int
	noSlice  bool
}

func NewQuplaSliceExpr(vi *QuplaSite, src string, offset, size int) *SliceExpr {
	noSlice := offset == 0 && size == vi.Size
	return &SliceExpr{
		ExpressionBase: NewExpressionBase(src),
		vi:             vi,
		offset:         offset,
		size:           size,
		sliceEnd:       offset + size,
		noSlice:        noSlice,
	}
}

func (e *SliceExpr) Size() int {
	if e == nil {
		return 0
	}
	return e.size
}

func (e *SliceExpr) Eval(frame *EvalFrame, result Trits) bool {
	restmp, null := e.vi.Eval(frame)
	if !null {
		if e.noSlice {
			copy(result, restmp)
		} else {
			copy(result, restmp[e.offset:e.sliceEnd])
		}
	}
	return null
}
