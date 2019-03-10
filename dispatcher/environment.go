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
	effectChan chan Trits // where all effects are sent
}

func NewEnvironment(disp *Dispatcher, name string) *Environment {
	ret := &Environment{
		dispatcher: disp,
		name:       name,
		joins:      make([]EntityInterface, 0),
		effectChan: make(chan Trits), // buffer to avoid deadlocks
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
	env.dispatcher.quantWG.Add(len(env.joins))
	logf(3, "---------------- ADD %v (env '%v')", len(env.joins), env.name)

	env.effectChan <- effect

}

// loop waits for effect in the environment and then process it
// null result mean nil
func (env *Environment) environmentListenToEffectsLoop() {
	logf(3, "environmentListenToEffectsLoop STARTED for environment '%v'", env.name)
	defer logf(3, "environmentListenToEffectsLoop STOPPED for environment '%v'", env.name)

	for effect := range env.effectChan {
		//  TODO in debug mode here wait for wave
		for _, entity := range env.joins {
			entity.Invoke(effect) // async
		}
	}
	// if the input channel (affect) is closed,
	// we have to close all join channels to stop listening routines
	for _, entity := range env.joins {
		entity.Stop()
	}
}
