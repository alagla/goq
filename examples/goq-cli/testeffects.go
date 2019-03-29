package main

import (
	"github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/cfg"
	"github.com/lunfardo314/goq/supervisor"
	"strconv"
	"time"
)

var postEffects = []struct {
	env       string
	decString string
}{
	{"fibInit", "10"},
}

func postEffectsToSupervisor(disp *supervisor.Supervisor) {
	Logf(0, "-----------------------")
	Logf(0, "Posting test effects to environments")

	var err error
	var val int
	for _, s := range postEffects {
		val, err = strconv.Atoi(s.decString)
		if err != nil {
			Logf(0, "error: %v", err)
			continue
		}
		t := trinary.IntToTrits(int64(val))
		start := time.Now()
		err = disp.PostEffect(s.env, t, 0)
		if err != nil {
			Logf(0, "error while starting quant with value '%v' and the environment '%v': %v",
				s.decString, s.env, err)
		} else {
			Logf(3, "Quant %v <- '%v' was finished in %v",
				s.decString, s.env, time.Since(start))
		}
	}
}
