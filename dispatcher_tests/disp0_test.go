package dispatcher_tests

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/dispatcher"
	"github.com/lunfardo314/goq/utils"
	"testing"
)

const postTimes0 = 100000

func TestPostEffect0(t *testing.T) {
	fmt.Printf("\nTest 0: posting %v effects to one mock environment\n", postTimes0)
	if err := dispatcher.CreateEnvironment(envName(0)); err != nil {
		t.Errorf("%v", err)
		return
	}
	entity := newMockEntity(0, -1)

	if err := dispatcher.Attach(entity, map[string]int{envName(0): 1}, nil); err != nil {
		t.Errorf("%v", err)
		return
	}

	start := utils.UnixMsNow()
	for i := 0; i < postTimes0; i++ {
		if err := dispatcher.PostEffect(envName(0), Trits{0}, 0); err != nil {
			t.Errorf("%v", err)
			return
		}
	}
	if durationSec := float64(utils.UnixMsNow()-start) / 1000; durationSec > 0.01 {
		fmt.Printf("Posting %v posts per second\n", int(postTimes0/durationSec))
	}

	core := entity.GetCore().(*mockEntityCore)
	dispatcher.CallWhenIdle(func() {
		if durationSec := float64(utils.UnixMsNow()-start) / 1000; durationSec > 0.01 {
			fmt.Printf("Processing speed %v waves per second\n", int(postTimes0/durationSec))
		}

		if core.state != 1 {
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

const postTimes1 = 100
const chainLen1 = 500

func TestPostEffect1(t *testing.T) {
	fmt.Printf("\nTest 1: posting %v effects to %v environments chained in a line\n", postTimes1, chainLen1)

	var prev *Entity
	cores := make([]*mockEntityCore, 0, chainLen1)

	for i := 0; i < chainLen1; i++ {
		// environment will be created when needed by attach

		entity := newMockEntity(i, -1)

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
		if err := dispatcher.PostEffect(envName(0), Trits{0}, 0); err != nil {
			t.Errorf("%v", err)
			return
		}
	}
	if durationSec := float64(utils.UnixMsNow()-start) / 1000; durationSec > 0.01 {
		fmt.Printf("Posting %v posts per second\n", int(postTimes1/durationSec))
	}

	dispatcher.CallWhenIdle(func() {
		if durationSec := float64(utils.UnixMsNow()-start) / 1000; durationSec > 0.01 {
			fmt.Printf("Processing speed %v waves per second\n", int(postTimes1*chainLen1/durationSec))
		}
		for i, core := range cores {
			if core.state != int64(i+1) {
				t.Errorf("failed with wrong state %v != expected %v", core.state, postTimes1+i)
			}
		}
		for i := 0; i < chainLen1; i++ {
			if err := dispatcher.DeleteEnvironment(envName(i)); err != nil {
				t.Errorf("failed while deleting environment '%v': %v", envName(i), err)
			}
		}
		if len(dispatcher.EnvironmentInfo()) != 0 {
			t.Errorf("expected 0 environments left")
		}
	})
}

// attach entities in a cycle of chainLen2.
// Posted effect should go in rounds until state of the mock enty reaches maxCount
// Ten mock entity return null value

const chainLen2 = 500
const maxCount = chainLen2 + 100000 // must be maxCount >= chainLen2 for test to be correct

func TestPostEffect2(t *testing.T) {
	fmt.Printf("\nTest 2: posting 1 effect to environment '%v'.\n%v environments connected in cycle. Max count: %v '\n",
		envName(0), chainLen2, maxCount)

	var prev *Entity
	cores := make([]*mockEntityCore, 0, chainLen2)

	// generating line chain
	var entity *Entity

	for i := 0; i < chainLen2; i++ {
		// environments created when needed by attach
		entity = newMockEntity(i, maxCount)

		if err := dispatcher.Attach(entity, map[string]int{envName(i): maxCount}, nil); err != nil {
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
	// connecting last will affect first
	if err := dispatcher.Attach(entity, nil, map[string]int{envName(0): 0}); err != nil {
		t.Errorf("%v", err)
		return
	}

	start := utils.UnixMsNow()
	if err := dispatcher.PostEffect(envName(0), Trits{0}, 0); err != nil {
		t.Errorf("%v", err)
		return
	}
	dispatcher.CallWhenIdle(func() {
		durationSec := float64(utils.UnixMsNow()-start) / 1000
		fmt.Printf("Processing speed %v waves per second\n", int(maxCount/durationSec))

		idxStop := -1
		for i, core := range cores {
			if core.state > maxCount {
				t.Errorf("state can't be > %v", maxCount)
			} else {
				if core.state == maxCount {
					idxStop = i
				}
			}
		}

		if idxStop == -1 {
			t.Errorf("failed with wrong state")
		} else {
			for i := 0; i < len(cores); i++ {
				idx := (idxStop + i + 1) % len(cores)
				expected := int64(maxCount - chainLen2 + i + 1)
				if cores[idx].state != expected {
					t.Errorf("failed with wrong state %v != expected %v", cores[idx].state, expected)
				}
			}

		}
		for i := 0; i < chainLen2; i++ {
			if err := dispatcher.DeleteEnvironment(envName(i)); err != nil {
				t.Errorf("failed while deleting environment '%v': %v", envName(i), err)
			}
		}
		if len(dispatcher.EnvironmentInfo()) != 0 {
			t.Errorf("expected 0 environments left")
		}
	})
}
