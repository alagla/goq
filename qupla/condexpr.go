package qupla

import (
	. "fmt"
	. "github.com/iotaledger/iota.go/trinary"
)

type CondExpr struct {
	ExpressionBase
}

func (e *CondExpr) Size() int {
	if e == nil {
		return 0
	}
	return e.subExpr[1].Size()
}

const (
	BOOL_TRUE  = 1
	BOOL_FALSE = -1
)

func (e *CondExpr) Eval(frame *EvalFrame, result Trits) bool {
	var buf [1]int8
	null := e.subExpr[0].Eval(frame, buf[:])
	if null {
		return true
	}
	// bool is 0/1
	switch buf[0] {
	case BOOL_TRUE:
		return e.subExpr[1].Eval(frame, result)
	case BOOL_FALSE:
		return e.subExpr[2].Eval(frame, result)
	default:
		return true
	}
	panic(Sprintf("trit value in cond Expr '%v'", e.source))
}

func (e *CondExpr) Copy() ExpressionInterface {
	return &CondExpr{
		ExpressionBase: e.copyBase(),
	}
}
