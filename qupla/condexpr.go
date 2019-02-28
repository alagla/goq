package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/quplayaml"
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
	if exprYAML.Else == nil {
		fmt.Printf("kuku")
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

func (e *QuplaCondExpr) Eval(proc ProcessorInterface, result Trits) bool {
	var buf [1]int8
	null := e.ifExpr.Eval(proc, buf[:])
	if null {
		return true
	}
	// bool is 0/1
	switch buf[0] {
	case 1:
		return e.thenExpr.Eval(proc, result)
	case 0:
		return e.elseExpr.Eval(proc, result)
	}
	panic("trit value in cond expr")
}
