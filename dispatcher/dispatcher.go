package dispatcher

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/utils"
	"sync"
	"time"
)

type Dispatcher struct {
	environments  map[string]*environment
	generalLock   *LockWithTimeout // controls environments, join, affect, modes
	timeout       time.Duration
	quantWG       sync.WaitGroup
	holdWaveWG    sync.WaitGroup
	releaseWaveWG sync.WaitGroup
	waveMode      bool // TODO
}

func NewDispatcher(lockTimeout time.Duration) *Dispatcher {
	return &Dispatcher{
		environments: make(map[string]*environment),
		generalLock:  NewAsyncLock(),
		timeout:      lockTimeout,
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

func (disp *Dispatcher) CreateEnvironment(name string) error {
	if !disp.generalLock.Acquire(disp.timeout) {
		return fmt.Errorf("request lock timeout: can't create environment")
	}
	defer disp.generalLock.Release()

	if disp.getEnvironment_(name) != nil {
		return fmt.Errorf("environment '%v' already exists", name)
	}
	disp.environments[name] = NewEnvironment(disp, name)
	return nil
}

// executes 'join' and 'affect' of the entity
func (disp *Dispatcher) Attach(entity *Entity, joins, affects []string) error {
	if !disp.generalLock.Acquire(disp.timeout) {
		return fmt.Errorf("acquire lock timeout: can't attach entity to environment")
	}
	defer disp.generalLock.Release()

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
	if !disp.generalLock.Acquire(disp.timeout) {
		return fmt.Errorf("request lock timeout: can't delete environment")
	}
	defer disp.generalLock.Release()

	env, ok := disp.environments[envName]
	if !ok {
		return fmt.Errorf("can't find environment '%v'", envName)
	}
	env.invalidate()
	delete(disp.environments, envName)
	logf(3, "deleted environment '%v'", envName)
	return nil
}

func (disp *Dispatcher) StartWave(envName string, effect Trits, onStop func()) error {
	if err := disp.startWave(envName, true, effect); err != nil {
		return err
	}
	disp.holdWaveWG.Wait()
	if onStop != nil {
		onStop()
	}
	return nil
}

func (disp *Dispatcher) StartQuant(envName string, effect Trits, onStop func()) error {
	if err := disp.startWave(envName, false, effect); err != nil {
		return err
	}
	disp.quantWG.Wait()
	if onStop != nil {
		onStop()
	}
	return nil
}

func (disp *Dispatcher) startWave(envName string, waveMode bool, effect Trits) error {
	env := disp.getEnvironment_(envName)
	if env == nil {
		return fmt.Errorf("startWave: can't find environment '%v'", envName)
	}
	if env.size == 0 {
		effect = Trits{0}
	} else {
		if int64(len(effect)) != env.size {
			if int64(len(effect)) > env.size {
				disp.generalLock.Release()
				return fmt.Errorf("startWave: trit vector '%v' is too long for the environment '%v', size = %v",
					utils.TritsToString(effect), envName, env.size)
			}
			effect = PadTrits(effect, int(env.size))
		}
	}

	disp.waveMode = waveMode
	if waveMode {
		disp.holdWaveWG.Add(1)
		disp.releaseWaveWG.Add(1)
	} else {
		disp.quantWG.Add(1)
	}
	env.effectChan <- effect
	return nil
}

func (disp *Dispatcher) Wave() error {
	if !disp.waveMode {
		return fmt.Errorf("not in wave mode")
	}
	disp.holdWaveWG.Wait()
	disp.releaseWaveWG.Done()
	disp.releaseWaveWG.Add(1)
	return nil
}

func (disp *Dispatcher) Value(envName string) (Trits, error) {

	env, ok := disp.environments[envName]
	if !ok {
		return nil, fmt.Errorf("can't find environment '%v'", envName)
	}
	return env.getValue(), nil
}

func (disp *Dispatcher) Values() map[string]Trits {

	ret := make(map[string]Trits)
	for name, env := range disp.environments {
		if env.value != nil {
			ret[name] = env.value
		}
	}
	return ret
}
