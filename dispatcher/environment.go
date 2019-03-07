package dispatcher

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/utils"
	"sync"
)

type Environment struct {
	sync.RWMutex
	name       string
	joins      []EntityInterface
	size       int64
	affectChan chan Trits // where all effects are sent
}

func NewEnvironment(name string) *Environment {
	ret := &Environment{
		name:       name,
		joins:      make([]EntityInterface, 0),
		affectChan: make(chan Trits, 1), // buffer to avoid deadlocks
	}
	go ret.AffectLoop()
	return ret
}

func (env *Environment) existsEntity(name string) bool {
	for _, ei := range env.joins {
		if ei.GetName() == name {
			return true
		}
	}
	return false
}

func (env *Environment) checkNewSize(size int64) error {
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
	env.Lock()
	defer env.Unlock()
	if env.existsEntity(entity.GetName()) {
		return fmt.Errorf("duplicated entity '%v' attempt to join to '%v'", entity.GetName(), env.name)
	}
	if err := env.checkNewSize(entity.InSize()); err != nil {
		return fmt.Errorf("error while joining entity '%v' to the environment '%v': %v",
			entity.GetName(), env.name, err)
	}
	env.joins = append(env.joins, entity)
	return nil
}

func (env *Environment) PostEffect(effect Trits) {
	env.affectChan <- effect
}

// loop waits for effect in the environment and then process it
// null result mean nil
func (env *Environment) AffectLoop() {
	for effect := range env.affectChan {
		if effect != nil {
			logf(1, "Value '%v' reached environment '%v'",
				TritsToString(effect), env.name)
			env.processEffect(effect)
		}
	}
}

func (env *Environment) processEffect(effect Trits) {
	env.RLock()
	defer env.RUnlock()
	for _, entity := range env.joins {
		go entity.Invoke(effect) // sync or async?
	}
}
