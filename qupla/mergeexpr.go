package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
)

type QuplaMergeExpr struct {
	QuplaExprBase
}

func (e *QuplaMergeExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.subexpr[0].Size()
}

func (e *QuplaMergeExpr) Eval(proc ProcessorInterface, result Trits) bool {
	null := proc.Eval(e.subexpr[0], result)
	if null {
		return proc.Eval(e.subexpr[1], result)
	}
	return false
}
