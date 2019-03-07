package dispatcher

import (
	"fmt"
	"github.com/lunfardo314/goq/cfg"
)

func logf(minVerbosity int, format string, args ...interface{}) {
	if cfg.Config.Verbosity >= minVerbosity {
		fmt.Printf(format+"\n", args...)
	}
}
