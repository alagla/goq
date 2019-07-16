package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/abra"
	cabra "github.com/lunfardo314/goq/abra/construct"
)

type ConcatExpr struct {
	ExpressionBase
	size     int
	offset   []int
	endSlice []int
}

func NewConcatExpression(src string, args []ExpressionInterface) *ConcatExpr {
	ret := &ConcatExpr{
		ExpressionBase: NewExpressionBase(src),
		offset:         make([]int, 0),
		endSlice:       make([]int, 0),
	}
	offset := 0
	for _, a := range args {
		ret.AppendSubExpr(a)
		ret.size += a.Size()
		ret.offset = append(ret.offset, offset)
		ret.endSlice = append(ret.endSlice, offset+a.Size())
		offset += a.Size()
	}
	return ret
}

func (e *ConcatExpr) Size() int {
	return e.size
}

func (e *ConcatExpr) Eval(frame *EvalFrame, result Trits) bool {
	for i, se := range e.subExpr {
		if se.Eval(frame, result[e.offset[i]:e.endSlice[i]]) {
			return true
		}
	}
	return false
}

func (e *ConcatExpr) Copy() ExpressionInterface {
	return &ConcatExpr{
		ExpressionBase: e.copyBase(),
		size:           e.size,
		offset:         e.offset,
		endSlice:       e.endSlice,
	}
}

func (e *ConcatExpr) GetAbraSite(branch *abra.Branch, codeUnit *abra.CodeUnit, lookupName string) *abra.Site {
	inputs := make([]*abra.Site, 0, len(e.subExpr))
	for _, se := range e.subExpr {
		s := se.GetAbraSite(branch, codeUnit, "")
		inputs = append(inputs, s)
	}
	concatBlock := cabra.GetConcatBlockForSize(codeUnit, e.Size())
	ret := cabra.NewKnotSite(e.Size(), lookupName, concatBlock, inputs...)
	return cabra.AddOrUpdateSite(branch, ret)
}
