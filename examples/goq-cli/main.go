package main

import (
	"flag"
	"github.com/c-bata/go-prompt"
	. "github.com/lunfardo314/goq/cfg"
	"github.com/lunfardo314/goq/qupla"
	. "github.com/lunfardo314/goq/readyaml"
	"github.com/lunfardo314/goq/supervisor"
	"time"
)

var svisor = supervisor.NewSupervisor("sv1", 1*time.Second)
var moduleYAML *QuplaModuleYAML
var module *qupla.QuplaModule

func execBatch(cmdlist []string) {
	for _, cmdline := range cmdlist {
		Logf(0, "Batch command: '%v'", cmdline)
		executor(cmdline)
	}
}

var startupCmd = []string{
	//"load modules/Fibonacci.yml exitonfail",
	//"load modules/Curl.yml exitonfail",
	"load modules/Examples.yml exitonfail",
	//"load modules/Qupla.yml exitonfail",
	//"trace listMap",
	//"run 113-115",
	//"run 0-1",
	"forabra",
	// "run all",
	//"run 0-10",
}

func main() {
	Logf(0, "Welcome to GOQ-CLI: a simple Qubic Supervisor in Go Command Line Interface ver %v", Config.Version)
	Logf(0, "Now is %v", time.Now())
	executor("dir")
	executor("verb")

	pnocli := flag.Bool("nocli", false, "bypass CLI")
	flag.Parse()
	if *pnocli {
		Logf(0, "Bypass CLI. Load module and run it.")
		execBatch(startupCmd)
		Logf(0, "sleep loop...")
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
