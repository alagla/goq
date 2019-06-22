package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/abra"
)

type SliceExpr struct {
	ExpressionBase
	site   *QuplaSite
	offset int
	size   int
	// precalcs to speed up Qupla interpretation
	sliceEnd int
	oneTrit  bool
}

func NewQuplaSliceExpr(site *QuplaSite, src string, offset, size int) *SliceExpr {
	return &SliceExpr{
		ExpressionBase: NewExpressionBase(src),
		site:           site,
		offset:         offset,
		size:           size,
		sliceEnd:       offset + size,
		oneTrit:        size == 1,
	}
}

func (e *SliceExpr) Site() *QuplaSite {
	return e.site
}

func (e *SliceExpr) Copy() ExpressionInterface {
	return &SliceExpr{
		ExpressionBase: NewExpressionBase(e.source),
		site:           e.site,
		offset:         e.offset,
		size:           e.size,
		sliceEnd:       e.sliceEnd,
		oneTrit:        e.oneTrit,
	}
}

func (e *SliceExpr) Size() int {
	if e == nil {
		return 0
	}
	return e.size
}

func (e *SliceExpr) Eval(frame *EvalFrame, result Trits) bool {
	restmp, null := e.site.Eval(frame)
	if !null {
		if e.oneTrit {
			result[0] = restmp[e.offset] // optimization ????
		} else {
			copy(result, restmp[e.offset:e.sliceEnd])
		}
	}
	return null
}

func (e *SliceExpr) GetAbraSiteForNonparamVar(branch *abra.Branch, codeUnit *abra.CodeUnit, vi *QuplaSite) *abra.Site {
	var ret *abra.Site
	lookupName := vi.GetAbraLookupName()
	ret = branch.FindSite(lookupName)
	if ret != nil {
		return ret
	}
	if vi.IsState {
		// state sites go in cycles.
		// This will be placeholder, will be resolved later
		return branch.AddUnfinishedStateSite(lookupName)
	}
	ret = e.site.Assign.GetAbraSite(branch, codeUnit, lookupName)
	return ret
}

func (e *SliceExpr) GetAbraSite(branch *abra.Branch, codeUnit *abra.CodeUnit, lookupName string) *abra.Site {
	var varsite *abra.Site
	if e.site.IsParam {
		varsite = branch.GetInputSite(e.site.Idx)
	} else {
		varsite = e.GetAbraSiteForNonparamVar(branch, codeUnit, e.site)
	}
	if e.offset == 0 && e.size == e.site.Size {
		// no actual slicing
		return varsite
	}
	// for actual slicing we have to have a slicing branch
	slicingBranchBlock := codeUnit.GetSlicingBranchBlock(e.site.Size, e.offset, e.size)
	var ret *abra.Site

	ret = abra.NewKnot(slicingBranchBlock, varsite).NewSite()
	ret.SetLookupName(lookupName)
	return branch.AddOrUpdateSite(ret)
}
