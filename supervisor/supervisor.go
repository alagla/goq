package supervisor

import (
	"fmt"
	"github.com/Workiva/go-datastructures/queue"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/cfg"
	"github.com/lunfardo314/goq/utils"
	"time"
)

// TODO dispose Supervisor ???

func (sv *Supervisor) getEnvironment_(name string) *environment {
	env, ok := sv.environments[name]
	if !ok {
		return nil
	}
	return env
}

func (sv *Supervisor) getOrCreateEnvironment(name string) *environment {
	ret := sv.getEnvironment_(name)
	if ret != nil {
		return ret
	}
	sv.environments[name] = newEnvironment(sv, name)
	return sv.environments[name]
}

func (sv *Supervisor) createEnvironment(name string) error {
	if sv.getEnvironment_(name) != nil {
		return fmt.Errorf("environment '%v' already exists", name)
	}
	sv.environments[name] = newEnvironment(sv, name)
	return nil
}

func (sv *Supervisor) resetCallCounters() {
	for _, env := range sv.environments {
		for _, joinInfo := range env.joins {
			joinInfo.count = 0
		}
	}
}

// perform a quant in sync way
// effect is sent to the channel of the environment
// and its waited until all waves settled down

func (sv *Supervisor) doQuant(env *environment, effect Trits) {
	// reset call counters for joined entites (needed to handle join limits)
	sv.resetCallCounters()

	sv.quantWG.Add(1)
	env.effectChan <- effect
	sv.quantWG.Wait()
}

type quantMsg struct {
	envName          string
	environment      *environment
	effect           Trits
	doNotStartBefore int64
}

// posts (internal) the effect to the main queue
// TODO not much use of the 'external' parameter

func (sv *Supervisor) postEffect(envName string, env *environment, effect Trits, delay int, external bool) error {
	n := envName
	if env != nil {
		n = env.name
	}
	LogDefer(5, func() {
		res := utils.TritsToString(effect)
		reslen := len(res)
		if reslen > 100 {
			res = res[:100] + "..."
		}
		Logf(5, "posted effect to Supervisor ->'%v', delay=%v external=%v: '%s' (len=%d)",
			n, delay, external, res, reslen)
	})

	return sv.queue.Put(&quantMsg{
		envName:          envName,
		environment:      env,
		effect:           effect,
		doNotStartBefore: sv.GetQuantCount() + int64(delay),
	})
}

// main Supervisor input loop
// one message read from the queue means one quant
// Supervisor is locked during processing of the quant
// It starts in locked state and the this loop doing 100 millisecond idle loops
// while queue is empty

func (sv *Supervisor) supervisorInputLoop() {
	var tmpItems []interface{}
	var msg *quantMsg
	var err error
	var env *environment

	for {
		tmpItems, err = sv.queue.Poll(1, 100*time.Millisecond)
		if err != nil {
			if err == queue.ErrTimeout {
				sv.setIdle(true) // unlock
				continue
			} else {
				panic(err)
			}
		}
		sv.setIdle(false) // lock

		msg = tmpItems[0].(*quantMsg)

		if msg.doNotStartBefore > sv.GetQuantCount() {
			// the effect is delayed: put it back to queue
			_ = sv.queue.Put(tmpItems[0])
			sv.incQuantCount()
			continue
		}

		// if environment is not given by pointer, find it by Name
		if msg.environment == nil {
			env = sv.getEnvironment_(msg.envName)
		} else {
			env = msg.environment
		}
		if env == nil || env.invalid {
			// environment can be invalid also in case it was deleted
			// from the Supervisor while some pending effects were still in the queue
			Logf(5, "supervisorInputLoop: can't find valid environment '%v'", msg.envName)
			continue
		}
		// perform a quant
		sv.doQuant(env, msg.effect)
		sv.incQuantCount()
	}
}

// Increases quant counter.
// Supervisor level quant counter is needed to handle affect delays and join limits for entities

func (sv *Supervisor) incQuantCount() {
	sv.quantCountMutex.Lock()
	defer sv.quantCountMutex.Unlock()
	sv.quantCount++
}

// setIdle is called from within Supervisor input loop
// Supervisor is idle only if there's no messages in the inpout queue and therefore
// it is not processing a quant.
// Within the quant Supervisor is locked for any calls from outside which change configuration of
// environments and entities

func (sv *Supervisor) setIdle(idle bool) {
	switch {
	case !sv.idle && idle:
		// is locked here
		// toggle to 'idle' state
		sv.idle = true
		sv.accessLock.release()
	case sv.idle && !idle:
		// is released here
		// toggle to 'busy' state
		sv.accessLock.acquire(-1)
		sv.idle = false
	}
}
