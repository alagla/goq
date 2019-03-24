package supervisor

import (
	"fmt"
	"github.com/Workiva/go-datastructures/queue"
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/utils"
	"sync"
	"time"
)

// TODO size checks when join/affect. Can be with different sizes
// TODO dispose supervisor

func (sv *Supervisor) incQuantCount() {
	sv.quantCountMutex.Lock()
	defer sv.quantCountMutex.Unlock()
	sv.quantCount++
}

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

func (sv *Supervisor) quantStart(env *environment, effect Trits, onQuantFinish func()) error {
	var err error
	if effect, err = env.adjustEffect(effect); err != nil {
		return err
	}

	sv.resetCallCounters()
	sv.quantWG.Add(1)
	env.effectChan <- effect
	go func() {
		env.supervisor.quantWG.Wait()
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

func (sv *Supervisor) postEffect(envName string, env *environment, effect Trits, delay int, external bool) error {
	dec, _ := utils.TritsToBigInt(effect)
	n := envName
	if env != nil {
		n = env.name
	}
	logf(5, "posted effect '%v' (%v) to supervisor, environment '%v', delay %v",
		utils.TritsToString(effect), dec, n, delay)

	return sv.queue.Put(&quantMsg{
		envName:          envName,
		environment:      env,
		effect:           effect,
		doNotStartBefore: sv.GetQuantCount() + int64(delay),
	})
}

func (sv *Supervisor) supervisorInputLoop() {
	var tmpItems []interface{}
	var msg *quantMsg
	var quantWG sync.WaitGroup
	var err error
	var env *environment

	for {
		tmpItems, err = sv.queue.Poll(1, 100*time.Millisecond)
		if err != nil {
			if err == queue.ErrTimeout {
				sv.setIdle(true)
				continue
			} else {
				panic(err)
			}
		}
		sv.setIdle(false)

		msg = tmpItems[0].(*quantMsg)
		//logf(5, "supervisorInputLoop: received %+v", msg)

		if msg.doNotStartBefore > sv.GetQuantCount() {
			// delayed: put it back to queue
			_ = sv.queue.Put(tmpItems[0])
			sv.incQuantCount()
			continue
		}

		if msg.environment == nil {
			env = sv.getEnvironment_(msg.envName)
		} else {
			env = msg.environment
		}
		if env == nil || env.invalid {
			logf(5, "supervisorInputLoop: can't find valid environment '%v'", msg.envName)
			continue
		}
		quantWG.Add(1)
		_ = sv.quantStart(env, msg.effect, func() {
			sv.incQuantCount()
			quantWG.Done()
		})
		quantWG.Wait()
	}
}

func (sv *Supervisor) setIdle(idle bool) {
	switch {
	case !sv.idle && idle:
		// is locked here
		sv.idle = true
		sv.accessLock.release()
	case sv.idle && !idle:
		// is released here
		sv.accessLock.acquire(-1)
		sv.idle = false
	}
}
