package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
)

type ValueExpr struct {
	ExpressionBase
	TritValue Trits
}

func NewValueExpr(t Trits) *ValueExpr {
	return &ValueExpr{
		TritValue: t,
	}
}

func (e *ValueExpr) InlineCopy(funExpr *FunctionExpr) ExpressionInterface {
	return &ValueExpr{
		ExpressionBase: e.inlineCopyBase(nil),
		TritValue:      e.TritValue,
	}
}

func (e *ValueExpr) Size() int {
	if e == nil {
		return 0
	}
	return len(e.TritValue)
}

func (e *ValueExpr) Eval(_ *EvalFrame, result Trits) bool {
	if e.TritValue == nil {
		return true
	}
	copy(result, e.TritValue)
	return false
}
