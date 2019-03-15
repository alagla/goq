package dispatcher

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/utils"
)

type environment struct {
	dispatcher *Dispatcher
	name       string
	invalid    bool
	joins      []*Entity
	affects    []*Entity
	size       int64
	effectChan chan Trits
	value      Trits // valid only between waves
}

func NewEnvironment(disp *Dispatcher, name string) *environment {
	ret := &environment{
		dispatcher: disp,
		name:       name,
		joins:      make([]*Entity, 0),
		affects:    make([]*Entity, 0),
		effectChan: make(chan Trits),
	}
	go ret.environmentLoop()
	return ret
}

//func (env *environment) Size() int64 {
//	return env.size
//}
//
func (env *environment) GetName() string {
	return env.name
}

func (env *environment) existsEntity_(name string) bool {
	for _, ei := range env.joins {
		if ei.name == name {
			return true
		}
	}
	return false
}

func (env *environment) checkNewSize(size int64) bool {
	if env.size != 0 {
		if env.size != size {
			return false
		}
	} else {
		env.size = size
	}
	return true
}

func (env *environment) join(entity *Entity) error {
	if !env.checkNewSize(entity.InSize()) {
		return fmt.Errorf("size mismach between joining entity '%v' and the environment '%v'",
			entity.name, env.name)
	}
	env.joins = append(env.joins, entity)
	entity.joinEnvironment(env)
	return nil
}

func (env *environment) affect(entity *Entity) error {
	if !env.checkNewSize(entity.OutSize()) {
		return fmt.Errorf("size mismach between affecting entity '%v' and the environment '%v'",
			entity.name, env.name)
	}
	env.affects = append(env.affects, entity)
	entity.affectEnvironment(env)
	return nil
}

func (env *environment) environmentLoop() {
	logf(4, "environment '%v': loop START", env.name)
	defer logf(4, "environment '%v': loop STOP", env.name)

	var wasSent bool
	for effect := range env.effectChan {
		dec, _ := TritsToBigInt(effect)
		logf(2, "Environment '%v' <- '%v' (%v)", env.name, TritsToString(effect), dec)
		if env.dispatcher.waveMode {
			env.setValue(effect)
			env.dispatcher.holdWaveWGDone("environmentLoop")
			env.dispatcher.releaseWaveWGWait("environmentLoop")
			env.setValue(nil) // environment value becomes invalid during wave
			env.dispatcher.holdWaveWGAdd(len(env.joins), "environmentLoop")
		}
		wasSent = false
		if effect != nil {
			if !env.dispatcher.waveMode {
				env.dispatcher.quantWGAdd(len(env.joins), "environmentLoop") // len(env.joins) new waves starts
			}
			// TODO change to dynamic select
			for _, entity := range env.joins {
				entity.inChan <- effect
				wasSent = true
			}
		}
		if !wasSent {
			// wave ends here
			env.setValue(effect)
		}
		if !env.dispatcher.waveMode {
			env.dispatcher.quantWGDone("environmentLoop")
		}
	}
}

// value is valid only outside quant and wave
func (env *environment) setValue(val Trits) {
	logf(3, "------ SET value env '%v' = '%v'", env.name, TritsToString(val))
	env.value = val
}

func (env *environment) getValue() Trits {
	return env.value
}

func (env *environment) invalidate() {
	if env.invalid {
		return
	}
	env.invalid = true
	close(env.effectChan)

	for _, entity := range env.joins {
		entity.stopListeningToEnvironment(env)
	}
	for _, entity := range env.affects {
		entity.stopAffectingEnvironment(env)
	}
}
