package dispatcher

import (
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/utils"
)

type EntityCore interface {
	Call(Trits, Trits) bool
}

type joinEntData struct {
	environment *environment
	limit       int
}

type affectEntData struct {
	environment *environment
	delay       int
}

type Entity struct {
	dispatcher *Dispatcher
	name       string
	inSize     int64
	outSize    int64
	affecting  []*affectEntData // list of affected environments where effects are sent
	joined     []*joinEntData   // list of environments which are being listened to
	inChan     chan Trits       // chan for incoming effects
	entityCore EntityCore       // function called for each effect
	terminal   bool             // can't affect environments
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

func (ent *Entity) joinEnvironment(env *environment, limit int) {
	ent.joined = append(ent.joined, &joinEntData{
		environment: env,
		limit:       limit,
	})
	ent.checkStart()
}

func (ent *Entity) affectEnvironment(env *environment, delay int) {
	ent.affecting = append(ent.affecting, &affectEntData{
		environment: env,
		delay:       delay,
	})
}

func (ent *Entity) stopListeningToEnvironment(env *environment) {
	newList := make([]*joinEntData, 0)
	for _, e := range ent.joined {
		if e.environment != env {
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
			for _, affectInfo := range ent.affecting {
				// TODO delay and limit
				affectInfo.environment.effectChan <- res
			}
		}
		ent.dispatcher.quantWG.Done()
	}
}
