package main

import (
	"encoding/json"
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path"
	"runtime"
	event "server/eventManager"
	home "server/homeManager"
	log "server/logger"
	webHandle "server/webHandle"
	"syscall"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
)

type settings struct {
	HostAddress        string // servers address
	RecordHistory      bool   // wether to save the events to a history file
	MaxHistory         int    // maximum number of events to store
	AllocateVMs        int    // number of javascript virtual machines to pre allocate
	QueueLen           int    // total number of events that can be held ready for processing
	BufferLen          int    // number of events that can be sent for concurrent processing should be less than QueueLen
	MaxPropertyHistory int    // maximum number of previous values to be saved per property
	// ScriptPath         string
	// PluginPath         string
	// PublicPath         string
}

func main() {

	log.Info("starting with", runtime.NumCPU(), "CPUs")
	done := make(chan bool, 1)

	// get users home folder
	// "publicPath": "/home/rich/data/Projects/go/auto-home/public/",
	// "pluginPath": "/home/rich/data/Projects/go/auto-home/plugins/",
	// "scriptPath":"./scripts/",
	profile, err := os.UserConfigDir()
	if err != nil {
		log.Panic("unable to get users home folder", err)
	}
	// homeDir := path.Join(profile, "auto-home", "system")
	homeDir := path.Join(profile, "auto-home")
	publicPath := path.Join(profile, "auto-home", "public")
	systemPath := path.Join(profile, "auto-home", "system")
	configPath := path.Join(profile, "auto-home", "config.json")

	jsonFile, err := os.Open(configPath)
	if err != nil {
		log.Error("unable to open", configPath, err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var conf settings
	json.Unmarshal(byteValue, &conf)

	evtMgr := event.NewManager(conf.QueueLen, conf.BufferLen)

	homeMgr := home.NewManager(conf.RecordHistory, conf.MaxHistory, conf.AllocateVMs, conf.MaxPropertyHistory, homeDir)

	www := webHandle.New(systemPath, publicPath, evtMgr, homeMgr, conf.HostAddress)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		// syscall.SIGQUIT,
	)

	go func() {
		s := <-sigc
		log.Info(" caught signal", s.String(), "closing connections, please wait")

		www.SaveSystem()
		www.Shutdown()
		evtMgr.Shutdown()
		homeMgr.SaveSystem()
		homeMgr.Shutdown()

		done <- true
	}()

	homeMgr.Start()

	www.LoadSystem()

	// homeMgr.runStartScript()

	// for _, v := range homeMgr.GetDevices() {
	// 	www.AddDeviceActionList(v.ActionId)
	// }

	// pass the trigger function "TriggerEvent" to the event loop, this allow us to keep some seperation of responsibilities
	go evtMgr.EventLoop(homeMgr)
	go evtMgr.EventManager()
	go StartServer(done, www, homeDir)
	go StartWebsite(www, homeDir)

	// TODO: start event manager, i think???
	//
	// server loops through looking for an id match
	// server recieves event from the client using a magic url
	//
	// start plugins??
	//

	// temporary, used for testing
	if true == false {
		time.Sleep(30 * time.Second)

		homeMgr.ReloadVMs()
		time.Sleep(30 * time.Second)

		homeMgr.SaveSystem()
	}

	// TEMPORARY: force close the program after timeout
	// time.Sleep(3000 * time.Second)

	// homeMgr.Shutdown()

	// done <- true

	<-done

}

// type Size interface {
// 	Size() int64
// }

func StartServer(done chan bool, handle *webHandle.Handler, homeDir string) {

	quicConf := &quic.Config{
		KeepAlivePeriod: 500 * time.Second,
		MaxIdleTimeout:  500 * time.Second,
	}

	server := http3.Server{
		Handler:    handle,
		Addr:       handle.Address,
		QuicConfig: quicConf,
	}

	log.Info("Starting server")

	err := server.ListenAndServeTLS(path.Join(homeDir, "cert.crt"), path.Join(homeDir, "cert.key"))

	if err != nil {
		log.Error(err)
	}
	defer server.CloseGracefully(5 * time.Second)

	done <- true // used to close the program
}

func StartWebsite(handle *webHandle.Handler, homeDir string) {

	server := &http.Server{
		Handler: handle,
		Addr:    handle.Address,
	}

	log.Info("Starting server")

	err := server.ListenAndServeTLS(path.Join(homeDir, "cert.crt"), path.Join(homeDir, "cert.key"))

	if err != nil {
		log.Error(err)
	}

	// done <- true // used to close the program
}
