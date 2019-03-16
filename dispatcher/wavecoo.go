package dispatcher

import (
	. "github.com/iotaledger/iota.go/trinary"
	"sync"
)

type WaveCoordinator struct {
	sync.RWMutex
	waves map[*environment]*waveResult
	chIn  chan *waveResult
}

type waveResult struct {
	environment *environment
	value       Trits
	wg          *sync.WaitGroup
}

func NewWaveCoordinator() *WaveCoordinator {
	ret := &WaveCoordinator{
		waves: make(map[*environment]*waveResult),
		chIn:  make(chan *waveResult),
	}
	go ret.loop()
	return ret
}

func (wcoo *WaveCoordinator) loop() {
	for wr := range wcoo.chIn {
		wcoo.Lock()
		if wr == nil {
			for _, r := range wcoo.waves {
				r.wg.Done()
			}
			wcoo.waves = make(map[*environment]*waveResult)
		} else {
			wcoo.waves[wr.environment] = wr
		}
		wcoo.Unlock()
	}
}

func (wcoo *WaveCoordinator) values() map[string]Trits {
	if wcoo == nil {
		return nil
	}
	wcoo.RLock()
	defer wcoo.RUnlock()
	ret := make(map[string]Trits)
	for _, wr := range wcoo.waves {
		ret[wr.environment.GetName()] = wr.value
	}
	return ret
}

func (wcoo *WaveCoordinator) nextWave() {
	wcoo.chIn <- nil
}
