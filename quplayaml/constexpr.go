package quplayaml

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	"strconv"
	"strings"
)

type ConstExpression interface {
	GetConstValue() int
}

type ConstValue struct {
	Value int64
	size  int64
}

func (e *ConstValue) Size() int64 {
	return 0 //todo
}

func (_ *ConstValue) Eval(_ ProcessorInterface, _ Trits) bool {
	return true // todo
}

func IsConstExpression(e interface{}) bool {
	switch e.(type) {
	case *QuplaConstExprYAML:
		return true
	case *QuplaConstTermYAML:
		return true
	case *QuplaConstTypeNameYAML:
		return true
	case *QuplaConstNumberYAML:
		return true
	case *ConstValue:
		return true
	}
	return false
}

func AnalyzeConstExpr(exprYAML *QuplaConstExprYAML, module ModuleInterface, scope FuncDefInterface) (*ConstValue, error) {
	var err error
	var lei, rei ExpressionInterface
	var ok bool
	if !strings.Contains("+-", exprYAML.Operator) {
		return nil, fmt.Errorf("wrong operator symbol %v", exprYAML.Operator)
	}
	if lei, err = module.AnalyzeExpression(exprYAML.Lhs, scope); err != nil {
		return nil, err
	}
	if rei, err = module.AnalyzeExpression(exprYAML.Rhs, scope); err != nil {
		return nil, err
	}
	if !IsConstExpression(lei) || !IsConstExpression(rei) {
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
	switch exprYAML.Operator {
	case "+":
		ret = NewConstValue(lv.Value + rv.Value)
	case "-":
		ret = NewConstValue(lv.Value - rv.Value)
	}
	return ret, nil
}

func AnalyzeConstTerm(exprYAML *QuplaConstTermYAML, module ModuleInterface, scope FuncDefInterface) (*ConstValue, error) {
	var err error
	var lei, rei ExpressionInterface
	var ok bool
	if !strings.Contains("+-", exprYAML.Operator) {
		return nil, fmt.Errorf("wrong operator symbol %v", exprYAML.Operator)
	}
	if lei, err = module.AnalyzeExpression(exprYAML.Lhs, scope); err != nil {
		return nil, err
	}
	if rei, err = module.AnalyzeExpression(exprYAML.Rhs, scope); err != nil {
		return nil, err
	}
	if !IsConstExpression(lei) || !IsConstExpression(rei) {
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
	switch exprYAML.Operator {
	case "*":
		ret = NewConstValue(lv.Value * rv.Value)
	case "/":
		if rv.Value != 0 {
			ret = NewConstValue(lv.Value / rv.Value)
		} else {
			return nil, fmt.Errorf("division by 0 in constant expression")
		}
	case "%":
		if rv.Value != 0 {
			ret = NewConstValue(lv.Value % rv.Value)
		} else {
			return nil, fmt.Errorf("division by 0 in constant expression")
		}
	}
	return ret, nil
}

func AnalyzeConstTypeName(exprYAML *QuplaConstTypeNameYAML, _ ModuleInterface, _ FuncDefInterface) (*ConstValue, error) {
	var err error
	var ret int
	if ret, err = strconv.Atoi(exprYAML.SizeString); err != nil {
		return nil, err
	}
	return NewConstValue(int64(ret)), nil
}

func AnalyzeConstNumber(exprYAML *QuplaConstNumberYAML, _ ModuleInterface, _ FuncDefInterface) (*ConstValue, error) {
	ret, err := strconv.Atoi(exprYAML.Value)
	if err != nil {
		return nil, err
	}
	return NewConstValue(int64(ret)), nil
}

func NewConstValue(value int64) *ConstValue {
	return &ConstValue{
		Value: value,
	}
}

func (e *ConstValue) GetTrits() Trits {
	t := IntToTrits(e.Value)
	if e.size == 0 {
		return t
	}
	if e.size == int64(len(t)) {
		return t
	}
	ret := make(Trits, 0, e.size)
	copy(ret, t)
	return ret
}

func GetConstValue(expr ExpressionInterface) (int64, error) {
	cv, ok := expr.(*ConstValue)
	if !ok {
		return 0, fmt.Errorf("not a constant value")
	}
	return cv.Value, nil
}
