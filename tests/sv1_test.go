package tests

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/supervisor"
	"testing"
)

func TestAffectDelay3(t *testing.T) {
	const delayOnStart = 5
	const delayAffect = 1000
	fmt.Printf("--------------\nSupervisor test 3: testing affect with delay on start = %v, delay affect = %v\n", delayOnStart, delayAffect)
	if !check0environments(t) {
		return
	}

	entity0, err := newMockEntity(0, -1)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	core0 := entity0.GetCore().(*mockEntityCore)
	if err := sv.Join(envName(0), entity0, 1); err != nil {
		t.Errorf("%v", err)
		return
	}

	entity1, err := newMockEntity(1, -1)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	core1 := entity1.GetCore().(*mockEntityCore)
	if err := sv.Join(envName(1), entity1, 1); err != nil {
		t.Errorf("%v", err)
		return
	}
	if err := sv.Affect(envName(1), entity0, delayAffect); err != nil {
		t.Errorf("%v", err)
		return
	}

	saveQuantCount := sv.GetQuantCount()
	if saveQuantCount < 0 {
		t.Errorf("GetQuantCount failed")
		return
	}
	if err := sv.PostEffect(envName(0), Trits{0}, delayOnStart); err != nil {
		t.Errorf("%v", err)
		return
	}
	sv.DoOnIdle(func() {
		expected0 := int64(1)
		if core0.state != expected0 {
			t.Errorf("failed with wrong state %v != expected %v", core0.state, expected0)
		}
		expected1 := int64(2)
		if core1.state != expected1 {
			t.Errorf("failed with wrong state %v != expected %v", core1.state, expected1)
		}
		expectedLastQuant := saveQuantCount + delayOnStart + delayAffect
		if core1.lastQuant != expectedLastQuant {
			fmt.Printf("failed with last quant %v != expected %v\n", core1.lastQuant, expectedLastQuant)
		}
		if err := sv.DeleteEnvironment(envName(0)); err != nil {
			t.Errorf("%v", err)
		}
		if err := sv.DeleteEnvironment(envName(1)); err != nil {
			t.Errorf("%v", err)
		}
	})
}

func TestJoinLimit4(t *testing.T) {
	const joinLimit = 1
	const chainLen = 20
	const maxCount = 22
	fmt.Printf("----------------\nSupervisor test 4: test join limit. Post 1 effect to %v environments connected in cycle. Max count = %v '. Join limit = %v\n",
		chainLen, maxCount, joinLimit)
	if !check0environments(t) {
		return
	}

	var prev *Entity
	cores := make([]*mockEntityCore, 0, chainLen)

	// generating line chain
	var entity *Entity
	var err error
	for i := 0; i < chainLen; i++ {
		// environments created when needed by attach
		entity, err = newMockEntity(i, maxCount)
		if err != nil {
			t.Errorf("%v", err)
			return
		}

		if err := sv.Join(envName(i), entity, joinLimit); err != nil {
			t.Errorf("%v", err)
			return
		}
		if prev != nil {
			if err := sv.Attach(prev, nil, map[string]int{envName(i): 0}); err != nil {
				t.Errorf("%v", err)
				return
			}
		}
		cores = append(cores, entity.GetCore().(*mockEntityCore))
		prev = entity
	}
	// connecting last will affect first
	if err := sv.Attach(entity, nil, map[string]int{envName(0): 0}); err != nil {
		t.Errorf("%v", err)
		return
	}
	saveQuantCount := sv.GetQuantCount()

	if err := sv.PostEffect(envName(0), Trits{0}, 0); err != nil {
		t.Errorf("%v", err)
		return
	}
	sv.DoOnIdle(func() {
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
				expected := int64(maxCount - chainLen + i + 1)
				if cores[idx].state != expected {
					t.Errorf("failed with wrong state %v != expected %v", cores[idx].state, expected)
				}
				numQuants := cores[idx].lastQuant + 1 - saveQuantCount
				//fmt.Printf("   '%v' state: %v, num quants: %v\n",
				//	cores[idx].name, cores[idx].state, numQuants)

				// valid only with joinLimit == 1
				// TODO for different combinations
				if joinLimit == 1 && cores[idx].state != int64(numQuants) {
					t.Errorf("failed with wrong numQuants %v != expected %v",
						numQuants, cores[idx].state)
				}
			}
		}
		for i := 0; i < chainLen; i++ {
			if err := sv.DeleteEnvironment(envName(i)); err != nil {
				t.Errorf("failed while deleting environment '%v': %v", envName(i), err)
			}
		}
	})
}
