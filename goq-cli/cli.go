package main

import (
	"github.com/c-bata/go-prompt"
	"os"
	"strings"
)

func executor(in string) {
	logf(2, "goq-cli cmd: '%v'", in)
	words := strings.Split(in, " ")
	if len(words) == 0 || words[0] == "" {
		return
	}
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
		CmdRun(words)
	case "wave":
		CmdWave(words)
	case "list":
		CmdList(words)
	case "functions":
		logf(0, "not implemented yet")
	case "runtime":
		CmdRuntime(words)
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
