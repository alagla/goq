package quplayaml

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
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
	ret.localVarIdx = vi.Idx
	if ret.localVarIdx < 0 {
		return nil, fmt.Errorf("can't find local variable '%v' in scope '%v'", exprYAML.Var, scope.GetName())
	}
	ret.varScope = scope.(*QuplaFuncDef)
	if ret.offset+ret.size > vi.Size {
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

func (e *QuplaSliceExpr) Eval(proc ProcessorInterface, result Trits) bool {
	tracef("%v sliceExpr '%v' offset = %v size = %v", proc.LevelPrefix(), e.localVarIdx, e.offset, e.size)
	null := proc.EvalVar(e.localVarIdx)
	if null {
		tracef("%v sliceExpr '%v' offset = %v size = %v resu = null",
			proc.LevelPrefix(), e.localVarIdx, e.offset, e.size)
		return true
	}
	copy(result, proc.Slice(e.offset, e.size))
	tracef("%v sliceExpr '%v' offset = %v size = %v result = '%v'",
		proc.LevelPrefix(), e.localVarIdx, e.offset, e.size, TritsToString(result))
	return false
}
