package main

import (
	"flag"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/lunfardo314/goq/cfg"
	"github.com/lunfardo314/goq/dispatcher"
	"github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/quplayaml/quplayaml"
	"strings"
	"time"
)

var dispatcherInstance = dispatcher.NewDispatcher(1 * time.Second)
var moduleYAML *QuplaModuleYAML
var module *qupla.QuplaModule

func logf(minVerbosity int, format string, args ...interface{}) {
	if cfg.Config.Verbosity < minVerbosity {
		return
	}
	prefix := fmt.Sprintf("%2d  %s", minVerbosity, strings.Repeat(" ", minVerbosity))
	fmt.Printf(prefix+format+"\n", args...)
}
func execBatch(cmdlist []string) {
	for _, cmdline := range cmdlist {
		logf(0, "Batch command: '%v'", cmdline)
		executor(cmdline)
	}
}

var startupCmd = []string{
	"load exitonfail",
	"run 0-1",
	// "run all",
	//"run 0-10",
}

func main() {
	logf(0, "Welcome to GOQ-CLI: a simple Qubic Dispatcher in Go Command Line Interface ver %v", cfg.Config.Version)
	logf(0, "Now is %v", time.Now())
	//logf(0, "Verbosity level is %v. Change with `verb <level>`", cfg.Config.Verbosity)
	//logf(0, "Use TAB to select suggestion")
	executor("dir")
	executor("verb")

	pnocli := flag.Bool("nocli", false, "bypass CLI")
	flag.Parse()
	if *pnocli {
		logf(0, "Bypass CLI. Load module and run it.")
		execBatch(startupCmd)
		logf(0, "sleep loop...")
		for {
			time.Sleep(10 * time.Second)
		}
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
