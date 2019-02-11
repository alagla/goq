package expr

import "fmt"

const (
	ExprType_CondExpr      = 0
	ExprType_LutExpr       = ExprType_CondExpr + 1
	ExprType_SliceExpr     = ExprType_LutExpr + 1
	ExprType_ValueExpr     = ExprType_SliceExpr + 1
	ExprType_FuncExpr      = ExprType_ValueExpr + 1
	ExprType_FieldExpr     = ExprType_FuncExpr + 1
	ExprType_ConstNumber   = ExprType_FieldExpr + 1
	ExprType_ConstTypeName = ExprType_ConstNumber + 1
	ExprType_ConstTerm     = ExprType_ConstTypeName + 1
	ExprType_ConstExpr     = ExprType_ConstTerm + 1
	ExprType_ConcatExpr    = ExprType_ConstExpr + 1
	ExprType_MergeExpr     = ExprType_ConcatExpr + 1
	ExprType_TypeExpr      = ExprType_MergeExpr + 1
)

type ExpressionInterface interface {
	Analyze() error
}

type QuplaExpression struct {
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
	//-----
	exprType      int
	theExpression ExpressionInterface
}

func (expr *QuplaExpression) Analyze() error {
	if expr == nil {
		return nil
	}
	var numCases int

	if expr.CondExpr != nil {
		expr.exprType = ExprType_CondExpr
		expr.theExpression = expr.CondExpr
		numCases++
	}
	if expr.LutExpr != nil {
		expr.exprType = ExprType_LutExpr
		expr.theExpression = expr.LutExpr
		numCases++
	}
	if expr.SliceExpr != nil {
		expr.exprType = ExprType_SliceExpr
		expr.theExpression = expr.SliceExpr
		numCases++
	}
	if expr.ValueExpr != nil {
		expr.exprType = ExprType_ValueExpr
		expr.theExpression = expr.ValueExpr
		numCases++
	}
	if expr.FuncExpr != nil {
		expr.exprType = ExprType_FuncExpr
		expr.theExpression = expr.FuncExpr
		numCases++
	}
	if expr.FieldExpr != nil {
		expr.exprType = ExprType_FieldExpr
		expr.theExpression = expr.FieldExpr
		numCases++
	}
	if expr.ConstNumber != nil {
		expr.exprType = ExprType_ConstNumber
		expr.theExpression = expr.ConstNumber
		numCases++
	}
	if expr.ConstTypeName != nil {
		expr.exprType = ExprType_ConstTypeName
		expr.theExpression = expr.ConstTypeName
		numCases++
	}
	if expr.ConstTerm != nil {
		expr.exprType = ExprType_ConstTerm
		expr.theExpression = expr.ConstTerm
		numCases++
	}
	if expr.ConstExpr != nil {
		expr.exprType = ExprType_ConstExpr
		expr.theExpression = expr.ConstExpr
		numCases++
	}
	if expr.ConcatExpr != nil {
		expr.exprType = ExprType_ConcatExpr
		expr.theExpression = expr.ConcatExpr
		numCases++
	}
	if expr.MergeExpr != nil {
		expr.exprType = ExprType_MergeExpr
		expr.theExpression = expr.MergeExpr
		numCases++
	}
	if expr.TypeExpr != nil {
		expr.exprType = ExprType_TypeExpr
		expr.theExpression = expr.TypeExpr
		numCases++
	}
	if numCases != 1 {
		return fmt.Errorf("internal error: must be exactly one expression case. Probably incorrect YAML")
	}
	return expr.theExpression.Analyze()
}
