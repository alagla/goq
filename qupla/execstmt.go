package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	"github.com/lunfardo314/goq/cfg"
	. "github.com/lunfardo314/goq/dispatcher"
	. "github.com/lunfardo314/goq/quplayaml"
	"github.com/lunfardo314/goq/utils"
	"math/big"
	"time"
)

type QuplaExecStmt struct {
	BaseEntity
	QuplaExprBase
	source        string
	isTest        bool
	isFloat       bool // needed for float comparison
	funcExpr      *QuplaFuncExpr
	valueExpected Trits
	module        *QuplaModule
	num           int
	duration      time.Duration
}

func AnalyzeExecStmt(execStmtYAML *QuplaExecStmtYAML, module *QuplaModule) error {
	var err error
	var expr ExpressionInterface
	var ok bool
	expr, err = module.factory.AnalyzeExpression(execStmtYAML.Expr, module, nil)
	if err != nil {
		return err
	}
	var funcExpr *QuplaFuncExpr
	if funcExpr, ok = expr.(*QuplaFuncExpr); !ok {
		return fmt.Errorf("top expression must be call to a function: '%v'", execStmtYAML.Source)
	}
	res := &QuplaExecStmt{
		QuplaExprBase: NewQuplaExprBase(execStmtYAML.Source),
		module:        module,
		funcExpr:      funcExpr,
	}

	res.isTest = execStmtYAML.Expected != nil
	var exprExpected ExpressionInterface
	if res.isTest {
		res.isFloat = execStmtYAML.IsFloat
		exprExpected, err = module.factory.AnalyzeExpression(execStmtYAML.Expected, module, nil)
		if err != nil {
			return err
		}
		// check sizes
		if err = MatchSizes(funcExpr, exprExpected); err != nil {
			return err
		}

		ve, ok := exprExpected.(*QuplaValueExpr)
		if !ok {
			return fmt.Errorf("test '%v': left hand side must be ValueExpr", res.GetSource())
		}
		res.valueExpected = ve.TritValue
		module.IncStat("numTest")
	} else {
		res.valueExpected = nil
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
	exp, err := utils.TritsToBigInt(ex.valueExpected)
	if err != nil {
		return false, err
	}
	logf(2, "Expected result '%v' (dec = %v)", utils.TritsToString(ex.valueExpected), exp)
	passed, _ = TritsEqual(ex.valueExpected, resExpr)
	if !passed && ex.isFloat && len(ex.valueExpected) > 1 && len(resExpr) > 1 {
		abs := ex.valueExpected[0] - resExpr[0]
		if abs < 0 {
			abs = -abs
		}
		passed, _ = TritsEqual(ex.valueExpected[1:], resExpr[1:])
		passed = passed && abs <= 1
	}
	if passed {
		logf(2, "Test #%v PASSED: '%v' Duration %v", ex.num, ex.GetSource(), time.Since(start))
	} else {
		d, _ = utils.TritsToBigInt(resExpr)
		logf(0, "Test #%v FAILED: '%v' Duration %v", ex.num, ex.GetSource(), time.Since(start))
		if cfg.Config.Verbosity < 2 {
			logf(0, "    Expected result '%v'", utils.TritsToString(ex.valueExpected))
			logf(0, "    Eval result '%v'", utils.TritsToString(resExpr))
		}
	}
	return passed, nil
}

// create temporary environment
// create two temporary entities
//   - one for the eval function expression itself, it will affect the environment
//   - another for the reaction of the function result: in case of eval ir prints result, in case of test it checks test
//   - invoke function expression
func (ex *QuplaExecStmt) Execute2() (bool, error) {
	env := NewEnvironment("$$tmp_environment$$")
	exprEntity := ex.newFuncExpressionEntity()
	resultEntity := ex.newResultEntity()
	if err := resultEntity.Join(env); err != nil {
		return false, err
	}
	if err := exprEntity.Affect(env); err != nil {
		return false, err
	}
	exprEntity.Invoke(nil)
	env.Stop()
	return false, nil
}

func (ex *QuplaExecStmt) newResultEntity() *BaseEntity {
	var callback func(Trits) Trits
	if ex.isTest {
		callback = func(result Trits) Trits {
			logf(0, "Executing '%v'. Eval result: '%v'",
				ex.source, utils.TritsToString(result))
			logf(0, "    expected result '%v'", utils.TritsToString(ex.valueExpected))

			if passed, _ := TritsEqual(result, ex.valueExpected); passed {
				logf(0, "    test PASSED. Duration %v", ex.duration)
			} else {
				logf(0, "    test FAILED. Duration %v", ex.duration)
			}
			return nil
		}
	} else {
		callback = func(result Trits) Trits {
			logf(0, "Executing '%v'. Eval result: '%v'. Duration %v",
				ex.source, utils.TritsToString(result), ex.duration)
			return nil
		}
	}
	return NewBaseEntity("$$eval_result_entity", ex.funcExpr.Size(), 0, callback)
}

// expression shouldn't have free variables
func (ex *QuplaExecStmt) newFuncExpressionEntity() *BaseEntity {
	effectCallback := func(_ Trits) Trits {
		start := time.Now()
		res := make(Trits, ex.funcExpr.Size(), ex.funcExpr.Size())
		null := ex.module.processor.Eval(ex.funcExpr, res)
		ex.duration = time.Since(start)
		if null {
			return nil
		}
		return res
	}
	return NewBaseEntity(ex.funcExpr.GetSource(), 0, ex.funcExpr.Size(), effectCallback)
}
