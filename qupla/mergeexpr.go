package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/abra"
	cabra "github.com/lunfardo314/goq/abra/construct"
)

type MergeExpr struct {
	ExpressionBase
}

func NewMergeExpression(src string, args []ExpressionInterface) *MergeExpr {
	ret := &MergeExpr{
		ExpressionBase: NewExpressionBase(src),
	}
	for _, a := range args {
		ret.AppendSubExpr(a)
	}
	return ret
}

func (e *MergeExpr) Size() int {
	if e == nil {
		return 0
	}
	return e.subExpr[0].Size()
}

func (e *MergeExpr) Copy() ExpressionInterface {
	return &MergeExpr{
		ExpressionBase: e.copyBase(),
	}
}

func (e *MergeExpr) Eval(frame *EvalFrame, result Trits) bool {
	for _, se := range e.subExpr {
		if !se.Eval(frame, result) {
			return false // return first not null
		}
	}
	return true // all nulls
}

func (e *MergeExpr) GetAbraSite(branch *abra.Branch, codeUnit *abra.CodeUnit, lookupName string) *abra.Site {
	inputs := make([]*abra.Site, 0, len(e.subExpr))
	for _, se := range e.subExpr {
		s := se.GetAbraSite(branch, codeUnit, "")
		inputs = append(inputs, s)
	}
	ret := cabra.NewMergeSite(e.Size(), lookupName, inputs...)
	return cabra.AddOrUpdateSite(branch, ret)
}
