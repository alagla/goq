package cfg

import (
	"fmt"
	"strings"
)

func Logf(minVerbosity int, format string, args ...interface{}) {
	if Config.Verbosity >= minVerbosity {
		prefix := fmt.Sprintf("%2d  %s", minVerbosity, strings.Repeat(" ", minVerbosity))
		fmt.Printf(prefix+format+"\n", args...)
	}
}

func LogDefer(minVerbosity int, fun func()) {
	if Config.Verbosity >= minVerbosity {
		fun()
	}
}
