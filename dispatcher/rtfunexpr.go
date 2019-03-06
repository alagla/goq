package dispatcher

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
)

type RunTimeFuncExpr struct {
	funDef FuncDefInterface
	argBuf Trits
}

type RTSlice struct {
	theSlice Trits
}

func NewRTSlice(buffer Trits, offset, size int64) Trits {
	return buffer[offset : offset+size]
}

func (s *RTSlice) GetSource() string {
	return "(runtime slice)"
}
func (s *RTSlice) Size() int64 {
	return int64(len(s.theSlice))
}

func (s *RTSlice) Eval(_ ProcessorInterface, res Trits) bool {
	copy(res, s.theSlice)
	return false
}

func (s *RTSlice) References(string) bool {
	return false
}
