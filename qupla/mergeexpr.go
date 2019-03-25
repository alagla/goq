package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
)

type MergeExpr struct {
	ExpressionBase
}

func (e *MergeExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.subexpr[0].Size()
}

func (e *MergeExpr) Eval(proc ProcessorInterface, result Trits) bool {
	null := proc.Eval(e.subexpr[0], result)
	if null {
		return proc.Eval(e.subexpr[1], result)
	}
	return false
}
