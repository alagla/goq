package dispatcher

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
)

type EntityInterface interface {
	GetName() string           // name of the entity
	OutSize() int64            // result size in trits
	InSize() int64             // concat arguments, total size in trits
	Join(*Environment) error   // join the environment = will be listening to the environment
	Stop()                     // stop listening to environments. Before GC
	Affect(*Environment) error // affect the environment = any results will be sent to the environments as effects
	Invoke(Trits)              // calls the entity with arguments
}

type BaseEntity struct {
	name           string
	inSize         int64
	outSize        int64
	affects        []*Environment // list of affected environments where effects are sent
	inChan         chan Trits     // chan for incoming effects
	joined         bool
	effectCallback func(Trits) Trits // function called for each effect
}

func NewBaseEntity(name string, inSize, outSize int64, effectCallback func(Trits) Trits) *BaseEntity {
	ret := &BaseEntity{
		name:           name,
		inSize:         inSize,
		outSize:        outSize,
		affects:        make([]*Environment, 0),
		inChan:         make(chan Trits, 1), // buffer to avoid deadlocks
		effectCallback: effectCallback,
	}
	return ret
}

// after that entity becomes invalid
// called by the environment only
func (ent *BaseEntity) Stop() {
	if ent.joined {
		close(ent.inChan)
	}
}

func (ent *BaseEntity) GetName() string {
	return ent.name
}

func (ent *BaseEntity) InSize() int64 {
	return ent.inSize
}

func (ent *BaseEntity) OutSize() int64 {
	return ent.outSize
}

func (ent *BaseEntity) Affect(env *Environment) error {
	if err := env.checkNewSize(ent.outSize); err != nil {
		return fmt.Errorf("error while registering affect, entity '%v': %v", ent.GetName(), err)
	}
	ent.affects = append(ent.affects, env)
	return nil
}

func (ent *BaseEntity) Join(env *Environment) error {
	if !ent.joined {
		ent.joined = true
		ent.inChan = make(chan Trits, 1) // buffer to avoid deadlocks
		go ent.loopEffects()             // start listening to incoming effects
	}
	return env.Join(ent)
}

func (ent *BaseEntity) Invoke(t Trits) {
	ent.inChan <- t
}

func (ent *BaseEntity) loopEffects() {
	var res Trits
	for args := range ent.inChan {
		res = ent.effectCallback(args)
		ent.postEffect(res)
	}
}

func (ent *BaseEntity) postEffect(effect Trits) {
	for _, env := range ent.affects {
		env.PostEffect(effect) // sync or async?
	}
}
