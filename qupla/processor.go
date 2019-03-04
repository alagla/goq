package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/utils"
	"strings"
)
import . "github.com/lunfardo314/goq/abstract"

type CallFrame struct {
	context  *QuplaFuncExpr // which function called
	parent   *CallFrame     // context where it was called
	buffer   Trits          // buffer to place all params and variables
	valueTag []uint8        // 0x01 bit mean 'evaluated', 0x02 bit means is not null
	result   Trits          // slice where to put result
}

type StackProcessor struct {
	curFrame      *CallFrame
	curStateParam Trits
	// aux
	levelFunc    int
	maxLevelFunc int
	level        int
	maxLevel     int
	numfuncall   int
	numvarcall   int
	trace        bool
	maxTraces    int
	curTraces    int
}

func NewStackProcessor() *StackProcessor {
	return &StackProcessor{}
}

func (proc *StackProcessor) Eval(expr ExpressionInterface, result Trits) bool {
	funExpr, isFunction := expr.(*QuplaFuncExpr)
	if isFunction {
		proc.numfuncall++
		proc.levelFunc++
		if proc.levelFunc > proc.maxLevelFunc {
			proc.maxLevelFunc = proc.levelFunc
		}
		proc.tracef("IN funExpr '%v'", funExpr.name)

		proc.curFrame = funExpr.NewCallFrame(proc.curFrame)
	} else {
		proc.level++
	}
	if proc.level > proc.maxLevel {
		proc.maxLevel = proc.level
	}
	null := expr.Eval(proc, result)
	if isFunction {
		proc.tracef("OUT funExpr '%v' null = %v res = '%v'", funExpr.name, null, utils.TritsToString(result))
		proc.levelFunc--
		proc.curFrame = proc.curFrame.parent
	} else {
		proc.level--
	}
	return null
}

func (proc *StackProcessor) Reset() {
	// not reset state var storage
	debugf("Proc stats: numfuncall = %v numvarcall = %v maxLevelFunc = %v maxLevel = %v",
		proc.numfuncall, proc.numvarcall, proc.maxLevelFunc, proc.maxLevel)
	proc.numfuncall = 0
	proc.numvarcall = 0
	proc.levelFunc = 0
	proc.level = 0
	proc.maxLevelFunc = 0
	proc.maxLevel = 0
}

func (proc *StackProcessor) EvalVar(idx int64) (Trits, bool) {
	proc.numvarcall++
	var null bool
	if proc.curFrame == nil {
		panic("variable can't be evaluated in nil context")
	}
	vi := proc.curFrame.context.funcDef.VarByIdx(idx)
	if vi == nil {
		panic("wrong var idx")
	}

	ret := proc.Slice(vi.Offset, vi.Size)

	if proc.curFrame.valueTag[vi.Idx]&0x01 != 0 {
		isNull := proc.curFrame.valueTag[vi.Idx]&0x02 != 0
		proc.tracef("EvalVar %v(%v) in '%v': already evaluated to '%v' null = %v",
			vi.Name, idx, proc.curFrame.context.funcDef.name, utils.TritsToString(ret), isNull)
		return ret, isNull
	}
	if vi.IsParam {
		expr := proc.curFrame.context.subexpr[vi.Idx]
		proc.tracef("EvalVar %v(idx=%v) in '%v' param=true expr = '%v'",
			vi.Name, idx, proc.curFrame.context.funcDef.name, expr.GetSource())
		saveCurFrame := proc.curFrame
		proc.curFrame = proc.curFrame.parent
		null = proc.Eval(expr, ret)
		proc.curFrame = saveCurFrame
	} else {
		proc.tracef("EvalVar %v (idx=%v) in '%v' param=false expr = '%v'",
			vi.Name, idx, proc.curFrame.context.funcDef.name, vi.Assign.GetSource())
		null = proc.Eval(vi.Assign, ret)
	}
	proc.tracef("Return EvalVar %v (idx=%v) in '%v': res = '%v' null = %v",
		vi.Name, idx, proc.curFrame.context.funcDef.name, utils.TritsToString(ret), null)

	proc.curFrame.valueTag[vi.Idx] |= 0x01 // mark evaluated
	if null {
		proc.curFrame.valueTag[vi.Idx] |= 0x02 // mark is null
	}
	if null {
		return nil, true
	}
	return ret, false
}

func (proc *StackProcessor) Slice(offset, size int64) Trits {
	return proc.curFrame.buffer[offset : offset+size]
}

func (proc *StackProcessor) SetTrace(trace bool, maxTraces int) {
	proc.trace = trace
	proc.maxTraces = maxTraces
	proc.curTraces = 0
}

func (proc *StackProcessor) LevelPrefix() string {
	r := strings.Repeat(".", proc.levelFunc) + strings.Repeat(" ", proc.level)
	return fmt.Sprintf("%4d: %s", proc.curTraces, r)
	//return fmt.Sprintf("%5d-%5d: "+r, proc.numfuncall, proc.numvarcall)
}

func (proc *StackProcessor) tracef(format string, args ...interface{}) {
	if !proc.trace {
		return
	}
	tracef("proc-> "+proc.LevelPrefix()+format, args...)
	proc.trace = proc.maxTraces <= 0 || proc.curTraces < proc.maxTraces
	proc.curTraces++
}
