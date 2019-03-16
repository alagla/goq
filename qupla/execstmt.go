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
	idx           int
	runResult     *execResultCallable
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
	return fmt.Sprintf("#%v-'%v'", ex.idx, ex.GetSource())
}

func (ex *QuplaExecStmt) GetIdx() int {
	return ex.idx
}

func (ex *QuplaExecStmt) HasState() bool {
	return ex.funcExpr.funcDef.hasState
}

func (ex *QuplaExecStmt) Execute(disp *Dispatcher) (bool, error) {
	return ex.ExecuteMulti(disp, 1)
}

func (ex *QuplaExecStmt) ExecuteMulti(disp *Dispatcher, repeat int) (bool, error) {
	if repeat < 1 {
		return false, fmt.Errorf("'repeat' parameter must be >1")
	}
	var err error
	if err = ex.prepareRun(disp); err != nil {
		return false, err
	}
	var t = Trits{0}
	envInName := ex.inEnvironmentName()
	for i := 0; i < repeat; i++ {
		err = disp.WaveStart(envInName, false, t)
		if err != nil {
			return false, err
		}
	}
	passed := ex.wrapUpRun(disp)
	return passed, err
}

func (ex *QuplaExecStmt) resultIsExpected(result Trits) bool {
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

func (ex *QuplaExecStmt) StartWave(disp *Dispatcher) error {
	var err error
	if err = ex.prepareRun(disp); err != nil {
		return err
	}
	var t = Trits{0}

	envInName := ex.inEnvironmentName()
	err = disp.WaveStart(envInName, true, t)
	return nil
}

// expression shouldn't have free variables
// only used to call from executables
type execEvalCallable struct {
	exec *QuplaExecStmt
}

func (ec *execEvalCallable) Call(_ Trits, res Trits) bool {
	null := ec.exec.module.processor.Eval(ec.exec.funcExpr, res)
	return null
}

func (ex *QuplaExecStmt) newEvalEntity(disp *Dispatcher) *Entity {
	name := fmt.Sprintf("#%v-EVAL_%v", ex.idx, ex.funcExpr.GetSource())
	return NewEntity(disp, name, 0, ex.funcExpr.Size(), &execEvalCallable{ex})
}

type execResultCallable struct {
	entity     *Entity
	start      time.Time
	lastRun    time.Time
	lastResult Trits
	called     int
	exec       *QuplaExecStmt
}

func (ec *execResultCallable) Call(effect Trits, _ Trits) bool {
	ec.lastRun = time.Now()
	ec.called++
	ec.lastResult = effect
	return false
}

func (ex *QuplaExecStmt) newResultEntity(disp *Dispatcher) (*Entity, *execResultCallable) {
	name := fmt.Sprintf("#%v-RESULT_%v", ex.idx, ex.funcExpr.GetSource())
	nowis := time.Now()
	core := &execResultCallable{
		start:   nowis,
		lastRun: nowis,
		exec:    ex,
	}
	ret := NewEntity(disp, name, ex.funcExpr.Size(), 0, core)
	core.entity = ret
	return ret, core
}

func (ex *QuplaExecStmt) inEnvironmentName() string {
	return "IN_EXE"
}

func (ex *QuplaExecStmt) outEnvironmentName() string {
	return "OUT_EXE"
}

func (ex *QuplaExecStmt) prepareRun(disp *Dispatcher) error {
	evalEntity := ex.newEvalEntity(disp)
	if err := disp.Attach(evalEntity, []string{ex.inEnvironmentName()}, []string{ex.outEnvironmentName()}); err != nil {
		return err
	}
	var resultEntity *Entity
	resultEntity, ex.runResult = ex.newResultEntity(disp)
	if err := disp.Attach(resultEntity, []string{ex.outEnvironmentName()}, nil); err != nil {
		return err
	}
	return nil
}

func (ex *QuplaExecStmt) wrapUpRun(disp *Dispatcher) bool {
	logf(0, "Executed '%v' %v times", ex.GetName(), ex.runResult.called)
	dur := ex.runResult.lastRun.Sub(ex.runResult.start)
	avgdur := int64(dur/time.Millisecond) / int64(ex.runResult.called)
	logf(0, "    eval result:     '%v'. Total duration %v, Average duration %v msec/run",
		utils.TritsToString(ex.runResult.lastResult), dur, avgdur)

	var passed bool
	if ex.isTest {
		logf(0, "    expected result: '%v'", utils.TritsToString(ex.valueExpected))
		if passed = ex.resultIsExpected(ex.runResult.lastResult); passed {
			logf(0, "    test PASSED")
		} else {
			logf(0, "    test FAILED")
		}
	}
	_ = disp.DeleteEnvironment(ex.inEnvironmentName())
	_ = disp.DeleteEnvironment(ex.outEnvironmentName())
	return passed
}
