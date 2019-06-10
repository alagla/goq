package main

import (
	"fmt"
	"github.com/lunfardo314/goq/analyzeyaml"
	. "github.com/lunfardo314/goq/cfg"
	"github.com/lunfardo314/goq/readyaml"
	"github.com/lunfardo314/goq/supervisor"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"
)

const webServerPort = 8000
const fname = "C:/Users/evaldas/Documents/proj/Go/src/github.com/lunfardo314/goq/examples/modules/GOL.yml"

func main() {
	Logf(0, "Starting GOL for GOQ example")
	currentDir, _ := os.Getwd()
	Logf(0, "Current dir = %v", currentDir)

	// load GOL Qupla module
	moduleYAML, err := readyaml.NewQuplaModuleFromYAML(fname)
	if err != nil {
		Logf(0, "Error while parsing YAML file: %v", err)
		moduleYAML = nil
		return
	}
	// analyze loaded module and produce interpretable IR
	module, succ := analyzeyaml.AnalyzeQuplaModule(fname, moduleYAML)
	if !succ {
		Logf(0, "Failed to lead module: %v", err)
		return
	}
	module.PrintStats()

	//module.SetTraceLevel(10, "gameOfLife")
	//module.SetTraceLevel(10, "golLoopRows")
	//module.SetTraceLevel(10, "golProcessRows")

	// create Qubic supervisor

	sv := supervisor.NewSupervisor("GOL supervisor", 2*time.Second)

	// attach (join, affect) environments of the module to the supervisor

	succ = module.AttachToSupervisor(sv)
	if !succ {
		Logf(0, "Failed to attach module to supervisor: %v", err)
		return
	}
	//traceEnvironment(sv, "GolView")
	//traceEnvironment(sv, "GolGen")
	//traceEnvironment(sv, "GolSend")

	printEnvironmentInfo(sv)

	// create GolOracle

	golOracle, err := NewGolOracle(sv)
	if err != nil {
		Logf(0, "error while creating GolOracle: %v", err)
		os.Exit(1)
	}

	// join the oracle to GolView environment
	// any effect posted to that environment (GolInfo Qupla type) will be sent
	// to the browser for display

	err = sv.Join("GolView", golOracle.GetEntity(), 1)
	if err != nil {
		Logf(0, "error while joining GolOracle to environment: %v", err)
		os.Exit(1)
	}

	// setting up and starting web server
	// webserver will call gWSServerHandle provided by oracle to initate session
	// the oracle will open WS with the browser and communicate directly

	staticFileRoot = path.Join(currentDir, "examples/gol")
	http.HandleFunc("/static/", staticFileHandler)
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/ws", golOracle.gWSServerHandle())

	Logf(0, "Web server will be running on port %d", webServerPort)
	panic(http.ListenAndServe(fmt.Sprintf(":%d", webServerPort), nil))
}

func traceEnvironment(sv *supervisor.Supervisor, env string) {
	printEffect, err := NewPrintEffectEntity(sv, env, 0, 0, 100)
	if err != nil {
		Logf(0, "error while creating printEffect entity: %v", err)
		os.Exit(1)
	}
	err = sv.Join(env, printEffect, 1)
	if err != nil {
		Logf(0, "error while joining pointEffect to '%v': %v", err, env)
		os.Exit(1)
	}
	Logf(0, "will be tracing environment '%v' with PrintEffectEntity", env)
}

var staticFileRoot string

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	loadStaticHTMLFile(w, "mainpage.html")
}

func staticFileHandler(w http.ResponseWriter, r *http.Request) {
	fname := r.URL.Path[len("/static/"):]
	loadStaticHTMLFile(w, fname)
}

func loadStaticHTMLFile(w http.ResponseWriter, fname string) {
	pathname := path.Join(staticFileRoot, fname)

	body, err := ioutil.ReadFile(pathname)
	if err != nil {
		Logf(0, "can't open static file: %v", err)
		_, _ = fmt.Fprintf(w, "can't open static file %v", fname)
	} else {
		_, _ = fmt.Fprint(w, string(body))
	}
}

func printEnvironmentInfo(sv *supervisor.Supervisor) {
	envinfo := sv.EnvironmentInfo()
	for env, ei := range envinfo {
		Logf(0, "Environment '%v'", env)
		Logf(0, "   Joined by:")
		for _, jo := range ei.JoinedEntities {
			Logf(0, "       '%v'", jo)
		}
		Logf(0, "   Affected by:")
		for _, af := range ei.AffectedBy {
			Logf(0, "       '%v'", af)
		}
	}
}

//
//var upgrader = websocket.Upgrader{
//	ReadBufferSize:  1024,
//	WriteBufferSize: 1024,
//}
//
//var conn *websocket.Conn
//
//func wsHandler(w http.ResponseWriter, r *http.Request) {
//	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
//	var err error
//	conn, err = upgrader.Upgrade(w, r, nil)
//	if err != nil {
//		fmt.Printf("wsHandler: %v\n", err)
//		return
//	}
//	if r.Header.Get("Origin") == "" {
//		http.Error(w, "Cross domain requests require Origin header", http.StatusBadRequest)
//		return
//	}
//	Logf(0, "Created websocket for %v", r.RemoteAddr)
//	doWebsocket(conn)
//}
//
//const (
//	golW = 50
//	golH = 50
//)
//
//func doWebsocket(conn *websocket.Conn) {
//	go writeWebsocket(conn)
//	go readWebsocket(conn)
//}
//
//func writeWebsocket(conn *websocket.Conn) {
//	coord := make([][]int, 0, 50)
//
//	var data []byte
//	var err error
//
//	for {
//		for i := 0; i < 50; i++ {
//			coord = append(coord, []int{rand.Intn(golW), rand.Intn(golH)})
//		}
//		data, err = json.Marshal(coord)
//		if err != nil {
//			Logf(0, "Marshal error: %v", err)
//			return
//		}
//		err = conn.WriteMessage(1, data)
//		if err != nil {
//			Logf(0, "Write socket error: %v", err)
//			return
//		}
//		coord = coord[0:0]
//		time.Sleep(2 * time.Second)
//	}
//}
//
//type clickCoord struct {
//	X int `json:"x"`
//	Y int `json:"y"`
//}
//
//func readWebsocket(conn *websocket.Conn) {
//	var data []byte
//	var err error
//	coord := clickCoord{}
//
//	for {
//		_, data, err = conn.ReadMessage()
//		if err != nil {
//			Logf(0, "Read socket error: %v", err)
//			return
//		}
//		err = json.Unmarshal(data, &coord)
//		if err != nil {
//			Logf(0, "Unmarshal error: %v data = '%v'", err, string(data))
//		} else {
//			Logf(0, "click coord received %+v", coord)
//		}
//	}
//}
