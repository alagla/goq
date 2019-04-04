package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/cfg"
)

type FunctionExpr struct {
	ExpressionBase
	FuncDef   *Function
	callIndex uint8
}

func NewFunctionExpr(src string, funcDef *Function, callIndex uint8) *FunctionExpr {
	return &FunctionExpr{
		ExpressionBase: NewExpressionBase(src),
		FuncDef:        funcDef,
		callIndex:      callIndex,
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
	//return e.FuncDef.RetExpr.Eval(&newFrame, result) // - avoid unnecessary call
	null := e.FuncDef.Eval(&newFrame, result)
	if !null {
		newFrame.SaveStateVariables()
	}
	return null
}

func (e *FunctionExpr) HasState() bool {
	return e.FuncDef.hasState || e.hasStateSubexpr()
}

func (e *FunctionExpr) Inline() ExpressionInterface {
	if !e.FuncDef.ZeroInternalSites() || e.FuncDef.isRecursive {
		return e
	}
	e.FuncDef.module.IncStat("numInlined")
	cfg.Logf(5, "inlined '%v' in %v", e.GetSource(), e.FuncDef.Name)
	ret := e.FuncDef.RetExpr.InlineCopy(e)
	return ret
}

func (e *FunctionExpr) InlineCopy(funExpr *FunctionExpr) ExpressionInterface {
	ret := &FunctionExpr{
		ExpressionBase: e.inlineCopyBase(funExpr),
		FuncDef:        e.FuncDef,
		callIndex:      e.callIndex,
	}
	return ret.Inline()
}
