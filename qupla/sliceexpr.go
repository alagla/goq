package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
)

type QuplaSliceExpr struct {
	QuplaExprBase
	LocalVarIdx int64
	VarScope    *QuplaFuncDef
	offset      int64
	size        int64
}

func NewQuplaSliceExpr(src string, offset, size int64) *QuplaSliceExpr {
	return &QuplaSliceExpr{
		QuplaExprBase: NewQuplaExprBase(src),
		offset:        offset,
		size:          size,
	}
}

func (e *QuplaSliceExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.size
}

func (e *QuplaSliceExpr) Eval(proc ProcessorInterface, result Trits) bool {
	restmp, null := proc.EvalVar(e.LocalVarIdx)
	if null {
		return true
	}
	copy(result, restmp[e.offset:e.offset+e.size])
	return false
}
