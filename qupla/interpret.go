package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/cfg"
	"github.com/lunfardo314/goq/utils"
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
	HasState() bool
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

func (frame *EvalFrame) getCallTrace() []uint8 {
	ret := make([]uint8, 0, 10)
	f := frame
	for ; f != nil; f = f.prev {
		ret = append(ret, f.context.callIndex)
	}
	return ret
}

func (vi *VarInfo) Eval(frame *EvalFrame) (Trits, bool) {
	result := frame.buffer[vi.Offset:vi.SliceEnd]
	null := false
	cached := false

	switch result[0] {
	case evaluatedToNull:
		null = true
		cached = true

	case notEvaluated:
		if vi.IsParam {
			// evaluated in the context of previous call
			if frame.context.subexpr[vi.Idx].Eval(frame.prev, result) {
				result[0] = evaluatedToNull
				null = true
			}
		} else {
			if vi.IsState {
				// for state variables (latches) we return value, retrieved from the key/value storage
				// at the module level.
				// the key is frame.getCallTrace(). It return all 0 f not present
				// but calculated value stays in the buffer
				result = frame.context.FuncDef.StateHashMap.getValue(frame.getCallTrace(), len(result))
			} else {
				if vi.Assign.Eval(frame, result) {
					result[0] = evaluatedToNull
					null = true
				}
			}
		}

	default: // evaluated, not null (must be valid trit, not checking)
		cached = true
	}

	if frame.context.FuncDef.traceLevel > 1 {
		var s string
		if cached {
			s = "cached value "
		} else {
			s = "evaluated value "
		}
		if null {
			s += "null"
		} else {
			bi, _ := utils.TritsToBigInt(result)
			s += fmt.Sprintf("%v, '%v'", bi, utils.TritsToString(result))
		}
		if vi.IsState {
			s += fmt.Sprintf(" (state with call trace '%v)' ",
				frame.getCallTrace())
		}
		Logf(frame.context.FuncDef.traceLevel, "trace var %v.%v: %v",
			frame.context.FuncDef.Name, vi.Name, s)
	}
	return result, null
}

func (frame *EvalFrame) SaveStateVariables() {
	if frame == nil || !frame.context.FuncDef.HasStateVariables {
		return
	}
	Logf(7, "SaveStateVariables for '%v'", frame.context.FuncDef.Name)
	var val Trits
	for _, vi := range frame.context.FuncDef.LocalVars {
		if !vi.IsState {
			continue
		}
		val = make(Trits, vi.Assign.Size(), vi.Assign.Size())
		if !vi.Assign.Eval(frame, val) {
			frame.context.FuncDef.StateHashMap.storeValue(frame.getCallTrace(), val)
		}
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
