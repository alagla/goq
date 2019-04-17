package tests

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	"github.com/lunfardo314/goq/cfg"
	. "github.com/lunfardo314/goq/supervisor"
	"testing"
	"time"
)

var sv *Supervisor

func init() {
	cfg.Config.Verbosity = 0
	sv = NewSupervisor(1 * time.Second)
	//pinline := flag.Bool("inline", false, "use inline call optimisation")
	//flag.Parse()
	//cfg.Config.OptimizeFunCallsInline = *pinline
}

type mockEntityCore struct {
	name      string
	state     int64
	maxWaves  int64
	entity    *Entity
	lastQuant int64
}

func newMockEntityCore(name string, maxWaves int64) *mockEntityCore {
	return &mockEntityCore{
		name:      name,
		maxWaves:  maxWaves,
		lastQuant: -1,
	}
}

func newMockEntity(id int, maxCount int64) (*Entity, error) {
	name := fmt.Sprintf("mock_%v", id)
	core := newMockEntityCore(name, maxCount)
	ret, err := sv.NewEntity(name, 81, 81, core)
	if err != nil {
		return nil, err
	}
	core.entity = ret
	return ret, nil
}

func (core *mockEntityCore) Call(args Trits, res Trits) bool {
	// + 1 to arg -> result -> state
	core.state = TritsToInt(args) + 1
	copy(res, IntToTrits(core.state))
	core.lastQuant = core.entity.GetQuantCount() // for testing only
	return core.maxWaves > 0 && core.state >= core.maxWaves
}

func envName(id int) string {
	return fmt.Sprintf("mock_environment_#%v", id)
}

func check0environments(t *testing.T) bool {
	if len(sv.EnvironmentInfo()) != 0 {
		t.Errorf("expected 0 environments, found %v", len(sv.EnvironmentInfo()))
		return false
	}
	return true
}
