package dispatcher

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/utils"
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

func (disp *Dispatcher) GetEnvironmentInfo(name string) (int64, bool) {
	ret, ok := disp.environments[name]
	if !ok {
		return 0, false
	}
	return ret.Size(), true
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
	return env, entity.JoinEnvironment(env)
}

func (disp *Dispatcher) Affect(envName string, entity EntityInterface) (*Environment, error) {
	env := disp.GetEnvironment(envName)
	return env, entity.AffectEnvironment(env)
}

func (disp *Dispatcher) PostEffect(envName string, effect Trits) error {
	env := disp.GetEnvironment(envName)
	if env == nil {
		return fmt.Errorf("can't find environment '%v'", envName)
	}

	if env.Size() != int64(len(effect)) {
		return fmt.Errorf("size mismatch while posting effect '%v' to the environment '%v', size must be %v",
			utils.TritsToString(effect), env.GetName(), env.Size())
	}
	env.PostEffect(effect)
	return nil
}
