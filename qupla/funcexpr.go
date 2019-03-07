package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	"github.com/lunfardo314/goq/dispatcher"
	"github.com/lunfardo314/goq/entities"
	. "github.com/lunfardo314/goq/quplayaml"
)

type QuplaFuncExpr struct {
	QuplaExprBase
	source  string
	funcDef *QuplaFuncDef
}

func NewQuplaFuncExpr(src string, funcDef FuncDefInterface) *QuplaFuncExpr {
	if fd, ok := funcDef.(*QuplaFuncDef); !ok {
		return nil
	} else {
		return &QuplaFuncExpr{
			QuplaExprBase: NewQuplaExprBase(src),
			funcDef:       fd,
		}
	}
}

func AnalyzeFuncExpr(exprYAML *QuplaFuncExprYAML, module ModuleInterface, scope FuncDefInterface) (*QuplaFuncExpr, error) {
	var err error
	ret := &QuplaFuncExpr{
		QuplaExprBase: NewQuplaExprBase(exprYAML.Source),
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

	for _, arg := range exprYAML.Args {
		if fe, err = module.AnalyzeExpression(arg, scope); err != nil {
			return nil, err
		}
		ret.AppendSubExpr(fe)
	}
	err = ret.funcDef.checkArgSizes(ret.subexpr)
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
		context:  e,
		parent:   parent,
		buffer:   make(Trits, e.funcDef.bufLen, e.funcDef.bufLen),
		valueTag: make([]uint8, numVars, numVars),
	}
}

func (e *QuplaFuncExpr) Eval(proc ProcessorInterface, result Trits) bool {
	return proc.Eval(e.funcDef.retExpr, result)
}

func (e *QuplaFuncExpr) References(funName string) bool {
	if e.funcDef.GetName() == funName {
		return true
	}
	return e.ReferencesSubExprs(funName)
}

// only expressions which can be calculated in nil context are valid
func (e *QuplaFuncExpr) NewEntity() dispatcher.EntityInterface {
	return entities.NewFunctionEntity(e.funcDef)
}
