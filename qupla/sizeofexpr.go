package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
)

type QuplaSizeofExpr struct {
	QuplaExprBase
	Value     int64
	TritValue Trits
}

func NewQuplaSizeofExpr(value int64, tritValue Trits) *QuplaSizeofExpr {
	return &QuplaSizeofExpr{
		Value:     value,
		TritValue: tritValue,
	}
}

func (e *QuplaSizeofExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return int64(len(e.TritValue))
}

func (e *QuplaSizeofExpr) Eval(_ ProcessorInterface, result Trits) bool {
	if e.TritValue == nil {
		return true
	}
	copy(result, e.TritValue)
	return false
}
