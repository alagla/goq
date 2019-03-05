package dispatcher

import (
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/abstract"
	. "github.com/lunfardo314/goq/utils"
	"sync"
	"time"
)

type Environment struct {
	sync.RWMutex
	name       string
	joins      []*Join
	affectChan chan Trits // where all effects are sent
}

type Join struct {
	environment *Environment
	function    FuncDefInterface
	inChan      chan Trits // each joined function is listening to this
}

func NewEnvironment(name string) *Environment {
	ret := &Environment{
		name:       name,
		joins:      make([]*Join, 0),
		affectChan: make(chan Trits),
	}
	go ret.AffectLoop()
	return ret
}

func (env *Environment) Join(fun FuncDefInterface) {
	env.Lock()
	defer env.Unlock()
	join := &Join{
		environment: env,
		function:    fun,
		inChan:      make(chan Trits),
	}
	env.joins = append(env.joins, join)
}

func (env *Environment) Affect(fun FuncDefInterface) {
	// todo inform function it must affect the environment
	//env.Lock()
	//defer env.Unlock()
	//affect := &Affect{
	//	environment: env,
	//	function: fun,
	//}
	//env.affects = append(env.affects, affect)
}

func (env *Environment) PostEffect(effect Trits) {
	env.affectChan <- effect
}

func (env *Environment) AffectLoop() {
	for effect := range env.affectChan {
		logf(3, "Value '%v' reached environment '%v'",
			TritsToString(effect), env.name)
		env.processEffect(effect)
	}
}

func (env *Environment) processEffect(effect Trits) {
	env.RLock()
	defer env.RUnlock()
	for _, join := range env.joins {
		go join.processEffect(effect)
	}
}

func (join *Join) processEffect(effect Trits) {
	for r := range join.inChan {
		logf(1, "Value '%v' triggered function %v in environment '%v'",
			TritsToString(r), join.function.GetName(), join.environment.name)
		// TODO calculate function
		time.Sleep(1 * time.Second)
	}
}
