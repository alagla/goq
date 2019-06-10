package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/cfg"
	. "github.com/lunfardo314/goq/supervisor"
	"github.com/lunfardo314/goq/utils"
	"net/http"
	"sync"
)

type WSOracle struct {
	sync.RWMutex // to access map
	entity       *Entity
	wsockets     map[Hash]*websocket.Conn
}

func (wsc *WSOracle) GetEntity() *Entity {
	return wsc.entity
}

func (wsc *WSOracle) GetWS(id Hash) *websocket.Conn {
	wsc.RLock()
	defer wsc.RUnlock()
	if ret, ok := wsc.wsockets[id]; ok {
		return ret
	}
	return nil
}

func (wsc *WSOracle) RegisterWS(id Hash, conn *websocket.Conn) error {
	Logf(0, "Registering wsocket for id = %v", id)
	wsc.Lock()
	defer wsc.Unlock()
	if _, ok := wsc.wsockets[id]; ok {
		return fmt.Errorf("duplicate websocket with id = '%v'", id)
	}
	wsc.wsockets[id] = conn
	return nil
}

func (wsc *WSOracle) RemoveWS(id Hash) error {
	Logf(0, "Removing wsocket for %v: %v", id)
	wsc.Lock()
	defer wsc.Unlock()
	if _, ok := wsc.wsockets[id]; ok {
		delete(wsc.wsockets, id)
		return nil
	}
	return fmt.Errorf("websocket with id = '%v' doesn't exist", id)
}

func NewWSOracle(sv *Supervisor) (*WSOracle, error) {
	core := &WSOracle{wsockets: make(map[Hash]*websocket.Conn)}
	ret, err := sv.NewEntity("wsServerOracle", 0, 0, core)
	if err != nil {
		return nil, err
	}
	core.entity = ret
	return core, nil
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (wsc *WSOracle) gWSServerHandle() func(http.ResponseWriter, *http.Request) {
	// will be called every time new connection arrives
	return func(w http.ResponseWriter, r *http.Request) {
		var conn *websocket.Conn
		var err error
		var id Trits
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, err = upgrader.Upgrade(w, r, nil)
		if err != nil {
			Logf(0, "wsHandler: %v", err)
			return
		}
		Logf(0, "Created websocket for %v", r.RemoteAddr)
		// create ID by hashing r.RemoteAddr
		remoteAddrTrits := utils.Bytes2Trits([]byte(r.RemoteAddr), 81*3)
		id, err = utils.KerlHash243(remoteAddrTrits)
		if err != nil {
			Logf(0, "Kerl error %v", err)
			return
		}
		id = id[:81]
		// register websocket
		idtrytes := MustTritsToTrytes(id)
		err = wsc.RegisterWS(idtrytes, conn)
		if err != nil {
			Logf(0, "error %v", err)
			return
		}
		// start listening to the connection and sending effects to specified environment
		go wsc.readWS(conn, id)
	}
}

const (
	TritSize = 1
	HashSize = 81
	HugeSize = HashSize
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
	NextGen bool `json:"nextGen"` // true next generation, false update map,
	X       int  `json:"x"`
	Y       int  `json:"y"`
}

const (
	golW = 81
	golH = 81
)

func (wsc *WSOracle) readWS(conn *websocket.Conn, idtrits Trits) {
	golMap := MustNewTritmap(golW, golH)
	var data []byte
	var err error
	var cell int8
	cmd := clickCmd{}
	id := MustTritsToTrytes(idtrits)

	for {
		_, data, err = conn.ReadMessage()
		if err != nil {
			Logf(0, "Read socket error: %v", err)
			_ = wsc.RemoveWS(id)
			return
		}
		err = json.Unmarshal(data, &cmd)
		if err != nil {
			Logf(0, "Unmarshal error: %v Data = '%v'", err, string(data))
			continue
		}
		if cmd.NextGen {
			Logf(0, "next gen click cmd received %+v", cmd)
			effect := golInfoStruct.ToTrits(map[string]Trits{
				"id":   idtrits,
				"grid": golMap.Data,
			},
			)
			err = wsc.entity.Supervisor.PostEffect("GolGen", effect, 0)
			if err != nil {
				Logf(0, "error: PostEffect returned '%v'", err)
			}
		} else {
			Logf(0, "update map click cmd received %+v", cmd)
			cell, err = golMap.Get(cmd.X, cmd.Y)
			if err != nil {
				Logf(0, "error: wrong coordinates %+v", cmd)
				continue
			}
			switch cell {
			case -1:
				_ = golMap.Put(cmd.X, cmd.Y, 1)
			case 0:
				_ = golMap.Put(cmd.X, cmd.Y, 1)
			case 1:
				_ = golMap.Put(cmd.X, cmd.Y, 0)
			default:
				Logf(0, "error: wrong cell content %v at %+v", cell, cmd)
			}
			effect := golInfoStruct.ToTrits(map[string]Trits{
				"id":   idtrits,
				"grid": golMap.Data,
			},
			)
			err = wsc.entity.Supervisor.PostEffect("GolView", effect, 0)
			if err != nil {
				Logf(0, "error: PostEffect returned '%v'", err)
			}
		}
	}
}

// request to display map

func (wsc *WSOracle) Call(effect Trits, _ Trits) bool {
	golInfo, err := golInfoStruct.Parse(effect)
	if err != nil {
		Logf(0, "error: golInfoStruct.Parse returned '%v'", err)
		return true
	}
	idtrits := golInfo["id"]
	id := MustTritsToTrytes(idtrits)
	conn := wsc.GetWS(id)
	if conn == nil {
		Logf(0, "can't find web socket with id = %v", id)
		return true
	}
	var data []byte
	data, err = json.Marshal(golInfo["grid"])
	//Logf(0, "+++++   %v", string(data))
	if err != nil {
		Logf(0, "Marshal error: %v", err)
		return true
	}
	err = conn.WriteMessage(1, data)
	if err != nil {
		Logf(0, "Write socket error: %v", err)
		return true
	}
	return true
}
