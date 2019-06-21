package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/abra"
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

func (e *SizeofExpr) Copy() ExpressionInterface {
	return &SizeofExpr{
		ExpressionBase: e.copyBase(),
		Value:          e.Value,
		TritValue:      e.TritValue,
	}
}

func (e *SizeofExpr) Size() int {
	if e == nil {
		return 0
	}
	return len(e.TritValue)
}

func (e *SizeofExpr) Eval(_ *EvalFrame, result Trits) bool {
	if e.TritValue == nil {
		return true
	}
	copy(result, e.TritValue)
	return false
}

func (e *SizeofExpr) GetAbraSite(branch *abra.Branch, codeUnit *abra.CodeUnit) *abra.Site {
	return branch.GetTritConstSite(codeUnit, e.TritValue)
}
