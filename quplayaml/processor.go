package quplayaml

import . "github.com/iotaledger/iota.go/trinary"
import . "github.com/lunfardo314/goq/abstract"

type CallFrame struct {
	context   *QuplaFuncExpr // which function called
	parent    *CallFrame     // context where it was called
	buffer    Trits          // buffer to place all params and variables
	evaluated []bool         // flag if respective variable was evaluated
	isNull    []bool         // flag if value was evaluated to null
	result    Trits          // slice where to put result
}

type StackProcessor struct {
	curFrame *CallFrame
}

func (proc *StackProcessor) Push(funExpr ExpressionInterface) {

}

func (proc *StackProcessor) Pull() {

}

func (proc *StackProcessor) Eval(expr ExpressionInterface, result Trits) bool {
	return true
}

func (proc *StackProcessor) EvalVar(idx int64) bool {
	return true
}

func (proc *StackProcessor) Slice(offset, size int64) Trits {
	return proc.curFrame.buffer[offset : offset+size]
}
