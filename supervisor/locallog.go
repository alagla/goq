package supervisor

import (
	"fmt"
	"github.com/lunfardo314/goq/cfg"
	"strings"
)

func logf(minVerbosity int, format string, args ...interface{}) {
	if cfg.Config.Verbosity >= minVerbosity {
		prefix := fmt.Sprintf("%2d  %s", minVerbosity, strings.Repeat(" ", minVerbosity))
		fmt.Printf(prefix+format+"\n", args...)
	}
}
