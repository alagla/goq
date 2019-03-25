package analyzeyaml

import (
	"fmt"
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/goq/readyaml"
	"strconv"
	"strings"
)

func AnalyzeConstExpr(exprYAML *QuplaConstExprYAML, module *QuplaModule, scope *Function) (ConstExpression, error) {
	var err error
	var lei, rei ExpressionInterface
	var ok bool
	if !strings.Contains("+-", exprYAML.Operator) {
		return nil, fmt.Errorf("wrong operator symbol %v", exprYAML.Operator)
	}
	if lei, err = AnalyzeExpression(exprYAML.Lhs, module, scope); err != nil {
		return nil, err
	}
	if rei, err = AnalyzeExpression(exprYAML.Rhs, module, scope); err != nil {
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

func AnalyzeConstTerm(exprYAML *QuplaConstTermYAML, module *QuplaModule, scope *Function) (ConstExpression, error) {
	var err error
	var lei, rei ExpressionInterface
	var ok bool
	if !strings.Contains("+-", exprYAML.Operator) {
		return nil, fmt.Errorf("wrong operator symbol %v", exprYAML.Operator)
	}
	if lei, err = AnalyzeExpression(exprYAML.Lhs, module, scope); err != nil {
		return nil, err
	}
	if rei, err = AnalyzeExpression(exprYAML.Rhs, module, scope); err != nil {
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

func AnalyzeConstNumber(exprYAML *QuplaConstNumberYAML, _ *QuplaModule, _ *Function) (ConstExpression, error) {
	ret, err := strconv.Atoi(exprYAML.Value)
	if err != nil {
		return nil, err
	}
	return NewConstValue("", ret), nil
}
