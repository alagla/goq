package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
)

type QuplaNullExpr struct {
	size int64
}

func IsNullExpr(e interface{}) bool {
	_, ok := e.(*QuplaNullExpr)
	return ok
}

func (e *QuplaNullExpr) Analyze(module *QuplaModule, scope *QuplaFuncDef) (ExpressionInterface, error) {
	return e, nil
}

func (e *QuplaNullExpr) Size() int64 {
	return e.size
}

func (e *QuplaNullExpr) Eval(_ ProcessorInterface, _ Trits) bool {
	return true
}

func (e *QuplaNullExpr) SetSize(size int64) {
	e.size = size
}
