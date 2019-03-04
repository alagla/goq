package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/quplayaml"
)

type QuplaCondExpr struct {
	QuplaExprBase
}

func AnalyzeCondExpr(exprYAML *QuplaCondExprYAML, module ModuleInterface, scope FuncDefInterface) (*QuplaCondExpr, error) {
	module.IncStat("numCond")

	ret := &QuplaCondExpr{
		QuplaExprBase: NewQuplaExprBase(exprYAML.Source),
	}
	if ifExpr, err := module.AnalyzeExpression(exprYAML.If, scope); err != nil {
		return nil, err
	} else {
		ret.AppendSubExpr(ifExpr)
	}
	if ret.subexpr[0].Size() != 1 {
		return nil, fmt.Errorf("condition size must be 1 trit, funDef %v: '%v'", scope.GetName(), ret.source)
	}
	if thenExpr, err := module.AnalyzeExpression(exprYAML.Then, scope); err != nil {
		return nil, err
	} else {
		ret.AppendSubExpr(thenExpr)
	}
	if elseExpr, err := module.AnalyzeExpression(exprYAML.Else, scope); err != nil {
		return nil, err
	} else {
		ret.AppendSubExpr(elseExpr)
	}
	if IsNullExpr(ret.subexpr[1]) && IsNullExpr(ret.subexpr[2]) {
		return nil, fmt.Errorf("can't be both branches null. Dunc def '%v': '%v'", scope.GetName(), ret.source)
	}
	if IsNullExpr(ret.subexpr[1]) {
		ret.subexpr[1].(*QuplaNullExpr).SetSize(ret.subexpr[1].Size())
	}
	if IsNullExpr(ret.subexpr[2]) {
		ret.subexpr[2].(*QuplaNullExpr).SetSize(ret.subexpr[1].Size())
	}
	return ret, nil
}

func (e *QuplaCondExpr) HasState() bool {
	return e.subexpr[0].HasState() || e.subexpr[1].HasState() || e.subexpr[2].HasState()
}

func (e *QuplaCondExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.subexpr[1].Size()
}

func (e *QuplaCondExpr) Eval(proc ProcessorInterface, result Trits) bool {
	var buf [1]int8
	null := proc.Eval(e.subexpr[0], buf[:])
	if null {
		return true
	}
	// bool is 0/1
	switch buf[0] {
	case 1:
		return proc.Eval(e.subexpr[1], result)
	case 0:
		return proc.Eval(e.subexpr[2], result)
	case -1:
		return true
	}
	panic(fmt.Sprintf("trit value in cond expr '%v'", e.source))
}
