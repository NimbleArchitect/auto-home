package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	event "server/eventManager"
	home "server/homeManager"
	webHandle "server/webHandle"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
)

type settings struct {
	HostAddress string
	PublicPath  string
}

func main() {

	log.Println("starting with", runtime.NumCPU(), "CPUs")
	done := make(chan bool, 1)

	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var conf settings
	json.Unmarshal(byteValue, &conf)

	evtMgr := event.NewManager(200, 50)
	homeMgr := home.NewManager()
	homeMgr.LoadSystem()

	www := webHandle.Handler{
		EventManager: evtMgr,
		HomeManager:  homeMgr,
		FsHandle:     http.FileServer(http.Dir(conf.PublicPath)),
		Address:      conf.HostAddress,
	}

	www.LoadSystem()
	// for _, v := range homeMgr.GetDevices() {
	// 	www.AddDeviceActionList(v.ActionId)
	// }

	// pass the trigger function "TriggerEvent" to the event loop, this allow us to keep some seperation of responsibilities
	go evtMgr.EventLoop(homeMgr)
	go evtMgr.EventManager()
	go StartServer(done, &www)
	go StartWebsite(&www)

	// TODO: start event manager, i think???
	//
	// server loops through looking for an id match
	// server recieves event from the client using a magic url
	//
	// start plugins??
	//

	time.Sleep(30 * time.Second)

	homeMgr.SaveSystem()
	// TEMPORARY: force close the program after timeout
	time.Sleep(3000 * time.Second)

	homeMgr.Shutdown()

	done <- true

	<-done

}

// type Size interface {
// 	Size() int64
// }

func StartServer(done chan bool, handle *webHandle.Handler) {
	quicConf := &quic.Config{
		KeepAlivePeriod: 60 * time.Second,
		MaxIdleTimeout:  600 * time.Second,
	}

	server := http3.Server{
		Handler:    handle,
		Addr:       handle.Address,
		QuicConfig: quicConf,
	}

	log.Println("Starting server")

	err := server.ListenAndServeTLS("cert.crt", "cert.key")

	if err != nil {
		log.Println(err)
	}
	defer server.CloseGracefully(30 * time.Second)

	done <- true // used to close the program
}

func StartWebsite(handle *webHandle.Handler) {

	server := &http.Server{
		Handler: handle,
		Addr:    handle.Address,
	}

	log.Println("Starting server")

	err := server.ListenAndServeTLS("cert.crt", "cert.key")

	if err != nil {
		log.Println(err)
	}

	// done <- true // used to close the program
}
