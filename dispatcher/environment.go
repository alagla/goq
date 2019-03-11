package dispatcher

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/utils"
	"sync"
)

type environment struct {
	sync.RWMutex
	dispatcher *Dispatcher
	name       string
	joins      []*Entity
	size       int64
	effectChan chan struct{} // signals about changed value are sent
	value      Trits         // valid only between waves
}

func NewEnvironment(disp *Dispatcher, name string) *environment {
	ret := &environment{
		dispatcher: disp,
		name:       name,
		joins:      make([]*Entity, 0),
		effectChan: make(chan struct{}),
	}
	go ret.environmentListenToEffectsLoop()
	return ret
}

//func (env *environment) Size() int64 {
//	return env.size
//}
//
func (env *environment) GetName() string {
	return env.name
}

func (env *environment) stop() {
	close(env.effectChan)
}

func (env *environment) existsEntity_(name string) bool {
	for _, ei := range env.joins {
		if ei.name == name {
			return true
		}
	}
	return false
}

func (env *environment) checkNewSize_(size int64) error {
	if env.size != 0 {
		if env.size != size {
			return fmt.Errorf("size mismatch in environment '%v'. Must be %v",
				env.name, env.size)
		}
	} else {
		env.size = size
	}
	return nil
}

func (env *environment) join(entity *Entity) error {
	if env.existsEntity_(entity.name) {
		return fmt.Errorf("duplicated entity '%v' attempt to join to '%v'", entity.name, env.name)
	}
	if err := env.checkNewSize_(entity.InSize()); err != nil {
		return fmt.Errorf("error while joining entity '%v' to the environment '%v': %v",
			entity.name, env.name, err)
	}
	env.joins = append(env.joins, entity)
	return nil
}

func (env *environment) postEffect(effect Trits) {
	if effect != nil {
		dec, _ := TritsToBigInt(effect)
		logf(2, "environment '%v' <- '%v' (%v)", env.name, TritsToString(effect), dec)
	} else {
		logf(2, "environment '%v' <- 'null'", env.name)
	}
	env.setNewValue(effect)
	env.dispatcher.quantWG.Add(len(env.joins))
	logf(4, "---------------- ADD %v (env '%v')", len(env.joins), env.name)

	env.effectChan <- struct{}{}
}

// loop waits for effect in the environment and then process it
// null result mean nil
func (env *environment) environmentListenToEffectsLoop() {
	logf(4, "environmentListenToEffectsLoop STARTED for environment '%v'", env.name)
	defer logf(4, "environmentListenToEffectsLoop STOPPED for environment '%v'", env.name)

	for range env.effectChan {
		if len(env.joins) == 0 {
			continue
		}
		// in wave-by-wave mode here waits
		env.dispatcher.waveCatchWG.Wait()
		env.dispatcher.waveReleaseWG.Wait()
		//  here starts new wave
		prev := env.setNewValue(nil) // environment value becomes invalid
		for _, entity := range env.joins {
			entity.invoke(prev) // async
		}
	}
	// if the input channel (affect) is closed,
	// we have to close all join channels to stop listening routines
	for _, entity := range env.joins {
		entity.stop()
	}
}

func (env *environment) setNewValue(val Trits) Trits {
	env.Lock()
	defer env.Unlock()
	//logf(3, "------ env '%v' set value to '%v'", env.name, TritsToString(val))
	saveValue := env.value
	env.value = val
	return saveValue
}

func (env *environment) GetValue() Trits {
	env.RLock()
	defer env.RUnlock()
	return env.value
}
