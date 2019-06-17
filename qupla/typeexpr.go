package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/abra"
)

type FieldExpr struct {
	Offset int
	Size   int
}
type TypeExpr struct {
	ExpressionBase
	size   int
	Fields []FieldExpr
}

func (e *TypeExpr) GenAbraSite(branch *abra.Branch, codeUnit *abra.CodeUnit) *abra.Site {
	panic("implement me")
}

func NewQuplaTypeExpr(src string, size int) *TypeExpr {
	return &TypeExpr{
		ExpressionBase: NewExpressionBase(src),
		size:           size,
		Fields:         make([]FieldExpr, 0, 5),
	}
}

func (e *TypeExpr) Copy() ExpressionInterface {
	return &TypeExpr{
		ExpressionBase: e.copyBase(),
		size:           e.size,
		Fields:         e.Fields,
	}
}

func (e *TypeExpr) Size() int {
	if e == nil {
		return 0
	}
	return e.size
}

func (e *TypeExpr) Eval(frame *EvalFrame, result Trits) bool {
	for idx, subExpr := range e.subExpr {
		if subExpr.Eval(frame, result[e.Fields[idx].Offset:e.Fields[idx].Offset+e.Fields[idx].Size]) {
			return true
		}
	}
	return false
}
