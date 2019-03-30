package supervisor

import (
	"fmt"
	"github.com/Workiva/go-datastructures/queue"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/cfg"
	"time"
)

// Public thread safe supervisor API

// Create new instance of the supervisor

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

// create new Entity

func (sv *Supervisor) NewEntity(name string, inSize, outSize int, core EntityCore) (*Entity, error) {
	if outSize < 1 || inSize < 0 {
		return nil, fmt.Errorf("must be: output size > 0, input size >= 0")
	}
	ret := &Entity{
		supervisor: sv,
		name:       name,
		inSize:     inSize,
		outSize:    outSize,
		affecting:  make([]*affectEntData, 0),
		joined:     make([]*environment, 0),
		core:       core,
	}
	return ret, nil
}

// return current quant count in thread safe manner
// quant count changing in async way therefore makes sense only
// when all quants are stopped
// Used mainly for tests

func (sv *Supervisor) GetQuantCount() int64 {
	sv.quantCountMutex.RLock()
	defer sv.quantCountMutex.RUnlock()
	return sv.quantCount
}

// creates new environment is doesn't exist another with the same name

func (sv *Supervisor) CreateEnvironment(name string) error {
	if !sv.accessLock.acquire(sv.timeout) {
		return fmt.Errorf("request lock timeout: can't create environment")
	}
	defer sv.accessLock.release()
	return sv.createEnvironment(name)
}

// executes several 'joins' and 'affects' of the entity to environments
// environments are created if necessary
// for 'joins' and 'affects' params inout is map[string]int with respective environments names as keys.
// int values of the map entry are interpreted as 'limit' for joins and 'delay' for affects

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

// shorter version for one 'join' only
func (sv *Supervisor) Join(envName string, entity *Entity, limit int) error {
	return sv.Attach(entity, map[string]int{envName: limit}, nil)
}

// shorter version for one 'affect' only
func (sv *Supervisor) Affect(envName string, entity *Entity, delay int) error {
	return sv.Attach(entity, nil, map[string]int{envName: delay})
}

// deletes environment from the supervisor.
//    - stops environment loop goroutine,
//    - "unjoins" and "unaffects" related entities (which may result in complete stop of the environment)
//    - marks environment as invalid

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
	Logf(5, "deleted environment '%v'", envName)
	return nil
}

// delete all environments. Useful when reloading module

func (sv *Supervisor) ClearEnvironments() error {
	if !sv.accessLock.acquire(sv.timeout) {
		return fmt.Errorf("request lock timeout: can't delete environment")
	}
	defer sv.accessLock.release()

	for _, env := range sv.environments {
		env.invalidate()
	}
	sv.environments = make(map[string]*environment)
	Logf(5, "supervisor: all environments were deleted")
	return nil
}

// Posts effect to the main queue

func (sv *Supervisor) PostEffect(envName string, effect Trits, delay int) error {
	return sv.postEffect(envName, nil, effect, delay, true)
}

// calls doFunct if supervisor becomes idle (= releases lock) within 'timeout'
// doFunc will be called upon release of the semaphore outside the locked section.
// The doFunc itself must take care about locking the supervisor if needed

func (sv *Supervisor) DoIfIdle(timeout time.Duration, doFunc func()) bool {
	if !sv.accessLock.acquire(timeout) {
		return false
	}
	sv.accessLock.release()
	doFunc()
	return true
}

// loops until supervisor becomes idle and calls doFunc
func (sv *Supervisor) DoOnIdle(doFunc func()) {
	for !sv.DoIfIdle(1*time.Second, doFunc) {
	}
}

type EnvironmentInfo struct {
	JoinedEntities []string
	AffectedBy     []string
}

// returns info about current configuration of the supervisor

func (sv *Supervisor) EnvironmentInfo() map[string]*EnvironmentInfo {
	if !sv.accessLock.acquire(sv.timeout) {
		return nil
	}
	defer sv.accessLock.release()

	ret := make(map[string]*EnvironmentInfo)

	for name, env := range sv.environments {
		envInfo := &EnvironmentInfo{
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
