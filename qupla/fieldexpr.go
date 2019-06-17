package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/abra"
)

// TODO with fields expressions
type QuplaFieldExpr struct {
	ExpressionBase
}

func (e *QuplaFieldExpr) GenAbraSite(branch *abra.Branch, codeUnit *abra.CodeUnit) *abra.Site {
	panic("implement me")
}

func NewFieldExpr(src string, condExpr ExpressionInterface) *QuplaFieldExpr {
	ret := &QuplaFieldExpr{
		ExpressionBase: NewExpressionBase(src),
	}
	ret.AppendSubExpr(condExpr)
	return ret
}

func (e *QuplaFieldExpr) Size() int {
	if e == nil {
		return 0
	}
	return e.subExpr[0].Size()
}

func (e *QuplaFieldExpr) Eval(_ *EvalFrame, _ Trits) bool {
	return true
}

func (e *QuplaFieldExpr) Copy() ExpressionInterface {
	return NewFieldExpr(e.source, e.subExpr[0])
}
