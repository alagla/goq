package dispatcher

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
)

type EntityInterface interface {
	GetName() string
	OutSize() int64 // result
	InSize() int64  // arguments
	Join(*Environment) error
	Affect(*Environment) error
	Call(Trits) Trits
}

type BaseEntity struct {
	name    string
	inSize  int64
	outSize int64
	affects []*Environment
}

func NewBaseEntity(name string, inSize, outSize int64) *BaseEntity {
	return &BaseEntity{
		name:    name,
		inSize:  inSize,
		outSize: outSize,
		affects: make([]*Environment, 0),
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
	if env.size == 0 {
		env.size = ent.outSize
	} else {
		if ent.outSize != env.size {
			return fmt.Errorf("size mismatch between environment '%v' and affecting entity '%v'",
				env.name, ent.GetName())
		}
	}
	ent.affects = append(ent.affects, env)
	return nil
}

func (ent *BaseEntity) Join(env *Environment) error {
	return env.Join(ent)
}

func (ent *BaseEntity) Call(_ Trits) Trits {
	return nil
}

type FunctionEntity struct {
	BaseEntity
	funDef FuncDefInterface
	inChan chan Trits
}

func NewFunctionEntity(funDef FuncDefInterface) *FunctionEntity {
	ret := &FunctionEntity{
		BaseEntity: *NewBaseEntity(funDef.GetName(), funDef.ArgSize(), funDef.Size()),
		funDef:     funDef,
		inChan:     make(chan Trits, 1), // buffer to avoid deadlocks
	}
	go ret.loopEffects()
	return ret
}

func (ent *FunctionEntity) Call(args Trits) Trits {
	return nil
}

func (ent *FunctionEntity) invoke(t Trits) {
	ent.inChan <- t
}

func (ent *FunctionEntity) loopEffects() {
	for t := range ent.inChan {
		_ = ent.Call(t)
	}
}
