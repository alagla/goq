package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/dispatcher"
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
	evalEntity    *Entity
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

func (ex *QuplaExecStmt) evalEnvironmentName() string {
	return fmt.Sprintf("$%v_IN", ex.GetIdx())
}

func (ex *QuplaExecStmt) attach(disp *Dispatcher, prev *QuplaExecStmt) error {
	ex.evalEntity = ex.newEvalEntity(disp)
	envJoin := map[string]int{ex.evalEnvironmentName(): 1}
	if err := disp.Attach(ex.evalEntity, envJoin, nil); err != nil {
		return err
	}
	if prev != nil {
		// chain mode: result of the previous affect input of the next
		envAffect := map[string]int{ex.evalEnvironmentName(): 0}
		if err := disp.Attach(prev.evalEntity, nil, envAffect); err != nil {
			return err
		}
	}
	return nil
}

func (ex *QuplaExecStmt) detach(disp *Dispatcher) error {
	return disp.DeleteEnvironment(ex.evalEnvironmentName())
}

func (ex *QuplaExecStmt) Run(disp *Dispatcher, repeat int) error {
	if repeat < 1 {
		return fmt.Errorf("'repeat' parameter must be >1")
	}
	var effect = Trits{0}
	envInName := ex.evalEnvironmentName()
	for i := 0; i < repeat; i++ {
		if err := disp.PostEffect(envInName, effect, 0); err != nil {
			return err
		}
	}
	return nil
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

// mock entities to run executables on dispatcher

type execEvalCore struct {
	exec              *QuplaExecStmt
	numRun            int
	numTestPassed     int
	totalDurationMsec uint64
	lastResult        Trits
}

func (ec *execEvalCore) Call(_ Trits, res Trits) bool {
	start := unixMsNow()
	null := ec.exec.module.processor.Eval(ec.exec.funcExpr, res)
	ec.numRun++
	ec.totalDurationMsec += unixMsNow() - start
	ec.lastResult = res
	if ec.exec.isTest && ec.exec.resultIsExpected(res) {
		ec.numTestPassed++
	}
	return null
}

func (ex *QuplaExecStmt) newEvalEntity(disp *Dispatcher) *Entity {
	name := fmt.Sprintf("#%v-EVAL_%v", ex.idx, ex.funcExpr.GetSource())
	core := &execEvalCore{exec: ex}
	return disp.NewEntity(EntityOpts{
		Name:    name,
		InSize:  0,
		OutSize: ex.funcExpr.Size(),
		Core:    core,
	})
}

type runSummary struct {
	isTest      bool
	testPassed  bool
	numRun      int
	lastResult  Trits
	avgDuration uint64
}

func (ex *QuplaExecStmt) GetRunResults() *runSummary {
	core := ex.evalEntity.GetCore().(*execEvalCore)
	var dur uint64
	if core.numRun != 0 {
		dur = core.totalDurationMsec / uint64(core.numRun)
	}
	return &runSummary{
		isTest:      ex.isTest,
		testPassed:  ex.isTest && core.numTestPassed == core.numRun,
		numRun:      core.numRun,
		lastResult:  core.lastResult,
		avgDuration: dur,
	}
}

func unixMsNow() uint64 {
	return uint64(time.Now().UnixNano()) / uint64(time.Millisecond)
}
