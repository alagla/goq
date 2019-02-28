package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/quplayaml"
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
	restmp, null := proc.EvalVar(e.localVarIdx)
	if null {
		tracef("%v sliceExpr '%v' result == null",
			proc.LevelPrefix(), e.localVarIdx)
		return true
	}
	numCopy := copy(result, restmp[e.offset:e.offset+e.size])
	if numCopy != len(result) {
		panic("wrong slice length 1")
	}
	tracef("%v sliceExpr '%v' offset = %v size = %v result = '%v'",
		proc.LevelPrefix(), e.localVarIdx, e.offset, e.size, TritsToString(result))
	return false
}
