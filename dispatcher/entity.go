package dispatcher

import (
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/utils"
)

type EntityCore interface {
	Call(Trits, Trits) bool
}

type Entity struct {
	dispatcher *Dispatcher
	name       string
	inSize     int64
	outSize    int64
	affecting  []*environment // list of affected environments where effects are sent
	joined     []*environment // list of environments which are being listened to
	inChan     chan Trits     // chan for incoming effects
	entityCore EntityCore     // function called for each effect
	terminal   bool           // can't affect environments
}

func (ent *Entity) GetName() string {
	return ent.name
}

func (ent *Entity) InSize() int64 {
	return ent.inSize
}

func (ent *Entity) OutSize() int64 {
	return ent.outSize
}

func (ent *Entity) affectEnvironment(env *environment) {
	ent.affecting = append(ent.affecting, env)
}

func (ent *Entity) joinEnvironment(env *environment) {
	ent.joined = append(ent.joined, env)
	ent.checkStart()
}

func (ent *Entity) stopAffectingEnvironment(env *environment) {
	tmpList := make([]*environment, 0)
	for _, e := range ent.affecting {
		if e != env {
			tmpList = append(tmpList, e)
		}
	}
}

func (ent *Entity) stopListeningToEnvironment(env *environment) {
	tmpList := make([]*environment, 0)
	for _, e := range ent.joined {
		if e != env {
			tmpList = append(tmpList, e)
		}
	}
	ent.joined = tmpList
	if ent.checkStop() {
		logf(5, "stopped entity '%v'", ent.GetName())
	}
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
		ent.inChan = make(chan Trits)
		go ent.entityLoop()
	}
}

func (ent *Entity) entityLoop() {
	logf(7, "entity '%v' loop STARTED", ent.name)
	defer logf(7, "entity '%v'loop STOPPED", ent.name)

	var null bool
	for effect := range ent.inChan {
		if effect == nil {
			panic("nil effect")
		}
		dec, _ := utils.TritsToBigInt(effect)
		logf(3, "effect '%v' (%v) -> entity '%v'", utils.TritsToString(effect), dec, ent.name)
		// calculate result
		res := make(Trits, ent.outSize)
		null = ent.entityCore.Call(effect, res)
		if !null {
			ent.dispatcher.quantWG.Add(len(ent.affecting))
			for _, env := range ent.affecting {
				env.effectChan <- res
			}
		}
		ent.dispatcher.quantWG.Done()
	}
}
