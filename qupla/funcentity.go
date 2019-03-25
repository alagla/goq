package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/supervisor"
)

type functionEntityCore struct {
	functionCalc *FunctionCalculator
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
