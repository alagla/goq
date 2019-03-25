package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
)

type FieldExpr struct {
	Offset int64
	Size   int64
}
type TypeExpr struct {
	ExpressionBase
	size   int64
	Fields []FieldExpr
}

func NewQuplaTypeExpr(src string, size int64) *TypeExpr {
	return &TypeExpr{
		ExpressionBase: NewExpressionBase(src),
		size:           size,
		Fields:         make([]FieldExpr, 0, 5),
	}
}

func (e *TypeExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.size
}

func (e *TypeExpr) Eval(frame *EvalFrame, result Trits) bool {
	for idx, subExpr := range e.subexpr {
		if subExpr.Eval(frame, result[e.Fields[idx].Offset:e.Fields[idx].Offset+e.Fields[idx].Size]) {
			return true
		}
	}
	return false
}
