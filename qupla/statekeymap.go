package qupla

import (
	. "github.com/iotaledger/iota.go/trinary"
	"hash/fnv"
	"sync"
)

type StateHashMap struct {
	sync.RWMutex
	theMap map[uint64]Trits
}

func newStateHashMap() *StateHashMap {
	return &StateHashMap{
		theMap: make(map[uint64]Trits),
	}
}

// using 64 bit hashing
//
// Possible hash collisions neglected !!!!

func hashCallTrace(callTrace []uint8) uint64 {
	h := fnv.New64()
	if _, err := h.Write(callTrace); err != nil {
		panic(err)
	}
	return h.Sum64()
}

func (hm *StateHashMap) getValue(key []uint8, nullSize int) Trits {
	hm.RLock()
	defer hm.RUnlock()

	hash := hashCallTrace(key)
	value, ok := hm.theMap[hash]
	if !ok {
		return PadTrits(Trits{0}, nullSize)
	}
	return value
}

func (hm *StateHashMap) storeValue(key []uint8, value Trits) {
	hm.Lock()
	defer hm.Unlock()

	hash := hashCallTrace(key)
	hm.theMap[hash] = value
}
