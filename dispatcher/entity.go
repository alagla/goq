package dispatcher

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/utils"
)

type CallableWithTrits interface {
	Call(Trits, Trits) bool
}

type Entity struct {
	dispatcher     *Dispatcher
	name           string
	inSize         int64
	outSize        int64
	affects        []*environment    // list of affected environments where effects are sent
	inChan         chan Trits        // chan for incoming effects
	effectCallable CallableWithTrits // function called for each effect
}

func NewEntity(disp *Dispatcher, name string, inSize, outSize int64, effectCallable CallableWithTrits) *Entity {
	ret := &Entity{
		dispatcher:     disp,
		name:           name,
		inSize:         inSize,
		outSize:        outSize,
		affects:        make([]*environment, 0),
		inChan:         make(chan Trits),
		effectCallable: effectCallable,
	}
	go ret.entityListenToEffectsLoop() // start listening to incoming effects
	return ret
}

func (ent *Entity) GetName() string {
	return ent.name
}

// after that entity becomes invalid
// called by the environment only
func (ent *Entity) stop() {
	close(ent.inChan)
}

func (ent *Entity) InSize() int64 {
	return ent.inSize
}

func (ent *Entity) OutSize() int64 {
	return ent.outSize
}

func (ent *Entity) affectEnvironment(env *environment) error {
	if err := env.checkNewSize_(ent.outSize); err != nil {
		return fmt.Errorf("error while registering affect, entity '%v': %v", ent.name, err)
	}
	ent.affects = append(ent.affects, env)
	return nil
}

func (ent *Entity) joinEnvironment(env *environment) error {
	return env.join(ent)
}

func (ent *Entity) invoke(t Trits) {
	ent.inChan <- t
}

func (ent *Entity) entityListenToEffectsLoop() {
	logf(4, "entityListenToEffectsLoop STARTED for entity '%v'", ent.name)
	defer logf(4, "entityListenToEffectsLoop STOPPED for entity '%v'", ent.name)

	res := make(Trits, ent.outSize)

	for effect := range ent.inChan {
		logf(2, "Entity '%v' <- '%v'", ent.name, utils.TritsToString(effect))
		// calculate result
		if !ent.effectCallable.Call(effect, res) {
			// is not null
			// mark it is done with entity
			// distribute result to affected environments
			for _, env := range ent.affects {
				env.postEffect(res)
			}
		}
		ent.dispatcher.quantWG.Done()
		logf(4, "---------------- DONE (entity '%v')", ent.name)
	}
}
