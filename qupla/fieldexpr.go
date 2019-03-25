package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
)

// TODO with fields expressions
type QuplaFieldExpr struct {
	ExpressionBase
	CondExpr ExpressionInterface
}

func (e *QuplaFieldExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.CondExpr.Size()
}

func (e *QuplaFieldExpr) Eval(_ ProcessorInterface, _ Trits) bool {
	return true
}
