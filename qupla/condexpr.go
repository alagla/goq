package qupla

import (
	. "fmt"
	. "github.com/iotaledger/iota.go/trinary"
)

type CondExpr struct {
	ExpressionBase
}

func (e *CondExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.subexpr[1].Size()
}

func (e *CondExpr) Eval(frame *EvalFrame, result Trits) bool {
	var buf [1]int8
	null := e.subexpr[0].Eval(frame, buf[:])
	if null {
		return true
	}
	// bool is 0/1
	switch buf[0] {
	case 1:
		return e.subexpr[1].Eval(frame, result)
	case 0:
		return e.subexpr[2].Eval(frame, result)
	case -1:
		return true
	}
	panic(Sprintf("trit value in cond expr '%v'", e.source))
}
