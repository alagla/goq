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

func (expr *QuplaExpressionWrapper) Unwarp() (ExpressionInterface, error) {
	if expr == nil {
		return nil, nil
	}
	var ret ExpressionInterface
	var numCases int

	if expr.CondExpr != nil {
		ret = expr.CondExpr
		numCases++
	}
	if expr.LutExpr != nil {
		ret = expr.LutExpr
		numCases++
	}
	if expr.SliceExpr != nil {
		ret = expr.SliceExpr
		numCases++
	}
	if expr.ValueExpr != nil {
		ret = expr.ValueExpr
		numCases++
	}
	if expr.FuncExpr != nil {
		ret = expr.FuncExpr
		numCases++
	}
	if expr.FieldExpr != nil {
		ret = expr.FieldExpr
		numCases++
	}
	if expr.ConstNumber != nil {
		ret = expr.ConstNumber
		numCases++
	}
	if expr.ConstTypeName != nil {
		ret = expr.ConstTypeName
		numCases++
	}
	if expr.ConstTerm != nil {
		ret = expr.ConstTerm
		numCases++
	}
	if expr.ConstExpr != nil {
		ret = expr.ConstExpr
		numCases++
	}
	if expr.ConcatExpr != nil {
		ret = expr.ConcatExpr
		numCases++
	}
	if expr.MergeExpr != nil {
		ret = expr.MergeExpr
		numCases++
	}
	if expr.TypeExpr != nil {
		ret = expr.TypeExpr
		numCases++
	}
	if numCases != 1 {
		return nil, fmt.Errorf("internal error: must be exactly one expression case")
	}
	return ret, nil
}
