package entities

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/supervisor"
	"github.com/lunfardo314/goq/utils"
	"time"
)

type printEffectEntityCore struct {
	entity                *Entity
	printEvery            int
	counter               int
	printNotMoreOftenThan time.Duration
	lastOutput            time.Time
}

func NewPrintEffectEntity(supervisor *Supervisor, name string, printEvery int, printNotMoreOftenThan time.Duration) (*Entity, error) {
	core := &printEffectEntityCore{
		printEvery:            printEvery,
		printNotMoreOftenThan: printNotMoreOftenThan,
	}
	ret, err := supervisor.NewEntity(name, 0, 0, core)
	core.entity = ret
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (e *printEffectEntityCore) Call(input Trits, _ Trits) bool {
	if (e.printEvery == 0 || e.counter%e.printEvery == 0) && time.Since(e.lastOutput) >= e.printNotMoreOftenThan {
		trits := utils.TritsToString(input)
		fmt.Printf("%9d %v.%v: '%v' (len=%v)", e.counter, e.entity.Supervisor.Name, e.entity.Name, trits, len(trits))
		e.lastOutput = time.Now()
	}
	e.counter++
	return true // does not affect any environment, does not produce any result
}
