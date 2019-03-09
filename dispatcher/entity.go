package dispatcher

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/utils"
)

type CallableWithTrits interface {
	Call(Trits, Trits) bool
}

type EntityInterface interface {
	Name() string                         // name of the entity
	OutSize() int64                       // result size in trits
	InSize() int64                        // concat arguments, total size in trits
	JoinEnvironment(*Environment) error   // join the environment = will be listening to the environment
	AffectEnvironment(*Environment) error // affect the environment = any results will be sent to the environments as effects
	Stop()                                // stop listening to environments. Before GC
	Invoke(Trits)                         // calls the entity with arguments
}

type BaseEntity struct {
	name           string
	inSize         int64
	outSize        int64
	affects        []*Environment    // list of affected environments where effects are sent
	inChan         chan Trits        // chan for incoming effects
	effectCallable CallableWithTrits // function called for each effect
}

func NewBaseEntity(name string, inSize, outSize int64, effectCallable CallableWithTrits) *BaseEntity {
	ret := &BaseEntity{
		name:           name,
		inSize:         inSize,
		outSize:        outSize,
		affects:        make([]*Environment, 0),
		inChan:         make(chan Trits),
		effectCallable: effectCallable,
	}
	go ret.loopEffects() // start listening to incoming effects
	return ret
}

// after that entity becomes invalid
// called by the environment only
func (ent *BaseEntity) Stop() {
	close(ent.inChan)
}

func (ent *BaseEntity) Name() string {
	return ent.name
}

func (ent *BaseEntity) InSize() int64 {
	return ent.inSize
}

func (ent *BaseEntity) OutSize() int64 {
	return ent.outSize
}

func (ent *BaseEntity) AffectEnvironment(env *Environment) error {
	if err := env.checkNewSize(ent.outSize); err != nil {
		return fmt.Errorf("error while registering affect, entity '%v': %v", ent.Name(), err)
	}
	ent.affects = append(ent.affects, env)
	return nil
}

func (ent *BaseEntity) JoinEnvironment(env *Environment) error {
	return env.Join(ent)
}

func (ent *BaseEntity) Invoke(t Trits) {
	ent.inChan <- t
}

func (ent *BaseEntity) loopEffects() {
	logf(3, "loopEffects STARTED for entity '%v'", ent.name)
	defer logf(3, "loopEffects STOPPED for entity '%v'", ent.name)

	res := make(Trits, ent.outSize)

	for args := range ent.inChan {
		logf(2, "Entity '%v' <- '%v'", ent.Name(), utils.TritsToString(args))
		if !ent.effectCallable.Call(args, res) {
			ent.postEffect(res)
		}
	}
}

func (ent *BaseEntity) postEffect(effect Trits) {
	for _, env := range ent.affects {
		env.PostEffect(effect) // sync or async?
	}
}
