package supervisor

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	"time"
)

var dispatcher = NewSupervisor(1 * time.Second)

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
	ret, err := dispatcher.NewEntity(name, 81, 81, core)
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
