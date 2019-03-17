package dispatcher

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	"sync"
	"time"
)

type Dispatcher struct {
	environments map[string]*environment
	generalLock  *LockWithTimeout // controls environments, join, affect, modes
	timeout      time.Duration
	waveCoo      *WaveCoordinator
	quantWG      sync.WaitGroup // released when quant ends
}

func NewDispatcher(lockTimeout time.Duration) *Dispatcher {
	return &Dispatcher{
		environments: make(map[string]*environment),
		generalLock:  NewAsyncLock(),
		timeout:      lockTimeout,
		waveCoo:      NewWaveCoordinator(),
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

func (disp *Dispatcher) WaveStart(envName string, waveMode bool, effect Trits, onQuantFinish func()) error {
	if disp.waveCoo.isWaveMode() {
		return fmt.Errorf("wave is already running")
	}
	env := disp.getEnvironment_(envName)
	if env == nil {
		return fmt.Errorf("startWave: can't find environment '%v'", envName)
	}
	var err error
	if effect, err = env.adjustEffect(effect); err != nil {
		return err
	}

	disp.waveCoo.setWaveMode(waveMode)
	disp.quantWG.Add(1)

	env.effectChan <- effect

	if onQuantFinish != nil {
		go func() {
			env.dispatcher.quantWG.Wait()
			onQuantFinish()
		}()
	}
	return nil
}

func (disp *Dispatcher) WaveNext() error {
	if !disp.waveCoo.isWaveMode() {
		return fmt.Errorf("not in wave mode")
	}
	disp.waveCoo.runWave()
	return nil
}

func (disp *Dispatcher) WaveRun() error {
	if !disp.waveCoo.isWaveMode() {
		return fmt.Errorf("not in wave mode")
	}
	disp.waveCoo.setWaveMode(false)
	disp.waveCoo.runWave()
	return nil
}

func (disp *Dispatcher) WaveValues() map[string]Trits {
	return disp.waveCoo.values()
}

func (disp *Dispatcher) IsWaveMode() bool {
	return disp.waveCoo.isWaveMode()
}
