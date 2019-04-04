package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
)

type SliceInline struct {
	ExpressionBase
	expr     ExpressionInterface
	offset   int
	size     int
	sliceEnd int
	noSlice  bool
	oneTrit  bool
}

func NewSliceInline(sliceExpr *SliceExpr, expr ExpressionInterface) *SliceInline {
	return &SliceInline{
		ExpressionBase: NewExpressionBase(sliceExpr.GetSource()),
		expr:           expr,
		offset:         sliceExpr.offset,
		size:           sliceExpr.size,
		sliceEnd:       sliceExpr.sliceEnd,
		noSlice:        sliceExpr.noSlice,
		oneTrit:        sliceExpr.oneTrit,
	}
}

func (e *SliceInline) InlineCopy(funExpr *FunctionExpr) ExpressionInterface {
	return &SliceInline{
		ExpressionBase: NewExpressionBase(e.GetSource()),
		expr:           e.expr.InlineCopy(funExpr),
		offset:         e.offset,
		size:           e.size,
		sliceEnd:       e.sliceEnd,
		noSlice:        e.noSlice,
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
	if e.noSlice {
		return e.expr.Eval(frame, result)
	}
	resTmp = make(Trits, e.expr.Size(), e.expr.Size())

	if e.expr.Eval(frame, resTmp) {
		return true
	}

	copy(result, resTmp[e.offset:e.sliceEnd])
	return false
}

func optimizeInlineSlicesExpr(expr ExpressionInterface) ExpressionInterface {
	inlineSlice, ok := expr.(*SliceInline)
	if !ok {
		subExpr := make([]ExpressionInterface, 0)
		for _, se := range expr.GetSubexpressions() {
			opt := optimizeInlineSlicesExpr(se)
			subExpr = append(subExpr, opt)
		}
		expr.SetSubexpressions(subExpr)
		return expr
	}
	if inlineSlice.noSlice {
		return inlineSlice.expr
	}
	valueExpr, ok := inlineSlice.expr.(*ValueExpr)
	if !ok {
		return inlineSlice
	}
	return NewValueExpr(valueExpr.TritValue[inlineSlice.offset:inlineSlice.sliceEnd])
}
