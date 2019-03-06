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
		affectChan: make(chan Trits),
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

func (env *Environment) Join(entity EntityInterface) error {
	env.Lock()
	defer env.Unlock()
	if env.existsEntity(entity.GetName()) {
		return fmt.Errorf("duplicated entity '%v' attempt to join to '%v'", entity.GetName(), env.name)
	}
	if env.size == 0 {
		env.size = entity.InSize()
	} else {
		if entity.InSize() != env.size {
			return fmt.Errorf("size mismatch between environment '%v', size = %v and joining entity '%v', size = %v",
				env.name, env.size, entity.GetName(), entity.InSize())
		}
	}
	env.joins = append(env.joins, entity)
	return nil
}

func (env *Environment) PostEffect(effect Trits) {
	env.affectChan <- effect
}

func (env *Environment) AffectLoop() {
	for effect := range env.affectChan {
		logf(1, "Value '%v' reached environment '%v'",
			TritsToString(effect), env.name)
		env.processEffect(effect)
	}
}

func (env *Environment) processEffect(effect Trits) {
	env.RLock()
	defer env.RUnlock()
	for _, entity := range env.joins {
		go env.calculateEntity(entity, effect)
	}
}

func (env *Environment) calculateEntity(entity EntityInterface, effect Trits) {
	logf(1, "Value '%v' triggered entity %v in environment '%v'",
		TritsToString(effect), entity.GetName(), env.name)
	// TODO calls function and then affects environment
}
