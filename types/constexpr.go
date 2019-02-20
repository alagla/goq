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
}

type QuplaConstTerm struct {
	Operator string                  `yaml:"operator"`
	LhsWrap  *QuplaExpressionWrapper `yaml:"lhs"`
	RhsWrap  *QuplaExpressionWrapper `yaml:"rhs"`
}

type QuplaConstTypeName struct {
	TypeName string `yaml:"typeName"` // not used
	Size     string `yaml:"size"`
}

type QuplaConstNumber struct {
	Value string `yaml:"value"`
}

type ConstValue struct {
	Value int64
	Trits trinary.Trits
}

func IsConstExpression(e ExpressionInterface) bool {
	if e == nil {
		return false //???????
	}
	switch e.(type) {
	case *QuplaConstExpr:
		return true
	case *QuplaConstTerm:
		return true
	case *QuplaConstTypeName:
		return true
	case *QuplaConstNumber:
		return true
	case *ConstValue:
		return true
	}
	return false
}

func (e *ConstValue) Analyze(module *QuplaModule) (ExpressionInterface, error) {
	return e, nil
}

func (e *QuplaConstExpr) Analyze(module *QuplaModule) (ExpressionInterface, error) {
	var err error
	var lei, rei ExpressionInterface
	var ok bool
	if !strings.Contains("+-", e.Operator) {
		return nil, fmt.Errorf("wrong operator symbol %v", e.Operator)
	}
	if lei, err = e.LhsWrap.Analyze(module); err != nil {
		return nil, err
	}
	if rei, err = e.RhsWrap.Analyze(module); err != nil {
		return nil, err
	}
	if !IsConstExpression(e.LhsWrap) || !IsConstExpression(e.RhsWrap) {
		return nil, fmt.Errorf("operands must be constant expression")
	}
	var rv, lv *ConstValue
	if lv, ok = lei.(*ConstValue); !ok {
		return nil, fmt.Errorf("inconsistency I")
	}
	if rv, ok = rei.(*ConstValue); !ok {
		return nil, fmt.Errorf("inconsistency II")
	}
	var ret *ConstValue
	switch e.Operator {
	case "+":
		ret = NewConstValue(lv.Value+rv.Value, 0)
	case "-":
		ret = NewConstValue(lv.Value-rv.Value, 0)
	}
	return ret, nil
}

func (e *QuplaConstTerm) Analyze(module *QuplaModule) (ExpressionInterface, error) {
	var err error
	var lei, rei ExpressionInterface
	var ok bool
	if !strings.Contains("*/%", e.Operator) {
		return nil, fmt.Errorf("wrong operator symbol %v", e.Operator)
	}
	if lei, err = e.LhsWrap.Analyze(module); err != nil {
		return nil, err
	}
	if rei, err = e.RhsWrap.Analyze(module); err != nil {
		return nil, err
	}
	if !IsConstExpression(e.LhsWrap) || !IsConstExpression(e.RhsWrap) {
		return nil, fmt.Errorf("operands must be constant expression")
	}
	var rv, lv *ConstValue
	if lv, ok = lei.(*ConstValue); !ok {
		return nil, fmt.Errorf("inconsistency I")
	}
	if rv, ok = rei.(*ConstValue); !ok {
		return nil, fmt.Errorf("inconsistency II")
	}
	var ret *ConstValue
	switch e.Operator {
	case "*":
		ret = NewConstValue(lv.Value*rv.Value, 0)
	case "/":
		if rv.Value != 0 {
			ret = NewConstValue(lv.Value/rv.Value, 0)
		} else {
			return nil, fmt.Errorf("division by 0 in constant expression")
		}
	case "%":
		if rv.Value != 0 {
			ret = NewConstValue(lv.Value%rv.Value, 0)
		} else {
			return nil, fmt.Errorf("division by 0 in constant expression")
		}
	}
	return ret, nil
}

func (e *QuplaConstTypeName) Analyze(module *QuplaModule) (ExpressionInterface, error) {
	var err error
	var ret int
	if ret, err = strconv.Atoi(e.Size); err != nil {
		return nil, err
	}
	return NewConstValue(int64(ret), 0), nil
}

func (e *QuplaConstNumber) Analyze(module *QuplaModule) (ExpressionInterface, error) {
	ret, err := strconv.Atoi(e.Value)
	if err != nil {
		return nil, err
	}
	return NewConstValue(int64(ret), 0), nil
}

func NewConstValue(value int64, size int) *ConstValue {
	t := trinary.IntToTrits(value)
	switch {
	case size <= len(t):
		return &ConstValue{
			Value: value,
			Trits: t,
		}
	case size > len(t):
		ret := make(trinary.Trits, size, size)
		copy(ret, t)
		return &ConstValue{
			Value: value,
			Trits: ret,
		}
	}
	panic("inconsistency")
}

func GetConstValue(expr ExpressionInterface) (int64, error) {
	cv, ok := expr.(*ConstValue)
	if !ok {
		return 0, fmt.Errorf("not a constant value")
	}
	return cv.Value, nil
}
