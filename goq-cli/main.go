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
	fmt.Printf(strings.Repeat(" ", minVerbosity)+format+"\n", args...)
}
func execBatch(cmdlist []string) {
	for _, cmdline := range cmdlist {
		executor(cmdline)
		time.Sleep(100 * time.Millisecond)
	}
}

var startupCmd = []string{
	"load",
	//"run",
	//"run 0-10",
	"wave 0",
	"wave status",
	"run",
	"wave status",
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
		execBatch(startupCmd)
		logf(0, "interrupt from keyboard...")
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
