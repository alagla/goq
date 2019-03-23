package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
)

type QuplaFieldExpr struct {
	QuplaExprBase
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
