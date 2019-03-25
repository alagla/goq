package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
)

type FunctionExpr struct {
	ExpressionBase
	source  string
	FuncDef *Function
}

func NewFunctionExpr(src string, funcDef *Function) *FunctionExpr {
	return &FunctionExpr{
		ExpressionBase: NewExpressionBase(src),
		FuncDef:        funcDef,
	}
}

func (e *FunctionExpr) Size() int {
	return e.FuncDef.Size()
}

func (e *FunctionExpr) References(funName string) bool {
	if e.FuncDef.Name == funName {
		return true
	}
	return e.ReferencesSubExprs(funName)
}

func (e *FunctionExpr) Eval(frame *EvalFrame, result Trits) bool {
	newFrame := newEvalFrame(e, frame)
	return e.FuncDef.RetExpr.Eval(&newFrame, result) // - avoid unnecessary call
	//return e.FuncDef.Eval(&newFrame, result)
}
