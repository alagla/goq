package main

import (
	"github.com/lunfardo314/goq/analyzeyaml"
	"github.com/lunfardo314/goq/cfg"
	"github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/goq/readyaml"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Commands:
//    verb
//    verb <level>
//    runtime
//    dir
//    dir <directory name>
//    load
//    load <module yaml file>
//    save
//    save <file to save as yaml>
//    list
//    list <substring>
//    run
//    run all
//    run <exec idx>
//    run <from exec idx>-<to exec idx>
//    repeat <exec idx> <repeat times>
//    chain
//    chain on|off

func CmdVerbosity(words []string) {
	if len(words) == 2 {
		v, err := strconv.Atoi(words[1])
		if err != nil || v < 0 {
			logf(0, "must be non-negative integer")
			return
		}
		cfg.Config.Verbosity = v
	}
	logf(0, "Verbosity level is %v", cfg.Config.Verbosity)
}

func CmdDir(words []string) {
	var err error
	if len(words) >= 2 {
		err = os.Chdir(words[1])
		if err != nil {
			logf(0, "%v", err)
		}
	}
	currentDir, err := os.Getwd()
	if err != nil {
		logf(0, "%v", err)
	} else {
		logf(0, "Current directory is %v", currentDir)
	}
}

const fnamedefault = "QuplaTests.yml"
const testoutdef = "echotest.yml"

func CmdLoadModule(words []string) {
	var err error

	fname := fnamedefault
	if len(words) >= 2 && words[1] != "exitonfail" {
		fname = words[1]
	}
	logf(0, "Loading QuplaYAML module form file %v", fname)
	moduleYAML, err = NewQuplaModuleFromYAML(fname)
	if err != nil {
		logf(0, "Error while parsing YAML file: %v", err)
		moduleYAML = nil
		return
	}
	var succ bool
	module, succ = analyzeyaml.AnalyzeQuplaModule(fname, moduleYAML)
	module.PrintStats()
	if succ {
		succ = module.AttachToSupervisor(svisor)
	}
	if succ {
		logf(0, "Module loaded successfully")
	} else {
		logf(0, "Failed to load module")
		module = nil
	}
	if !succ && len(words) == 2 && words[1] == "exitonfail" {
		os.Exit(1)
	}
}

func CmdSaveModule(words []string) {
	if moduleYAML == nil {
		logf(0, "Error: module is not loaded")
		return
	}
	fname := testoutdef
	if len(words) == 2 {
		fname = words[1]
	}
	logf(0, "Writing Qupla module to YAML file %v", fname)

	if err := moduleYAML.WriteToFile(fname); err != nil {
		logf(0, "%v", err)
	} else {
		logf(0, "Successfully saved Qupla module to %v", fname)
	}
}

func logExecs(list []*qupla.ExecStmt) {
	logf(0, "Found %v executable statements:", len(list))
	for _, ex := range list {
		logf(0, "   #%v:  %v", ex.GetIdx(), ex.GetSource())
	}
}

func CmdList(words []string) {
	if moduleYAML == nil {
		logf(0, "Error: module was not loaded")
		return
	}
	substr := ""
	if len(words) == 2 {
		substr = words[1]
	}
	execs := module.FindExecs(substr)
	logExecs(execs)
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
		if err := module.RunExecs(svisor, currentExecIdx, currentExecIdx, chainMode); err != nil {
			logf(0, "%v", err)
			currentExecIdx = 0
			return
		}
		currentExecIdx++

	case len(words) == 2 && words[1] == "all":
		start := time.Now()
		if err := module.RunExecs(svisor, -1, -1, chainMode); err != nil {
			logf(0, "%v", err)
			return
		}
		currentExecIdx = 0
		logf(0, "Duration %v", time.Since(start))

	case len(words) == 2 && stringIsInt(words[1]):
		idx, _ := strconv.Atoi(words[1])
		if err := module.RunExecs(svisor, idx, idx, chainMode); err != nil {
			logf(0, "%v", err)
			currentExecIdx = 0
			return
		}
		currentExecIdx = idx + 1
		return

	case len(words) == 2 && !stringIsInt(words[1]):
		split := strings.Split(words[1], "-")
		if len(split) != 2 {
			logf(0, "wrong commend")
			return
		}

		fromIdx, _ := strconv.Atoi(split[0])
		toIdx, _ := strconv.Atoi(split[1])
		if err := module.RunExecs(svisor, fromIdx, toIdx, chainMode); err != nil {
			logf(0, "%v", err)
			currentExecIdx = 0
			return
		}
		currentExecIdx = toIdx + 1
		return
	}
}

var chainMode = false

func CmdChain(words []string) {
	if len(words) == 2 {
		chainMode = words[1] == "on"
	}
	if chainMode {
		logf(0, "chain mode is ON. Executable statements will be linked in a chain of environments")
	} else {
		logf(0, "chain mode is OFF. Executable statements will be run in unlinked environments")
	}
}

func CmdRuntime(_ []string) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	memAllocMB := math.Round(100*(float64(mem.Alloc/1024)/1024)) / 100
	m := "not loaded"
	if module != nil {
		m = module.GetName()
	}
	logf(0, "Module: %v", m)
	logf(0, "Memory allocated: %vM", memAllocMB)
	logf(0, "Number of goroutines: %v", runtime.NumGoroutine())
}

func CmdRepeat(words []string) {
	if len(words) != 3 || !stringIsInt(words[1]) || !stringIsInt(words[2]) {
		logf(0, "usage: 'repeat <exec idx> <numrepeat>'")
	}
	idx, _ := strconv.Atoi(words[1])
	repeat, _ := strconv.Atoi(words[2])
	if repeat < 1 {
		logf(0, "wrong number of repeats'")
		return
	}
	if err := module.RunExec(svisor, idx, repeat); err != nil {
		logf(0, "%v", err)
	}
}
