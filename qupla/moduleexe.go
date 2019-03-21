package qupla

import (
	"fmt"
	"github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/cfg"
	"github.com/lunfardo314/goq/dispatcher"
	"time"
)

func (module *QuplaModule) attachExecs(disp *dispatcher.Dispatcher, fromIdx int, toIdx int, chain bool) []*QuplaExecStmt {
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
	switch {
	case cfg.Config.ExecEvals && cfg.Config.ExecTests:
		logf(2, "Attaching to dispatcher: evals and tests")
	case cfg.Config.ExecEvals && !cfg.Config.ExecTests:
		logf(2, "Attaching to dispatcher: evals only")
	case !cfg.Config.ExecEvals && cfg.Config.ExecTests:
		logf(2, "Attaching to dispatcher: Etests only")
	case !cfg.Config.ExecEvals && !cfg.Config.ExecTests:
		logf(2, "Attaching to dispatcher: wrong config values, assume tests only")
	}
	if fromIdx < 0 && toIdx < 0 {
		logf(2, "Index range: ALL (total %v)", len(module.execs))
	} else {
		logf(2, "Index range: %v - %v", fromIdx, toIdx)
	}

	ret := make([]*QuplaExecStmt, 0)

	var exec *QuplaExecStmt
	var err error
	var prev *QuplaExecStmt
	for idx := fromIdx; idx <= toIdx; idx++ {
		exec = module.execs[idx]
		if exec.HasState() {
			logf(2, "skipped '%v'", exec.GetName())
			continue
		}
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

func (module *QuplaModule) detachExecs(disp *dispatcher.Dispatcher, execs []*QuplaExecStmt) error {
	for _, e := range execs {
		if err := e.detach(disp); err != nil {
			return err
		}
	}
	return nil
}

func (module *QuplaModule) runAttachedExecs(disp *dispatcher.Dispatcher, execs []*QuplaExecStmt, chain bool) error {
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

func (module *QuplaModule) RunExecs(disp *dispatcher.Dispatcher, fromIdx int, toIdx int, chain bool) error {
	attachedExecs := module.attachExecs(disp, fromIdx, toIdx, chain)
	logf(0, "Total executables in the module: %v", len(module.execs))
	logf(0, "Skipped: %v", len(module.execs)-len(attachedExecs))
	logf(0, "Start running executables: %v", len(attachedExecs))
	start := time.Now()
	if err := module.runAttachedExecs(disp, attachedExecs, chain); err != nil {
		return err
	}

	onFinish := func() {
		_ = module.detachExecs(disp, attachedExecs)
		logf(0, "Finished running execs, chain mode = %v. Duration %v", chain, time.Since(start))
	}

	for !disp.CallIfIdle(5*time.Second, onFinish) {
	}

	reportRunResults(attachedExecs)
	return nil
}

func reportRunResults(execs []*QuplaExecStmt) {
	logf(0, "Run summary:")
	logf(0, "Executed %v executables", len(execs))
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
			}
			numTest++
		} else {
			logf(2, "evaluated %v %v time{s}. Avg duration: %v msec", ex.GetName(), summ.numRun, summ.avgDuration)
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
