package entities

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
	"sync"
)

type Tritmap struct {
	sync.RWMutex
	Width  int
	Height int
	data   Trits
}

func NewTritmap(width, height int) (*Tritmap, error) {
	if width*height < 1 {
		return nil, fmt.Errorf("NewTritmap: must be at least 1x1")
	}
	return &Tritmap{
		Width:  width,
		Height: height,
		data:   make(Trits, width*height),
	}, nil
}

func (m *Tritmap) Get(x, y int) (int8, error) {
	m.RLock()
	defer m.RUnlock()

	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return 0, fmt.Errorf("Tritmap.Get: out of bounds")
	}
	return m.data[x*m.Height+y], nil
}

func (m *Tritmap) Put(x, y int, value int8) error {
	m.Lock()
	defer m.Unlock()

	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return fmt.Errorf("Tritmap.Get: out of bounds")
	}
	m.data[x*m.Height+y] = value
	return nil
}
