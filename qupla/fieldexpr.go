package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/quplayaml/go"
)

type QuplaFieldExpr struct {
	QuplaExprBase
	condExpr ExpressionInterface
}

func AnalyzeFieldExpr(exprYAML *QuplaFieldExprYAML, module ModuleInterface, scope FuncDefInterface) (*QuplaFieldExpr, error) {
	var err error
	module.IncStat("numFieldExpr")
	ret := &QuplaFieldExpr{}
	ret.condExpr, err = module.AnalyzeExpression(exprYAML.CondExpr, scope)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (e *QuplaFieldExpr) Size() int64 {
	if e == nil {
		return 0
	}
	return e.condExpr.Size()
}

func (e *QuplaFieldExpr) Eval(_ ProcessorInterface, _ Trits) bool {
	return true
}
