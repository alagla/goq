package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
)

type FunctionExpr struct {
	ExpressionBase
	source  string
	FuncDef *Function
}

func NewQuplaFuncExpr(src string, funcDef *Function) *FunctionExpr {
	return &FunctionExpr{
		ExpressionBase: NewExpressionBase(src),
		FuncDef:        funcDef,
	}
}

func (e *FunctionExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.FuncDef.Size()
}

func (e *FunctionExpr) NewCallFrame(parent *CallFrame) *CallFrame {
	numVars := len(e.FuncDef.LocalVars)
	return &CallFrame{
		context:  e,
		parent:   parent,
		buffer:   make(Trits, e.FuncDef.BufLen, e.FuncDef.BufLen),
		valueTag: make([]uint8, numVars, numVars),
	}
}

func (e *FunctionExpr) Eval(proc ProcessorInterface, result Trits) bool {
	return proc.Eval(e.FuncDef.RetExpr, result)
}

func (e *FunctionExpr) References(funName string) bool {
	if e.FuncDef.Name == funName {
		return true
	}
	return e.ReferencesSubExprs(funName)
}
