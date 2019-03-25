package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
)

type SizeofExpr struct {
	ExpressionBase
	Value     int64
	TritValue Trits
}

func NewQuplaSizeofExpr(value int64, tritValue Trits) *SizeofExpr {
	return &SizeofExpr{
		Value:     value,
		TritValue: tritValue,
	}
}

func (e *SizeofExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return int64(len(e.TritValue))
}

func (e *SizeofExpr) Eval(_ ProcessorInterface, result Trits) bool {
	if e.TritValue == nil {
		return true
	}
	copy(result, e.TritValue)
	return false
}
