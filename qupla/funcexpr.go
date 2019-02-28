package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/quplayaml"
)

type QuplaFuncExpr struct {
	name    string
	funcDef *QuplaFuncDef
	args    []ExpressionInterface
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

func (e *QuplaFuncExpr) Eval(proc ProcessorInterface, result Trits) bool {
	tracef("%vfuncExpr '%v'", proc.LevelPrefix(), e.name)
	null := proc.Eval(e.funcDef.retExpr, result)
	tracef("%vfuncExpr '%v' res = '%v' null = %v", proc.LevelPrefix(), e.name, TritsToString(result), null)
	return null
}
