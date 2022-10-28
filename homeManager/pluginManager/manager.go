package pluginManager

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"sync"

	"github.com/dop251/goja"
)

type response struct {
	Method string
	Id     int
	Data   interface{}
}

// var plugins map[string]Caller
type Plugin struct {
	lock sync.RWMutex
	list map[string]*PluginConnector
}

func (p *Plugin) Add(name string, connector *PluginConnector) {
	p.lock.Lock()
	if p.list == nil {
		p.list = make(map[string]*PluginConnector)
	}
	// fmt.Println("1>> add plugin", name)
	p.list[name] = connector
	p.lock.Unlock()
}

func (p *Plugin) All() map[string]*PluginConnector {
	p.lock.RLock()
	out := p.list
	p.lock.RUnlock()
	return out
}

func Start(sockAddr string, wg *sync.WaitGroup, plugins *Plugin, jsCallBack func(string, string, *goja.Object)) {
	if err := os.RemoveAll(sockAddr); err != nil {
		log.Fatal(err)
	}

	l, e := net.Listen("unix", sockAddr)
	if e != nil {
		log.Fatal("listen error:", e)
	}

	for {
		incoming, _ := l.Accept()

		con := PluginConnector{
			c:          incoming,
			lock:       sync.Mutex{},
			wait:       map[int]chan bool{},
			funcList:   make(map[string]*Caller),
			plug:       plugins,
			jsCallBack: jsCallBack,
			wg:         wg,
		}
		completed := make(chan bool, 1)
		con.wait[0] = completed
		go con.handle()
		// wait for plugin to self register
		// TODO: this should be wrapped in a select so it timesout
		<-completed
		// now we have registered we add to the global plugin list
		plugins.Add(con.name, &con)
	}
}

func makeError(id int, err error) []byte {
	var msg string

	ok := false
	if err == nil {
		ok = true
	} else {
		msg = err.Error()
	}

	resp := response{
		Method: "result",
		Id:     id,
		Data: result{
			Ok:      ok,
			Message: msg,
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		log.Println("json error", err)
	}

	return data
}
