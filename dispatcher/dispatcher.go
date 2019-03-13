package dispatcher

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/utils"
	"sync"
)

type Dispatcher struct {
	sync.RWMutex  // access to environment structure, joins, affect. Except values
	environments  map[string]*environment
	running       bool // is within quant
	quantWG       sync.WaitGroup
	holdWaveWG    sync.WaitGroup
	releaseWaveWG sync.WaitGroup
	isWaveMode    bool
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		environments: make(map[string]*environment),
	}
}

func (disp *Dispatcher) getEnvironment_(name string) *environment {
	env, ok := disp.environments[name]
	if !ok {
		return nil
	}
	return env
}

func (disp *Dispatcher) getOrCreateEnvironment_(name string) *environment {
	ret := disp.getEnvironment_(name)
	if ret != nil {
		return ret
	}
	disp.environments[name] = NewEnvironment(disp, name)
	return disp.environments[name]
}

func (disp *Dispatcher) SetWaveMode(mode bool) {
	disp.Lock()
	defer disp.Unlock()
	disp.isWaveMode = mode
}

func (disp *Dispatcher) IsWaveMode() bool {
	disp.RLock()
	defer disp.RUnlock()
	return disp.isWaveMode
}

func (disp *Dispatcher) CreateEnvironment(name string) error {
	disp.Lock()
	defer disp.Unlock()

	if disp.getEnvironment_(name) != nil {
		return fmt.Errorf("environment '%v' already exists", name)
	}
	disp.environments[name] = NewEnvironment(disp, name)
	return nil
}

// executes 'join' and 'affect' of the entity
func (disp *Dispatcher) Attach(entity *Entity, joins, affects []string) error {
	disp.Lock()
	defer disp.Unlock()

	for _, envName := range joins {
		env := disp.getOrCreateEnvironment_(envName)
		if err := env.join(entity); err != nil {
			return err
		}
	}
	for _, envName := range affects {
		env := disp.getOrCreateEnvironment_(envName)
		if err := env.affect(entity); err != nil {
			return err
		}
	}
	return nil
}

func (disp *Dispatcher) DeleteEnvironment(envName string) error {
	disp.Lock()
	defer disp.Unlock()
	env, ok := disp.environments[envName]
	if !ok {
		return fmt.Errorf("can't find environment '%v'", envName)
	}
	env.invalidate()
	delete(disp.environments, envName)
	logf(3, "deleted environment '%v'", envName)
	return nil
}

func (disp *Dispatcher) IsRunning() bool {
	disp.RLock()
	defer disp.RUnlock()
	return disp.running
}
func (disp *Dispatcher) RunQuant(envName string, effect Trits, async bool) error {
	disp.Lock()
	disp.running = true

	env := disp.getEnvironment_(envName)
	if env == nil {
		disp.Unlock()
		return fmt.Errorf("RunQuant: can't find environment '%v'", envName)
	}
	if env.size == 0 {
		effect = Trits{0}
	} else {
		if int64(len(effect)) != env.size {
			if int64(len(effect)) > env.size {
				disp.Unlock()
				return fmt.Errorf("RunQuant: trit vector '%v' is too long for the environment '%v', size = %v",
					utils.TritsToString(effect), envName, env.size)
			}
			effect = PadTrits(effect, int(env.size))
		}
	}
	if disp.isWaveMode {
		disp.releaseWaveWG.Add(1)
	}
	env.postEffect(effect)

	if async {
		go disp.finishQuant()
	} else {
		disp.finishQuant()
	}
	return nil
}

func (disp *Dispatcher) finishQuant() {
	disp.quantWG.Wait()
	logf(3, "---------------- Quant finished")
	disp.running = false
	disp.Unlock()
}

func (disp *Dispatcher) Value(envName string) (Trits, error) {
	disp.RLock()
	defer disp.RUnlock()

	env, ok := disp.environments[envName]
	if !ok {
		return nil, fmt.Errorf("can't find environment '%v'", envName)
	}
	return env.getValue(), nil
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

func (disp *Dispatcher) Wave() error {
	// TODO locking
	if !disp.isWaveMode {
		return fmt.Errorf("Not in 'wave' mode")
	}
	if !disp.running {
		return fmt.Errorf("Quant is not started")
	}
	disp.releaseWaveWG.Done()
	return nil
}
