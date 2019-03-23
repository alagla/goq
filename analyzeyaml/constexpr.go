package analyzeyaml

import (
	"fmt"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/quplayaml/quplayaml"
	"strconv"
	"strings"
)

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
