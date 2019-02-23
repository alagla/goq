package program

import (
	"github.com/iotaledger/iota.go/trinary"
)

// TODO
type Processor struct {
	curOffset int64
	buffer    trinary.Trits
}

func NewProcessor(width int64) *Processor {
	return &Processor{
		buffer: make(trinary.Trits, width, width),
	}
}

func (proc *Processor) Slice(offset, size int64) trinary.Trits {
	return proc.buffer[proc.curOffset+offset : size]
}

func (proc *Processor) Eval(expr ExpressionInterface, offset int64) bool {
	proc.curOffset += offset
	ret := expr.Eval(proc)
	proc.curOffset -= offset
	return ret
}

func (proc *Processor) Trit(offset int64) int8 {
	return proc.buffer[proc.curOffset+offset]
}
