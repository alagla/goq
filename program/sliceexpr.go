package program

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
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

func (e *QuplaSliceExpr) Eval(callFrame *CallFrame, result Trits) bool {
	res, null := callFrame.EvalVar(e.localVarIdx) // must be big enough to fit whole result
	if null {
		return true
	}
	copy(result, res[e.Offset:e.Offset+e.SliceSize])
	return false
}
