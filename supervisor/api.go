package supervisor

import (
	"fmt"
	"github.com/Workiva/go-datastructures/queue"
	. "github.com/iotaledger/iota.go/trinary"
	"time"
)

// Public supervisor API

func NewSupervisor(lockTimeout time.Duration) *Supervisor {
	ret := &Supervisor{
		queue:        queue.New(5),
		environments: make(map[string]*environment),
		accessLock:   newSema(),
		timeout:      lockTimeout,
	}
	ret.accessLock.acquire(-1) // supervisor starts locked
	go ret.supervisorInputLoop()
	return ret
}

type EntityOpts struct {
	Name    string
	InSize  int64
	OutSize int64
	Core    EntityCore
}

func (sv *Supervisor) NewEntity(opt EntityOpts) *Entity {
	ret := &Entity{
		supervisor: sv,
		name:       opt.Name,
		inSize:     opt.InSize,
		outSize:    opt.OutSize,
		affecting:  make([]*affectEntData, 0),
		joined:     make([]*environment, 0),
		core:       opt.Core,
	}
	return ret
}

func (sv *Supervisor) GetQuantCount() int64 {
	sv.quantCountMutex.RLock()
	defer sv.quantCountMutex.RUnlock()
	return sv.quantCount
}

func (sv *Supervisor) CreateEnvironment(name string) error {
	if !sv.accessLock.acquire(sv.timeout) {
		return fmt.Errorf("request lock timeout: can't create environment")
	}
	defer sv.accessLock.release()
	return sv.createEnvironment(name)
}

// executes 'join' and 'affect' of the entity
func (sv *Supervisor) Attach(entity *Entity, joins, affects map[string]int) error {
	if !sv.accessLock.acquire(sv.timeout) {
		return fmt.Errorf("acquire lock timeout: can't attach entity to environment")
	}
	defer sv.accessLock.release()

	for envName, limit := range joins {
		env := sv.getOrCreateEnvironment(envName)
		if err := env.join(entity, limit); err != nil {
			return err
		}
	}
	for envName, delay := range affects {
		env := sv.getOrCreateEnvironment(envName)
		if err := env.affect(entity, delay); err != nil {
			return err
		}
	}
	return nil
}

func (sv *Supervisor) Join(envName string, entity *Entity, limit int) error {
	return sv.Attach(entity, map[string]int{envName: limit}, nil)
}

func (sv *Supervisor) Affect(envName string, entity *Entity, delay int) error {
	return sv.Attach(entity, nil, map[string]int{envName: delay})
}

func (sv *Supervisor) DeleteEnvironment(envName string) error {
	if !sv.accessLock.acquire(sv.timeout) {
		return fmt.Errorf("request lock timeout: can't delete environment")
	}
	defer sv.accessLock.release()

	env, ok := sv.environments[envName]
	if !ok {
		return fmt.Errorf("can't find environment '%v'", envName)
	}
	env.invalidate()
	delete(sv.environments, envName)
	logf(5, "deleted environment '%v'", envName)
	return nil
}

func (sv *Supervisor) PostEffect(envName string, effect Trits, delay int) error {
	return sv.postEffect(envName, nil, effect, delay, true)
}

// calls doFunct if supervisor becomes idle within 'timeout'
// doFunct will be called upon release of the semaphore in the same goroutine.
// The doFunct itself must take care about locking the supervisor if needed
func (sv *Supervisor) DoIfIdle(timeout time.Duration, doFunct func()) bool {
	if !sv.accessLock.acquire(timeout) {
		return false
	}
	sv.accessLock.release()
	doFunct()
	return true
}

func (sv *Supervisor) DoOnIdle(doFunct func()) {
	for !sv.DoIfIdle(1*time.Second, doFunct) {
	}
}

type EnvironmentInfo struct {
	Size           int64
	JoinedEntities []string
	AffectedBy     []string
}

func (sv *Supervisor) EnvironmentInfo() map[string]*EnvironmentInfo {
	if !sv.accessLock.acquire(sv.timeout) {
		return nil
	}
	defer sv.accessLock.release()

	ret := make(map[string]*EnvironmentInfo)

	for name, env := range sv.environments {
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
	return ent.supervisor.GetQuantCount()
}
