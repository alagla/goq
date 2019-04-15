package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
)

type SliceInline struct {
	ExpressionBase
	Offset   int
	size     int
	SliceEnd int
	NoSlice  bool
	oneTrit  bool
}

func NewSliceInline(sliceExpr *SliceExpr, expr ExpressionInterface) *SliceInline {
	ret := &SliceInline{
		ExpressionBase: NewExpressionBase(sliceExpr.GetSource()),
		Offset:         sliceExpr.offset,
		size:           sliceExpr.size,
		SliceEnd:       sliceExpr.sliceEnd,
		NoSlice:        sliceExpr.noSlice,
		oneTrit:        sliceExpr.oneTrit,
	}
	ret.AppendSubExpr(expr)
	return ret
}

func (e *SliceInline) Copy() ExpressionInterface {
	ret := &SliceInline{
		ExpressionBase: e.copyBase(),
		Offset:         e.Offset,
		size:           e.size,
		SliceEnd:       e.SliceEnd,
		NoSlice:        e.NoSlice,
		oneTrit:        e.oneTrit,
	}
	return ret
}

func (e *SliceInline) Size() int {
	if e == nil {
		return 0
	}
	return e.size
}

func (e *SliceInline) Eval(frame *EvalFrame, result Trits) bool {
	var resTmp Trits
	if e.NoSlice {
		return e.subExpr[0].Eval(frame, result)
	}
	resTmp = make(Trits, e.subExpr[0].Size(), e.subExpr[0].Size())

	if e.subExpr[0].Eval(frame, resTmp) {
		return true
	}

	copy(result, resTmp[e.Offset:e.SliceEnd])
	return false
}
