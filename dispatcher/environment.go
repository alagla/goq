package dispatcher

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/utils"
	"sync"
)

type environment struct {
	dispatcher *Dispatcher
	name       string
	invalid    bool
	joins      []*Entity
	affects    []*Entity
	size       int64
	effectChan chan Trits
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

func (env *environment) adjustEffect(effect Trits) (Trits, error) {
	if env.size == 0 {
		effect = Trits{0}
	} else {
		if int64(len(effect)) != env.size {
			if int64(len(effect)) > env.size {
				return nil, fmt.Errorf("trit vector '%v' is too long for the environment '%v', size = %v",
					TritsToString(effect), env.name, env.size)
			}
			effect = PadTrits(effect, int(env.size))
		}
	}
	return effect, nil
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

// main loop of the environment
func (env *environment) environmentLoop() {
	logf(4, "environment '%v': loop START", env.name)
	defer logf(4, "environment '%v': loop STOP", env.name)

	for effect := range env.effectChan {
		if effect == nil {
			panic("nil effect")
		}
		dec, _ := TritsToBigInt(effect)
		logf(2, "Environment '%v' <- '%v' (%v)", env.name, TritsToString(effect), dec)
		env.dispatcher.quantWG.Add(len(env.joins))
		env.waitWave(effect)
		for _, entity := range env.joins {
			entity.inChan <- effect
		}
		env.dispatcher.quantWG.Done()
	}
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

func (env *environment) waitWave(value Trits) {
	if !env.dispatcher.waveCoo.isWaveMode() {
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	env.dispatcher.waveCoo.chIn <- &waveCmd{
		environment: env,
		value:       value,
		wg:          &wg,
	}
	wg.Wait()
	return
}
