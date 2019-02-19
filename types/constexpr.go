package types

import (
	"fmt"
	"strings"
)

type QuplaConstExpr struct {
	Operator string                  `yaml:"operator"`
	LhsWrap  *QuplaExpressionWrapper `yaml:"lhs"`
	RhsWrap  *QuplaExpressionWrapper `yaml:"rhs"`
	//---
	lhsExpr ExpressionInterface
	rhsExpr ExpressionInterface
}

type QuplaConstTerm struct {
	Operator string                  `yaml:"operator"`
	LhsWrap  *QuplaExpressionWrapper `yaml:"lhs"`
	RhsWrap  *QuplaExpressionWrapper `yaml:"rhs"`
	//---
	lhsExpr ExpressionInterface
	rhsExpr ExpressionInterface
}

type QuplaConstTypeName struct {
	TypeName string `yaml:"typeName"`
	Size     string `yaml:"size"`
}

type QuplaConstNumber struct {
	Value string `yaml:"value"`
}

func (e *QuplaConstExpr) Analyze(module *QuplaModule) error {
	var err error
	if !strings.Contains("+-", e.Operator) {
		return fmt.Errorf("wrong operator symbol %v", e.Operator)
	}
	if e.lhsExpr, err = e.LhsWrap.Unwarp(); err != nil {
		return err
	}
	if e.rhsExpr, err = e.RhsWrap.Unwarp(); err != nil {
		return err
	}
	if err := e.rhsExpr.Analyze(module); err != nil {
		return err
	}
	return e.rhsExpr.Analyze(module)
}

func (e *QuplaConstTerm) Analyze(module *QuplaModule) error {
	var err error
	if !strings.Contains("*/%", e.Operator) {
		return fmt.Errorf("wrong operator symbol %v", e.Operator)
	}
	if e.lhsExpr, err = e.LhsWrap.Unwarp(); err != nil {
		return err
	}
	if e.rhsExpr, err = e.RhsWrap.Unwarp(); err != nil {
		return err
	}
	if err := e.rhsExpr.Analyze(module); err != nil {
		return err
	}
	return e.rhsExpr.Analyze(module)
}

func (e *QuplaConstTypeName) Analyze(module *QuplaModule) error {
	return nil
}

func (e *QuplaConstNumber) Analyze(module *QuplaModule) error {
	return nil
}
