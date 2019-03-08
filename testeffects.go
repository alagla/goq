package main

import (
	"github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/dispatcher"
	"github.com/lunfardo314/goq/utils"
	"strconv"
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
		size, ok := disp.GetEnvironmentInfo(s.env)
		if !ok {
			logf(0, "can't find environment '%v'")
			continue
		}
		if len(t) > int(size) {
			logf(0, "Trit vector '%v' is too long for the environment '%v'",
				utils.TritsToString(t), s.env)
			continue
		}
		t = trinary.PadTrits(t, int(size))
		err = disp.PostEffect(s.env, t)
		if err != nil {
			logf(0, "error while posting value '%v' to the environment '%v': %v",
				s.decString, s.env, err)
		} else {
			logf(0, "trit vector '%v' ('%v') was posted to the environment '%v'",
				utils.TritsToString(t), s.decString, s.env)
		}
	}

}
