package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/cfg"
	. "github.com/lunfardo314/goq/supervisor"
	"sync"
	"time"
)

// file contain general definitions for GolOracle

const (
	golW = 81
	golH = 81
)

type GolOracle struct {
	sync.RWMutex   // to access connection map
	entity         *Entity
	golConnections map[Hash]*golConnection
}

type golConnection struct {
	conn         *websocket.Conn
	themap       Trits
	timeLastCall time.Time
}

func NewGolOracle(sv *Supervisor) (*GolOracle, error) {
	core := &GolOracle{golConnections: make(map[Hash]*golConnection)}
	ret, err := sv.NewEntity("wsServerOracle", 0, 0, core)
	if err != nil {
		return nil, err
	}
	core.entity = ret
	return core, nil
}

func (gol *GolOracle) GetEntity() *Entity {
	return gol.entity
}

func (gol *GolOracle) GetConnectionById(id Hash) *websocket.Conn {
	gol.RLock()
	defer gol.RUnlock()
	if ret, ok := gol.golConnections[id]; ok {
		return ret.conn
	}
	return nil
}

func (gol *GolOracle) RegisterGolConnection(id Hash, conn *websocket.Conn) error {
	Logf(0, "Registering wsocket for id = %v", id)
	gol.Lock()
	defer gol.Unlock()
	if _, ok := gol.golConnections[id]; ok {
		return fmt.Errorf("duplicate websocket with id = '%v'", id)
	}
	gol.golConnections[id] = &golConnection{
		conn:   conn,
		themap: make(Trits, golH*golW),
	}
	return nil
}

func (gol *GolOracle) RemoveGolConnection(id Hash) error {
	Logf(0, "Removing connection for %v: %v", id)
	gol.Lock()
	defer gol.Unlock()
	if _, ok := gol.golConnections[id]; ok {
		delete(gol.golConnections, id)
		return nil
	}
	return fmt.Errorf("websocket with id = '%v' doesn't exist", id)
}

func (gol *GolOracle) CopyMap(id Hash) (Trits, error) {
	gol.RLock()
	defer gol.RUnlock()
	if c, ok := gol.golConnections[id]; ok {
		ret := make(Trits, len(c.themap))
		copy(ret, c.themap)
		return ret, nil
	}
	return nil, fmt.Errorf("websocket with id = '%v' doesn't exist", id)
}

func (gol *GolOracle) SetMap(id Hash, themap Trits) error {
	gol.Lock()
	defer gol.Unlock()
	if c, ok := gol.golConnections[id]; ok {
		c.themap = themap
		return nil
	}
	return fmt.Errorf("websocket with id = '%v' doesn't exist", id)
}

func (gol *GolOracle) CallStart(id Hash) error {
	gol.Lock()
	defer gol.Unlock()
	if c, ok := gol.golConnections[id]; ok {
		c.timeLastCall = time.Now()
		return nil
	}
	return fmt.Errorf("websocket with id = '%v' doesn't exist", id)
}

func (gol *GolOracle) CallDuration(id Hash) time.Duration {
	gol.RLock()
	defer gol.RUnlock()
	if c, ok := gol.golConnections[id]; ok {
		return time.Since(c.timeLastCall)
	}
	return 0
}

func getCell(themap Trits, x, y int) (int8, error) {
	if x < 0 || x >= golW || y < 0 || y >= golH {
		return 0, fmt.Errorf("Tritmap.Get: out of bounds")
	}
	return themap[y*golW+x], nil
}

func putCell(themap Trits, x, y int, value int8) error {
	if x < 0 || x >= golW || y < 0 || y >= golH {
		return fmt.Errorf("Tritmap.Get: out of bounds")
	}
	themap[y*golW+x] = value
	return nil
}

const (
	TritSize = 1
	HashSize = 243
	HugeSize = 81
	GolSize  = HugeSize * HugeSize
)

var golInfoStruct = QStruct{
	QField{"id", HashSize},
	QField{"grid", GolSize},
	QField{"signature", HashSize},
	QField{"address", HashSize},
	QField{"cmd", TritSize},
}

type clickCmd struct {
	Cmd int `json:"cmd"` // 0 - update map, 1 - next generation, 2 - clear, 3 - randomize, 4 - radomize with gliders
	X   int `json:"x"`
	Y   int `json:"y"`
}
