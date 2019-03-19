package dispatcher

import (
	"fmt"
	"github.com/Workiva/go-datastructures/queue"
	. "github.com/iotaledger/iota.go/trinary"
	"sync"
	"time"
)

// TODO size checks when join/affect. Can be with different sizes

type Dispatcher struct {
	queue        *queue.Queue
	quantCount   uint64
	environments map[string]*environment
	generalLock  *LockWithTimeout // controls environments, join, affect, modes
	timeout      time.Duration
	waveCoo      *WaveCoordinator
	quantWG      sync.WaitGroup // released when quant ends
}

func NewDispatcher(lockTimeout time.Duration) *Dispatcher {
	ret := &Dispatcher{
		queue:        queue.New(5),
		environments: make(map[string]*environment),
		generalLock:  NewAsyncLock(),
		timeout:      lockTimeout,
		waveCoo:      NewWaveCoordinator(),
	}
	go ret.dispatcherInputLoop()
	return ret
}

type EntityOpts struct {
	Name     string
	InSize   int64
	OutSize  int64
	Core     EntityCore
	Terminal bool // can't affect environments (doesn't produce any result, always returns null
}

func (disp *Dispatcher) NewEntity(opt EntityOpts) *Entity {
	ret := &Entity{
		dispatcher: disp,
		name:       opt.Name,
		inSize:     opt.InSize,
		outSize:    opt.OutSize,
		affecting:  make([]*affectEntData, 0),
		joined:     make([]*joinEntData, 0),
		entityCore: opt.Core,
		terminal:   opt.Terminal,
	}
	return ret
}

func (disp *Dispatcher) getQuantCount() uint64 {
	return disp.quantCount
}

func (disp *Dispatcher) incQuantCount() {
	disp.quantCount++
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
	disp.environments[name] = newEnvironment(disp, name, false)
	return disp.environments[name]
}

func (disp *Dispatcher) createEnvironment(name string, builtin bool) error {
	if disp.getEnvironment_(name) != nil {
		return fmt.Errorf("environment '%v' already exists", name)
	}
	disp.environments[name] = newEnvironment(disp, name, builtin)
	return nil
}

func (disp *Dispatcher) CreateEnvironment(name string) error {
	if !disp.generalLock.Acquire(disp.timeout) {
		return fmt.Errorf("request lock timeout: can't create environment")
	}
	defer disp.generalLock.Release()
	return disp.createEnvironment(name, false)
}

// executes 'join' and 'affect' of the entity
func (disp *Dispatcher) Attach(entity *Entity, joins, affects map[string]int) error {
	if !disp.generalLock.Acquire(disp.timeout) {
		return fmt.Errorf("acquire lock timeout: can't attach entity to environment")
	}
	defer disp.generalLock.Release()

	for envName, limit := range joins {
		env := disp.getOrCreateEnvironment_(envName)
		if err := env.join(entity, limit); err != nil {
			return err
		}
	}
	for envName, delay := range affects {
		env := disp.getOrCreateEnvironment_(envName)
		if err := env.affect(entity, delay); err != nil {
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
	logf(5, "deleted environment '%v'", envName)
	return nil
}

// sends effect to the environment and thus asynchronously starts first wave in the quant.
// parameter 'waveMode' controls how the process ends:
//    waveMode = true:
//    	process stops after first wave is completed.
//      All the intermediate environment values are stored in wave coordinator
//      To continue there are two valid ways: WaveNext and WaveRun
//  	Usually used only in the debug mode
//    waveMode = false
//      process stops at the end of the quant

func (disp *Dispatcher) QuantStart(envName string, effect Trits, waveMode bool, onQuantFinish func()) error {
	env := disp.getEnvironment_(envName)
	if env == nil || env.invalid {
		return fmt.Errorf("can't find environment '%v'", envName)
	}
	return disp.quantStart(env, effect, waveMode, onQuantFinish)
}

func (disp *Dispatcher) quantStart(env *environment, effect Trits, waveMode bool, onQuantFinish func()) error {
	if disp.waveCoo.isWaveMode() {
		return fmt.Errorf("wave is already running")
	}
	var err error
	if effect, err = env.adjustEffect(effect); err != nil {
		return err
	}

	disp.waveCoo.setWaveMode(waveMode)
	disp.quantWG.Add(1)

	env.effectChan <- effect

	go func() {
		env.dispatcher.quantWG.Wait()
		disp.waveCoo.setWaveMode(false)
		if onQuantFinish != nil {
			onQuantFinish()
		}
	}()
	return nil
}

// if in waveMode, continues to the next wave and stops

func (disp *Dispatcher) WaveNext() error {
	if !disp.waveCoo.isWaveMode() {
		return fmt.Errorf("not in wave mode")
	}
	disp.waveCoo.runWave()
	return nil
}

// if in waveMode, continues to the next mode and stops at the end of the quant

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

type EnvironmentStatus struct {
	Size           int64
	JoinedEntities []string
	AffectedBy     []string
}

func (disp *Dispatcher) EnvironmentInfo() map[string]*EnvironmentStatus {
	ret := make(map[string]*EnvironmentStatus)

	for name, env := range disp.environments {
		envInfo := &EnvironmentStatus{
			Size:           env.size,
			JoinedEntities: make([]string, 0, len(env.joins)),
			AffectedBy:     make([]string, 0, len(env.affects)),
		}
		for _, joinData := range env.joins {
			envInfo.JoinedEntities = append(envInfo.JoinedEntities,
				fmt.Sprintf("%v(%v)", joinData.entity.GetName(), joinData.limit))
		}
		for _, affectData := range env.affects {
			envInfo.AffectedBy = append(envInfo.AffectedBy,
				fmt.Sprintf("%v(%v)", affectData.entity.GetName(), affectData.delay))
		}
		ret[name] = envInfo
	}
	return ret
}
