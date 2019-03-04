package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	"github.com/lunfardo314/goq/cfg"
	. "github.com/lunfardo314/goq/quplayaml"
	"github.com/lunfardo314/goq/utils"
	"math/big"
	"time"
)

type QuplaExecStmt struct {
	QuplaExprBase
	source string
	isTest bool
	//expr         ExpressionInterface
	funcExpr     *QuplaFuncExpr
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
	var expr ExpressionInterface
	var ok bool
	expr, err = module.factory.AnalyzeExpression(execStmtYAML.Expr, module, nil)
	if err != nil {
		return err
	}
	if res.funcExpr, ok = expr.(*QuplaFuncExpr); !ok {
		return fmt.Errorf("top expression must be call to a function: '%v'", execStmtYAML.Source)
	}
	res.isTest = execStmtYAML.Expected != nil
	if res.isTest {
		res.exprExpected, err = module.factory.AnalyzeExpression(execStmtYAML.Expected, module, nil)
		if err != nil {
			return err
		}
		// check sizes
		if err = MatchSizes(res.funcExpr, res.exprExpected); err != nil {
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

func (ex *QuplaExecStmt) HasState() bool {
	return ex.funcExpr.funcDef.hasState
}

func (ex *QuplaExecStmt) Execute() (bool, error) {
	//ex.module.processor.SetTrace(ex.num == 0, 0)
	logf(2, "Running #%v: '%v'", ex.num, ex.GetSource())

	start := time.Now()

	if !cfg.Config.ExecEvals && !cfg.Config.ExecTests {
		return true, nil
	}
	resExpr := make(Trits, ex.funcExpr.Size(), ex.funcExpr.Size())
	null := ex.module.processor.Eval(ex.funcExpr, resExpr)

	var minVerb int
	var d *big.Int
	if !ex.isTest {
		if cfg.Config.ExecEvals {
			minVerb = 0
		} else {
			minVerb = 2
		}
		logf(minVerb, "Executing #%v: '%v'. Duration %v", ex.num, ex.GetSource(), time.Since(start))
		if null {
			logf(minVerb, "Eval result is null")
		} else {
			d, _ = utils.TritsToBigInt(resExpr)
			logf(minVerb, "Eval result is '%v' (dec = %v)", utils.TritsToString(resExpr), d)
		}
		return true, nil
	}
	if !cfg.Config.ExecTests {
		return true, nil
	}
	if null {
		logf(0, "Expression result is null. Test #%v FAILED: '%v'", ex.num, ex.GetSource())
		return false, nil
	}
	passed := false
	resExpected := make(Trits, ex.funcExpr.Size(), ex.exprExpected.Size())
	null = ex.module.processor.Eval(ex.exprExpected, resExpected)

	exp, err := utils.TritsToBigInt(resExpected)
	if err != nil {
		return false, err
	}
	logf(2, "Expected result '%v' (dec = %v)", utils.TritsToString(resExpected), exp)
	passed, _ = TritsEqual(resExpected, resExpr)
	if passed {
		logf(2, "Test #%v PASSED: '%v' Duration %v", ex.num, ex.GetSource(), time.Since(start))
	} else {
		d, _ = utils.TritsToBigInt(resExpr)
		logf(0, "Test #%v FAILED: '%v' Duration %v", ex.num, ex.GetSource(), time.Since(start))
		if cfg.Config.Verbosity < 2 {
			logf(0, "    Expected result '%v'", utils.TritsToString(resExpected))
			logf(0, "    Eval result '%v'", utils.TritsToString(resExpr))
		}
	}
	return passed, nil
}
