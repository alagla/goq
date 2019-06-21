package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/abra"
)

type NullExpr struct {
	ExpressionBase
	size int
}

func (e *NullExpr) GetAbraSite(branch *abra.Branch, codeUnit *abra.CodeUnit) *abra.Site {
	panic("implement me")
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

func (e *NullExpr) Copy() ExpressionInterface {
	return &NullExpr{
		ExpressionBase: e.copyBase(),
		size:           e.size,
	}
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
