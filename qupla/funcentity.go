package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/supervisor"
)

type functionEntityCore struct {
	functionCalc *FunctionCalculator
}

type FunctionCalculator struct {
	function *Function
	frame    *EvalFrame
}

func NewFunctionEntity(sv *Supervisor, function *Function) (*Entity, error) {
	return sv.NewEntity(EntityOpts{
		Name:    function.Name,
		InSize:  function.ArgSize(),
		OutSize: function.Size(),
		Core:    &functionEntityCore{functionCalc: NewFunctionCalculator(function)},
	})
}

func (fc *functionEntityCore) Call(args Trits, result Trits) bool {
	return fc.functionCalc.evalWithArgs(args, result)
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
	// Go copy does exactly what Abra specs prescribes:
	// cut longer inout and pad with zero shorter than the input size
	copy(fc.frame.buffer, args)
	return fc.function.Eval(fc.frame, result)
}
