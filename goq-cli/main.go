package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/lunfardo314/goq/cfg"
	"os"
	"strconv"
	"strings"
	"time"
)

func executor(in string) {
	words := strings.Split(in, " ")
	if len(words) == 0 {
		return
	}
	logf(2, "Your input: %v\n", words)
	switch words[0] {
	case "exit", "quit":
		logf(0, "Bye!")
		os.Exit(0)
	case "verb":
		CmdVerbosity(words)
	case "load":
		logf(0, "not implemented yet")
	case "functions":
		logf(0, "not implemented yet")
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

	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefixTextColor(prompt.LightGray),
		prompt.OptionPrefix(">>> "),
	)
	p.Run()
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
