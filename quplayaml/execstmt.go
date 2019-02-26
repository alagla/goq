package quplayaml

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
)

type QuplaExecStmt struct {
	isTest       bool
	expr         ExpressionInterface
	exprExpected ExpressionInterface
}

func AnalyzeExecStmt(execStmtYAML *QuplaExecStmtYAML, module *QuplaModule) error {
	res := &QuplaExecStmt{}
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
	funcExpr, ok := ex.expr.(*QuplaFuncExpr)
	if !ok {
		return fmt.Errorf("must be call to function")
	}
	res := make(Trits, funcExpr.Size(), funcExpr.Size())
	null := funcExpr.Eval(nil, res)
	if null {
		debugf("result is null")
	} else {
		debugf("result is not null")
	}
	return nil
}
