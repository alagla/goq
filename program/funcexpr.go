package program

import . "github.com/iotaledger/iota.go/trinary"

type QuplaFuncExpr struct {
	Name     string                    `yaml:"name"`
	ArgsWrap []*QuplaExpressionWrapper `yaml:"args"`
	//---
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

func (e *QuplaFuncExpr) Analyze(module *QuplaModule, scope *QuplaFuncDef) (ExpressionInterface, error) {
	var err error
	e.funcDef, err = module.FindFuncDef(e.Name)
	if err != nil {
		return nil, err
	}
	var fe ExpressionInterface
	module.IncStat("numFuncExpr")

	e.args = make([]ExpressionInterface, 0, len(e.ArgsWrap))
	for _, arg := range e.ArgsWrap {
		if fe, err = arg.Analyze(module, scope); err != nil {
			return nil, err
		}
		e.args = append(e.args, fe)
	}
	err = e.funcDef.checkArgSizes(e.args)
	return e, err
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
	frame := e.NewCallFrame(parentFrame)
	return e.funcDef.retExpr.Eval(frame, result)
}

func (frame *CallFrame) EvalVar(idx int) (Trits, bool) {
	null := false
	resOffset := frame.context.funcDef.localVars[idx].offset
	resSize := frame.context.funcDef.localVars[idx].size
	resultSlice := frame.buffer[resOffset : resOffset+resSize]
	if frame.evaluated[idx] {
		if frame.isNull[idx] {
			null = true
		}
	} else {
		var f *CallFrame
		if idx < frame.context.funcDef.numParams {
			// variable is parameter
			f = frame.parent
		} else {
			// variable is assigned
			f = frame
		}
		null = frame.context.funcDef.localVars[idx].expr.Eval(f, resultSlice)
		if null {
			resultSlice = nil
		}
	}
	return resultSlice, null
}
