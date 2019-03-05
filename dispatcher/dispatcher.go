package dispatcher

import (
	"sync"
)

type Dispatcher struct {
	sync.RWMutex
	environments map[string]*Environment
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		environments: make(map[string]*Environment),
	}
}

func (disp *Dispatcher) GetEnvironment(name string) *Environment {
	disp.RLock()
	defer disp.RUnlock()
	_, ok := disp.environments[name]
	if !ok {
		disp.environments[name] = NewEnvironment(name)
	}
	return disp.environments[name]

}

func (disp *Dispatcher) Join(envName string, entity EntityInterface) (*Environment, error) {
	env := disp.GetEnvironment(envName)
	return env, entity.Join(env)
}

func (disp *Dispatcher) Affect(envName string, entity EntityInterface) (*Environment, error) {
	env := disp.GetEnvironment(envName)
	return env, entity.Affect(env)
}
