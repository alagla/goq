package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
)

type NullExpr struct {
	ExpressionBase
	size int64
}

func IsNullExpr(e interface{}) bool {
	_, ok := e.(*NullExpr)
	return ok
}

func (e *NullExpr) Size() int64 {
	return e.size
}

func (e *NullExpr) Eval(_ ProcessorInterface, _ Trits) bool {
	return true
}

func (e *NullExpr) SetSize(size int64) {
	e.size = size
}
