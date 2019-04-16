package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
)

type SliceExpr struct {
	ExpressionBase
	site     *QuplaSite
	offset   int
	size     int
	sliceEnd int
	noSlice  bool
	oneTrit  bool
}

func NewQuplaSliceExpr(site *QuplaSite, src string, offset, size int) *SliceExpr {
	noSlice := offset == 0 && size == site.Size
	return &SliceExpr{
		ExpressionBase: NewExpressionBase(src),
		site:           site,
		offset:         offset,
		size:           size,
		sliceEnd:       offset + size,
		noSlice:        noSlice,
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
		noSlice:        e.noSlice,
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
			if e.noSlice {
				copy(result, restmp)
			} else {
				copy(result, restmp[e.offset:e.sliceEnd])
			}
		}
	}
	return null
}
