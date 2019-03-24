package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
)

type ExpressionInterface1 interface {
	Eval(frame *EvalFrame, result Trits) bool
}

type growingBuffer struct {
	arr Trits
}

const segmentSize = 1024

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

type EvalFrame struct {
	buffer   *growingBuffer
	offset   int
	size     int
	scope    *FunctionExpr
	valueTag []uint8 // 0x01 bit mean 'evaluated', 0x02 bit means is not null TODO optimize
}

func (frame *EvalFrame) EvalVar(idx int) (Trits, bool) {
	switch frame.valueTag[idx] & 0x03 {
	case 1: // evaluated, null
		return nil, true

	case 3: // evaluated, not null
		vi := frame.scope.FuncDef.VarByIdx(int64(idx))
		return frame.Slice(int(vi.Offset), int(vi.Size)), false

	case 0: // not evaluated
		vi := frame.scope.FuncDef.VarByIdx(int64(idx))
		//res := frame.Slice(int(vi.Offset), int(vi.Size))
		if vi.IsParam {
			//if frame.scope.subexpr[vi.Idx].Eval(frame, res){
			//	return nil, true
			//}
		} else {
			//if vi.Assign.Eval(frame, res){
			//	return nil, true
			//}
		}
		return frame.Slice(int(vi.Offset), int(vi.Size)), false

	default:
		panic("wrong var tag")
	}
}

func (frame *EvalFrame) Slice(offset, size int) Trits {
	return frame.buffer.arr[frame.offset+offset : frame.offset+offset+size]
}
