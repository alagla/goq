package qupla

import (
	"fmt"
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
	if _, ok := res.expr.(*QuplaFuncExpr); !ok {
		return fmt.Errorf("top expression must be call to a function: '%v'", execStmtYAML.Source)
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
	if res.num == 6 {
		fmt.Printf("kuku\n")
	}
	return nil
}

func (ex *QuplaExecStmt) HasState() bool {
	return ex.expr.HasState()
}

func (ex *QuplaExecStmt) Execute() (time.Duration, bool, error) {
	//ex.module.processor.SetTrace(ex.num == 0, 0)
	debugf("running #%v: '%v'", ex.num, ex.GetSource())

	start := time.Now()

	resExpr := make(Trits, ex.expr.Size(), ex.expr.Size())
	null := ex.module.processor.Eval(ex.expr, resExpr)

	if null {
		debugf("eval result is null")
		if ex.isTest {
			return time.Since(start), false, nil
		}
	} else {
		d, _ := utils.TritsToBigInt(resExpr)
		debugf("eval result '%v' (dec = %v) ", utils.TritsToString(resExpr), d)
	}
	passed := false
	if ex.isTest {
		resExpected := make(Trits, ex.expr.Size(), ex.exprExpected.Size())
		null = ex.module.processor.Eval(ex.exprExpected, resExpected)

		exp, err := utils.TritsToBigInt(resExpected)
		if err != nil {
			return time.Since(start), false, err
		}
		debugf("expected result '%v' (dec = %v)", utils.TritsToString(resExpected), exp)
		passed, _ = TritsEqual(resExpected, resExpr)
	}
	return time.Since(start), passed, nil
}
