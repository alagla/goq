package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
)

type FunctionExpr struct {
	ExpressionBase
	source    string
	FuncDef   *Function
	callIndex uint8
}

func NewFunctionExpr(src string, funcDef *Function) *FunctionExpr {
	return &FunctionExpr{
		ExpressionBase: NewExpressionBase(src),
		FuncDef:        funcDef,
		callIndex:      funcDef.NextCallIndex(),
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
	//if frame != nil && strings.Contains(e.FuncDef.Name, "arcLeaf"){
	//	logf(0, "+++++++++++++++++++ KUKU")
	//}
	newFrame := newEvalFrame(e, frame)
	//return e.FuncDef.RetExpr.Eval(&newFrame, result) // - avoid unnecessary call
	null := e.FuncDef.Eval(&newFrame, result)
	if !null {
		frame.SaveStateVariables()
	}
	return null
}

func (e *FunctionExpr) HasState() bool {
	return e.FuncDef.hasState || e.hasStateSubexpr()
}
