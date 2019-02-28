package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/quplayaml"
)

type QuplaLutExpr struct {
	argExpr []ExpressionInterface
	lutDef  *QuplaLutDef
}

func AnalyzeLutExpr(exprYAML *QuplaLutExprYAML, module ModuleInterface, scope FuncDefInterface) (*QuplaLutExpr, error) {
	var err error
	var ae ExpressionInterface
	var li LUTInterface
	var ok bool
	ret := &QuplaLutExpr{}
	li, err = module.FindLUTDef(exprYAML.Name)
	if err != nil {
		return nil, err
	}
	ret.lutDef, ok = li.(*QuplaLutDef)
	if !ok {
		return nil, fmt.Errorf("inconsistency with types")
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

func (e *QuplaLutExpr) Eval(proc ProcessorInterface, result Trits) bool {
	tracef("%v lutExpr '%v'", proc.LevelPrefix(), e.lutDef.name)
	var buf [3]int8 // no more than 3 inputs
	for i, a := range e.argExpr {
		if proc.Eval(a, buf[i:i+1]) {
			tracef("%v--- eval lutExpr '%v'. Arg %v is null", proc.LevelPrefix(), e.lutDef.name, i)
			return true
		}
	}
	lutArg := buf[:e.lutDef.inputSize]
	tracef("%v LUT lookup '%v' args '%v'", proc.LevelPrefix(), e.lutDef.name, TritsToString(lutArg))
	null := e.lutDef.Lookup(result, lutArg)
	tracef("%v lutExpr '%v' res = '%v' null = %v", proc.LevelPrefix(), e.lutDef.name, TritsToString(result), null)
	return null
}
