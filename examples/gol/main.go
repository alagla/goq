package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	. "github.com/lunfardo314/goq/cfg"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"time"
)

const webServerPort = 8000

func main() {
	Logf(0, "Starting GOL for GOQ example")
	runWebServer(webServerPort)
}

var staticFileRoot string

func runWebServer(port int) {
	currentDir, _ := os.Getwd()
	staticFileRoot = path.Join(currentDir, "examples/gol")
	Logf(0, "Current dir = %v", currentDir)
	Logf(0, "Web server is running on port %d", port)
	http.HandleFunc("/static/", staticFileHandler)
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/ws", wsHandler)
	panic(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var conn *websocket.Conn

func wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	var err error
	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("wsHandler: %v\n", err)
		return
	}
	if r.Header.Get("Origin") == "" {
		http.Error(w, "Cross domain requests require Origin header", http.StatusBadRequest)
		return
	}
	Logf(0, "Created websocket for %v", r.RemoteAddr)
	doWebsocket(conn)
}

const (
	golW = 50
	golH = 50
)

func doWebsocket(conn *websocket.Conn) {
	go writeWebsocket(conn)
	go readWebsocket(conn)
}

func writeWebsocket(conn *websocket.Conn) {
	coord := make([][]int, 0, 50)

	var data []byte
	var err error

	for {
		for i := 0; i < 50; i++ {
			coord = append(coord, []int{rand.Intn(golW), rand.Intn(golH)})
		}
		data, err = json.Marshal(coord)
		if err != nil {
			Logf(0, "Marshal error: %v", err)
			return
		}
		err = conn.WriteMessage(1, data)
		if err != nil {
			Logf(0, "Write socket error: %v", err)
			return
		}
		coord = coord[0:0]
		time.Sleep(2 * time.Second)
	}
}

type clickCoord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func readWebsocket(conn *websocket.Conn) {
	var data []byte
	var err error
	coord := clickCoord{}

	for {
		_, data, err = conn.ReadMessage()
		if err != nil {
			Logf(0, "Read socket error: %v", err)
			return
		}
		err = json.Unmarshal(data, &coord)
		if err != nil {
			Logf(0, "Unmarshal error: %v data = '%v'", err, string(data))
		} else {
			Logf(0, "click coord received %+v", coord)
		}
	}
}
