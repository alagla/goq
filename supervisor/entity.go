package supervisor

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/cfg"
	"github.com/lunfardo314/goq/utils"
)

func (ent *Entity) joinEnvironment(env *environment, limit int) {
	ent.joined = append(ent.joined, env)
	ent.checkStart()
}

func (ent *Entity) affectEnvironment(env *environment, delay int) {
	ent.affecting = append(ent.affecting, &affectEntData{
		environment: env,
		delay:       delay,
	})
}

func (ent *Entity) stopListeningToEnvironment(env *environment) {
	newList := make([]*environment, 0)
	for _, e := range ent.joined {
		if e != env {
			newList = append(newList, e)
		}
	}
	ent.joined = newList
	if ent.checkStop() {
		Logf(5, "stopped entity '%v'", ent.Name)
	}
}

func (ent *Entity) stopAffectingEnvironment(env *environment) {
	newList := make([]*affectEntData, 0)
	for _, e := range ent.affecting {
		if e.environment != env {
			newList = append(newList, e)
		}
	}
	ent.affecting = newList
}

func (ent *Entity) checkStop() bool {
	ret := false
	if ent.inChan != nil && len(ent.joined) == 0 {
		c := ent.inChan
		ent.inChan = nil
		close(c)
		ret = true
	}
	return ret
}

func (ent *Entity) checkStart() {
	if ent.inChan == nil && len(ent.joined) != 0 {
		ent.inChan = make(chan entityMsg)
		go ent.entityLoop()
	}
}

// main loop of the entity

func (ent *Entity) entityLoop() {
	var null bool
	for msg := range ent.inChan {
		if msg.effect == nil {
			panic("entity loop: nil effect")
		}

		LogDefer(3, func() {
			Logf(3, "effect '%v' (%v) -> entity '%v'",
				utils.TritsToString(msg.effect), utils.MustTritsToBigInt(msg.effect), ent.Name)
		})

		// calculate result by calling entity core
		// memory allocation is inevitable this effect will travel along channels
		// to environments and other entities
		result := make(Trits, ent.outSize)
		null = ent.core.Call(msg.effect, result)

		if !null {
			if msg.lastWithinLimit {
				// if effect message is marked as last in the quant for this entity, it means
				// join limit is reached for this entity.
				// The effect is postponed to new quant by resending effect to the main queue of the Supervisor
				for _, affectInfo := range ent.affecting {
					_ = ent.Supervisor.postEffect("", affectInfo.environment, result, affectInfo.delay, false)
				}
			} else {
				// otherwise the effect can be processed with current quant and therefore is
				// posted to the each environment of affected environments with respective delay
				for _, affectInfo := range ent.affecting {
					// if there's no delay, effect is sent directly to the environment's channel
					// in the current quant
					if affectInfo.delay == 0 {
						ent.Supervisor.quantWG.Add(1)
						affectInfo.environment.effectChan <- result
					} else {
						// otherwise effect is posted to the main input queue
						_ = ent.Supervisor.postEffect("", affectInfo.environment, result, affectInfo.delay, false)
					}
				}
			}
		}
		ent.Supervisor.quantWG.Done()
	}
}

func (ent *Entity) sendEffect(effect Trits, lastWithinLimit bool) {
	ent.inChan <- entityMsg{
		effect:          effect,
		lastWithinLimit: lastWithinLimit,
	}
}
