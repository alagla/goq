package supervisor

import (
	. "github.com/iotaledger/iota.go/trinary"
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
		logf(5, "stopped entity '%v'", ent.name)
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

func (ent *Entity) entityLoop() {
	//logf(7, "entity '%v' loop STARTED", ent.name)
	//defer logf(7, "entity '%v'loop STOPPED", ent.name)

	var null bool
	for msg := range ent.inChan {
		if msg.effect == nil {
			panic("entity loop: nil effect")
		}
		dec, _ := utils.TritsToBigInt(msg.effect)
		logf(3, "effect '%v' (%v) -> entity '%v'", utils.TritsToString(msg.effect), dec, ent.name)

		// calculate result
		result := make(Trits, ent.outSize)
		null = ent.core.Call(msg.effect, result)

		if !null {
			if msg.lastWithinLimit {
				// postpone to new quant
				for _, affectInfo := range ent.affecting {
					_ = ent.supervisor.postEffect("", affectInfo.environment, result, affectInfo.delay, false)
				}
			} else {
				for _, affectInfo := range ent.affecting {
					if affectInfo.delay == 0 {
						ent.supervisor.quantWG.Add(1)
						affectInfo.environment.effectChan <- result
					} else {
						_ = ent.supervisor.postEffect("", affectInfo.environment, result, affectInfo.delay, false)
					}
				}
			}
		}
		ent.supervisor.quantWG.Done()
	}
}

func (ent *Entity) sendEffect(effect Trits, lastWithinLimit bool) {
	ent.inChan <- entityMsg{
		effect:          effect,
		lastWithinLimit: lastWithinLimit,
	}
}
