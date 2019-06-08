package supervisor

import (
	"fmt"
	"github.com/Workiva/go-datastructures/queue"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/cfg"
	"time"
)

// Public thread safe Supervisor API

// Creates new instance of the Supervisor

func NewSupervisor(name string, lockTimeout time.Duration) *Supervisor {
	ret := &Supervisor{
		Name:         name,
		queue:        queue.New(5),
		environments: make(map[string]*environment),
		accessLock:   newSema(),
		timeout:      lockTimeout,
	}
	ret.accessLock.acquire(-1) // Supervisor starts locked
	go ret.supervisorInputLoop()
	return ret
}

// create new Entity
// Params:
//    - Name, used only for tracing
//    - inSize, expected size of input trit vector. O means any not nil
//    - outSize, size of output trit vector
//    - core, an object which implement EntityCore interface. It calculates output trits from inputs

func (sv *Supervisor) NewEntity(name string, inSize, outSize int, core EntityCore) (*Entity, error) {
	if outSize < 1 || inSize < 0 {
		return nil, fmt.Errorf("must be: output size > 0, input size >= 0")
	}
	ret := &Entity{
		Supervisor: sv,
		Name:       name,
		inSize:     inSize,
		outSize:    outSize,
		affecting:  make([]*affectEntData, 0),
		joined:     make([]*environment, 0),
		core:       core,
	}
	return ret, nil
}

// returns current quant count.
// Quant count is changing in async way therefore makes sense only
// when Supervisor is in idle state, i e input queue is empty
// Used when delay is posted with effect and for testing

func (sv *Supervisor) GetQuantCount() int64 {
	sv.quantCountMutex.RLock()
	ret := sv.quantCount
	sv.quantCountMutex.RUnlock()
	return ret
}

// creates new environment if doesn't exist another with the same Name

func (sv *Supervisor) CreateEnvironment(name string) error {
	if !sv.accessLock.acquire(sv.timeout) {
		return fmt.Errorf("request lock timeout: can't create environment")
	}
	ret := sv.createEnvironment(name)
	sv.accessLock.release()
	return ret
}

// executes several 'joins' and 'affects' of the entity to environments
// Environments are created if necessary
// Params 'joins' and 'affects' expect map[string]int with respective environments names as keys.
// int values of the map entry are interpreted as 'limit' for joins and 'delay' for affects

func (sv *Supervisor) Attach(entity *Entity, joins, affects map[string]int) error {
	if !sv.accessLock.acquire(sv.timeout) {
		return fmt.Errorf("acquire lock timeout: can't attach entity to environment")
	}

	for envName, limit := range joins {
		env := sv.getOrCreateEnvironment(envName)
		if err := env.join(entity, limit); err != nil {
			sv.accessLock.release()
			return err
		}
	}
	for envName, delay := range affects {
		env := sv.getOrCreateEnvironment(envName)
		if err := env.affect(entity, delay); err != nil {
			sv.accessLock.release()
			return err
		}
	}
	sv.accessLock.release()
	return nil
}

// shorter version for one 'join' only: entity is joined (subscribes) to the environment 'envName'
// Entity automatically starts it's input loop/goroutine upon first join

func (sv *Supervisor) Join(envName string, entity *Entity, limit int) error {
	return sv.Attach(entity, map[string]int{envName: limit}, nil)
}

// shorter version for one 'affect' only: entity starts to 'affect' (post results to) the environment 'envName'

func (sv *Supervisor) Affect(envName string, entity *Entity, delay int) error {
	return sv.Attach(entity, nil, map[string]int{envName: delay})
}

// deletes environment from the Supervisor.
//    - stops environment loop goroutine,
//    - "unjoins" and "unaffects" related entities.
//    - marks environment as invalid. Any effects for this environment which may remain in the queue will be ignored
// After environment is deleted, some entities may become completely detached from the Supervisor.
// Input loop and goroutine of such entity is automatically stopped.

func (sv *Supervisor) DeleteEnvironment(envName string) error {
	if !sv.accessLock.acquire(sv.timeout) {
		return fmt.Errorf("request lock timeout: can't delete environment")
	}

	env, ok := sv.environments[envName]
	if !ok {
		sv.accessLock.release()
		return fmt.Errorf("can't find environment '%v'", envName)
	}
	env.invalidate()
	delete(sv.environments, envName)
	Logf(5, "deleted environment '%v'", envName)

	sv.accessLock.release()
	return nil
}

// delete all environments. Useful when reloading a module

func (sv *Supervisor) ClearEnvironments() error {
	if !sv.accessLock.acquire(sv.timeout) {
		return fmt.Errorf("request lock timeout: can't delete environment")
	}

	for _, env := range sv.environments {
		env.invalidate()
	}
	sv.environments = make(map[string]*environment)
	Logf(5, "Supervisor: all environments were deleted")
	sv.accessLock.release()
	return nil
}

// Posts effect to the main queue. This can be done by any code, not only from entity

func (sv *Supervisor) PostEffect(envName string, effect Trits, delay int) error {
	return sv.postEffect(envName, nil, effect, delay, true)
}

// executes `doFunct` if Supervisor becomes idle (= releases lock) within 'timeout'
// `doFunc` will be called upon release of the semaphore outside of the locked section.
// Note, that there's no guarantee that Supervisor will be idle (unlocked) during execution of `doFunc`

func (sv *Supervisor) DoIfIdle(timeout time.Duration, doFunc func()) bool {
	if !sv.accessLock.acquire(timeout) {
		return false
	}
	sv.accessLock.release()
	doFunc()
	return true
}

// loops until Supervisor becomes idle and calls doFunc
func (sv *Supervisor) DoOnIdle(doFunc func()) {
	for !sv.DoIfIdle(1*time.Second, doFunc) {
	}
}

type EnvironmentInfo struct {
	JoinedEntities []string
	AffectedBy     []string
}

// returns info about current configuration of the Supervisor

func (sv *Supervisor) EnvironmentInfo() map[string]*EnvironmentInfo {
	if !sv.accessLock.acquire(sv.timeout) {
		return nil
	}

	ret := make(map[string]*EnvironmentInfo)

	for name, env := range sv.environments {
		envInfo := &EnvironmentInfo{
			JoinedEntities: make([]string, 0, len(env.joins)),
			AffectedBy:     make([]string, 0, len(env.affects)),
		}
		for _, joinData := range env.joins {
			envInfo.JoinedEntities = append(envInfo.JoinedEntities,
				fmt.Sprintf("%v(%v)", joinData.entity.Name, joinData.limit))
		}
		for _, ent := range env.affects {
			envInfo.AffectedBy = append(envInfo.AffectedBy, ent.Name)
		}
		ret[name] = envInfo
	}
	sv.accessLock.release()
	return ret
}

// Entity API

func (ent *Entity) GetCore() EntityCore {
	return ent.core
}

// for calls from within entity core. For testing/debugging
func (ent *Entity) GetQuantCount() int64 {
	return ent.Supervisor.GetQuantCount()
}
