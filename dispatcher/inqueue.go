package dispatcher

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/utils"
	"sync"
)

type quantMsg struct {
	env              *environment
	effect           Trits
	doNotStartBefore uint64
}

func (disp *Dispatcher) postEffect(env *environment, effect Trits, delay int) error {
	dec, _ := utils.TritsToBigInt(effect)
	logf(5, "posted effect '%v' (%v) to dispatcher, environment '%v', delay %v",
		utils.TritsToString(effect), dec, env.GetName(), delay)

	return disp.queue.Put(&quantMsg{
		env:              env,
		effect:           effect,
		doNotStartBefore: disp.getQuantCount() + uint64(delay),
	})
}

func (disp *Dispatcher) PostEffect(envName string, effect Trits, delay int) error {
	if !disp.generalLock.Acquire(disp.timeout) {
		return fmt.Errorf("request lock timeout: can't create environment")
	}
	defer disp.generalLock.Release()

	env := disp.getEnvironment_(envName)
	if env == nil || env.invalid {
		return fmt.Errorf("can't find environment '%v'", envName)
	}
	return disp.postEffect(env, effect, delay)
}

func (disp *Dispatcher) dispatcherInputLoop() {
	var tmpItems []interface{}
	var msg *quantMsg
	var err error
	var wg sync.WaitGroup

	for {
		tmpItems, err = disp.queue.Get(1)
		if err != nil {
			panic(err)
		}
		msg = tmpItems[0].(*quantMsg)
		if msg.doNotStartBefore > disp.getQuantCount() {
			// put it back to queue
			if err = disp.queue.Put(msg); err != nil {
				panic(err)
			}
			disp.incQuantCount()
			continue
		}

		wg.Add(1)
		err = disp.quantStart(msg.env, msg.effect, false, func() {
			disp.incQuantCount()
			wg.Done()
		})
		wg.Wait()
	}

}
