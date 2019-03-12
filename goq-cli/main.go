package main

import (
	"flag"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/lunfardo314/goq/cfg"
	"github.com/lunfardo314/goq/dispatcher"
	"github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/quplayaml/quplayaml"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func executor(in string) {
	words := strings.Split(in, " ")
	if len(words) == 0 || words[0] == "" {
		return
	}
	logf(2, "----- Command '%v'", words[0])
	start := time.Now()
	defer logf(2, "----- Command '%v'. Duration %v", words[0], time.Since(start))

	switch words[0] {
	case "exit", "quit":
		logf(0, "Bye!")
		os.Exit(0)
	case "verb":
		CmdVerbosity(words)
	case "load":
		CmdLoadModule(words)
	case "save":
		CmdSaveModule(words)
	case "run":
		CmdRunExecs(words)
	case "functions":
		logf(0, "not implemented yet")
	case "runtime":
		CmdRuntime(words)
	case "module":
		logf(0, "not implemented yet")
	default:
		logf(0, "unknown command")
	}
}

func completer(in prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{}
	//if in.GetWordBeforeCursor() == ""{
	//	return []prompt.Suggest{}
	//}
	//
	//switch strings.Trim(in.GetWordBeforeCursorWithSpace(), " "){
	//case "verb":
	//	return []prompt.Suggest{
	//		{Text: "0", Description: "normal"},
	//		{Text: "1", Description: "verbose"},
	//		{Text: "2", Description: "debug"},
	//		{Text: "3", Description: "trace"},
	//	}
	//}
	//s := []prompt.Suggest{
	//	{Text: "exit", Description: "Exit goq-cli"},
	//	{Text: "verb", Description: "Change verbosity level to 0,1,2,3"},
	//	{Text: "load", Description: "Load Qupla module"},
	//	{Text: "module", Description: "Current module info"},
	//	{Text: "functions", Description: "list functions of the current module"},
	//	{Text: "help", Description: "list goq dispatcher commands"},
	//}
	//return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func main() {
	logf(0, "goq-cli: GOQ (Qubic Dispatcher in Go) Command Line Interface ver %v", cfg.Config.Version)
	logf(0, "Now is %v", time.Now())
	logf(0, "Verbosity is %v", cfg.Config.Verbosity)
	logf(0, "Use TAB to select suggestion")

	pnocli := flag.Bool("nocli", false, "bypass CLI")
	flag.Parse()
	if *pnocli {
		logf(0, "Bypass CLI. Load module and run it.")
		CmdLoadModule(nil)
		CmdRunExecs(nil)
	} else {
		p := prompt.New(
			executor,
			completer,
			prompt.OptionPrefixTextColor(prompt.LightGray),
			prompt.OptionPrefix(">>> "),
		)
		p.Run()
	}
}

func logf(minVerbosity int, format string, args ...interface{}) {
	if cfg.Config.Verbosity < minVerbosity {
		return
	}
	fmt.Printf(format+"\n", args...)
}

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

var moduleYAML *QuplaModuleYAML
var module *qupla.QuplaModule

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

func CmdRunExecs(_ []string) {
	if module == nil {
		logf(0, "Error: module not loaded")
		return
	}
	disp := dispatcher.NewDispatcher()
	module.AttachToDispatcher(disp)
	module.Execute(disp)
	postEffectsToDispatcher(disp)
}

func CmdRuntime(_ []string) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	memAllocMB := math.Round(100*(float64(mem.Alloc/1024)/1024)) / 100
	logf(0, "Memory allocated: %vM", memAllocMB)
	logf(0, "Number of goroutines: %v", runtime.NumGoroutine())
}
