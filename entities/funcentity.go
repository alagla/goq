package entities

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/dispatcher"
)

type FunctionEntity struct {
	BaseEntity
}

func NewFunctionEntity(funDef FuncDefInterface) *FunctionEntity {
	effectCallback := func(args Trits) Trits {
		expr, err := funDef.NewExpressionWithArgs(args)
		if err != nil {
			panic(err)
		}
		var proc ProcessorInterface // todo
		res := make(Trits, funDef.Size(), funDef.Size())
		null := proc.Eval(expr, res)
		if null {
			return nil
		}
		return res
	}
	return &FunctionEntity{
		BaseEntity: *NewBaseEntity(funDef.GetName(), funDef.ArgSize(), funDef.Size(), effectCallback),
	}
}
