package funcentity

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/dispatcher"
)

type FunctionEntity struct {
	BaseEntity
	funDef FuncDefInterface
}

func NewFunctionEntity(funDef FuncDefInterface) *FunctionEntity {
	effectCallback := func(args Trits) Trits {
		return callFundef(funDef, args)
	}
	ret := &FunctionEntity{
		BaseEntity: *NewBaseEntity(funDef.GetName(), funDef.ArgSize(), funDef.Size(), effectCallback),
		funDef:     funDef,
	}
	return ret
}

func callFundef(funDef FuncDefInterface, args Trits) Trits {
	return nil
}
