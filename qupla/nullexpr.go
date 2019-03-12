package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/quplayaml/go"
)

type QuplaNullExpr struct {
	QuplaExprBase
	size int64
}

func IsNullExpr(e interface{}) bool {
	_, ok := e.(*QuplaNullExpr)
	return ok
}

func AnalyzeNullExpr(_ *QuplaNullExprYAML, module ModuleInterface, _ FuncDefInterface) (*QuplaNullExpr, error) {
	module.IncStat("nullExpr")
	return &QuplaNullExpr{}, nil
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
