package types

import (
	"fmt"
)

type QuplaExpressionWrapper struct {
	CondExpr      *QuplaCondExpr      `yaml:"CondExpr,omitempty"`
	LutExpr       *QuplaLutExpr       `yaml:"LutExpr,omitempty"`
	SliceExpr     *QuplaSliceExpr     `yaml:"SliceExpr,omitempty"`
	ValueExpr     *QuplaValueExpr     `yaml:"ValueExpr,omitempty"`
	FuncExpr      *QuplaFuncExpr      `yaml:"FuncExpr,omitempty"`
	FieldExpr     *QuplaFieldExpr     `yaml:"FieldExpr,omitempty"`
	ConstNumber   *QuplaConstNumber   `yaml:"ConstNumber,omitempty"`
	ConstTypeName *QuplaConstTypeName `yaml:"ConstTypeName,omitempty"`
	ConstTerm     *QuplaConstTerm     `yaml:"ConstTerm,omitempty"`
	ConstExpr     *QuplaConstExpr     `yaml:"ConstExpr,omitempty"`
	ConcatExpr    *QuplaConcatExpr    `yaml:"ConcatExpr,omitempty"`
	MergeExpr     *QuplaMergeExpr     `yaml:"MergeExpr,omitempty"`
	TypeExpr      *QuplaTypeExpr      `yaml:"TypeExpr,omitempty"`
}

func (e *QuplaExpressionWrapper) Analyze(module *QuplaModule, scope *QuplaFuncDef) (ExpressionInterface, error) {
	if e == nil {
		return nil, nil
	}
	ret, err := e.unwarp()
	if err != nil {
		return nil, err
	}
	if ret != nil {
		if ret, err = ret.Analyze(module, scope); err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func (e *QuplaExpressionWrapper) unwarp() (ExpressionInterface, error) {
	var ret ExpressionInterface
	var numCases int

	if e.CondExpr != nil {
		ret = e.CondExpr
		numCases++
	}
	if e.LutExpr != nil {
		ret = e.LutExpr
		numCases++
	}
	if e.SliceExpr != nil {
		ret = e.SliceExpr
		numCases++
	}
	if e.ValueExpr != nil {
		ret = e.ValueExpr
		numCases++
	}
	if e.FuncExpr != nil {
		ret = e.FuncExpr
		numCases++
	}
	if e.FieldExpr != nil {
		ret = e.FieldExpr
		numCases++
	}
	if e.ConstNumber != nil {
		ret = e.ConstNumber
		numCases++
	}
	if e.ConstTypeName != nil {
		ret = e.ConstTypeName
		numCases++
	}
	if e.ConstTerm != nil {
		ret = e.ConstTerm
		numCases++
	}
	if e.ConstExpr != nil {
		ret = e.ConstExpr
		numCases++
	}
	if e.ConcatExpr != nil {
		ret = e.ConcatExpr
		numCases++
	}
	if e.MergeExpr != nil {
		ret = e.MergeExpr
		numCases++
	}
	if e.TypeExpr != nil {
		ret = e.TypeExpr
		numCases++
	}
	if numCases != 1 {
		return nil, fmt.Errorf("internal error: must be exactly one expression case")
	}
	return ret, nil
}
