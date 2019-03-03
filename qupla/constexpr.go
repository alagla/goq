package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/quplayaml"
	"strconv"
	"strings"
)

type ConstExpression interface {
	ExpressionInterface
	GetConstValue() int64
	GetConstName() string
}

type ConstValue struct {
	QuplaExprBase
	value int64
	name  string
	size  int64
}

func (e *ConstValue) Size() int64 {
	return 0 //todo
}

func (_ *ConstValue) Eval(_ ProcessorInterface, _ Trits) bool {
	return true // todo
}

func (e *ConstValue) GetConstValue() int64 {
	return e.value
}

func (e *ConstValue) GetConstName() string {
	return e.name
}

func AnalyzeConstExpr(exprYAML *QuplaConstExprYAML, module ModuleInterface, scope FuncDefInterface) (ConstExpression, error) {
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
	var rv, lv ConstExpression
	if lv, ok = lei.(ConstExpression); !ok {
		return nil, fmt.Errorf("must be constant expression")
	}
	if rv, ok = rei.(ConstExpression); !ok {
		return nil, fmt.Errorf("must be constant expression")
	}
	var ret ConstExpression
	switch exprYAML.Operator {
	case "+":
		ret = NewConstValue("", lv.GetConstValue()+rv.GetConstValue())
	case "-":
		ret = NewConstValue("", lv.GetConstValue()-rv.GetConstValue())
	}
	return ret, nil
}

func AnalyzeConstTerm(exprYAML *QuplaConstTermYAML, module ModuleInterface, scope FuncDefInterface) (ConstExpression, error) {
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
	var rv, lv ConstExpression
	if lv, ok = lei.(ConstExpression); !ok {
		return nil, fmt.Errorf("inconsistency I")
	}
	if rv, ok = rei.(ConstExpression); !ok {
		return nil, fmt.Errorf("inconsistency II")
	}
	var ret ConstExpression
	switch exprYAML.Operator {
	case "*":
		ret = NewConstValue("", lv.GetConstValue()*rv.GetConstValue())
	case "/":
		if rv.GetConstValue() != 0 {
			ret = NewConstValue("", lv.GetConstValue()/rv.GetConstValue())
		} else {
			return nil, fmt.Errorf("division by 0 in constant expression")
		}
	case "%":
		if rv.GetConstValue() != 0 {
			ret = NewConstValue("", lv.GetConstValue()%rv.GetConstValue())
		} else {
			return nil, fmt.Errorf("division by 0 in constant expression")
		}
	}
	return ret, nil
}

func AnalyzeConstNumber(exprYAML *QuplaConstNumberYAML, _ ModuleInterface, _ FuncDefInterface) (ConstExpression, error) {
	ret, err := strconv.Atoi(exprYAML.Value)
	if err != nil {
		return nil, err
	}
	return NewConstValue("", int64(ret)), nil
}

func NewConstValue(name string, value int64) ConstExpression {
	return &ConstValue{
		name:  name,
		value: value,
	}
}

func (e *ConstValue) GetTrits() Trits {
	t := IntToTrits(e.value)
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
	ce, ok := expr.(ConstExpression)
	if !ok {
		return 0, fmt.Errorf("not a constant value")
	}
	return ce.GetConstValue(), nil
}
