package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
)

type QuplaValueExpr struct {
	QuplaExprBase
	TritValue Trits
}

func NewQuplaValueExpr(t Trits) *QuplaValueExpr {
	return &QuplaValueExpr{
		TritValue: t,
	}
}
func (e *QuplaValueExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return int64(len(e.TritValue))
}

func (e *QuplaValueExpr) Eval(_ ProcessorInterface, result Trits) bool {
	if e.TritValue == nil {
		return true
	}
	copy(result, e.TritValue)
	return false
}
