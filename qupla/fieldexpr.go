package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
)

// TODO with fields expressions
type QuplaFieldExpr struct {
	ExpressionBase
	CondExpr ExpressionInterface
}

func (e *QuplaFieldExpr) Size() int {
	if e == nil {
		return 0
	}
	return e.CondExpr.Size()
}

func (e *QuplaFieldExpr) Eval(_ *EvalFrame, _ Trits) bool {
	return true
}

func (e *QuplaFieldExpr) InlineCopy(funExpr *FunctionExpr) ExpressionInterface {
	return &QuplaFieldExpr{
		ExpressionBase: e.inlineCopyBase(funExpr),
		CondExpr:       e.CondExpr.InlineCopy(funExpr),
	}
}
