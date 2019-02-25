package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/quplayaml"
)

type QuplaLutExpr struct {
	argExpr []ExpressionInterface
	lutDef  *QuplaLutDef
}

func AnalyzeLutExpr(exprYAML *QuplaLutExprYAML, module *QuplaModule, scope *QuplaFuncDef) (*QuplaLutExpr, error) {
	var err error
	var ae ExpressionInterface
	ret := &QuplaLutExpr{}
	ret.lutDef, err = module.FindLUTDef(exprYAML.Name)
	if err != nil {
		return nil, err
	}
	module.IncStat("numLUTExpr")

	ret.argExpr = make([]ExpressionInterface, 0, len(exprYAML.Args))
	for _, a := range exprYAML.Args {
		ae, err = module.AnalyzeExpression(a, scope)
		if err != nil {
			return nil, err
		}
		if err = RequireSize(ae, 1); err != nil {
			return nil, fmt.Errorf("LUT expression with '%v': %v", ret.lutDef.name, err)
		}
		ret.argExpr = append(ret.argExpr, ae)
	}
	if ret.lutDef.inputSize != len(ret.argExpr) {
		return nil, fmt.Errorf("num arg doesnt't match input dimension of the LUT %v", ret.lutDef.name)
	}
	return ret, nil
}

func (e *QuplaLutExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.lutDef.Size()
}

func (e *QuplaLutExpr) Eval(callFrame *CallFrame, result Trits) bool {
	tracef("eval var lutExpr '%v', frame = %v", e.lutDef.name, callFrame.context.name)

	var null bool
	var buf [3]int8 // no more than 3 inputs
	for i, a := range e.argExpr {
		null = a.Eval(callFrame, buf[i:i+1])
		if null {
			return true
		}
	}
	return e.lutDef.Lookup(result, buf[:e.lutDef.inputSize])
}
