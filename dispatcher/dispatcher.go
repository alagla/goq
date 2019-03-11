package dispatcher

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/utils"
	"sync"
)

type Dispatcher struct {
	sync.RWMutex
	environments  map[string]*Environment
	quantWG       sync.WaitGroup
	waveStopWG    sync.WaitGroup
	waveReleaseWG sync.WaitGroup
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		environments: make(map[string]*Environment),
	}
}

func (disp *Dispatcher) SetEnvironmentSize(envName string, size int64) error {
	disp.Lock()
	defer disp.Unlock()

	env, ok := disp.environments[envName]
	if !ok {
		return fmt.Errorf("no such environment: '%v'", envName)
	}
	env.size = size
	return nil
}

func (disp *Dispatcher) GetOrCreateEnvironment_(name string) *Environment {
	_, ok := disp.environments[name]
	if !ok {
		disp.environments[name] = NewEnvironment(disp, name)
	}
	return disp.environments[name]
}

func (disp *Dispatcher) Join(envName string, entity EntityInterface) (*Environment, error) {
	disp.Lock()
	defer disp.Unlock()

	env := disp.GetOrCreateEnvironment_(envName)
	return env, entity.JoinEnvironment(env)
}

func (disp *Dispatcher) Affect(envName string, entity EntityInterface) (*Environment, error) {
	disp.Lock()
	defer disp.Unlock()

	env := disp.GetOrCreateEnvironment_(envName)
	return env, entity.AffectEnvironment(env)
}

func (disp *Dispatcher) DeleteEnvironment(envName string) error {
	disp.Lock()
	defer disp.Unlock()
	env, ok := disp.environments[envName]
	if !ok {
		return fmt.Errorf("can't find environment '%v'", envName)
	}
	env.stop()
	delete(disp.environments, envName)
	logf(3, "deleted environment '%v'", envName)
	return nil
}

func (disp *Dispatcher) DoQuant(envName string, effect Trits) error {
	// only one quant at a time
	disp.quantWG.Wait()

	// No changes to dispatcher state within a quant
	disp.Lock()
	defer disp.Unlock()

	if effect == nil {
		return fmt.Errorf("DoQuant: effect is nil")
	}
	env := disp.GetOrCreateEnvironment_(envName)
	if env == nil {
		return fmt.Errorf("DoQuant: can't find environment '%v'", envName)
	}
	size := int(env.Size())
	if size == 0 {
		effect = Trits{0}
	} else {
		if len(effect) != size {
			if len(effect) > size {
				return fmt.Errorf("DoQuant: trit vector '%v' is too long for the environment '%v', size = %v",
					utils.TritsToString(effect), envName, size)
			}
			effect = PadTrits(effect, size)
		}
	}
	env.PostEffect(effect)

	// wait for quant to finish
	disp.quantWG.Wait()

	return nil
}

func (disp *Dispatcher) Value(envName string) (Trits, error) {
	disp.RLock()
	defer disp.RUnlock()

	env, ok := disp.environments[envName]
	if !ok {
		return nil, fmt.Errorf("can't find environment '%v'", envName)
	}
	return env.GetValue(), nil
}

func (disp *Dispatcher) Values() map[string]Trits {
	disp.RLock()
	defer disp.RUnlock()

	ret := make(map[string]Trits)
	for name, env := range disp.environments {
		if env.value != nil {
			ret[name] = env.value
		}
	}
	return ret
}
