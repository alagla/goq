package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/abra"
	"sort"
)

type FieldExpr struct {
	Offset int
	Size   int
}
type TypeExpr struct {
	ExpressionBase
	size   int
	Fields []*FieldExpr
}

func NewQuplaTypeExpr(src string, size int) *TypeExpr {
	return &TypeExpr{
		ExpressionBase: NewExpressionBase(src),
		size:           size,
		Fields:         make([]*FieldExpr, 0, 5),
	}
}

func (e *TypeExpr) Copy() ExpressionInterface {
	return &TypeExpr{
		ExpressionBase: e.copyBase(),
		size:           e.size,
		Fields:         e.Fields,
	}
}

func (e *TypeExpr) Size() int {
	if e == nil {
		return 0
	}
	return e.size
}

func (e *TypeExpr) Eval(frame *EvalFrame, result Trits) bool {
	for idx, subExpr := range e.subExpr {
		if subExpr.Eval(frame, result[e.Fields[idx].Offset:e.Fields[idx].Offset+e.Fields[idx].Size]) {
			return true
		}
	}
	return false
}

type fieldExprPair struct {
	expr  ExpressionInterface
	field *FieldExpr
}

func (e *TypeExpr) GetAbraSite(branch *abra.Branch, codeUnit *abra.CodeUnit, lookupName string) *abra.Site {
	// sort field expression by field offset
	sortedByOffset := make([]*fieldExprPair, len(e.Fields))
	for i, fe := range e.Fields {
		sortedByOffset[i] = &fieldExprPair{
			expr:  e.GetSubExpr(i),
			field: fe,
		}
	}
	sort.SliceStable(sortedByOffset, func(i, j int) bool {
		return sortedByOffset[i].field.Offset < sortedByOffset[j].field.Offset
	})
	inputs := make([]*abra.Site, len(sortedByOffset))
	for i, fi := range sortedByOffset {
		inputs[i] = fi.expr.GetAbraSite(branch, codeUnit, "")
	}
	concatBranch := codeUnit.GetConcatBlockForSize(e.Size())
	ret := abra.NewKnot(concatBranch, inputs...).NewSite()
	ret.SetLookupName(lookupName)
	return branch.AddOrUpdateSite(ret)
}
