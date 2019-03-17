package dispatcher

import (
	"fmt"
	"github.com/lunfardo314/goq/cfg"
	"strings"
)

func logf(minVerbosity int, format string, args ...interface{}) {
	if cfg.Config.Verbosity >= minVerbosity {
		fmt.Printf(strings.Repeat(" ", minVerbosity)+format+"\n", args...)
	}
}
