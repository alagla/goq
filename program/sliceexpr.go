package program

import (
	"fmt"
)

type QuplaSliceExpr struct {
	Var           string                  `yaml:"var"`
	Offset        int64                   `yaml:"offset"`
	SliceSize     int64                   `yaml:"size"`
	StartExprWrap *QuplaExpressionWrapper `yaml:"start,omitempty"` // not used
	EndExprWrap   *QuplaExpressionWrapper `yaml:"end,omitempty"`   // not used
	//----
	localVarIdx int
	varScope    *QuplaFuncDef
}

func (e *QuplaSliceExpr) Analyze(module *QuplaModule, scope *QuplaFuncDef) (ExpressionInterface, error) {
	var err error

	if e.localVarIdx, err = scope.FindVarIdx(e.Var, module); err != nil {
		return nil, err
	}
	module.IncStat("numSliceExpr")
	if e.localVarIdx < 0 {
		return nil, fmt.Errorf("can't find local variable '%v' in scope '%v'", e.Var, scope.GetName())
	}
	e.varScope = scope
	if e.Offset+e.SliceSize > scope.VarByIdx(e.localVarIdx).size {
		return nil, fmt.Errorf("wrong offset/size for the slice of '%v'", e.Var)
	}
	return e, nil
}

func (e *QuplaSliceExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.SliceSize
}

func (e *QuplaSliceExpr) Eval(proc *Processor) bool {
	null := proc.Eval(e.varScope.VarByIdx(e.localVarIdx).expr, 0)
	if null {
		return true
	}
	copy(proc.Slice(0, e.SliceSize), proc.Slice(e.Offset, e.SliceSize))
	return true
}
