package dispatcher

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/utils"
)

type Environment struct {
	//sync.RWMutex
	dispatcher *Dispatcher
	name       string
	joins      []EntityInterface
	size       int64
	effectChan chan struct{} // signals about changed value are sent
	value      Trits         // valid only between waves
}

func NewEnvironment(disp *Dispatcher, name string) *Environment {
	ret := &Environment{
		dispatcher: disp,
		name:       name,
		joins:      make([]EntityInterface, 0),
		effectChan: make(chan struct{}), // buffer to avoid deadlocks
	}
	go ret.environmentListenToEffectsLoop()
	return ret
}

func (env *Environment) Size() int64 {
	return env.size
}

func (env *Environment) GetName() string {
	return env.name
}

func (env *Environment) stop() {
	close(env.effectChan)
}

func (env *Environment) existsEntity_(name string) bool {
	for _, ei := range env.joins {
		if ei.Name() == name {
			return true
		}
	}
	return false
}

func (env *Environment) checkNewSize_(size int64) error {
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

func (env *Environment) Join(entity EntityInterface) error {
	if env.existsEntity_(entity.Name()) {
		return fmt.Errorf("duplicated entity '%v' attempt to join to '%v'", entity.Name(), env.name)
	}
	if err := env.checkNewSize_(entity.InSize()); err != nil {
		return fmt.Errorf("error while joining entity '%v' to the environment '%v': %v",
			entity.Name(), env.name, err)
	}
	env.joins = append(env.joins, entity)
	return nil
}

func (env *Environment) PostEffect(effect Trits) {
	if effect != nil {
		dec, _ := TritsToBigInt(effect)
		logf(2, "Environment '%v' <- '%v' (%v)", env.name, TritsToString(effect), dec)
	} else {
		logf(2, "Environment '%v' <- 'null'", env.name)
	}
	env.setNewValue(effect)
	env.dispatcher.quantWG.Add(len(env.joins))
	logf(4, "---------------- ADD %v (env '%v')", len(env.joins), env.name)

	env.effectChan <- struct{}{}

}

// loop waits for effect in the environment and then process it
// null result mean nil
func (env *Environment) environmentListenToEffectsLoop() {
	logf(4, "environmentListenToEffectsLoop STARTED for environment '%v'", env.name)
	defer logf(4, "environmentListenToEffectsLoop STOPPED for environment '%v'", env.name)

	for range env.effectChan {
		if len(env.joins) == 0 {
			continue
		}
		// in wave-by-wave mode here waits
		env.dispatcher.waveStopWG.Wait()
		env.dispatcher.waveReleaseWG.Wait()
		//  here starts new wave
		prev := env.setNewValue(nil) // environment value becomes invalid
		for _, entity := range env.joins {
			entity.Invoke(prev) // async
		}
	}
	// if the input channel (affect) is closed,
	// we have to close all join channels to stop listening routines
	for _, entity := range env.joins {
		entity.Stop()
	}
}

func (env *Environment) setNewValue(val Trits) Trits {
	//logf(3, "------ env '%v' set value to '%v'", env.name, TritsToString(val))
	saveValue := env.value
	env.value = val
	return saveValue
}

func (env *Environment) GetValue() Trits {
	return env.value
}
