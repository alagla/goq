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
	waveLock      *LockWithTimeout // waves are locked
	timeout       time.Duration
	quantWG       sync.WaitGroup
	holdWaveWG    sync.WaitGroup
	releaseWaveWG sync.WaitGroup
	running       bool // quant is running. Set by RunQuant under generak and wave locks

	isWaveMode bool
	modeMutex  *sync.RWMutex
}

func NewDispatcher(lockTimeout time.Duration) *Dispatcher {
	return &Dispatcher{
		environments: make(map[string]*environment),
		generalLock:  NewAsyncLock(),
		waveLock:     NewAsyncLock(),
		timeout:      lockTimeout,
		modeMutex:    &sync.RWMutex{},
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

func (disp *Dispatcher) SetWaveMode(mode bool) error {
	if !disp.generalLock.Acquire(disp.timeout) {
		return fmt.Errorf("can't change wave/quant mode within quant")
	}
	defer disp.generalLock.Release()
	disp.modeMutex.Lock()
	defer disp.modeMutex.Unlock()
	disp.isWaveMode = mode
	return nil
}

func (disp *Dispatcher) IsWaveMode() bool {
	disp.modeMutex.RLock()
	defer disp.modeMutex.RUnlock()
	return disp.isWaveMode
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

func (disp *Dispatcher) RunQuant(envName string, effect Trits, async bool) error {
	if !disp.generalLock.Acquire(disp.timeout) {
		return fmt.Errorf("request dispatcher lock timeout: can't run quant")
	}
	env := disp.getEnvironment_(envName)
	if env == nil {
		disp.generalLock.Release()
		return fmt.Errorf("RunQuant: can't find environment '%v'", envName)
	}
	if env.size == 0 {
		effect = Trits{0}
	} else {
		if int64(len(effect)) != env.size {
			if int64(len(effect)) > env.size {
				disp.generalLock.Release()
				return fmt.Errorf("RunQuant: trit vector '%v' is too long for the environment '%v', size = %v",
					utils.TritsToString(effect), envName, env.size)
			}
			effect = PadTrits(effect, int(env.size))
		}
	}
	// quant is always started with a wave
	if !disp.waveLock.Acquire(disp.timeout) {
		disp.generalLock.Release()
		return fmt.Errorf("request wave lock timeout: can't run a quant")
	}
	disp.running = true
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
	disp.waveLock.Release()
	disp.generalLock.Release()
}

func (disp *Dispatcher) Value(envName string) (Trits, error) {
	// environments only have values outside waves
	disp.waveLock.Acquire(disp.timeout)
	defer disp.waveLock.Release()

	env, ok := disp.environments[envName]
	if !ok {
		return nil, fmt.Errorf("can't find environment '%v'", envName)
	}
	return env.getValue(), nil
}

func (disp *Dispatcher) Values() map[string]Trits {
	disp.waveLock.Acquire(disp.timeout)
	defer disp.waveLock.Release()

	ret := make(map[string]Trits)
	for name, env := range disp.environments {
		if env.value != nil {
			ret[name] = env.value
		}
	}
	return ret
}

func (disp *Dispatcher) Wave() error {
	// if wave lock is acquired, that means disp.running is consistent
	if !disp.waveLock.Acquire(disp.timeout) {
		return fmt.Errorf("can't acquire wave lock")
	}

	if !disp.isWaveMode {
		return fmt.Errorf("not in in 'wave' mode")
	}
	if !disp.running {
		return fmt.Errorf("can't run wave outside a quant")
	}
	disp.releaseWaveWG.Done()
	return nil
}
