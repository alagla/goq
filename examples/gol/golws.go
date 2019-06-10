package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/cfg"
	"github.com/lunfardo314/goq/utils"
	"net/http"
)

// file contains WS processing code of the GolOracle

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// web server handle to start session

func (gol *GolOracle) gWSServerHandle() func(http.ResponseWriter, *http.Request) {
	// will be called every time new connection arrives
	return func(w http.ResponseWriter, r *http.Request) {
		var conn *websocket.Conn
		var err error
		var idtrits Trits
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, err = upgrader.Upgrade(w, r, nil)
		if err != nil {
			Logf(0, "wsHandler: %v", err)
			return
		}
		Logf(0, "Created websocket for %v", r.RemoteAddr)
		// create ID by hashing r.RemoteAddr
		remoteAddrTrits := utils.Bytes2Trits([]byte(r.RemoteAddr), 81*3)
		idtrits, err = utils.KerlHash243(remoteAddrTrits)
		if err != nil {
			Logf(0, "Kerl error %v", err)
			return
		}
		idtrits = idtrits[:81]
		// register websocket
		id := MustTritsToTrytes(idtrits)
		err = gol.RegisterGolConnection(id, conn)
		if err != nil {
			Logf(0, "error %v", err)
			return
		}
		// start listening to the connection and sending effects to specified environment
		go gol.readWS(id, idtrits, conn)
	}
}

// WS reading loop. It wait for mouse clicks to update the map which contains current generation
// of the GOL

func (gol *GolOracle) readWS(id string, idtrits Trits, conn *websocket.Conn) {
	var data []byte
	var err error
	var cell int8
	var themapCopy Trits
	cmd := clickCmd{}

	for {
		_, data, err = conn.ReadMessage()
		if err != nil {
			Logf(0, "Read socket error: %v", err)
			_ = gol.RemoveGolConnection(id)
			return
		}
		err = json.Unmarshal(data, &cmd)
		if err != nil {
			Logf(0, "Unmarshal error: %v Data = '%v'", err, string(data))
			continue
		}
		themapCopy, err = gol.CopyMap(id)
		if err != nil {
			Logf(0, "error while copying map: %v ", err)
			continue
		}
		if cmd.NextGen {
			Logf(0, "next gen click cmd received %+v from '%v'", cmd, id)
			effect := golInfoStruct.ToTrits(map[string]Trits{
				"id":   idtrits,
				"grid": themapCopy,
			},
			)
			err = gol.entity.Supervisor.PostEffect("GolGen", effect, 0)
			if err != nil {
				Logf(0, "error: PostEffect returned '%v'", err)
			}
		} else {
			Logf(0, "update map click cmd received %+v from '%v'", cmd, id)
			cell, err = getCell(themapCopy, cmd.X, cmd.Y)
			if err != nil {
				Logf(0, "error: wrong coordinates %+v", cmd)
				continue
			}
			switch cell {
			case -1:
				_ = putCell(themapCopy, cmd.X, cmd.Y, 1)
			case 0:
				_ = putCell(themapCopy, cmd.X, cmd.Y, 1)
			case 1:
				_ = putCell(themapCopy, cmd.X, cmd.Y, 0)
			default:
				Logf(0, "error: wrong cell content %v at %+v", cell, cmd)
			}
			effect := golInfoStruct.ToTrits(map[string]Trits{
				"id":   idtrits,
				"grid": themapCopy,
			},
			)
			err = gol.entity.Supervisor.PostEffect("GolView", effect, 0)
			if err != nil {
				Logf(0, "error: PostEffect returned '%v'", err)
			}
		}
	}
}

// main function of the GolOracle implementing Qubic entity interface
// it is called every time effect arrives to the environment GolOracle has been joined
// usually it is environment GolView where Gol maps are posted for display
// also map is stored in the oracle as a current generation which later may be
// modified by the user

func (gol *GolOracle) Call(effect Trits, _ Trits) bool {
	golInfo, err := golInfoStruct.Parse(effect)
	if err != nil {
		Logf(0, "error: golInfoStruct.Parse returned '%v'", err)
		return true
	}
	idtrits := golInfo["id"]
	id := MustTritsToTrytes(idtrits)
	id = id[:27]
	conn := gol.GetConnectionById(id)
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
	err = gol.SetMap(id, golInfo["grid"])
	if err != nil {
		Logf(0, "error while saving map for id = %v", id)
	}
	return true
}
