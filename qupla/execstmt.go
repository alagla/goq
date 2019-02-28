package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/quplayaml"
)

type QuplaExecStmt struct {
	isTest       bool
	expr         ExpressionInterface
	exprExpected ExpressionInterface
	module       *QuplaModule
}

func AnalyzeExecStmt(execStmtYAML *QuplaExecStmtYAML, module *QuplaModule) error {
	res := &QuplaExecStmt{
		module: module,
	}
	var err error
	res.expr, err = module.factory.AnalyzeExpression(execStmtYAML.Expr, module, nil)
	if err != nil {
		return err
	}
	res.isTest = execStmtYAML.Expected != nil
	if res.isTest {
		res.exprExpected, err = module.factory.AnalyzeExpression(execStmtYAML.Expected, module, nil)
		if err != nil {
			return err
		}
		// check sizes
		if err = MatchSizes(res.expr, res.exprExpected); err != nil {
			return err
		}
		module.IncStat("numTest")
	} else {
		res.exprExpected = nil
		module.IncStat("numEval")
	}
	module.AddExec(res)
	return nil
}

func (ex *QuplaExecStmt) Execute() error {
	resExpr := make(Trits, ex.expr.Size(), ex.expr.Size())
	null := ex.module.processor.Eval(ex.expr, resExpr)
	if null {
		debugf("eval result is null")
	} else {
		debugf("eval result = '%v'", TritsToString(resExpr))
	}
	if ex.isTest {
		resExpected := make(Trits, ex.expr.Size(), ex.exprExpected.Size())
		null = ex.module.processor.Eval(ex.exprExpected, resExpected)
		debugf("expected result is '%v'", TritsToString(resExpected))
		if eq, _ := TritsEqual(resExpected, resExpr); eq {
			debugf("Test passed")
		} else {
			debugf("Test failed")
		}
	}
	return nil
}
