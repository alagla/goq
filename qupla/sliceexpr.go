package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/quplayaml"
)

type QuplaSliceExpr struct {
	localVarIdx int
	varScope    *QuplaFuncDef
	offset      int64
	size        int64
}

func AnalyzeSliceExpr(exprYAML *QuplaSliceExprYAML, module *QuplaModule, scope *QuplaFuncDef) (*QuplaSliceExpr, error) {
	var err error
	ret := &QuplaSliceExpr{
		offset: exprYAML.Offset,
		size:   exprYAML.SliceSize,
	}
	if ret.localVarIdx, err = scope.FindVarIdx(exprYAML.Var, module); err != nil {
		return nil, err
	}
	module.IncStat("numSliceExpr")
	if ret.localVarIdx < 0 {
		return nil, fmt.Errorf("can't find local variable '%v' in scope '%v'", exprYAML.Var, scope.GetName())
	}
	ret.varScope = scope
	if ret.offset+ret.size > scope.VarByIdx(ret.localVarIdx).size {
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
