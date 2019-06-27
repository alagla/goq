package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/abra"
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

func (e *FunctionExpr) Copy() ExpressionInterface {
	return &FunctionExpr{
		ExpressionBase: e.copyBase(),
		FuncDef:        e.FuncDef,
		callIndex:      e.callIndex,
	}
}

func (e *FunctionExpr) GetAbraSite(branch *abra.Branch, codeUnit *abra.CodeUnit, lookupName string) *abra.Site {
	branchBlock := e.FuncDef.GetAbraBranchBlock(codeUnit)
	inputs := make([]*abra.Site, e.FuncDef.NumParams)
	for i, se := range e.GetSubexpressions() {
		inputs[i] = se.GetAbraSite(branch, codeUnit, "")
	}
	ret := abra.NewKnot(branchBlock, inputs...).NewSite(e.Size())
	ret.SetLookupName(lookupName)
	return branch.AddOrUpdateSite(ret)
}
