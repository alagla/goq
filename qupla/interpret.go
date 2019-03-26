package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
)

type VarInfo struct {
	Name     string
	Analyzed bool
	Idx      int
	Offset   int
	Size     int
	SliceEnd int
	IsState  bool
	IsParam  bool
	Assign   ExpressionInterface
}

type EvalFrame struct {
	prev    *EvalFrame
	buffer  Trits
	context *FunctionExpr
}

type ExpressionInterface interface {
	Size() int
	Eval(*EvalFrame, Trits) bool
	References(string) bool
}

const (
	notEvaluated    = int8(100)
	evaluatedToNull = int8(101)
)

func newEvalFrame(expr *FunctionExpr, prev *EvalFrame) EvalFrame {
	ret := EvalFrame{
		prev:    prev,
		buffer:  make(Trits, expr.FuncDef.BufLen, expr.FuncDef.BufLen),
		context: expr,
	}
	for _, vi := range expr.FuncDef.LocalVars {
		ret.buffer[vi.Offset] = notEvaluated
	}
	return ret
}

func (vi *VarInfo) Eval(frame *EvalFrame) (Trits, bool) {
	result := frame.buffer[vi.Offset:vi.SliceEnd]
	switch result[0] {
	case evaluatedToNull:
		//logf(10, "evalVar evaluated NULL idx = %v in = '%v'", idx, frame.context.FuncDef.Name)
		return nil, true

	case notEvaluated:
		//logf(10, "evalVar NOT evaluated, idx = %v in '%v'", idx, frame.context.FuncDef.Name)
		if vi.IsParam {
			// evaluated in the context of previous call
			if frame.context.subexpr[vi.Idx].Eval(frame.prev, result) {
				result[0] = evaluatedToNull
				return nil, true
			}
		} else {
			if vi.Assign.Eval(frame, result) {
				result[0] = evaluatedToNull
				return nil, true
			}
		}
		return result, false

	default: // evaluated, not null (must be valid trit, not checking)
		//logf(10, "evalVar evaluated NOT NULL idx = %v in '%v'", idx, frame.context.FuncDef.Name)
		return result, false
	}
}

func MatchSizes(e1, e2 ExpressionInterface) error {
	s1 := e1.Size()
	s2 := e2.Size()

	if s1 != s2 {
		return fmt.Errorf("sizes doesn't match: %v != %v", s1, s2)
	}
	return nil
}

func RequireSize(e ExpressionInterface, size int) error {
	s := e.Size()

	if s != size {
		return fmt.Errorf("sizes doesn't match: required %v != %v", size, s)
	}
	return nil
}

//type growingBuffer struct {
//	arr Trits
//}
//
//const segmentSize = 1024
//
//func newGrowingBuffer(size int) *growingBuffer {
//	alloc := (size/segmentSize + 1) * segmentSize
//	return &growingBuffer{
//		arr: make(Trits, alloc, alloc),
//	}
//}
//
//func (b *growingBuffer) growTo(size int) *growingBuffer {
//	if size <= len(b.arr) {
//		return b // no need to grow
//	}
//	alloc := (size/segmentSize + 1) * segmentSize
//	newArr := make(Trits, alloc, alloc)
//	copy(newArr, b.arr)
//	b.arr = newArr
//	return b
//}
