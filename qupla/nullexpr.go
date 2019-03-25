package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
)

type NullExpr struct {
	ExpressionBase
	size int
}

func NewNullExpr(size int) *NullExpr {
	return &NullExpr{
		ExpressionBase: NewExpressionBase(""),
		size:           size,
	}
}

func IsNullExpr(e interface{}) bool {
	_, ok := e.(*NullExpr)
	return ok
}

func (e *NullExpr) Size() int {
	return e.size
}

func (e *NullExpr) Eval(_ *EvalFrame, _ Trits) bool {
	return true
}

func (e *NullExpr) SetSize(size int) {
	e.size = size
}
