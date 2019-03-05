package dispatcher

import (
	. "github.com/lunfardo314/goq/abstract"
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

func (disp *Dispatcher) Join(envName string, fun FuncDefInterface) *Environment {
	env := disp.GetEnvironment(envName)
	env.Join(fun)
	return env
}

func (disp *Dispatcher) Affect(envName string, fun FuncDefInterface) *Environment {
	env := disp.GetEnvironment(envName)
	env.Affect(fun)
	return env
}
