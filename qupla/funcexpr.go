package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
)

type QuplaFuncExpr struct {
	QuplaExprBase
	source  string
	FuncDef *QuplaFuncDef
}

func NewQuplaFuncExpr(src string, funcDef *QuplaFuncDef) *QuplaFuncExpr {
	return &QuplaFuncExpr{
		QuplaExprBase: NewQuplaExprBase(src),
		FuncDef:       funcDef,
	}
}

func (e *QuplaFuncExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.FuncDef.Size()
}

func (e *QuplaFuncExpr) NewCallFrame(parent *CallFrame) *CallFrame {
	numVars := len(e.FuncDef.LocalVars)
	return &CallFrame{
		context:  e,
		parent:   parent,
		buffer:   make(Trits, e.FuncDef.BufLen, e.FuncDef.BufLen),
		valueTag: make([]uint8, numVars, numVars),
	}
}

func (e *QuplaFuncExpr) Eval(proc ProcessorInterface, result Trits) bool {
	return proc.Eval(e.FuncDef.RetExpr, result)
}

func (e *QuplaFuncExpr) References(funName string) bool {
	if e.FuncDef.Name == funName {
		return true
	}
	return e.ReferencesSubExprs(funName)
}
