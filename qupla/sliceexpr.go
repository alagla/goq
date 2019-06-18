package qupla

import (
	"fmt"
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

func (e *SliceExpr) GenAbraSite(branch *abra.Branch, codeUnit *abra.CodeUnit) *abra.Site {
	if e.site.IsParam {
		// for inputs we take site not by name but by index
		return branch.InputSites[e.site.Idx]
	}
	var ret *abra.Site
	// for other sites (body and state) we use variable name_offset_size
	lookupName := fmt.Sprintf("qupla_site_%s_%d_%d", e.site.Name, e.offset, e.size)
	ret = branch.FindBodySite(lookupName)
	if ret != nil {
		return ret
	}
	if e.offset == 0 && e.size == e.site.Size {
		// no actual slicing
		// generate new site
		ret = e.site.Assign.GenAbraSite(branch, codeUnit).SetLookupName(lookupName)
		return ret
	}
	// for actual slicing we have to have a slicing branch
	slicingBranchBlock := codeUnit.GetSlicingBranch(e.offset, e.size)
	if e.site.IsParam {
		ret = abra.NewKnot(slicingBranchBlock, branch.InputSites[e.site.Idx]).NewSite(lookupName)
	} else {
		input := e.site.Assign.GenAbraSite(branch, codeUnit)
		ret = abra.NewKnot(slicingBranchBlock, input).NewSite(lookupName)
	}
	return ret
}
