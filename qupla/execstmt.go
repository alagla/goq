package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/quplayaml"
	"github.com/lunfardo314/goq/utils"
	"time"
)

type QuplaExecStmt struct {
	QuplaExprBase
	source       string
	isTest       bool
	expr         ExpressionInterface
	exprExpected ExpressionInterface
	module       *QuplaModule
	num          int
}

func AnalyzeExecStmt(execStmtYAML *QuplaExecStmtYAML, module *QuplaModule) error {
	res := &QuplaExecStmt{
		QuplaExprBase: NewQuplaExprBase(execStmtYAML.Source),
		module:        module,
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

	//ex.module.processor.SetTrace(ex.num == 19, 0)

	debugf("-------------")
	debugf("running #%v: '%v'", ex.num, ex.GetSource())

	start := time.Now()

	resExpr := make(Trits, ex.expr.Size(), ex.expr.Size())
	null := ex.module.processor.Eval(ex.expr, resExpr)

	debugf("Duration: %v", time.Since(start))

	if null {
		debugf("eval result is null")
		if ex.isTest {
			debugf("Test FAILED")
			return nil
		}
	} else {
		d, _ := utils.TritsToBigInt(resExpr)
		debugf("eval result dec = %v trits = '%v' ", d, utils.TritsToString(resExpr))
	}
	if ex.isTest {
		resExpected := make(Trits, ex.expr.Size(), ex.exprExpected.Size())
		null = ex.module.processor.Eval(ex.exprExpected, resExpected)

		exp, err := utils.TritsToBigInt(resExpected)
		if err != nil {
			return err
		}
		debugf("expected dec = %v", exp)
		if eq, _ := TritsEqual(resExpected, resExpr); eq {
			debugf("Test PASSED")
		} else {
			debugf("Test FAILED")
		}
	}
	return nil
}
