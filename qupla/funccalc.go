package qupla

import . "github.com/iotaledger/iota.go/trinary"

type FunctionCalculator struct {
	function *Function
	frame    *EvalFrame
}

func NewFunctionCalculator(function *Function) *FunctionCalculator {
	mockExpr := function.NewFuncExpressionTemplate()
	frame := newEvalFrame(mockExpr, nil)
	return &FunctionCalculator{
		function: function,
		frame:    &frame,
	}
}

func (fc *FunctionCalculator) evalWithArgs(args Trits, result Trits) bool {
	copy(fc.frame.buffer, args)
	return fc.function.Eval(fc.frame, result)
}
