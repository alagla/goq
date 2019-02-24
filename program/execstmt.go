package program

import (
	"fmt"
	"github.com/iotaledger/iota.go/trinary"
)

type QuplaExecStmt struct {
	ExprWrap     *QuplaExpressionWrapper `yaml:"expr"`
	ExpectedWrap *QuplaExpressionWrapper `yaml:"expected,omitempty"`
	//---
	isTest       bool
	expr         ExpressionInterface
	exprExpected ExpressionInterface
}

func (ex *QuplaExecStmt) Execute() error {
	funcExpr, ok := ex.expr.(*QuplaFuncExpr)
	if !ok {
		return fmt.Errorf("must be call to function")
	}
	res := make(trinary.Trits, funcExpr.Size(), funcExpr.Size())
	null := funcExpr.Eval(nil, res)
	if null {
		debugf("result is null")
	} else {
		debugf("result is not null")
	}
	return nil
}
