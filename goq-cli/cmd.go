package main

import (
	"github.com/lunfardo314/goq/cfg"
	"github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/quplayaml/quplayaml"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// Commands:
// 		load [<module file>]
//      save [<module file>]
// 		verb [<level>]
//      list [execs <search string>]
//      runtime
//      run  [idxFrom [- idxTo]]

func CmdVerbosity(words []string) {
	if len(words) == 1 {
		logf(0, "current verbosity level is %v", cfg.Config.Verbosity)
		return
	}
	v, err := strconv.Atoi(words[1])
	if err != nil || v < 0 {
		logf(0, "must be non-negative integers")
		return
	}
	cfg.Config.Verbosity = v
	logf(0, "verbosity was set to %v", cfg.Config.Verbosity)

}

const fname = "C:/Users/evaldas/Documents/proj/Java/github.com/qupla/src/main/resources/Qupla.yml"
const testout = "C:/Users/evaldas/Documents/proj/site_data/tmp/echotest.yml"

func CmdLoadModule(words []string) {
	var err error

	logf(0, "Loading module form file %v", fname)
	moduleYAML, err = NewQuplaModuleFromYAML(fname)
	if err != nil {
		logf(0, "Error while parsing YAML file: %v", err)
		moduleYAML = nil
		return
	}
	var succ bool
	module, succ = qupla.AnalyzeQuplaModule("single_module", moduleYAML, &qupla.ExpressionFactoryFromYAML{})
	module.PrintStats()
	if succ {
		succ = module.AttachToDispatcher(dispatcherInstance)
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
		if err := module.RunExecs(dispatcherInstance, currentExecIdx, currentExecIdx, chainMode); err != nil {
			logf(0, "%v", err)
			currentExecIdx = 0
			return
		}
		currentExecIdx++

	case len(words) == 2 && words[1] == "all":
		if err := module.RunExecs(dispatcherInstance, -1, -1, chainMode); err != nil {
			logf(0, "%v", err)
			return
		}
		currentExecIdx = 0

	case len(words) == 2 && stringIsInt(words[1]):
		idx, _ := strconv.Atoi(words[1])
		if err := module.RunExecs(dispatcherInstance, idx, idx, chainMode); err != nil {
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
		if err := module.RunExecs(dispatcherInstance, fromIdx, toIdx, chainMode); err != nil {
			logf(0, "%v", err)
			currentExecIdx = 0
			return
		}
		currentExecIdx = toIdx + 1
		return
	}
}

var chainMode = false

func CmdMode(words []string) {
	if len(words) == 1 {
		if chainMode {
			logf(0, "chain mode is ON")
		} else {
			logf(0, "chain mode is OFF")
		}
		return
	}
	if len(words) == 2 && words[1] == "on" {
		chainMode = true
		logf(0, "chain mode is ON")
	} else {
		chainMode = false
		logf(0, "chain mode is OFF")
	}
}

//var waveModeON = false
//
//func CmdWave(words []string) {
//	if module == nil {
//		logf(0, "error: module not loaded")
//		return
//	}
//	if len(words) == 1 {
//		if waveModeON {
//			logf(0, "Wave mode is ON")
//		} else {
//			logf(0, "Wave mode is OFF")
//		}
//		return
//	}
//
//	switch words[1] {
//	case "on":
//		waveModeON = true
//		logf(0, "wave mode is ON")
//	case "off":
//		waveModeON = false
//		logf(0, "wave mode is OFF")
//
//	case "start":
//		if len(words) != 4 {
//			logf(0, "error: wrong command")
//			return
//		}
//		if dispatcherInstance.IsWaveMode() {
//			logf(0, "   quant is already running")
//			return
//		}
//		effectDec, err := strconv.Atoi(words[2])
//		if err != nil {
//			logf(0, "error: effect must be decimal integers")
//			return
//		}
//		effectTrits := trinary.IntToTrits(int64(effectDec))
//		envName := words[3]
//		err = dispatcherInstance.QuantStart(envName, effectTrits, waveModeON, func() {
//			logf(0, " -------- quant finished")
//			waveModeON = false
//		})
//		if err != nil {
//			logf(0, "%v", err)
//		}
//	case "next":
//		if !dispatcherInstance.IsWaveMode() {
//			logf(0, "error: quant wasn't started: can't continue with the wave")
//		}
//		if err := dispatcherInstance.WaveNext(); err != nil {
//			logf(0, "error: %v", err)
//			return
//		}
//		time.Sleep(100 * time.Millisecond)
//	case "run":
//		if !dispatcherInstance.IsWaveMode() {
//			logf(0, "   can't continue: quant not running")
//			return
//		}
//		if err := dispatcherInstance.WaveRun(); err != nil {
//			logf(0, "%v", err)
//		}
//	case "status":
//		listValues()
//	}
//}
//
//func listValues() {
//	vDict := dispatcherInstance.WaveValues()
//	if len(vDict) == 0 {
//		logf(0, "   wave is empty")
//	} else {
//		names := make([]string, 0, len(vDict))
//		for n := range vDict {
//			names = append(names, n)
//		}
//		sort.Strings(names)
//		logf(0, "   environment values:")
//		for _, name := range names {
//			logf(0, "    %v: '%v'", name, utils.TritsToString(vDict[name]))
//		}
//	}
//}

func CmdRuntime(_ []string) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	memAllocMB := math.Round(100*(float64(mem.Alloc/1024)/1024)) / 100
	logf(0, "Memory allocated: %vM", memAllocMB)
	logf(0, "Number of goroutines: %v", runtime.NumGoroutine())
}

//func CmdStatus(_ []string) {
//	eInfo := dispatcherInstance.EnvironmentInfo()
//	logf(0, "Dispatcher status:")
//	logf(1, "Found %v environments", len(eInfo))
//
//	names := make([]string, 0, len(eInfo))
//	for n := range eInfo {
//		names = append(names, n)
//	}
//	sort.Strings(names)
//
//	for _, name := range names {
//		envStatus := eInfo[name]
//		logf(2, "%v (size = %v):", name, envStatus.Size)
//		entStr := ""
//		for _, entName := range envStatus.AffectedBy {
//			if entStr != "" {
//				entStr += ", "
//			}
//			entStr += entName
//		}
//		logf(4, "Affected by %v entities: %v", len(envStatus.AffectedBy), entStr)
//
//		entStr = ""
//		for _, entName := range envStatus.JoinedEntities {
//			if entStr != "" {
//				entStr += ", "
//			}
//			entStr += entName
//		}
//		logf(4, "Joined %v entities: %v", len(envStatus.JoinedEntities), entStr)
//
//	}
//}
