package quplayaml

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
)

type QuplaSliceExpr struct {
	localVarIdx int64
	varScope    *QuplaFuncDef
	offset      int64
	size        int64
}

func AnalyzeSliceExpr(exprYAML *QuplaSliceExprYAML, module ModuleInterface, scope FuncDefInterface) (*QuplaSliceExpr, error) {
	var err error
	ret := &QuplaSliceExpr{
		offset: exprYAML.Offset,
		size:   exprYAML.SliceSize,
	}
	module.IncStat("numSliceExpr")
	var vi *VarInfo
	if vi, err = scope.GetVarInfo(exprYAML.Var, module); err != nil {
		return nil, err
	}
	ret.localVarIdx = vi.idx
	if ret.localVarIdx < 0 {
		return nil, fmt.Errorf("can't find local variable '%v' in scope '%v'", exprYAML.Var, scope.GetName())
	}
	ret.varScope = scope.(*QuplaFuncDef)
	if ret.offset+ret.size > vi.size {
		return nil, fmt.Errorf("wrong offset/size for the slice of '%v'", exprYAML.Var)
	}
	return ret, nil
}

func (e *QuplaSliceExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.size
}

func (e *QuplaSliceExpr) Eval(callFrame *CallFrame, result Trits) bool {
	res, null := callFrame.EvalVar(e.localVarIdx) // must be big enough to fit whole result
	if null {
		return true
	}
	copy(result, res[e.offset:e.offset+e.size])
	return false
}
