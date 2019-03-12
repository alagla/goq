package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/dispatcher"
	"github.com/lunfardo314/goq/utils"
	. "github.com/lunfardo314/quplayaml/quplayaml"
	"time"
)

type QuplaExecStmt struct {
	Entity
	QuplaExprBase
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

func (ex *QuplaExecStmt) GetName() string {
	return fmt.Sprintf("#%v-'%v'", ex.num, ex.GetSource())
}

func (ex *QuplaExecStmt) HasState() bool {
	return ex.funcExpr.funcDef.hasState
}

// create temporary environment
// create two temporary entities
//   - one for the eval function expression itself, affect the environment
//   - another for the reaction of the function result: in case of eval ir prints result, in case of test it checks test
//   - post effect to the environment
func (ex *QuplaExecStmt) Execute(disp *Dispatcher) (bool, error) {
	envInName := "ENV_IN$$" + ex.GetName() + "$$"
	envOutName := "ENV_OUT$$" + ex.GetName() + "$$"
	var err error
	//if err = disp.CreateEnvironment(envInName); err != nil {
	//	return false, err
	//}
	//if err = disp.CreateEnvironment(envOutName); err != nil {
	//	return false, err
	//}
	//
	exprEntity := ex.newEvalEntity(disp)
	if err = disp.Attach(exprEntity, []string{envInName}, []string{envOutName}); err != nil {
		return false, nil
	}

	var t = Trits{0}
	var result Trits

	if err = disp.RunQuant(envInName, t, false); err != nil {
		return false, err
	}
	if result, err = disp.Value(envOutName); err != nil {
		return false, err
	}
	logf(0, "Executing %v", ex.GetName())
	logf(0, "    eval result:     '%v'. Duration %v", utils.TritsToString(result), ex.duration)

	var passed bool
	if ex.isTest {
		logf(0, "    expected result: '%v'", utils.TritsToString(ex.valueExpected))
		if passed, _ = TritsEqual(result, ex.valueExpected); passed {
			logf(0, "    test PASSED")
		} else {
			logf(0, "    test FAILED")
		}
	}

	//logf(0, "environment values after quant:")
	//printTritMap(disp.Values())

	_ = disp.DeleteEnvironment(envInName)
	_ = disp.DeleteEnvironment(envOutName)

	return passed, err
}

// expression shouldn't have free variables
// only used to call from executables
type execEvalCallable struct {
	exec *QuplaExecStmt
}

func (ec *execEvalCallable) Call(_ Trits, res Trits) bool {
	start := time.Now()
	null := ec.exec.module.processor.Eval(ec.exec.funcExpr, res)
	ec.exec.duration = time.Since(start)
	return null
}

func (ex *QuplaExecStmt) newEvalEntity(disp *Dispatcher) *Entity {
	return NewEntity(disp, "EVAL_"+ex.funcExpr.GetSource(), 0, ex.funcExpr.Size(), &execEvalCallable{ex})
}
