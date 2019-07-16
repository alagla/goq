package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/abra"
	cabra "github.com/lunfardo314/goq/abra/construct"
)

type CondExpr struct {
	ExpressionBase
}

func (e *CondExpr) Size() int {
	if e == nil {
		return 0
	}
	return e.subExpr[1].Size()
}

const (
	BOOL_TRUE  = 1
	BOOL_FALSE = -1
)

func (e *CondExpr) Eval(frame *EvalFrame, result Trits) bool {
	var buf [1]int8
	null := e.subExpr[0].Eval(frame, buf[:])
	if null {
		return true
	}
	// bool is 0/1
	switch buf[0] {
	case BOOL_TRUE:
		return e.subExpr[1].Eval(frame, result)
	case BOOL_FALSE:
		return e.subExpr[2].Eval(frame, result)
	default:
		return true
	}
	//panic(Sprintf("trit value in cond Expr '%v'", e.source))
}

func (e *CondExpr) Copy() ExpressionInterface {
	return &CondExpr{
		ExpressionBase: e.copyBase(),
	}
}

func (e *CondExpr) GetAbraSite(branch *abra.Branch, codeUnit *abra.CodeUnit, lookupName string) *abra.Site {
	condSite := e.subExpr[0].GetAbraSite(branch, codeUnit, "")

	trueSite := e.subExpr[1].GetAbraSite(branch, codeUnit, "")
	nullifyTrueBlock := cabra.GetNullifyBranchBlock(codeUnit, e.subExpr[1].Size(), true)
	nullifiedTrueSite := cabra.NewKnotSite(e.subExpr[1].Size(), "", nullifyTrueBlock, condSite, trueSite)
	if _, ok := e.subExpr[2].(*NullExpr); ok {
		// if second arg is null, then first arg is stored as final site
		nullifiedTrueSite.LookupName = lookupName
		ret := cabra.AddOrUpdateSite(branch, nullifiedTrueSite)
		return ret
	}

	cabra.AddOrUpdateSite(branch, nullifiedTrueSite)

	falseSite := e.subExpr[2].GetAbraSite(branch, codeUnit, "")
	nullifyFalseBlock := cabra.GetNullifyBranchBlock(codeUnit, e.subExpr[2].Size(), false)
	nullifiedFalseSite := cabra.NewKnotSite(e.subExpr[2].Size(), "", nullifyFalseBlock, condSite, falseSite)
	nullifiedFalseSite = cabra.AddOrUpdateSite(branch, nullifiedFalseSite)

	ret := cabra.NewMergeSite(e.subExpr[2].Size(), lookupName, nullifiedTrueSite, nullifiedFalseSite)
	return cabra.AddOrUpdateSite(branch, ret)
}
