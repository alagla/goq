package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/quplayaml"
)

type QuplaCondExpr struct {
	QuplaExprBase
	ifExpr   ExpressionInterface
	thenExpr ExpressionInterface
	elseExpr ExpressionInterface
}

func AnalyzeCondExpr(exprYAML *QuplaCondExprYAML, module ModuleInterface, scope FuncDefInterface) (*QuplaCondExpr, error) {
	var err error
	module.IncStat("numCond")

	ret := &QuplaCondExpr{
		QuplaExprBase: NewQuplaExprBase(exprYAML.Source),
	}
	if ret.ifExpr, err = module.AnalyzeExpression(exprYAML.If, scope); err != nil {
		return nil, err
	}
	if ret.ifExpr.Size() != 1 {
		return nil, fmt.Errorf("condition size must be 1 trit, funDef %v: '%v'", scope.GetName(), ret.source)
	}
	if ret.thenExpr, err = module.AnalyzeExpression(exprYAML.Then, scope); err != nil {
		return nil, err
	}
	if ret.elseExpr, err = module.AnalyzeExpression(exprYAML.Else, scope); err != nil {
		return nil, err
	}
	if IsNullExpr(ret.thenExpr) && IsNullExpr(ret.elseExpr) {
		return nil, fmt.Errorf("can't be both branches null. Dunc def '%v': '%v'", scope.GetName(), ret.source)
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

func (e *QuplaCondExpr) Eval(proc ProcessorInterface, result Trits) bool {
	var buf [1]int8
	null := proc.Eval(e.ifExpr, buf[:])
	if null {
		return true
	}
	// bool is 0/1
	switch buf[0] {
	case 1:
		return proc.Eval(e.thenExpr, result)
	case 0:
		return proc.Eval(e.elseExpr, result)
	}
	panic(fmt.Sprintf("trit value in cond expr '%v'", e.source))
}
