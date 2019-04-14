package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
)

type SliceInline struct {
	ExpressionBase
	Expr     ExpressionInterface
	Offset   int
	size     int
	SliceEnd int
	NoSlice  bool
	oneTrit  bool
}

func NewSliceInline(sliceExpr *SliceExpr, expr ExpressionInterface) *SliceInline {
	return &SliceInline{
		ExpressionBase: NewExpressionBase(sliceExpr.GetSource()),
		Expr:           expr,
		Offset:         sliceExpr.offset,
		size:           sliceExpr.size,
		SliceEnd:       sliceExpr.sliceEnd,
		NoSlice:        sliceExpr.noSlice,
		oneTrit:        sliceExpr.oneTrit,
	}
}

func (e *SliceInline) InlineCopy(funExpr *FunctionExpr) ExpressionInterface {
	return &SliceInline{
		ExpressionBase: NewExpressionBase(e.GetSource()),
		Expr:           e.Expr.InlineCopy(funExpr),
		Offset:         e.Offset,
		size:           e.size,
		SliceEnd:       e.SliceEnd,
		NoSlice:        e.NoSlice,
	}
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
		return e.Expr.Eval(frame, result)
	}
	resTmp = make(Trits, e.Expr.Size(), e.Expr.Size())

	if e.Expr.Eval(frame, resTmp) {
		return true
	}

	copy(result, resTmp[e.Offset:e.SliceEnd])
	return false
}
