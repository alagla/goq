package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
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

func (e *ConcatExpr) InlineCopy(funExpr *FunctionExpr) ExpressionInterface {
	return &ConcatExpr{
		ExpressionBase: e.inlineCopyBase(funExpr),
		size:           e.size,
		offset:         e.offset,
		endSlice:       e.endSlice,
	}
}
