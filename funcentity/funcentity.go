package funcentity

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/dispatcher"
)

type FunctionEntity struct {
	BaseEntity
	funDef     FuncDefInterface
	argbuf     Trits
	expression ExpressionInterface
}

func NewFunctionEntity(funDef FuncDefInterface) *FunctionEntity {
	effectCallback := func(args Trits) Trits {
		return callFundef(funDef, args)
	}
	argbuf := make(Trits, funDef.ArgSize(), funDef.ArgSize())
	expr, err := funDef.NewExpressionWithArgs(argbuf)
	if err != nil {
		panic(err)
	}
	return &FunctionEntity{
		BaseEntity: *NewBaseEntity(funDef.GetName(), funDef.ArgSize(), funDef.Size(), effectCallback),
		funDef:     funDef,
		argbuf:     argbuf,
		expression: expr,
	}
}

func callFundef(funDef FuncDefInterface, args Trits) Trits {
	return nil
}
