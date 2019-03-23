package analyzeyaml

import (
	"fmt"
	"github.com/lunfardo314/goq/cfg"
	"github.com/op/go-logging"
	"strings"
)

var (
	localLog   *logging.Logger
	localTrace bool
)

const (
	logInfoPrefix  = ""
	logDebugPrefix = ""
	logErrorPrefix = "ERRO "
)

func errorf(format string, args ...interface{}) {
	if localLog != nil {
		localLog.Errorf(format, args...)
	} else {
		fmt.Printf(logErrorPrefix+format+"\n", args...)
	}
}

func debugf(format string, args ...interface{}) {
	if localLog != nil {
		localLog.Debugf(format, args...)
	} else {
		fmt.Printf(logDebugPrefix+format+"\n", args...)
	}
}

func tracef(format string, args ...interface{}) {
	if !localTrace {
		return
	}
	logf(3, format, args...)
}

func infof(format string, args ...interface{}) {
	if localLog != nil {
		localLog.Infof(format, args...)
	} else {
		fmt.Printf(logInfoPrefix+format+"\n", args...)
	}
}

func logf(minVerbosity int, format string, args ...interface{}) {
	if cfg.Config.Verbosity >= minVerbosity {
		prefix := fmt.Sprintf("%2d  %s", minVerbosity, strings.Repeat(" ", minVerbosity))
		fmt.Printf(prefix+format+"\n", args...)
	}
}
