package eval

import (
	"github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/program"
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
	return proc.buffer[offset:size]
}

func (proc *Processor) Eval(expr program.ExpressionInterface, offset int64) bool {

}
