package dispatcher

import (
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/utils"
)

type EntityCore interface {
	Call(Trits, Trits) bool
}

type affectEntData struct {
	environment *environment
	delay       int
}

type entityMsg struct {
	effect          Trits
	lastWithinLimit bool
}

type Entity struct {
	dispatcher *Dispatcher
	name       string
	inSize     int64
	outSize    int64
	affecting  []*affectEntData // list of affected environments where effects are sent
	joined     []*environment   // list of environments which are being listened to
	inChan     chan entityMsg   // chan for incoming effects
	core       EntityCore       // function called for each effect
	terminal   bool             // can't affect environments
}

func (ent *Entity) GetName() string {
	return ent.name
}

func (ent *Entity) GetCore() EntityCore {
	return ent.core
}

func (ent *Entity) InSize() int64 {
	return ent.inSize
}

func (ent *Entity) OutSize() int64 {
	return ent.outSize
}

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
		logf(5, "stopped entity '%v'", ent.GetName())
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
	logf(7, "entity '%v' loop STARTED", ent.name)
	defer logf(7, "entity '%v'loop STOPPED", ent.name)

	var null bool
	for msg := range ent.inChan {
		if msg.effect == nil {
			panic("nil effect")
		}
		dec, _ := utils.TritsToBigInt(msg.effect)
		logf(3, "effect '%v' (%v) -> entity '%v'", utils.TritsToString(msg.effect), dec, ent.name)
		// calculate result
		res := make(Trits, ent.outSize)
		null = ent.core.Call(msg.effect, res)
		if !null {
			if msg.lastWithinLimit {
				// postpone to new quant
				for _, affectInfo := range ent.affecting {
					_ = ent.dispatcher.postEffect("", affectInfo.environment, res, 0, false)
				}
			} else {
				for _, affectInfo := range ent.affecting {
					if affectInfo.delay == 0 {
						ent.dispatcher.quantWG.Add(1)
						affectInfo.environment.effectChan <- res
					} else {
						_ = ent.dispatcher.postEffect("", affectInfo.environment, res, affectInfo.delay, false)
					}
				}
			}
		}
		ent.dispatcher.quantWG.Done()
	}
}

func (ent *Entity) sendEffect(effect Trits, lastWithinLimit bool) {
	ent.inChan <- entityMsg{
		effect:          effect,
		lastWithinLimit: lastWithinLimit,
	}
}
