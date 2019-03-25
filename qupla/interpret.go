package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
)

type VarInfo struct {
	Name     string
	Analyzed bool
	Idx      int64
	Offset   int64
	Size     int64
	IsState  bool
	IsParam  bool
	Assign   ExpressionInterface
}
type EvalFrame struct {
	prev     *EvalFrame
	buffer   *growingBuffer
	offset   int
	size     int
	scope    *FunctionExpr
	valueTag []uint8 // 0x01 bit mean 'evaluated', 0x02 bit means is not null TODO optimize
}

type ExpressionInterface interface {
	GetSource() string
	Size() int64
	Eval(*EvalFrame, Trits) bool
	References(string) bool
}

type growingBuffer struct {
	arr Trits
}

const segmentSize = 1024

func newEvalFrame(expr *FunctionExpr, prev *EvalFrame) EvalFrame {
	size := int(expr.FuncDef.BufLen)
	var buf *growingBuffer
	offset := 0
	if prev == nil {
		buf = newGrowingBuffer(size)
	} else {
		buf = prev.buffer.growTo(prev.offset + prev.size + size)
		offset = prev.offset + prev.size
	}
	numVars := len(expr.FuncDef.LocalVars)
	return EvalFrame{
		prev:     prev,
		buffer:   buf,
		offset:   offset,
		size:     size,
		scope:    expr,
		valueTag: make([]uint8, numVars, numVars),
	}
}

func newGrowingBuffer(size int) *growingBuffer {
	alloc := (size/segmentSize + 1) * segmentSize
	return &growingBuffer{
		arr: make(Trits, alloc, alloc),
	}
}

func (b *growingBuffer) growTo(size int) *growingBuffer {
	if size <= len(b.arr) {
		return b // no need to grow
	}
	alloc := (size/segmentSize + 1) * segmentSize
	newArr := make(Trits, alloc, alloc)
	copy(newArr, b.arr)
	b.arr = newArr
	return b
}

func (frame *EvalFrame) EvalVar(idx int) (Trits, bool) {
	switch frame.valueTag[idx] & 0x03 {
	case 0x01: // evaluated, null
		logf(10, "evalVar evaluated NULL idx = %v in = '%v'", idx, frame.scope.FuncDef.Name)
		return nil, true

	case 0x03: // evaluated, not null
		logf(10, "evalVar evaluated NOT NULL idx = %v in '%v'", idx, frame.scope.FuncDef.Name)
		vi, err := frame.scope.FuncDef.VarByIdx(int64(idx))
		if err != nil {
			panic(err)
		}
		return frame.Slice(int(vi.Offset), int(vi.Size)), false

	case 0x00: // not evaluated
		logf(10, "evalVar NOT evaluated, idx = %v in '%v'", idx, frame.scope.FuncDef.Name)
		vi, err := frame.scope.FuncDef.VarByIdx(int64(idx))
		if err != nil {
			panic(err)
		}
		res := frame.Slice(int(vi.Offset), int(vi.Size))
		if vi.IsParam {
			if frame.scope.subexpr[vi.Idx].Eval(frame.prev, res) {
				frame.valueTag[idx] = 0x01 // evaluated null
				return nil, true
			}
		} else {
			if vi.Assign.Eval(frame, res) {
				frame.valueTag[idx] = 0x01 // evaluated null
				return nil, true
			}
		}
		frame.valueTag[idx] = 0x03 // evaluated not null
		return frame.Slice(int(vi.Offset), int(vi.Size)), false

	default:
		panic("wrong var tag")
	}
}

func (frame *EvalFrame) Slice(offset, size int) Trits {
	if frame == nil {
		panic("attempt to take slice from nil frame")
	}
	return frame.buffer.arr[frame.offset+offset : frame.offset+offset+size]
}

func MatchSizes(e1, e2 ExpressionInterface) error {
	s1 := e1.Size()
	s2 := e2.Size()

	if s1 != s2 {
		return fmt.Errorf("sizes doesn't match: %v != %v", s1, s2)
	}
	return nil
}

func RequireSize(e ExpressionInterface, size int64) error {
	s := e.Size()

	if s != size {
		return fmt.Errorf("sizes doesn't match: required %v != %v", size, s)
	}
	return nil
}
