package dispatcher_tests

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/dispatcher"
	"github.com/lunfardo314/goq/utils"
	"testing"
	"time"
)

var dispatcher = NewDispatcher(1 * time.Second)

type mockEntityCore struct {
	name  string
	state int64
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
	core.state += TritsToInt(args)
	return false
}

func envName(id int) string {
	return fmt.Sprintf("mock_environment_#%v", id)
}

const postTimes0 = 1000000

func TestPostEffect0(t *testing.T) {
	fmt.Printf("Test 0: posting %v effects to one mock environment\n", postTimes0)
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
	for i := 0; i < postTimes0; i++ {
		if err := dispatcher.PostEffect(envName(0), Trits{1}, 0); err != nil {
			t.Errorf("%v", err)
			return
		}
	}
	durationSec := float64(utils.UnixMsNow()-start) / 1000
	fmt.Printf("Posting %v posts per second\n", int(postTimes0/durationSec))

	core := entity.GetCore().(*mockEntityCore)
	dispatcher.CallWhenIdle(func() {
		durationSec := float64(utils.UnixMsNow()-start) / 1000
		fmt.Printf("Processing speed %v effects per second\n", int(postTimes0/durationSec))

		if core.state != postTimes0 {
			t.Errorf("failed with wrong state %v != expected %v", core.state, postTimes0)
		}
		if err := dispatcher.DeleteEnvironment(envName(0)); err != nil {
			t.Errorf("failed while deleting environment: %v", err)
		}
		if len(dispatcher.EnvironmentInfo()) != 0 {
			t.Errorf("expected 0 environments left")
		}
	})
}

const postTimes1 = 100000
const chainLen = 50

func TestPostEffect1(t *testing.T) {
	fmt.Printf("Test 1: posting %v effects to %v environments, chained with 'affect'\n", postTimes1, chainLen)

	var prev *Entity
	cores := make([]*mockEntityCore, 0, chainLen)

	for i := 0; i < chainLen; i++ {
		if err := dispatcher.CreateEnvironment(envName(i)); err != nil {
			t.Errorf("%v", err)
			return
		}
		entity := newMockEntity(i)

		if err := dispatcher.Attach(entity, map[string]int{envName(i): 1}, nil); err != nil {
			t.Errorf("%v", err)
			return
		}
		if prev != nil {
			if err := dispatcher.Attach(prev, nil, map[string]int{envName(i): 0}); err != nil {
				t.Errorf("%v", err)
				return
			}
		}
		cores = append(cores, entity.GetCore().(*mockEntityCore))
		prev = entity
	}
	start := utils.UnixMsNow()
	for i := 0; i < postTimes1; i++ {
		if err := dispatcher.PostEffect(envName(0), Trits{1}, 0); err != nil {
			t.Errorf("%v", err)
			return
		}
	}
	durationSec := float64(utils.UnixMsNow()-start) / 1000
	fmt.Printf("Posting %v posts per second\n", int(postTimes1/durationSec))

	dispatcher.CallWhenIdle(func() {
		durationSec := float64(utils.UnixMsNow()-start) / 1000
		fmt.Printf("Processing speed %v effects per second\n", int(postTimes1*chainLen/durationSec))

		for i, core := range cores {
			if core.state != int64(postTimes1*(i+1)) {
				t.Errorf("failed with wrong state %v != expected %v", core.state, postTimes1+i)
			}
		}
		for i := 0; i < chainLen; i++ {
			if err := dispatcher.DeleteEnvironment(envName(i)); err != nil {
				t.Errorf("failed while deleting environment '%v': %v", envName(i), err)
			}
		}
		if len(dispatcher.EnvironmentInfo()) != 0 {
			t.Errorf("expected 0 environments left")
		}
	})
}
