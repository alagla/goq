package analyzeyaml

import (
	. "fmt"
	. "github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/quplayaml/quplayaml"
)

func AnalyzeCondExpr(exprYAML *QuplaCondExprYAML, module *QuplaModule, scope *Function) (*CondExpr, error) {
	module.IncStat("numCond")

	ret := &CondExpr{
		ExpressionBase: NewExpressionBase(exprYAML.Source),
	}
	if ifExpr, err := AnalyzeExpression(exprYAML.If, module, scope); err != nil {
		return nil, err
	} else {
		ret.AppendSubExpr(ifExpr)
	}
	if ret.NumSubExpr() != 1 {
		return nil, Errorf("condition size must be 1 trit, funDef %v: '%v'", scope.Name, ret.GetSource())
	}
	if thenExpr, err := AnalyzeExpression(exprYAML.Then, module, scope); err != nil {
		return nil, err
	} else {
		ret.AppendSubExpr(thenExpr)
	}
	if elseExpr, err := AnalyzeExpression(exprYAML.Else, module, scope); err != nil {
		return nil, err
	} else {
		ret.AppendSubExpr(elseExpr)
	}
	s1 := ret.GetSubExpr(1)
	s2 := ret.GetSubExpr(2)
	if IsNullExpr(s1) && IsNullExpr(s2) {
		return nil, Errorf("can't be both branches null. Dunc def '%v': '%v'", scope.Name, ret.GetSource())
	}
	if IsNullExpr(s1) {
		s1.(*NullExpr).SetSize(s1.Size())
	}
	if IsNullExpr(s2) {
		s2.(*NullExpr).SetSize(s1.Size())
	}
	return ret, nil
}
