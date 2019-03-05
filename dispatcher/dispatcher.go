package dispatcher

import (
	. "github.com/lunfardo314/goq/abstract"
	"sync"
)

type environment struct {
	joins   []FuncDefInterface
	affects []FuncDefInterface
}

type Dispatcher struct {
	sync.RWMutex
	environments map[string]*environment
}

func StartDispatcher() *Dispatcher {
	ret := &Dispatcher{}
	return ret
}

func (disp *Dispatcher) Stop() {

}

func (disp *Dispatcher) AddJoin(module ModuleInterface, envName string, fun FuncDefInterface) error {
	disp.Lock()
	defer disp.Unlock()

	name := module.GetName() + "::" + envName
	if env, ok := disp.environments[name]; ok {
		env.joins = append(env.joins, fun)
	}
	return nil
}

func (disp *Dispatcher) AddAffect(module ModuleInterface, envName string, fun FuncDefInterface) error {
	disp.Lock()
	defer disp.Unlock()

	name := module.GetName() + "::" + envName
	if env, ok := disp.environments[name]; ok {
		env.affects = append(env.affects, fun)
	}
	return nil
}
