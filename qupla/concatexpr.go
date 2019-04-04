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
	//var ret int
	//for _, se := range e.subExpr{
	//	ret += se.Size()
	//}
	//return ret
}

func (e *ConcatExpr) Eval(frame *EvalFrame, result Trits) bool {
	for i, se := range e.subExpr {
		if se.Eval(frame, result[e.offset[i]:e.endSlice[i]]) {
			return true
		}
	}
	return false

	//null := e.subExpr[0].Eval(frame, result)
	//if null {
	//	return true
	//}
	//return e.subExpr[1].Eval(frame, result[e.subExpr[0].Size():])
}

func (e *ConcatExpr) InlineCopy(funExpr *FunctionExpr) ExpressionInterface {
	return &ConcatExpr{
		ExpressionBase: e.inlineCopyBase(funExpr),
		size:           e.size,
		offset:         e.offset,
		endSlice:       e.endSlice,
	}
}

func optimizeConcatExpr(expr ExpressionInterface) ExpressionInterface {
	_, ok := expr.(*ConcatExpr)
	subExpr := make([]ExpressionInterface, 0)
	if !ok {
		for _, se := range expr.GetSubexpressions() {
			subExpr = append(subExpr, optimizeConcatExpr(se))
		}
		expr.SetSubexpressions(subExpr)
		return expr
	}
	for _, se := range expr.GetSubexpressions() {
		oe := optimizeConcatExpr(se)
		if ce, ok := oe.(*ConcatExpr); ok {
			for _, e := range ce.subExpr {
				subExpr = append(subExpr, e)
			}
		} else {
			subExpr = append(subExpr, oe)
		}
	}
	return NewConcatExpression("optimized", subExpr)
}
