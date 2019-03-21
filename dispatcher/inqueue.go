package dispatcher

import (
	"github.com/Workiva/go-datastructures/queue"
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/utils"
	"sync"
	"time"
)

type quantMsg struct {
	envName          string
	environment      *environment
	effect           Trits
	doNotStartBefore uint64
}

func (disp *Dispatcher) postEffect(envName string, env *environment, effect Trits, delay int, external bool) error {
	dec, _ := utils.TritsToBigInt(effect)
	n := envName
	if env != nil {
		n = env.GetName()
	}
	logf(5, "posted effect '%v' (%v) to dispatcher, environment '%v', delay %v",
		utils.TritsToString(effect), dec, n, delay)

	return disp.queue.Put(&quantMsg{
		envName:          envName,
		environment:      env,
		effect:           effect,
		doNotStartBefore: disp.getQuantCount() + uint64(delay),
	})
}

func (disp *Dispatcher) PostEffect(envName string, effect Trits, delay int) error {
	return disp.postEffect(envName, nil, effect, delay, true)
}

func (disp *Dispatcher) dispatcherInputLoop() {
	var tmpItems []interface{}
	var msg *quantMsg
	var quantWG sync.WaitGroup
	var err error
	var env *environment

	for {
		tmpItems, err = disp.queue.Poll(1, 100*time.Millisecond)
		if err != nil {
			if err == queue.ErrTimeout {
				disp.setIdle(true)
				continue
			} else {
				panic(err)
			}
		}
		disp.setIdle(false)

		msg = tmpItems[0].(*quantMsg)
		//logf(5, "dispatcherInputLoop: received %+v", msg)

		if msg.doNotStartBefore > disp.getQuantCount() {
			// delyed: put it back to queue
			_ = disp.queue.Put(msg)
			disp.incQuantCount()
			continue
		}

		if msg.environment == nil {
			env = disp.getEnvironment_(msg.envName)
		} else {
			env = msg.environment
		}
		if env == nil || env.invalid {
			logf(5, "dispatcherInputLoop: can't find valid environment '%v'", msg.envName)
			continue
		}
		quantWG.Add(1)
		_ = disp.quantStart(env, msg.effect, false, func() {
			disp.incQuantCount()
			quantWG.Done()
		})
		quantWG.Wait()
	}
}

func (disp *Dispatcher) setIdle(idle bool) {
	switch {
	case !disp.idle && idle:
		// is locked here
		disp.idle = true
		disp.environmentLock.Release()
	case disp.idle && !idle:
		// is released here
		disp.environmentLock.Acquire(-1)
		disp.idle = false
	}
}

// calls callback if dispatcher becomes idle within 'timeout'
// callback will be called after release of the semaphore.
// the callback must take care about locking the dispatcher if needed
func (disp *Dispatcher) CallIfIdle(timeout time.Duration, callback func()) bool {
	if !disp.environmentLock.Acquire(timeout) {
		return false
	}
	disp.environmentLock.Release()
	callback()
	return true
}
