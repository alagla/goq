package entities

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/dispatcher"
)

type functionCallable struct {
	funcDef FuncDefInterface
	proc    ProcessorInterface
}

func (fc *functionCallable) Call(args Trits, res Trits) bool {
	expr, err := fc.funcDef.NewExpressionWithArgs(args)
	if err != nil {
		panic(err)
	}
	return fc.proc.Eval(expr, res)
}

func NewFunctionEntity(funcDef FuncDefInterface, proc ProcessorInterface) *BaseEntity {
	return NewBaseEntity(funcDef.GetName(), funcDef.ArgSize(), funcDef.Size(),
		&functionCallable{funcDef, proc})
}