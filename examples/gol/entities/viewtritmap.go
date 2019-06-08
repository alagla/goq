package entities

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/iotaledger/iota.go/consts"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/cfg"
	. "github.com/lunfardo314/goq/supervisor"
	"sync"
)

type viewTritmapEntityCore struct {
	sync.RWMutex
	entity     *Entity
	id         Hash
	themap     *Tritmap
	conn       *websocket.Conn
	updateChan chan Trits
	running    bool
}

type GolInfo struct {
	id        Trits
	grid      Trits
	signature Trits
	address   Trits
	cmd       int8
}

func NewViewTritmapEntity(supervisor *Supervisor, conn *websocket.Conn, name string, width, height int) (*Entity, error) {
	themap, err := NewTritmap(width, height)
	if err != nil {
		return nil, err
	}
	core := &viewTritmapEntityCore{
		conn:       conn,
		themap:     themap,
		updateChan: make(chan Trits),
	}
	ret, err := supervisor.NewEntity(name, 0, 0, core)
	core.entity = ret
	if err != nil {
		return nil, err
	}
	_ := supervisor.CreateEnvironment("GolGen")
	core.running = true
	go core.updateLoop()
	go core.inputLoop()
	return ret, nil
}

func (e *viewTritmapEntityCore) parseEffect(effect Trits) *GolInfo {
	return &GolInfo{
		id:   effect[:consts.HashTrinarySize],
		grid: effect[consts.HashTrinarySize : e.themap.Width*e.themap.Height],
		// signature:
		// address:
		// cmd:
	}
}

// when updated map is sent to the environment, the entity will react with Call
func (e *viewTritmapEntityCore) Call(input Trits, _ Trits) bool {
	e.RLock()
	defer e.RUnlock()
	if e.running {
		e.updateChan <- input
	}
	return true // does not affect any environment, does not produce any result
}

func (e *viewTritmapEntityCore) updateLoop() {
	defer func() {
		e.Lock()
		defer e.Unlock()
		e.running = false
		close(e.updateChan)
	}()

	var data []byte
	var err error
	var golInfo *GolInfo
	var id Hash
	// loop until first error
	for golInfoTrits := range e.updateChan {
		golInfo = e.parseEffect(golInfoTrits)
		id = MustTritsToTrytes(golInfo.id)
		if id != e.id {
			continue
		}
		data, err = json.Marshal(e.parseEffect(golInfo.grid))
		if err != nil {
			Logf(0, "Marshal error: %v", err)
			return
		}
		err = e.conn.WriteMessage(1, data)
		if err != nil {
			Logf(0, "Write socket error: %v", err)
			return
		}
	}
}

// on each mouse message from the remote client, updates
// the map and posts GolInfo effect to the environment GolGen

type mouseEvent struct {
	X   int  `json:"x"`
	Y   int  `json:"y"`
	Run bool `json:"run"`
}

func (e *viewTritmapEntityCore) inputLoop() {
	var data []byte
	var err error
	var event mouseEvent
	for {
		_, data, err = e.conn.ReadMessage()
		if err != nil {
			Logf(0, "Read socket error: %v", err)
			return
		}
		err = json.Unmarshal(data, &event)
		if err != nil {
			Logf(0, "Unmarshal error: %v data = '%v'", err, string(data))
		} else {
			Logf(0, "click coord received %+v", event)
		}
		e.updateTritmap(&event)
		// form GolInfo and post to GolGen envronment
	}
}

func (e *viewTritmapEntityCore) updateTritmap(event *mouseEvent) {

}
