package dispatcher_tests

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/dispatcher"
	"github.com/lunfardo314/goq/utils"
	"math"
	"testing"
	"time"
)

var dispatcher = NewDispatcher(1 * time.Second)

type mockEntityCore struct {
	name  string
	state int
}

func newMockEntityCore(name string) *mockEntityCore {
	return &mockEntityCore{
		name: name,
	}
}

func newMockEntity(id int) *Entity {
	name := fmt.Sprintf("mock_%v", id)
	return dispatcher.NewEntity(EntityOpts{
		Name:    name,
		InSize:  9,
		OutSize: 9,
		Core:    newMockEntityCore(name),
	})
}

func (core *mockEntityCore) Call(args Trits, res Trits) bool {
	resTmp := AddTrits(args, Trits{1, 0})
	copy(res, resTmp)
	core.state++
	return false
}

func envName(id int) string {
	return fmt.Sprintf("mock_environment_#%v", id)
}

const postTimes = 100000

func TestPostEffect0(t *testing.T) {
	if err := dispatcher.CreateEnvironment(envName(0)); err != nil {
		t.Errorf("%v", err)
		return
	}
	entity := newMockEntity(0)

	if err := dispatcher.Attach(entity, map[string]int{envName(0): 1}, nil); err != nil {
		t.Errorf("%v", err)
		return
	}

	start := utils.UnixMsNow()
	for i := 0; i < postTimes; i++ {
		if err := dispatcher.PostEffect(envName(0), Trits{0}, 0); err != nil {
			t.Errorf("%v", err)
			return
		}
	}
	durationSec := float64(utils.UnixMsNow()-start) / 1000
	fmt.Printf("Speed %v posts per second", math.Round(postTimes/durationSec))

	core := entity.GetCore().(*mockEntityCore)
	timeout := 5 * time.Second
	if !dispatcher.CallIfIdle(timeout, func() {
		if core.state != postTimes {
			t.Errorf("failed with wrong state %v != expected 1", core.state)
		}
		if err := dispatcher.DeleteEnvironment(envName(0)); err != nil {
			t.Errorf("failed while deleting environment: %v", err)
		}
		if len(dispatcher.EnvironmentInfo()) != 0 {
			t.Errorf("expected 0 environments left")
		}
	}) {
		t.Errorf("failed with timeout of %v", timeout)
	}
}
