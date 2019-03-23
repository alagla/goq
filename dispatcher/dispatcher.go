package dispatcher

import (
	"fmt"
	"github.com/Workiva/go-datastructures/queue"
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/utils"
	"sync"
	"time"
)

// TODO size checks when join/affect. Can be with different sizes
// TODO dispose dispatcher
// TODO rename to 'supervisor' ??

func (disp *Dispatcher) incQuantCount() {
	disp.quantCountMutex.Lock()
	defer disp.quantCountMutex.Unlock()
	disp.quantCount++
}

func (disp *Dispatcher) getEnvironment_(name string) *environment {
	env, ok := disp.environments[name]
	if !ok {
		return nil
	}
	return env
}

func (disp *Dispatcher) getOrCreateEnvironment(name string) *environment {
	ret := disp.getEnvironment_(name)
	if ret != nil {
		return ret
	}
	disp.environments[name] = newEnvironment(disp, name)
	return disp.environments[name]
}

func (disp *Dispatcher) createEnvironment(name string) error {
	if disp.getEnvironment_(name) != nil {
		return fmt.Errorf("environment '%v' already exists", name)
	}
	disp.environments[name] = newEnvironment(disp, name)
	return nil
}

func (disp *Dispatcher) resetCallCounters() {
	for _, env := range disp.environments {
		for _, joinInfo := range env.joins {
			joinInfo.count = 0
		}
	}
}

func (disp *Dispatcher) quantStart(env *environment, effect Trits, onQuantFinish func()) error {
	var err error
	if effect, err = env.adjustEffect(effect); err != nil {
		return err
	}

	disp.resetCallCounters()
	disp.quantWG.Add(1)
	env.effectChan <- effect
	go func() {
		env.dispatcher.quantWG.Wait()
		if onQuantFinish != nil {
			onQuantFinish()
		}
	}()
	return nil
}

type quantMsg struct {
	envName          string
	environment      *environment
	effect           Trits
	doNotStartBefore int64
}

func (disp *Dispatcher) postEffect(envName string, env *environment, effect Trits, delay int, external bool) error {
	dec, _ := utils.TritsToBigInt(effect)
	n := envName
	if env != nil {
		n = env.name
	}
	logf(5, "posted effect '%v' (%v) to dispatcher, environment '%v', delay %v",
		utils.TritsToString(effect), dec, n, delay)

	return disp.queue.Put(&quantMsg{
		envName:          envName,
		environment:      env,
		effect:           effect,
		doNotStartBefore: disp.GetQuantCount() + int64(delay),
	})
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

		if msg.doNotStartBefore > disp.GetQuantCount() {
			// delayed: put it back to queue
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
		_ = disp.quantStart(env, msg.effect, func() {
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
		disp.accessLock.release()
	case disp.idle && !idle:
		// is released here
		disp.accessLock.acquire(-1)
		disp.idle = false
	}
}
