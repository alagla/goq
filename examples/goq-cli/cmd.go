package main

import (
	"github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/analyzeyaml"
	. "github.com/lunfardo314/goq/cfg"
	"github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/goq/readyaml"
	"github.com/lunfardo314/goq/utils"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
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
//    lexe
//    lexe <filter substring>
//    lfun
//    lfun <filter substring>
//    lenv
//    trace stop
//    trace <filter substring>
//    trace <filter substring> <traceLevel>
//    run
//    run all
//    run <exec idx>
//    run <from exec idx>-<to exec idx>
//    repeat <exec idx> <repeat times>
//    chain
//    chain on|off
//    post <effect decimal> <environment>

func CmdVerbosity(words []string) {
	if len(words) == 2 {
		v, err := strconv.Atoi(words[1])
		if err != nil || v < 0 {
			Logf(0, "must be non-negative integer")
			return
		}
		Config.Verbosity = v
	}
	Logf(0, "Verbosity level is %v", Config.Verbosity)
}

func CmdDir(words []string) {
	var err error
	if len(words) >= 2 {
		err = os.Chdir(words[1])
		if err != nil {
			Logf(0, "%v", err)
		}
	}
	currentDir, err := os.Getwd()
	if err != nil {
		Logf(0, "%v", err)
	} else {
		Logf(0, "Current directory is %v", currentDir)
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
	Logf(0, "Loading QuplaYAML module form file %v", fname)
	moduleYAML, err = NewQuplaModuleFromYAML(fname)
	if err != nil {
		Logf(0, "Error while parsing YAML file: %v", err)
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
		Logf(0, "Module loaded successfully")
	} else {
		Logf(0, "Failed to load module")
		module = nil
	}
	if !succ && len(words) == 2 && words[1] == "exitonfail" {
		os.Exit(1)
	}
}

func CmdSaveModule(words []string) {
	if moduleYAML == nil {
		Logf(0, "Error: module is not loaded")
		return
	}
	fname := testoutdef
	if len(words) == 2 {
		fname = words[1]
	}
	Logf(0, "Writing Qupla module to YAML file %v", fname)

	if err := moduleYAML.WriteToFile(fname); err != nil {
		Logf(0, "%v", err)
	} else {
		Logf(0, "Successfully saved Qupla module to %v", fname)
	}
}

func logExecs(list []*qupla.ExecStmt) {
	for _, ex := range list {
		Logf(0, "   #%v:  %v", ex.GetIdx(), ex.GetSource())
	}
	Logf(0, "Found %v executable statements:", len(list))
}

func CmdLexe(words []string) {
	if moduleYAML == nil {
		Logf(0, "Error: module was not loaded")
		return
	}
	substr := ""
	if len(words) == 2 {
		substr = words[1]
	}
	execs := module.FindExecs(substr)
	logExecs(execs)
}

func CmdInline(words []string) {
	if len(words) == 2 {
		switch words[1] {
		case "on":
			Config.OptimizeInline = true
		default:
			Config.OptimizeInline = false
		}
	}
	if Config.OptimizeInline {
		Logf(0, "Inline call optimization is ON")
	} else {
		Logf(0, "Inline call optimization is OFF")
	}
}

func CmdTrace(words []string) {
	if moduleYAML == nil {
		Logf(0, "Error: module was not loaded")
		return
	}
	if len(words) < 2 {
		Logf(0, "usage: trace stop|<filter substr> [1|2]")
		return
	}
	if words[1] == "stop" {
		module.SetTraceLevel(0, "")
		Logf(0, "all tracing stopped")
		return
	}
	traceLevel := 1
	var err error
	if len(words) == 3 {
		traceLevel, err = strconv.Atoi(words[2])
		if err != nil {
			Logf(0, "wrong command: %v", err)
			return
		}
	}
	funcs := module.SetTraceLevel(traceLevel, words[1])
	Logfuncs(funcs)
	Logf(0, "Set trace level = %v", traceLevel)
}

func Logfuncs(list []*qupla.Function) {
	for _, fun := range list {
		Logf(0, "   %v", fun.Name)
	}
	Logf(0, "Found %v functions:", len(list))
}

func CmdLfun(words []string) {
	if moduleYAML == nil {
		Logf(0, "Error: module was not loaded")
		return
	}
	substr := ""
	if len(words) == 2 {
		substr = words[1]
	}
	funcs := module.FindFuncs(substr)
	Logfuncs(funcs)
}

func CmdPassparam(words []string) {
	if moduleYAML == nil {
		Logf(0, "Error: module was not loaded")
		return
	}
	substr := ""
	if len(words) == 2 {
		substr = words[1]
	}
	funcs := module.FindFuncs(substr)
	lfun := make([]*qupla.Function, 0)
	for _, f := range funcs {
		if f.ZeroInternalSites() {
			lfun = append(lfun, f)
		}
	}
	Logfuncs(lfun)
}

func CmdLenv(words []string) {
	if moduleYAML == nil {
		Logf(0, "Error: module was not loaded")
		return
	}
	for env := range module.Environments {
		Logf(0, "    %v", env)
	}
	Logf(0, "   Total %v environments found", len(module.Environments))
}

func stringIsInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

var currentExecIdx = 0

func CmdRun(words []string) {
	if module == nil {
		Logf(0, "Error: module not loaded")
		return
	}
	switch {
	case len(words) == 1:
		if err := module.RunExecs(svisor, currentExecIdx, currentExecIdx, chainMode); err != nil {
			Logf(0, "%v", err)
			currentExecIdx = 0
			return
		}
		currentExecIdx++

	case len(words) == 2 && words[1] == "all":
		start := time.Now()
		if err := module.RunExecs(svisor, -1, -1, chainMode); err != nil {
			Logf(0, "%v", err)
			return
		}
		currentExecIdx = 0
		Logf(0, "Duration %v", time.Since(start))

	case len(words) == 2 && stringIsInt(words[1]):
		idx, _ := strconv.Atoi(words[1])
		if err := module.RunExecs(svisor, idx, idx, chainMode); err != nil {
			Logf(0, "%v", err)
			currentExecIdx = 0
			return
		}
		currentExecIdx = idx + 1
		return

	case len(words) == 2 && !stringIsInt(words[1]):
		split := strings.Split(words[1], "-")
		if len(split) != 2 {
			Logf(0, "wrong commend")
			return
		}

		fromIdx, _ := strconv.Atoi(split[0])
		toIdx, _ := strconv.Atoi(split[1])
		if err := module.RunExecs(svisor, fromIdx, toIdx, chainMode); err != nil {
			Logf(0, "%v", err)
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
		Logf(0, "chain mode is ON. Executable statements will be linked in a chain of environments")
	} else {
		Logf(0, "chain mode is OFF. Executable statements will be run in unlinked environments")
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
	Logf(0, "Module: %v", m)
	Logf(0, "Memory allocated: %vM", memAllocMB)
	Logf(0, "Number of goroutines: %v", runtime.NumGoroutine())
}

func CmdRepeat(words []string) {
	if len(words) != 3 || !stringIsInt(words[1]) || !stringIsInt(words[2]) {
		Logf(0, "usage: 'repeat <exec idx> <numrepeat>'")
	}
	idx, _ := strconv.Atoi(words[1])
	repeat, _ := strconv.Atoi(words[2])
	if repeat < 1 {
		Logf(0, "wrong number of repeats'")
		return
	}
	if err := module.RunExec(svisor, idx, repeat); err != nil {
		Logf(0, "%v", err)
	}
}

func CmdPost(words []string) {
	if module == nil {
		Logf(0, "Error: module not loaded")
		return
	}
	if len(words) != 3 {
		Logf(0, "Usage: post <effect decimal> <environment>")
		return
	}
	dec, err := strconv.Atoi(words[1])
	if err != nil {
		Logf(0, "Usage: post <effect decimal> <environment>")
		return
	}
	effect := trinary.IntToTrits(int64(dec))
	Logf(0, "Posting effect %v, '%v' to environment '%v'",
		dec, utils.TritsToString(effect), words[2])

	err = svisor.PostEffect(words[2], effect, 0)
	if err != nil {
		Logf(0, "error: %v", err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	svisor.DoOnIdle(func() {
		wg.Done()
	})
	wg.Wait()
}
