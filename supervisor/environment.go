package supervisor

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/cfg"
	. "github.com/lunfardo314/goq/utils"
)

func newEnvironment(disp *Supervisor, name string) *environment {
	ret := &environment{
		supervisor: disp,
		name:       name,
		joins:      make([]*joinEnvData, 0),
		affects:    make([]*Entity, 0),
		effectChan: make(chan Trits),
	}
	go ret.environmentLoop()
	return ret
}

func (env *environment) join(entity *Entity, limit int) error {
	env.joins = append(env.joins, &joinEnvData{
		entity: entity,
		limit:  limit,
	})
	entity.joinEnvironment(env, limit)
	return nil
}

func (env *environment) affect(entity *Entity, delay int) error {
	env.affects = append(env.affects, entity)
	entity.affectEnvironment(env, delay)
	return nil
}

// main loop of the environment
func (env *environment) environmentLoop() {
	Logf(7, "environment '%v': loop START", env.name)
	defer Logf(7, "environment '%v': loop STOP", env.name)

	for effect := range env.effectChan {
		if effect == nil {
			panic("nil effect")
		}
		// TODO optimize logging in the loop
		Logf(3, "effect '%v' (%v) -> environment '%v'",
			TritsToString(effect), MustTritsToBigInt(effect), env.name)
		env.supervisor.quantWG.Add(len(env.joins))
		for _, joinData := range env.joins {
			joinData.count++
			joinData.entity.sendEffect(effect, joinData.count == joinData.limit)
		}
		env.supervisor.quantWG.Done()
	}
}

func (env *environment) invalidate() {
	if env.invalid {
		return
	}
	env.invalid = true
	close(env.effectChan)

	for _, joinData := range env.joins {
		joinData.entity.stopListeningToEnvironment(env)
	}
	for _, entity := range env.affects {
		entity.stopAffectingEnvironment(env)
	}
}
