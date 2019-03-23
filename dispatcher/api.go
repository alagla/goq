package dispatcher

import (
	"fmt"
	"github.com/Workiva/go-datastructures/queue"
	. "github.com/iotaledger/iota.go/trinary"
	"time"
)

// Dispatcher API

func NewDispatcher(lockTimeout time.Duration) *Dispatcher {
	ret := &Dispatcher{
		queue:        queue.New(5),
		environments: make(map[string]*environment),
		accessLock:   newSema(),
		timeout:      lockTimeout,
	}
	ret.accessLock.acquire(-1)
	go ret.dispatcherInputLoop()
	return ret
}

type EntityOpts struct {
	Name    string
	InSize  int64
	OutSize int64
	Core    EntityCore
}

func (disp *Dispatcher) NewEntity(opt EntityOpts) *Entity {
	ret := &Entity{
		dispatcher: disp,
		name:       opt.Name,
		inSize:     opt.InSize,
		outSize:    opt.OutSize,
		affecting:  make([]*affectEntData, 0),
		joined:     make([]*environment, 0),
		core:       opt.Core,
	}
	return ret
}

func (disp *Dispatcher) GetQuantCount() int64 {
	disp.quantCountMutex.RLock()
	defer disp.quantCountMutex.RUnlock()
	return disp.quantCount
}

func (disp *Dispatcher) CreateEnvironment(name string) error {
	if !disp.accessLock.acquire(disp.timeout) {
		return fmt.Errorf("request lock timeout: can't create environment")
	}
	defer disp.accessLock.release()
	return disp.createEnvironment(name, false)
}

// executes 'join' and 'affect' of the entity
func (disp *Dispatcher) Attach(entity *Entity, joins, affects map[string]int) error {
	if !disp.accessLock.acquire(disp.timeout) {
		return fmt.Errorf("acquire lock timeout: can't attach entity to environment")
	}
	defer disp.accessLock.release()

	for envName, limit := range joins {
		env := disp.getOrCreateEnvironment(envName)
		if err := env.join(entity, limit); err != nil {
			return err
		}
	}
	for envName, delay := range affects {
		env := disp.getOrCreateEnvironment(envName)
		if err := env.affect(entity, delay); err != nil {
			return err
		}
	}
	return nil
}

func (disp *Dispatcher) Join(envName string, entity *Entity, limit int) error {
	return disp.Attach(entity, map[string]int{envName: limit}, nil)
}

func (disp *Dispatcher) Affect(envName string, entity *Entity, delay int) error {
	return disp.Attach(entity, nil, map[string]int{envName: delay})
}

func (disp *Dispatcher) DeleteEnvironment(envName string) error {
	if !disp.accessLock.acquire(disp.timeout) {
		return fmt.Errorf("request lock timeout: can't delete environment")
	}
	defer disp.accessLock.release()

	env, ok := disp.environments[envName]
	if !ok {
		return fmt.Errorf("can't find environment '%v'", envName)
	}
	env.invalidate()
	delete(disp.environments, envName)
	logf(5, "deleted environment '%v'", envName)
	return nil
}

func (disp *Dispatcher) PostEffect(envName string, effect Trits, delay int) error {
	return disp.postEffect(envName, nil, effect, delay, true)
}

// calls doFunct if dispatcher becomes idle within 'timeout'
// doFunct will be called upon release of the semaphore in the same goroutine.
// The doFunct itself must take care about locking the dispatcher if needed
func (disp *Dispatcher) DoIfIdle(timeout time.Duration, doFunct func()) bool {
	if !disp.accessLock.acquire(timeout) {
		return false
	}
	disp.accessLock.release()
	doFunct()
	return true
}

func (disp *Dispatcher) DoOnIdle(doFunct func()) {
	for !disp.DoIfIdle(1*time.Second, doFunct) {
	}
}

type EnvironmentInfo struct {
	Size           int64
	JoinedEntities []string
	AffectedBy     []string
}

func (disp *Dispatcher) EnvironmentInfo() map[string]*EnvironmentInfo {
	if !disp.accessLock.acquire(disp.timeout) {
		return nil
	}
	defer disp.accessLock.release()

	ret := make(map[string]*EnvironmentInfo)

	for name, env := range disp.environments {
		envInfo := &EnvironmentInfo{
			Size:           env.size,
			JoinedEntities: make([]string, 0, len(env.joins)),
			AffectedBy:     make([]string, 0, len(env.affects)),
		}
		for _, joinData := range env.joins {
			envInfo.JoinedEntities = append(envInfo.JoinedEntities,
				fmt.Sprintf("%v(%v)", joinData.entity.name, joinData.limit))
		}
		for _, ent := range env.affects {
			envInfo.AffectedBy = append(envInfo.AffectedBy, ent.name)
		}
		ret[name] = envInfo
	}
	return ret
}

// Entity API

func (ent *Entity) GetCore() EntityCore {
	return ent.core
}

// for calls from within entity core. For debugging
func (ent *Entity) GetQuantCount() int64 {
	return ent.dispatcher.GetQuantCount()
}
