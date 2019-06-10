package main

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
)

type Tritmap struct {
	Width  int
	Height int
	Data   Trits
}

func NewTritmap(width, height int) (*Tritmap, error) {
	if width*height < 1 {
		return nil, fmt.Errorf("NewTritmap: must be at least 1x1")
	}
	return &Tritmap{
		Width:  width,
		Height: height,
		Data:   make(Trits, width*height),
	}, nil
}

func MustNewTritmap(width, height int) *Tritmap {
	ret, err := NewTritmap(width, height)
	if err != nil {
		panic(err)
	}
	return ret
}

func (m *Tritmap) Get(x, y int) (int8, error) {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return 0, fmt.Errorf("Tritmap.Get: out of bounds")
	}
	return m.Data[y*m.Height+x], nil
}

func (m *Tritmap) Put(x, y int, value int8) error {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return fmt.Errorf("Tritmap.Get: out of bounds")
	}
	m.Data[y*m.Height+x] = value
	return nil
}
