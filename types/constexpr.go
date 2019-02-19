package types

import (
	"fmt"
	"github.com/iotaledger/iota.go/trinary"
	"strconv"
	"strings"
)

type ConstExpression interface {
	GetConstValue() int
}

type QuplaConstExpr struct {
	Operator string                  `yaml:"operator"`
	LhsWrap  *QuplaExpressionWrapper `yaml:"lhs"`
	RhsWrap  *QuplaExpressionWrapper `yaml:"rhs"`
	//---
	lhsExpr ConstExpression
	rhsExpr ConstExpression
}

type QuplaConstTerm struct {
	Operator string                  `yaml:"operator"`
	LhsWrap  *QuplaExpressionWrapper `yaml:"lhs"`
	RhsWrap  *QuplaExpressionWrapper `yaml:"rhs"`
	//---
	lhsExpr ConstExpression
	rhsExpr ConstExpression
}

type QuplaConstTypeName struct {
	TypeName string `yaml:"typeName"` // not used
	Size     string `yaml:"size"`
	//---
	size int
}

type QuplaConstNumber struct {
	Value string `yaml:"value"`
	//--
	value int
}

type ConstValue struct {
	Trits trinary.Trits
}

func ToConstExpression(e ExpressionInterface) (ConstExpression, bool) {
	switch e.(type) {
	case *QuplaConstExpr:
		return e.(ConstExpression), true
	case *QuplaConstTerm:
		return e.(ConstExpression), true
	case *QuplaConstTypeName:
		return e.(ConstExpression), true
	case *QuplaConstNumber:
		return e.(ConstExpression), true
	}
	return nil, false
}

func (e *QuplaConstExpr) Analyze(module *QuplaModule) error {
	var err error
	var ei ExpressionInterface
	var ok bool
	if !strings.Contains("+-", e.Operator) {
		return fmt.Errorf("wrong operator symbol %v", e.Operator)
	}
	ei, err = e.LhsWrap.Unwarp()
	if err = ei.Analyze(module); err != nil {
		return err
	}
	if e.lhsExpr, ok = ToConstExpression(ei); !ok {
		return fmt.Errorf("must be const expression")
	}
	ei, err = e.RhsWrap.Unwarp()
	if err = ei.Analyze(module); err != nil {
		return err
	}
	if e.rhsExpr, ok = ToConstExpression(ei); !ok {
		return fmt.Errorf("must be const expression")
	}
	return nil
}

func (e *QuplaConstTerm) Analyze(module *QuplaModule) error {
	var err error
	var ei ExpressionInterface
	var ok bool
	if !strings.Contains("*/%", e.Operator) {
		return fmt.Errorf("wrong operator symbol %v", e.Operator)
	}
	ei, err = e.LhsWrap.Unwarp()
	if err = ei.Analyze(module); err != nil {
		return err
	}
	if e.lhsExpr, ok = ToConstExpression(ei); !ok {
		return fmt.Errorf("must be const expression")
	}
	ei, err = e.RhsWrap.Unwarp()
	if err = ei.Analyze(module); err != nil {
		return err
	}
	if e.rhsExpr, ok = ToConstExpression(ei); !ok {
		return fmt.Errorf("must be const expression")
	}
	return nil
}

func (e *QuplaConstTypeName) Analyze(module *QuplaModule) error {
	var err error
	if e.size, err = strconv.Atoi(e.Size); err != nil {
		return err
	}
	return nil
}

func (e *QuplaConstNumber) Analyze(module *QuplaModule) error {
	var err error
	e.value, err = strconv.Atoi(e.Value)
	return err
}

func (e *QuplaConstExpr) GetConstValue() int {
	lv := e.lhsExpr.GetConstValue()
	rv := e.rhsExpr.GetConstValue()
	switch e.Operator {
	case "+":
		return lv + rv
	case "-":
		return lv - rv
	}
	panic("bad operator")
}

func (e *QuplaConstTerm) GetConstValue() int {
	lv := e.lhsExpr.GetConstValue()
	rv := e.rhsExpr.GetConstValue()
	switch e.Operator {
	case "*":
		return lv * rv
	case "/":
		if rv == 0 {
			panic("divide by zero")
		}
		return lv / rv
	case "%":
		if rv == 0 {
			panic("divide by zero")
		}
		return lv % rv
	}
	panic("bad operator")
}

func (e *QuplaConstTypeName) GetConstValue() int {
	return e.size
}

func (e *QuplaConstNumber) GetConstValue() int {
	return e.value
}

func GetTritValue(e ConstExpression, size int) *ConstValue {
	v := e.GetConstValue()
	t := trinary.IntToTrits(int64(v))
	switch {
	case size < len(t):
		panic(fmt.Errorf("value doesn't fit into %v trits", size))
	case size == len(t):
		return &ConstValue{
			Trits: t,
		}
	case size > len(t):
		ret := make(trinary.Trits, size, size)
		copy(ret, t)
		return &ConstValue{
			Trits: ret,
		}
	}
	panic("inconsistency")
}
