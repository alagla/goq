package qupla

import (
	"fmt"
	"github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/supervisor"
	"github.com/lunfardo314/goq/utils"
	"time"
)

func (module *QuplaModule) AttachExecs(disp *supervisor.Supervisor, fromIdx int, toIdx int, chain bool) []*ExecStmt {
	if len(module.execs) == 0 {
		logf(0, "No executables to execute")
		return nil
	}
	if fromIdx < 0 || fromIdx >= len(module.execs) {
		fromIdx = 0
	}
	if toIdx < 0 || toIdx >= len(module.execs) {
		toIdx = len(module.execs) - 1
	}
	if fromIdx < 0 || fromIdx > toIdx {
		logf(0, "Wrong range of indices: from %v to %v", fromIdx, toIdx)
		return nil
	}
	ret := make([]*ExecStmt, 0)

	var exec *ExecStmt
	var err error
	var prev *ExecStmt
	for idx := fromIdx; idx <= toIdx; idx++ {
		exec = module.execs[idx]
		if err = exec.attach(disp, prev); err != nil {
			logf(0, "can't attach executable '%v'", exec.GetName())
		} else {
			ret = append(ret, exec)
			if chain {
				prev = exec
			}
		}
	}
	return ret
}

func (module *QuplaModule) detachExecs(disp *supervisor.Supervisor, execs []*ExecStmt) error {
	for _, e := range execs {
		if err := e.detach(disp); err != nil {
			return err
		}
	}
	return nil
}

func (module *QuplaModule) runAttachedExecs(disp *supervisor.Supervisor, execs []*ExecStmt, chain bool) error {
	if len(execs) == 0 {
		return fmt.Errorf("No executables to execute")
	}
	effect := trinary.Trits{1}
	if chain {
		return disp.PostEffect(execs[0].evalEnvironmentName(), effect, 0)
	}
	for _, exec := range execs {
		if err := disp.PostEffect(exec.evalEnvironmentName(), effect, 0); err != nil {
			logf(0, "%v", err)
		}
	}
	return nil
}

func (module *QuplaModule) RunExec(disp *supervisor.Supervisor, idx int, repeat int) error {
	if module.ExecByIdx(idx) == nil {
		return fmt.Errorf("can't find executable statement #%v", idx)
	}
	attachedExecs := module.AttachExecs(disp, idx, idx, false)
	if len(attachedExecs) != 1 {
		return fmt.Errorf("inconsistency")
	}
	logf(0, "Running %v times: '%v'", repeat, attachedExecs[0].GetName())

	start := time.Now()
	for i := 0; i < repeat; i++ {
		if err := disp.PostEffect(attachedExecs[0].evalEnvironmentName(), trinary.Trits{0}, 0); err != nil {
			return err
		}
	}
	var duration time.Duration
	onFinish := func() {
		_ = module.detachExecs(disp, attachedExecs)
		duration = time.Since(start)
		logf(0, "Stop")
	}

	for !disp.DoIfIdle(5*time.Second, onFinish) {
	}

	reportRunResults(attachedExecs, duration)
	return nil
}

func (module *QuplaModule) RunExecs(disp *supervisor.Supervisor, fromIdx int, toIdx int, chain bool) error {
	attachedExecs := module.AttachExecs(disp, fromIdx, toIdx, chain)

	logf(0, "Running executable statements with indices between %v and %v", fromIdx, toIdx)
	logf(0, "   total in the module: %v", len(module.execs))
	logf(0, "   running: %v", len(attachedExecs))
	cmode := "OFF"
	if chain {
		cmode = "ON"
	}
	logf(0, "Chain mode is %v", cmode)
	start := time.Now()
	if err := module.runAttachedExecs(disp, attachedExecs, chain); err != nil {
		return err
	}

	var duration time.Duration
	onFinish := func() {
		_ = module.detachExecs(disp, attachedExecs)
		duration = time.Since(start)
		logf(0, "Stop")
	}

	for !disp.DoIfIdle(5*time.Second, onFinish) {
	}

	reportRunResults(attachedExecs, duration)
	return nil
}

func reportRunResults(execs []*ExecStmt, duration time.Duration) {
	logf(0, "Run summary:")
	logf(0, "   Executed %v executable statements in %v", len(execs), duration)
	numTest := 0
	numEvals := 0
	numTestsPassed := 0
	var summ *runSummary
	for _, ex := range execs {
		summ = ex.GetRunResults()
		if summ.isTest {
			if summ.testPassed {
				logf(2, "evaluated %v %v time{s}. Avg duration: %v msec", ex.GetName(), summ.numRun, summ.avgDuration)
				logf(2, "     test PASSED")
			} else {
				logf(0, "evaluated %v %v time{s}. Avg duration: %v msec", ex.GetName(), summ.numRun, summ.avgDuration)
				logf(0, "     test FAILED")
				logf(0, "     Result %v != expected %v",
					utils.ReprTrits(summ.lastResult), utils.ReprTrits(ex.expected))
			}
			numTest++
		} else {
			logf(2, "evaluated %v %v time{s}. Avg duration: %v msec", ex.GetName(), summ.numRun, summ.avgDuration)
			bi, _ := utils.TritsToBigInt(summ.lastResult)
			logf(2, "   result: %v, '%v'", bi, utils.TritsToString(summ.lastResult))
			numEvals++
		}
		if summ.testPassed {
			numTestsPassed++
		}

	}
	logf(0, "Total evals: %v", numEvals)
	logf(0, "Total test: %v", numTest)
	percPassed := "n/a"
	if numTest != 0 {
		percPassed = fmt.Sprintf("%v%%", (numTestsPassed*100)/numTest)
	}
	logf(0, "Tests passed: %v (%v)", numTestsPassed, percPassed)
}
