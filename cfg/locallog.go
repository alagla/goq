package cfg

import (
	"fmt"
	"sort"
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

func LogStats(minVerbosity int, stats map[string]int) {
	tmpKeys := make([]string, 0)
	for k := range stats {
		tmpKeys = append(tmpKeys, k)
	}
	sort.Strings(tmpKeys)
	for _, key := range tmpKeys {
		Logf(minVerbosity, "      %v: %v", key, stats[key])
	}
}

func AddStats(dst, add map[string]int) {
	for k := range add {
		if _, ok := dst[k]; !ok {
			dst[k] = 0
		}
		dst[k] += add[k]
	}
}
