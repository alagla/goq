package main

import (
	"github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/dispatcher"
	"strconv"
	"time"
)

var postEffects = []struct {
	env       string
	decString string
}{
	{"fibInit", "10"},
}

func postEffectsToDispatcher(disp *dispatcher.Dispatcher) {
	logf(0, "-----------------------")
	logf(0, "Posting test effects to environments")

	var err error
	var val int
	for _, s := range postEffects {
		val, err = strconv.Atoi(s.decString)
		if err != nil {
			logf(0, "error: %v", err)
			continue
		}
		t := trinary.IntToTrits(int64(val))
		start := time.Now()
		err = disp.DoQuant(s.env, t)
		if err != nil {
			logf(0, "error while starting quant with value '%v' and the environment '%v': %v",
				s.decString, s.env, err)
		} else {
			logf(3, "Quant %v <- '%v' was finished in %v",
				s.decString, s.env, time.Since(start))
		}
	}
}
