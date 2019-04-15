package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/supervisor"
	"time"
)

type ExecStmt struct {
	ExpressionBase
	isTest     bool
	isFloat    bool // needed for float comparison
	expected   Trits
	Expr       ExpressionInterface // *FunctionExpr
	module     *QuplaModule
	idx        int
	evalEntity *Entity
}

func NewExecStmt(src string, expr ExpressionInterface, isTest, isFloat bool, expected Trits, module *QuplaModule) *ExecStmt {
	return &ExecStmt{
		ExpressionBase: NewExpressionBase(src),
		isTest:         isTest,
		isFloat:        isFloat,
		expected:       expected,
		Expr:           expr,
		module:         module,
	}
}

func (ex *ExecStmt) GetName() string {
	return fmt.Sprintf("#%v '%v'", ex.idx, ex.GetSource())
}

func (ex *ExecStmt) GetIdx() int {
	return ex.idx
}

func (ex *ExecStmt) HasState() bool {
	return ex.Expr.HasState()
}

func (ex *ExecStmt) evalEnvironmentName() string {
	return fmt.Sprintf("$%v_IN", ex.GetIdx())
}

func (ex *ExecStmt) attach(disp *Supervisor, prev *ExecStmt) error {
	var err error
	if ex.evalEntity, err = ex.newEvalEntity(disp); err != nil {
		return err
	}
	envJoin := map[string]int{ex.evalEnvironmentName(): 5}
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

func (ex *ExecStmt) detach(disp *Supervisor) error {
	return disp.DeleteEnvironment(ex.evalEnvironmentName())
}

func (ex *ExecStmt) Run(disp *Supervisor, repeat int) error {
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

func (ex *ExecStmt) resultIsExpected(result Trits) bool {
	passed, _ := TritsEqual(result, ex.expected)
	if passed {
		return true
	}
	if len(result) != len(ex.expected) {
		return false
	}
	if !ex.isFloat {
		return false
	}
	dif0 := result[0] - ex.expected[0]
	if dif0 < 0 {
		dif0 = -dif0
	}
	if dif0 > 1 {
		return false
	}
	if len(result) == 1 {
		return true
	}
	passed, _ = TritsEqual(result[1:], ex.expected[1:])
	return passed

}

// mock entities to run executables on dispatcher

type execEvalCore struct {
	exec              *ExecStmt
	numRun            int
	numTestPassed     int
	totalDurationMsec uint64
	lastResult        Trits
}

func (ec *execEvalCore) Call(_ Trits, result Trits) bool {
	start := unixMsNow()
	null := ec.exec.Expr.Eval(nil, result)
	ec.numRun++
	ec.totalDurationMsec += unixMsNow() - start
	ec.lastResult = result
	if ec.exec.isTest {
		if ec.exec.resultIsExpected(result) {
			ec.numTestPassed++
		}
	}
	return null
}

func (ex *ExecStmt) newEvalEntity(sv *Supervisor) (*Entity, error) {
	name := fmt.Sprintf("#%v-EVAL", ex.idx)
	core := &execEvalCore{exec: ex}
	return sv.NewEntity(name, 0, ex.Expr.Size(), core)
}

type runSummary struct {
	isTest      bool
	testPassed  bool
	numRun      int
	lastResult  Trits
	avgDuration uint64
}

func (ex *ExecStmt) GetRunResults() *runSummary {
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
