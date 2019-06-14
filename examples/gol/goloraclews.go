package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	. "github.com/iotaledger/iota.go/trinary"
	. "github.com/lunfardo314/goq/cfg"
	"github.com/lunfardo314/goq/utils"
	"math/rand"
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
	cmd := userMouseCmd{}

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
		_ = gol.CallStart(id)
		switch cmd.Cmd {
		case 0:
			gol.updateMapCmd(id, idtrits, cmd.Coord)
		case 1:
			gol.nextGenCmd(id, idtrits)
		case 2:
			gol.clearCmd(id, idtrits)
		case 3:
			gol.randomizeCmd(id, idtrits)
		case 4:
			gol.randomizeGlidersCmd(id, idtrits)
		default:
			Logf(0, "wrong cmd from browser")
			continue
		}
	}
}

func (gol *GolOracle) nextGenCmd(id string, idtrits Trits) {
	Logf(0, "next gen click cmd received '%v'", id)
	themapCopy, err := gol.CopyMap(id)
	if err != nil {
		Logf(0, "error while copying map: %v ", err)
		return
	}
	effect := golInfoStruct.ToTrits(map[string]Trits{
		"id":   idtrits,
		"grid": themapCopy,
	},
	)
	err = gol.entity.Supervisor.PostEffect("GolGen", effect, 0)
	if err != nil {
		Logf(0, "error: PostEffect returned '%v'", err)
	}
}

func (gol *GolOracle) updateMapCmd(id string, idtrits Trits, coo []coord) {
	Logf(0, "update map cmd received from '%v': %v", id, coo)
	themapCopy, err := gol.CopyMap(id)
	if err != nil {
		Logf(0, "error while copying map: %v ", err)
		return
	}
	for _, c := range coo {
		var cell int8
		cell, err = getCell(themapCopy, c.X, c.Y)
		if err != nil {
			Logf(0, "error: %v", err)
		}
		switch cell {
		case -1:
			_ = putCell(themapCopy, c.X, c.Y, 1)
		case 0:
			_ = putCell(themapCopy, c.X, c.Y, 1)
		case 1:
			_ = putCell(themapCopy, c.X, c.Y, 0)
		default:
			Logf(0, "error: wrong cell content %v at x = %d y = %d", cell, c.X, c.Y)
		}
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

func (gol *GolOracle) clearCmd(id string, idtrits Trits) {
	Logf(0, "clear map click cmd received form '%v'", id)
	effect := golInfoStruct.ToTrits(map[string]Trits{
		"id":   idtrits,
		"grid": make(Trits, golH*golW),
	},
	)
	err := gol.entity.Supervisor.PostEffect("GolView", effect, 0)
	if err != nil {
		Logf(0, "error: PostEffect returned '%v'", err)
	}
}

func (gol *GolOracle) randomizeCmd(id string, idtrits Trits) {
	Logf(0, "ramdomize map click cmd received form '%v'", id)
	themapCopy, err := gol.CopyMap(id)
	if err != nil {
		Logf(0, "error while copying map: %v ", err)
		return
	}

	for i := 0; i < 50; i++ {
		rx := int(rand.NormFloat64()*golW/5 + golW/2)
		if rx < 0 {
			rx = 0
		}
		if rx >= golW {
			rx = golW - 1
		}
		ry := int(rand.NormFloat64()*golH/5 + golH/2)
		if ry < 0 {
			ry = 0
		}
		if ry >= golW {
			ry = golW - 1
		}
		_ = putCell(themapCopy, rx, ry, 1)
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

var gliders = [][]struct{ x, y int }{
	{{-1, 0}, {0, 1}, {1, 1}, {1, 0}, {1, -1}},
	{{0, -1}, {-1, 0}, {-1, 1}, {0, 1}, {1, 1}},
	{{-1, -1}, {0, -1}, {-1, 0}, {1, 0}, {-1, 1}},
	{{-1, -1}, {0, -1}, {1, -1}, {1, 0}, {0, 1}},
}

func putGlider(themap Trits, x, y int) {
	rnd := rand.Intn(4)
	for _, offset := range gliders[rnd] {
		_ = putCell(themap, x+offset.x, y+offset.y, 1)
	}
}

func (gol *GolOracle) randomizeGlidersCmd(id string, idtrits Trits) {
	Logf(0, "ramdomize with gliders map click cmd received form '%v'", id)
	themapCopy, err := gol.CopyMap(id)
	if err != nil {
		Logf(0, "error while copying map: %v ", err)
		return
	}

	for i := 0; i < 10; i++ {
		rx := rand.Intn(golW)
		ry := rand.Intn(golH)
		putGlider(themapCopy, rx, ry)
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
	data := []byte(utils.TritsToString(golInfo["grid"]))
	err = conn.WriteMessage(1, data)
	if err != nil {
		Logf(0, "Write socket error: %v", err)
		return true
	}
	err = gol.SetMap(id, golInfo["grid"])
	if err != nil {
		Logf(0, "error while saving map for id = %v", id)
	}
	Logf(0, "Call duration on GolView: %v", gol.CallDuration(id))
	return true
}
