package main

import (
	"flag"
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

const (
	webServerPort   = 8000
	golCodeLocation = "C:/Users/evaldas/Documents/proj/Go/src/github.com/lunfardo314/goq/examples/modules/GOL.yml"
	siteRoot        = "./examples/gol/static"
)

func main() {
	cfgStr := flag.String("cfg", "gol.yml", "Location of the configuration YAML")
	flag.Parse()
	Logf(0, "Starting GOL for GOQ example")
	currentDir, _ := os.Getwd()
	Logf(0, "Current dir is %v", currentDir)

	err := readConfig(*cfgStr)
	if err != nil {
		Logf(0, "error while reading config file %v: %v", *cfgStr, err)
		Logf(0, "using defaults")

		ConfigGol.WebServerPort = webServerPort
		ConfigGol.SiteRoot = siteRoot
		ConfigGol.GolCodeLocation = golCodeLocation
	}

	// load GOL Qupla module
	Logf(0, "Loading GOL code from %v", ConfigGol.GolCodeLocation)
	moduleYAML, err := readyaml.NewQuplaModuleFromYAML(ConfigGol.GolCodeLocation)
	if err != nil {
		Logf(0, "Error while parsing YAML file: %v", err)
		moduleYAML = nil
		return
	}
	// analyze loaded module and produce interpretable IR
	module, succ := analyzeyaml.AnalyzeQuplaModule(ConfigGol.GolCodeLocation, moduleYAML)
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

	//printEnvironmentInfo(sv)

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

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(ConfigGol.SiteRoot))))
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/ws", golOracle.gWSServerHandle())

	Logf(0, "Web server will be running on port %d", ConfigGol.WebServerPort)
	panic(http.ListenAndServe(fmt.Sprintf(":%d", ConfigGol.WebServerPort), nil))
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

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	loadStaticHTMLFile(w, "mainpage.html")
}

//func staticFileHandler(w http.ResponseWriter, r *http.Request) {
//	golCodeLocation := r.URL.Path[len("/static/"):]
//	loadStaticHTMLFile(w, golCodeLocation)
//}
//
func loadStaticHTMLFile(w http.ResponseWriter, fname string) {
	pathname := path.Join(ConfigGol.SiteRoot, fname)
	Logf(0, "load static file: %v", pathname)

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
