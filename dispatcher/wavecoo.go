package dispatcher

import (
	. "github.com/iotaledger/iota.go/trinary"
	"sync"
)

type WaveCoordinator struct {
	sync.RWMutex
	waves    map[*environment]*waveCmd
	chIn     chan *waveCmd
	waveMode bool
}

type waveCmd struct {
	environment *environment
	value       Trits
	wg          *sync.WaitGroup
}

func NewWaveCoordinator() *WaveCoordinator {
	ret := &WaveCoordinator{
		waves: make(map[*environment]*waveCmd),
		chIn:  make(chan *waveCmd),
	}
	go ret.loop()
	return ret
}

func (wcoo *WaveCoordinator) loop() {
	for wr := range wcoo.chIn {
		wcoo.Lock()
		if wcoo.waveMode {
			if wr == nil {
				for _, r := range wcoo.waves {
					r.wg.Done()
				}
				wcoo.waves = make(map[*environment]*waveCmd)
			} else {
				wcoo.waves[wr.environment] = wr
			}
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

func (wcoo *WaveCoordinator) runWave() {
	if wcoo == nil {
		return
	}
	wcoo.chIn <- nil
}

func (wcoo *WaveCoordinator) setWaveMode(mode bool) {
	if wcoo == nil {
		return
	}
	wcoo.Lock()
	defer wcoo.Unlock()
	wcoo.waveMode = mode
}

func (wcoo *WaveCoordinator) isWaveMode() bool {
	if wcoo == nil {
		return false
	}
	wcoo.RLock()
	defer wcoo.RUnlock()
	return wcoo.waveMode
}
