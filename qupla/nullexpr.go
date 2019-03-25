package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
)

type NullExpr struct {
	ExpressionBase
	size int64
}

func NewNullExpr(size int64) *NullExpr {
	return &NullExpr{
		ExpressionBase: NewExpressionBase(""),
		size:           size,
	}
}

func IsNullExpr(e interface{}) bool {
	_, ok := e.(*NullExpr)
	return ok
}

func (e *NullExpr) Size() int64 {
	return e.size
}

func (e *NullExpr) Eval(_ *EvalFrame, _ Trits) bool {
	return true
}

func (e *NullExpr) SetSize(size int64) {
	e.size = size
}
