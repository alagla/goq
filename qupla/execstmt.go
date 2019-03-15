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

func (ex *QuplaExecStmt) GetIdx() int {
	return ex.num
}

func (ex *QuplaExecStmt) HasState() bool {
	return ex.funcExpr.funcDef.hasState
}

func (ex *QuplaExecStmt) Execute(disp *Dispatcher) (bool, error) {
	return ex.ExecuteMulti(disp, 1)
}

// create temporary environment
// create two temporary entities
//   - one for the eval function expression itself, affect the environment
//   - another for the reaction of the function result: in case of eval ir prints result, in case of test it checks test
//   - post effect to the environment
func (ex *QuplaExecStmt) ExecuteMulti(disp *Dispatcher, repeat int) (bool, error) {
	if repeat < 1 {
		return false, fmt.Errorf("'repeat' parameter must be >1")
	}
	envInName := "ENV_IN$$" + ex.GetName() + "$$"
	envOutName := "ENV_OUT$$" + ex.GetName() + "$$"
	var err error
	exprEntity := ex.newEvalEntity(disp)
	if err = disp.Attach(exprEntity, []string{envInName}, []string{envOutName}); err != nil {
		return false, nil
	}

	var t = Trits{0}
	var result Trits

	start := time.Now()
	for i := 0; i < repeat; i++ {
		err = disp.StartQuant(envInName, t, func() { logf(3, "%v ++++++++++ Done with %v\n", i, envInName) })
		if err != nil {
			return false, err
		}
	}

	if result, err = disp.Value(envOutName); err != nil {
		return false, err
	}
	logf(0, "Executing %v. Repeat %v times", ex.GetName(), repeat)
	dur := time.Since(start)
	avgdur := int64(dur/time.Millisecond) / int64(repeat)
	logf(0, "    eval result:     '%v'. Total duration %v, Average duration %v msec/run",
		utils.TritsToString(result), dur, avgdur)

	var passed bool
	if ex.isTest {
		logf(0, "    expected result: '%v'", utils.TritsToString(ex.valueExpected))
		if passed = ex.ResultIsExpected(result); passed {
			logf(0, "    test PASSED")
		} else {
			logf(0, "    test FAILED")
		}
	}
	_ = disp.DeleteEnvironment(envInName)
	_ = disp.DeleteEnvironment(envOutName)

	return passed, err
}

func (ex *QuplaExecStmt) ResultIsExpected(result Trits) bool {
	passed, _ := TritsEqual(result, ex.valueExpected)
	if passed {
		return true
	}
	if len(result) != len(ex.valueExpected) {
		return false
	}
	if !ex.isFloat {
		return false
	}
	dif0 := result[0] - ex.valueExpected[0]
	if dif0 < 0 {
		dif0 = -dif0
	}
	if dif0 > 1 {
		return false
	}
	if len(result) == 1 {
		return true
	}
	passed, _ = TritsEqual(result[1:], ex.valueExpected[1:])
	return passed

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
	name := fmt.Sprintf("#%v-EVAL_%v", ex.num, ex.funcExpr.GetSource())
	return NewEntity(disp, name, 0, ex.funcExpr.Size(), &execEvalCallable{ex})
}
