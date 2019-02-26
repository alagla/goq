package quplayaml

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
)

type QuplaCondExpr struct {
	ifExpr   ExpressionInterface
	thenExpr ExpressionInterface
	elseExpr ExpressionInterface
}

func AnalyzeCondExpr(exprYAML *QuplaCondExprYAML, module ModuleInterface, scope FuncDefInterface) (*QuplaCondExpr, error) {
	var err error
	module.IncStat("numCond")

	ret := &QuplaCondExpr{}
	if ret.ifExpr, err = module.AnalyzeExpression(exprYAML.If, scope); err != nil {
		return nil, err
	}
	if ret.ifExpr.Size() != 1 {
		return nil, fmt.Errorf("condition size must be 1 trit: scope %v", scope.GetName())
	}
	if ret.thenExpr, err = module.AnalyzeExpression(exprYAML.Then, scope); err != nil {
		return nil, err
	}
	if ret.elseExpr, err = module.AnalyzeExpression(exprYAML.Else, scope); err != nil {
		return nil, err
	}
	if IsNullExpr(ret.thenExpr) && IsNullExpr(ret.elseExpr) {
		return nil, fmt.Errorf("can't be both branches null: scope %v", scope.GetName())
	}
	if IsNullExpr(ret.thenExpr) {
		ret.thenExpr.(*QuplaNullExpr).SetSize(ret.elseExpr.Size())
	}
	if IsNullExpr(ret.elseExpr) {
		ret.elseExpr.(*QuplaNullExpr).SetSize(ret.thenExpr.Size())
	}
	return ret, nil
}
func (e *QuplaCondExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.thenExpr.Size()
}

func (e *QuplaCondExpr) Eval(callFrame *CallFrame, result Trits) bool {
	var buf [1]int8
	null := e.ifExpr.Eval(callFrame, buf[:])
	if null {
		return true
	}
	switch buf[0] {
	case 1:
		return e.thenExpr.Eval(callFrame, result)
	case -1:
		return e.elseExpr.Eval(callFrame, result)
	}
	panic("trit value")
}
