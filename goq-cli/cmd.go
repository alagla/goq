package main

import (
	"github.com/lunfardo314/goq/cfg"
	"github.com/lunfardo314/goq/qupla"
	"github.com/lunfardo314/goq/utils"
	. "github.com/lunfardo314/quplayaml/quplayaml"
	"math"
	"runtime"
	"strconv"
	"strings"
)

func CmdVerbosity(words []string) {
	if len(words) == 1 {
		logf(0, "current verbosity level is %v", cfg.Config.Verbosity)
		return
	}
	if len(words) != 2 {
		logf(0, "usage: verb [0|1|2|3]")
	}
	var v int
	v, err := strconv.Atoi(words[1])
	if err != nil || v < 0 || v > 2 {
		logf(0, "usage: verb [0|1|2|3]")
		return
	}
	cfg.Config.Verbosity = v
	logf(0, "verbosity was set to %v", cfg.Config.Verbosity)

}

const fname = "C:/Users/evaldas/Documents/proj/Java/github.com/qupla/src/main/resources/Qupla.yml"
const testout = "C:/Users/evaldas/Documents/proj/site_data/tmp/echotest.yml"

func CmdLoadModule(_ []string) {
	var err error

	logf(0, "Loading module form file %v", fname)
	moduleYAML, err = NewQuplaModuleFromYAML(fname)
	if err != nil {
		logf(0, "Error while parsing YAML file: %v", err)
		moduleYAML = nil
		return
	}
	logf(0, "Module '%v' loaded successfully", moduleYAML.Name)
	logf(0, "Analyzing module")

	var succ bool
	module, succ = qupla.AnalyzeQuplaModule("single_module", moduleYAML, &qupla.ExpressionFactoryFromYAML{})
	module.PrintStats()
	if succ {
		module.AttachToDispatcher(dispatcherInstance)
		logf(0, "Module analyzed succesfully")
	} else {
		logf(0, "Failed to analyze module")
		module = nil
	}
}

func CmdSaveModule(_ []string) {
	if moduleYAML == nil {
		logf(0, "Error: module was not loaded")
		return
	}
	logf(0, "Writing Qupla module to YAML file %v", testout)

	if err := moduleYAML.WriteToFile(testout); err != nil {
		logf(0, "Error occured: %v", err)
	} else {
		logf(0, "Succesfully saved Qupla module")
	}
}

func logExecs(list []*qupla.QuplaExecStmt) {
	logf(0, "Found %v executables:", len(list))
	for _, ex := range list {
		logf(0, "   #%v:  %v", ex.GetIdx(), ex.GetSource())
	}
}

func CmdList(words []string) {
	if moduleYAML == nil {
		logf(0, "Error: module was not loaded")
		return
	}
	target := "execs"
	if len(words) >= 2 {
		target = words[1]
	}
	substr := ""
	if len(words) >= 3 {
		substr = words[2]
	}
	switch target {
	case "execs":
		execs := module.FindExecs(substr)
		logExecs(execs)
		return
	}
}

func stringIsInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

var currentExecIdx = 0

func CmdRun(words []string) {
	if module == nil {
		logf(0, "Error: module not loaded")
		return
	}
	switch {
	case len(words) == 1:
		if dispatcherInstance.IsWaveMode() {
			if err := dispatcherInstance.WaveRun(); err != nil {
				logf(0, "%v", err)
			}
		} else {
			module.Execute(dispatcherInstance, currentExecIdx, currentExecIdx)
		}
		currentExecIdx++

	case len(words) == 2 && words[1] == "all":
		if dispatcherInstance.IsWaveMode() {
			logf(0, "Already running #%v", currentExecIdx)
			return
		}
		// run all executables
		module.Execute(dispatcherInstance, currentExecIdx, -1)

	case len(words) == 2 && stringIsInt(words[1]):
		// run specific executable in quant mode
		if dispatcherInstance.IsWaveMode() {
			logf(0, "Already running #%v", currentExecIdx)
			return
		}
		idx, _ := strconv.Atoi(words[1])
		module.Execute(dispatcherInstance, idx, idx)
		currentExecIdx = idx + 1
		return

	case len(words) >= 2 && !stringIsInt(words[1]):
		if dispatcherInstance.IsWaveMode() {
			logf(0, "Already running #%v", currentExecIdx)
			return
		}
		nn := strings.Split(words[1], "-")
		if len(nn) != 2 {
			return
		}
		var nfrom, nto int
		var err error
		nfrom, err = strconv.Atoi(nn[0])
		if err != nil {
			logf(0, "%v", err)
			return
		}
		nto, err = strconv.Atoi(nn[1])
		if err != nil {
			logf(0, "%v", err)
			return
		}
		if nfrom > nto {
			logf(0, "wrong index range")
			return
		}

		module.Execute(dispatcherInstance, nfrom, nto)
		currentExecIdx = nto + 1
		return

	case len(words) == 3 && stringIsInt(words[1]) && stringIsInt(words[2]):
		// run specific executable in quant mode
		idx, _ := strconv.Atoi(words[1])
		exec := module.ExecByIdx(idx)
		if exec == nil {
			logf(0, "Can't find executable #%v", idx)
			return
		}
		num, _ := strconv.Atoi(words[2])
		_, err := exec.ExecuteMulti(dispatcherInstance, num)
		if err != nil {
			logf(0, "Error: %v", err)
		}
		return
	}
}

func CmdWave(words []string) {
	if module == nil {
		logf(0, "Error: module not loaded")
		return
	}
	if len(words) != 2 {
		logf(0, "Wrong commend")
		return

	}
	switch words[1] {
	case "next":
		if !dispatcherInstance.IsWaveMode() {
			logf(0, "quant wasn't started: can't continue with the wave")
		}
		if err := dispatcherInstance.WaveNext(); err != nil {
			logf(0, "%v", err)
			return
		}
	case "status":
		logf(0, "Wave mode = %v", dispatcherInstance.IsWaveMode())
		listValues()

	case "run":
		if err := dispatcherInstance.WaveRun(); err != nil {
			logf(0, "%v", err)
		}

	default:
		idx, err := strconv.Atoi(words[1])
		if err != nil {
			logf(0, "Wrong command: %v", err)
			return
		}
		exec := module.ExecByIdx(idx)
		if exec == nil {
			logf(0, "Can't find executable #%v", idx)
			return
		}
		if dispatcherInstance.IsWaveMode() {
			logf(0, "use 'wave next' or 'wave cancel' commands to continue")
			return
		}
		if err := exec.StartWave(dispatcherInstance); err != nil {
			logf(0, "%v", err)
		}
	}
}

func listValues() {
	logf(0, "Not nil values:")
	vDict := dispatcherInstance.WaveValues()
	for name, val := range vDict {
		logf(0, "%v: '%v'", name, utils.TritsToString(val))
	}
}

func CmdRuntime(_ []string) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	memAllocMB := math.Round(100*(float64(mem.Alloc/1024)/1024)) / 100
	logf(0, "Memory allocated: %vM", memAllocMB)
	logf(0, "Number of goroutines: %v", runtime.NumGoroutine())
}
