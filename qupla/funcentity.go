package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/supervisor"
)

type functionEntityCore struct {
	funcDef *QuplaFuncDef
	proc    ProcessorInterface
}

func (fc *functionEntityCore) Call(args Trits, res Trits) bool {
	expr, err := fc.funcDef.NewExpressionWithArgs(args)
	if err != nil {
		panic(err)
	}
	return fc.proc.Eval(expr, res)
}

func NewFunctionEntity(disp *Supervisor, funcDef *QuplaFuncDef, proc ProcessorInterface) (*Entity, error) {
	return disp.NewEntity(EntityOpts{
		Name:    funcDef.Name,
		InSize:  funcDef.ArgSize(),
		OutSize: funcDef.Size(),
		Core:    &functionEntityCore{funcDef, proc},
	})
}
