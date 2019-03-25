package qupla

import . "github.com/iotaledger/iota.go/trinary"

type FunctionCalculator struct {
	exprTemplate *FunctionExpr
	frame        *EvalFrame
}

func NewFunctionCalculator(function *Function) *FunctionCalculator {
	expr := function.NewFuncExpressionTemplate()
	frame := newEvalFrame(expr, nil)
	return &FunctionCalculator{
		exprTemplate: expr,
		frame:        &frame,
	}
}

func (fc *FunctionCalculator) evalWithArgs(args Trits, result Trits) bool {
	copy(fc.frame.buffer.arr, args)
	for i := range fc.frame.valueTag {
		fc.frame.valueTag[i] = 0x03 // evaluated not null
	}
	return fc.exprTemplate.FuncDef.Eval(fc.frame, result)
}
