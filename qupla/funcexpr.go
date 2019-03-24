package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/abstract"
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

func (e *FunctionExpr) References(funName string) bool {
	if e.FuncDef.Name == funName {
		return true
	}
	return e.ReferencesSubExprs(funName)
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

func (e *FunctionExpr) Eval(proc abstract.ProcessorInterface, result Trits) bool {
	return proc.Eval(e.FuncDef.RetExpr, result)
}

func (e *FunctionExpr) newEvalFrame(prev *EvalFrame) EvalFrame {
	size := int(e.FuncDef.BufLen)
	if prev != nil {
		return EvalFrame{
			buffer:   prev.buffer.growTo(prev.offset + prev.size + size),
			offset:   prev.offset + prev.size,
			size:     size,
			scope:    e,
			valueTag: make([]uint8, e.FuncDef.NumParams, e.FuncDef.NumParams),
		}
	}
	return EvalFrame{
		buffer: newGrowingBuffer(size),
		offset: 0,
		size:   size,
	}
}
