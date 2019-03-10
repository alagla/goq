package qupla

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/dispatcher"
	. "github.com/lunfardo314/goq/quplayaml"
	"github.com/lunfardo314/goq/utils"
	"time"
)

type QuplaExecStmt struct {
	BaseEntity
	QuplaExprBase
	isTest        bool
	isFloat       bool // needed for float comparison
	funcExpr      *QuplaFuncExpr
	valueExpected Trits
	module        *QuplaModule
	num           int
	duration      time.Duration
	passed        bool
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

// create temporary environment
// create two temporary entities
//   - one for the eval function expression itself, affect the environment
//   - another for the reaction of the function result: in case of eval ir prints result, in case of test it checks test
//   - post effect to the environment
func (ex *QuplaExecStmt) Execute(disp *Dispatcher) (bool, error) {
	envInName := "ENV_IN$$" + ex.GetSource() + "$$"
	envOutName := "ENV_OUT$$" + ex.GetSource() + "$$"

	disp.GetOrCreateEnvironment_(envInName)
	//if err := disp.SetEnvironmentSize(envInName, 1); err != nil{
	//	return false, err
	//}
	exprEntity := ex.newEvalEntity(disp)
	if _, err := disp.Join(envInName, exprEntity); err != nil {
		return false, err
	}

	disp.GetOrCreateEnvironment_(envOutName)
	if _, err := disp.Affect(envOutName, exprEntity); err != nil {
		return false, err
	}

	resultEntity := ex.newEvalResultEntity(disp)
	if _, err := disp.Join(envOutName, resultEntity); err != nil {
		return false, err
	}
	var t = Trits{0}
	err := disp.DoQuant(envInName, t)

	_ = disp.DeleteEnvironment(envInName)
	_ = disp.DeleteEnvironment(envOutName)

	return ex.passed, err
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

func (ex *QuplaExecStmt) newEvalEntity(disp *Dispatcher) *BaseEntity {
	return NewBaseEntity(disp, ex.funcExpr.GetSource(), 0, ex.funcExpr.Size(), &execEvalCallable{ex})
}

type execEvalResultCallable struct {
	exec *QuplaExecStmt
}

func (ec *execEvalResultCallable) Call(result Trits, _ Trits) bool {
	logf(0, "Executing '%v' ", ec.exec.source)
	logf(0, "    eval result:     '%v'. Duration %v", utils.TritsToString(result), ec.exec.duration)

	if ec.exec.isTest {
		logf(0, "    expected result: '%v'", utils.TritsToString(ec.exec.valueExpected))

		if ec.exec.passed, _ = TritsEqual(result, ec.exec.valueExpected); ec.exec.passed {
			logf(0, "    test PASSED")
		} else {
			logf(0, "    test FAILED")
		}
	}
	return true
}

func (ex *QuplaExecStmt) newEvalResultEntity(disp *Dispatcher) *BaseEntity {
	ec := &execEvalResultCallable{exec: ex}
	return NewBaseEntity(disp, ex.funcExpr.GetSource(), ex.funcExpr.Size(), 0, ec)
}
