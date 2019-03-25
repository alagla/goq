package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
)

type ValueExpr struct {
	ExpressionBase
	TritValue Trits
}

func NewQuplaValueExpr(t Trits) *ValueExpr {
	return &ValueExpr{
		TritValue: t,
	}
}
func (e *ValueExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return int64(len(e.TritValue))
}

func (e *ValueExpr) Eval(_ ProcessorInterface, result Trits) bool {
	if e.TritValue == nil {
		return true
	}
	copy(result, e.TritValue)
	return false
}
