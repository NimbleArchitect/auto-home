package home

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"server/deviceManager"
	"server/globals"
	"server/groupManager"
	"server/homeManager/clientConnector"
	js "server/homeManager/js"
	"server/homeManager/pluginManager"
	log "server/logger"
	"sync"
	"time"

	"github.com/dop251/goja"
)

type Manager struct {
	// homeManager
	RecordHistory   bool              // true = capture the history and save it to the history file, false = dont save history
	isExternalEvent bool              // internal flag used to signal if the event happened from outside the house
	eventHistory    *historyProcessor //TODO: finish history capture
	// devices       map[string]Device
	devices *deviceManager.Manager
	hubs    map[string]Hub // hubs store a list of device references
	// events  event.Manager
	clientConnection *clientConnector.Manager // client communication channel
	groups           *groupManager.Manager
	// actions       map[string]Action

	configPath    string
	timeoutWindow map[string]map[string]int64 // property repeat timeout

	MaxPropertyHistory int
	MaxVMCount         int                // maximum number of VMs to start
	activeVMs          []*js.JavascriptVM // list of active initalised VM's
	chActiveVM         chan int           // stores a list of usable VM id's that are ready for use
	vmWaits            []*sync.WaitGroup

	compiledScripts js.CompiledScripts    // list of pre compiled script code
	plugins         *pluginManager.Plugin // manages the plugins and connestions exspose a caller object to allow the server to run remote functions

	pluginPath        string    // path to plugin excutables
	chStartupComplete chan bool // recieves true when the home server has started enough tostart processing scripts
	scriptPath        string    // path to javascript files

	globals *globals.Global // items stored/referenced here can be used across VM instances
}

// type lockClient struct {
// 	lock sync.RWMutex
// 	id   string
// }

type eventHistory struct {
	Deviceid   string                   `json:"deviceid"`
	Timestamp  time.Time                `json:"timestamp"`
	Properties []map[string]interface{} `json:"properties"`
}

func NewManager(recordHistory bool, maxEventHistory int, preAllocateVMs int, maxPropertyHistory int, homePath string) *Manager {

	eventProc := historyProcessor{
		lock: sync.RWMutex{},
		max:  maxEventHistory,
	}

	clientMgr := clientConnector.NewManager()

	systemPath := path.Join(homePath, "system")
	deviceMgr := deviceManager.New(maxPropertyHistory, systemPath, clientMgr)

	globalItems := globals.New()
	// p := pluginManager.Plugin{}

	m := Manager{
		configPath:      path.Join(homePath, "system"),
		RecordHistory:   recordHistory,
		isExternalEvent: false, // by default events are considured internal to the house
		eventHistory:    &eventProc,
		MaxVMCount:      preAllocateVMs,
		chActiveVM:      make(chan int, preAllocateVMs),
		scriptPath:      path.Join(homePath, "scripts"),
		// plugins:            &p,
		devices:            deviceMgr,
		groups:             groupManager.New(systemPath),
		MaxPropertyHistory: maxPropertyHistory,
		chStartupComplete:  make(chan bool, 1),
		pluginPath:         path.Join(homePath, "plugins"),
		clientConnection:   clientMgr,

		globals: globalItems,
	}

	return &m
}

func (m *Manager) Start() {
	pluginsList := pluginManager.Plugin{}

	m.LoadSystem()
	m.StartPlugins(&pluginsList)

	// TODO: need a channel to signal when the plugins have finished loading

	// go func() {
	ok := <-m.chStartupComplete
	// wait for the startup complete signal
	m.initVMs(&pluginsList)
	// then run the start script
	if ok {
		m.runStartScript()
	}
	// }()
}

func (m *Manager) DeviceWindow(deviceId string) map[string]int64 {
	if timeout, ok := m.timeoutWindow[deviceId]; ok {
		return timeout
	}

	return make(map[string]int64)
}

func (m *Manager) SaveSystem() {

	// log.Println("saving system configuration")

	m.devices.Save()

	file, err := json.Marshal(m.hubs)
	if err != nil {
		log.Error("unable to serialize hubs", err)
	}

	err = os.WriteFile(path.Join(m.configPath, "hubs.json"), file, 0640)
	if err != nil {
		log.Error("unable to write jubs.json", err)
	}

	m.groups.Save()
}

func (m *Manager) LoadSystem() {
	// var window map[string]map[string]int64

	log.Info("loading system configuration")

	m.devices.Load()

	file, err := os.ReadFile(path.Join(m.configPath, "hubs.json"))
	if !errors.Is(err, os.ErrNotExist) {
		if err != nil {
			log.Panic("unable to read hubs.json ", err)
		}
		err = json.Unmarshal(file, &m.hubs)
		if err != nil {
			log.Panic("unable to read previous system state ", err)
		}
	}

	m.groups.Load()

	// file, err = os.ReadFile("groups.json")
	// if !errors.Is(err, os.ErrNotExist) {
	// 	if err != nil {
	// 		log.Panic("unable to read groups.json ", err)
	// 	}
	// 	err = json.Unmarshal(file, &m.groups)
	// 	if err != nil {
	// 		log.Panic("unable to read previous system state ", err)
	// 	}
	// }

	file, err = os.ReadFile(path.Join(m.configPath, "window.json"))
	if !errors.Is(err, os.ErrNotExist) {
		if err != nil {
			log.Panic("unable to read window.json ", err)
		}
		err = json.Unmarshal(file, &m.timeoutWindow)
		if err != nil {
			log.Panic("unable to read previous system state ", err)
		}
	}
}

func (m *Manager) initVMs(plugs *pluginManager.Plugin) {

	m.compiledScripts = js.LoadAllScripts(m.scriptPath)

	for i := 0; i < m.MaxVMCount; i++ {
		// start the preallocated javascrip VMs
		// js, err := m.compiledScripts.NewVM(m.plugins)
		js, err := m.compiledScripts.NewVM(plugs, m.globals)
		if err == nil {
			// Set the group properties, this is not expected to change very often so we run it during vm creation
			iterator := m.groups.Iterate()
			for iterator.Next() {
				v := iterator.Get()
				js.SetGroup(v.Id, v.Name, v.Groups, v.Devices, v)
			}

			m.activeVMs = append(m.activeVMs, js)
			m.chActiveVM <- i
		}
	}

	if len(m.activeVMs) != m.MaxVMCount {
		log.Error("error unable to start enough javascript instances")
	}

}

func (m *Manager) ReloadVMs() {

	log.Info("reloading javascript VMs, please wait")
	for {
		count := 0
		for _, v := range m.activeVMs {
			if v != nil {
				count++
				break
			}
		}
		if count == 0 {
			break
		}

		_, id := m.GetNextVM()
		m.activeVMs[id] = nil
	}

	m.activeVMs = nil

	m.initVMs(m.plugins)

	log.Info("reload complete")

}

func (m *Manager) PushVMID(id int) {

	log.Info("release VM id:", id)
	m.chActiveVM <- id
}

func (m *Manager) GetNextVM() (*js.JavascriptVM, int) {

	tryagain := true

	hastried := 0
	// TODO: finish ME
	for tryagain {
		select {
		case id := <-m.chActiveVM:
			if len(m.activeVMs) >= id {
				log.Infof("selected javascript VM #%d\n", id)
				return m.activeVMs[id], id
			}
			tryagain = false
		case <-time.After(time.Second * 5):
			// TODO: this warning needs to be visible in the UI
			log.Warning("WARN: not enough javascript VMs avaliable for use")
			tryagain = true
			hastried++
		}
		// if hastried > 5 {
		// 	log.Panicln("wont clear")
		// }
	}

	return nil, 0
}

// Trigger is called one at a time with the deviceid
func (m *Manager) Trigger(id int, deviceid string, timestamp time.Time, props []map[string]interface{}) error {
	log.Info("start Trigger:", id)

	//get next avaliable vm
	vm, id := m.GetNextVM()
	// once we have finished we make sure to add the vm id back to the channel list ready for next use
	defer m.PushVMID(id)

	if vm == nil {
		log.Error("invalid javascript vm")
	} else {

		// TODO: need to write a proper check function that disables the history recording when an external event happens
		//  external events can be the doorbell being pressed or an external door being opened
		//  for now I have hard coded two external events so I can test the history recording works as expected
		if deviceid == "door-bell" || deviceid == "front-door" {
			// setting isExternalEvent to false will disable the history recording while an external event is in play
			m.isExternalEvent = true
			// re-enable recording
			defer func() { m.isExternalEvent = false }()
		}

		// TODO: somewhere I need to validate the properties so I only save valid states
		log.Debug("state:", m.devices)

		// save the current state of all devices
		vm.SaveDeviceState(m.devices)

		// lookup changes, trigger change notifications, what am I supposed
		//  to trigger and how am I supposed to trigger it???

		event := eventHistory{
			Deviceid:   deviceid,
			Timestamp:  timestamp,
			Properties: props,
		}

		// lookup device, trigger device scripts
		// dev := m.devices[deviceid]

		// process the event
		devList := m.verifyMap2jsDevice(deviceid, timestamp, props)
		// moved the go func from event loop to here as there isnt an easy way to copy an interface and it was causing a strange race bug

		vm.Updater = m
		vm.Process(deviceid, timestamp, devList)

		vm.Wait()

		//now we have finished processing save the event to our internal history list
		m.eventHistory.Add(event)

		// only record the history to a file if the user wants recording (RecordHistory) and if we are not processing an external event
		if m.RecordHistory && !m.isExternalEvent {
			// save history to file, we do this after processing the event so we have a quicker response to the event
			fileData, err := json.Marshal(event)
			if err != nil {
				log.Error("unable to serialize event", err)
			}
			var f *os.File

			if f, err = os.OpenFile(path.Join(m.configPath, "history.json"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640); err != nil {
				log.Error("unable to open file", err)
			} else {
				defer f.Close()
			}

			_, err = f.Write(append(fileData, []byte("\n")...))
			if err != nil {
				log.Error("unable to write history.json", err)
			}

		}
		// fmt.Println("finish Trigger:", id)

	}

	return nil
}

func (m *Manager) verifyMap2jsDevice(deviceid string, timestamp time.Time, props []map[string]interface{}) js.JSPropsList {

	newdev := js.NewJSDevice()

	dev, ok := m.devices.Device(deviceid)
	if !ok {
		return newdev
	}

	for _, prop := range props {
		rawName, ok := prop["name"]
		if !ok {
			log.Error("recieved property without a name")
			continue
		}
		name := rawName.(string)
		if val, ok := prop["type"]; ok {
			log.Infof("processing %s property: %s", val.(string), name)
			switch val.(string) {
			case "switch":
				swi, err := js.MapToJsSwitch(prop)
				if err != nil {
					log.Error("map error", err)
					continue
				}
				newdev.AddSwitch(name, swi)

			case "dial":
				dial, err := js.MapToJsDial(prop)
				if err != nil {
					log.Error("map error", err)
					continue
				}

				p := dev.Dial(name)
				// check min and max are within range
				if dial.Value > p.Max {
					dial.Value = p.Max
				}
				if dial.Value < p.Min {
					dial.Value = p.Min
				}

				// check we are outside of our repeat window
				if dev.DialWindow(name, timestamp) {
					newdev.AddDial(name, dial)
				}

			case "button":
				button, err := js.MapToJsButton(prop)
				if err != nil {
					log.Error("map error", err)
					continue
				}
				newdev.AddButton(name, button)

			case "text":
				text, err := js.MapToJsText(prop)
				if err != nil {
					log.Error("map error", err)
					continue
				}
				newdev.AddText(name, text)

			default:
				log.Error("unknown property type")
			}
		}
	}

	return newdev
}

func (m *Manager) Shutdown() {
	m.clientConnection.CloseAll()
	time.Sleep(1 * time.Second)
}

func (m *Manager) runStartScript() {

	// called during startup to run the server onstart function

	vm, id := m.GetNextVM()
	defer m.PushVMID(id)

	svr := "homeserver"
	_, err := vm.RunJS(svr, js.StrOnStart, goja.Undefined())
	if err != nil {
		log.Error(err)
	}
}
