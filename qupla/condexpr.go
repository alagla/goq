package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
)

type QuplaCondExpr struct {
	QuplaExprBase
}

func (e *QuplaCondExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.subexpr[1].Size()
}

func (e *QuplaCondExpr) Eval(proc ProcessorInterface, result Trits) bool {
	var buf [1]int8
	null := proc.Eval(e.subexpr[0], buf[:])
	if null {
		return true
	}
	// bool is 0/1
	switch buf[0] {
	case 1:
		return proc.Eval(e.subexpr[1], result)
	case 0:
		return proc.Eval(e.subexpr[2], result)
	case -1:
		return true
	}
	panic(fmt.Sprintf("trit value in cond expr '%v'", e.source))
}
