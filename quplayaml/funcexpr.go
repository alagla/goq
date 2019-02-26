package quplayaml

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
)

type QuplaFuncExpr struct {
	name    string
	funcDef *QuplaFuncDef
	args    []ExpressionInterface
}

type CallFrame struct {
	context   *QuplaFuncExpr // which function called
	parent    *CallFrame     // context where it was called
	buffer    Trits          // buffer to place all params and variables
	evaluated []bool         // flag if respective variable was evaluated
	isNull    []bool         // flag if value was evaluated to null
}

func AnalyzeFuncExpr(exprYAML *QuplaFuncExprYAML, module ModuleInterface, scope FuncDefInterface) (*QuplaFuncExpr, error) {
	var err error
	ret := &QuplaFuncExpr{
		name: exprYAML.Name,
	}
	var fdi FuncDefInterface
	fdi, err = module.FindFuncDef(exprYAML.Name)
	if err != nil {
		return nil, err
	}
	var ok bool
	if ret.funcDef, ok = fdi.(*QuplaFuncDef); !ok {
		return nil, fmt.Errorf("inconsistency with types in %v", exprYAML.Name)
	}

	var fe ExpressionInterface
	module.IncStat("numFuncExpr")

	ret.args = make([]ExpressionInterface, 0, len(exprYAML.Args))
	for _, arg := range exprYAML.Args {
		if fe, err = module.AnalyzeExpression(arg, scope); err != nil {
			return nil, err
		}
		ret.args = append(ret.args, fe)
	}
	err = ret.funcDef.checkArgSizes(ret.args)
	return ret, err
}

func (e *QuplaFuncExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.funcDef.Size()
}

func (e *QuplaFuncExpr) NewCallFrame(parent *CallFrame) *CallFrame {
	numVars := len(e.funcDef.localVars)
	return &CallFrame{
		context:   e,
		parent:    parent,
		buffer:    make(Trits, e.funcDef.bufLen, e.funcDef.bufLen),
		evaluated: make([]bool, numVars, numVars),
		isNull:    make([]bool, numVars, numVars),
	}
}

func (e *QuplaFuncExpr) Eval(parentFrame *CallFrame, result Trits) bool {
	tracef("eval funcExpr '%v'", e.name)

	frame := e.NewCallFrame(parentFrame)
	return e.funcDef.retExpr.Eval(frame, result)
}

func (frame *CallFrame) EvalVar(idx int64) (Trits, bool) {
	tracef("eval var funcExpr '%v', idx = %v", frame.context.name, idx)
	null := false
	resOffset := frame.context.funcDef.localVars[idx].offset
	resSize := frame.context.funcDef.localVars[idx].size
	resultSlice := frame.buffer[resOffset : resOffset+resSize]
	if frame.evaluated[idx] {
		if frame.isNull[idx] {
			null = true
		}
	} else {
		if idx < frame.context.funcDef.numParams {
			// variable is parameter
			null = frame.context.args[idx].Eval(frame.parent, resultSlice)
		} else {
			// variable is assigned
			null = frame.context.funcDef.localVars[idx].expr.Eval(frame, resultSlice)
		}
		if null {
			resultSlice = nil
		}
	}
	return resultSlice, null
}
