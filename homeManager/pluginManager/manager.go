package pluginManager

import (
	"encoding/json"
	"server/logger"

	"net"
	"os"
	"sync"
)

var debugLevel int

type response struct {
	Method string
	Id     int
	Data   interface{}
}

type Manager struct {
	SockAddr  string
	Plugins   *Plugin
	WaitGroup *sync.WaitGroup
}

func (m *Manager) Start(jsCallBack func(string, string, map[string]interface{})) {
	debugLevel = logger.GetDebugLevel()

	log := logger.New(&debugLevel)

	if err := os.RemoveAll(m.SockAddr); err != nil {
		log.Error(err)
		os.Exit(1)
	}

	l, e := net.Listen("unix", m.SockAddr)
	if e != nil {
		log.Error("listen error:", e)
		os.Exit(1)
	}

	for {
		incoming, err := l.Accept()
		if err != nil {
			log.Error("unable to accept connection:", err)
		}

		con := PluginConnector{
			c:            incoming,
			lock:         sync.Mutex{},
			responseWait: map[int]*chan result{},
			funcList:     make(map[string]*Caller),
			plug:         m.Plugins,
			jsCallBack:   jsCallBack,
			wg:           m.WaitGroup,
		}
		nextId := con.WaitAdd()
		go con.handle()
		// wait for plugin to self register
		// TODO: this should be wrapped in a select so it timesout
		con.WaitOn(nextId)
		// now we have registered we add to the global plugin list
		m.Plugins.Add(con.name, &con)
	}
}

func makeError(id int, err error) []byte {
	var msg string

	log := logger.New(&debugLevel)

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
		log.Error("json error", err)
	}

	return data
}

func makeResponse(id int, singleArg interface{}) []byte {
	log := logger.New(&debugLevel)

	arg := make(map[string]interface{})
	arg["0"] = singleArg

	resp := response{
		Method: "result",
		Id:     id,
		Data: result{
			Ok:   true,
			Data: arg,
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		log.Error("json error", err)
	}

	return data
}
