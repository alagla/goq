package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
)

type ConcatExpr struct {
	ExpressionBase
}

func (e *ConcatExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.subexpr[0].Size() + e.subexpr[1].Size()
}

func (e *ConcatExpr) Eval(proc ProcessorInterface, result Trits) bool {
	null := proc.Eval(e.subexpr[0], result)
	if null {
		return true
	}
	return proc.Eval(e.subexpr[1], result[e.subexpr[0].Size():])
}
