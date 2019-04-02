package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
)

type ConcatExpr struct {
	ExpressionBase
}

func (e *ConcatExpr) Size() int {
	if e == nil {
		return 0
	}
	return e.subExpr[0].Size() + e.subExpr[1].Size()
}

func (e *ConcatExpr) Eval(frame *EvalFrame, result Trits) bool {
	null := e.subExpr[0].Eval(frame, result)
	if null {
		return true
	}
	return e.subExpr[1].Eval(frame, result[e.subExpr[0].Size():])
}

func (e *ConcatExpr) InlineCopy(funExpr *FunctionExpr) ExpressionInterface {
	return &ConcatExpr{
		ExpressionBase: e.inlineCopyBase(funExpr),
	}
}
