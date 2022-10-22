package home

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"path"
	"sync"
	"time"

	"server/deviceManager"
	"server/groupManager"
	js "server/homeManager/js"

	"github.com/dop251/goja"
)

type Manager struct {
	RecordHistory bool
	eventHistory  *historyProcessor //TODO: finish history capture
	// devices       map[string]Device
	devices *deviceManager.Manager
	hubs    map[string]Hub
	// events  event.Manager
	actionChannel map[string]actionsChannel
	groups        *groupManager.Manager
	// actions       map[string]Action

	configPath    string
	timeoutWindow map[string]map[string]int64

	MaxPropertyHistory int
	MaxVMCount         int
	activeVMs          []*js.JavascriptVM
	chActiveVM         chan int

	compiledScripts js.CompiledScripts
	plugins         map[string]*rpc.Client

	pluginPath        string
	chStartupComplete chan bool
	scriptPath        string
}

// type lockClient struct {
// 	lock sync.RWMutex
// 	id   string
// }

type eventHistory struct {
	Deviceid   string
	Timestamp  time.Time
	Properties []map[string]interface{}
}

func NewManager(recordHistory bool, maxEventHistory int, preAllocateVMs int, maxPropertyHistory int, homePath string) *Manager {
	eventProc := historyProcessor{
		lock: sync.RWMutex{},
		max:  maxEventHistory,
	}

	systemPath := path.Join(homePath, "system")
	deviceMgr := deviceManager.New(maxPropertyHistory, systemPath)

	m := Manager{
		configPath:         path.Join(homePath, "system"),
		RecordHistory:      recordHistory,
		eventHistory:       &eventProc,
		MaxVMCount:         preAllocateVMs,
		chActiveVM:         make(chan int, preAllocateVMs),
		scriptPath:         path.Join(homePath, "scripts"),
		plugins:            make(map[string]*rpc.Client),
		devices:            deviceMgr,
		groups:             groupManager.New(systemPath),
		MaxPropertyHistory: maxPropertyHistory,
		chStartupComplete:  make(chan bool, 1),
		pluginPath:         path.Join(homePath, "plugins"),
	}

	return &m
}

func (m *Manager) Start() {
	m.LoadSystem()
	m.StartPlugins()
	m.initVMs()

	// TODO: need a channel to signal when the plugins have finished loading

	go func() {
		// wait for the startup complete signal
		ok := <-m.chStartupComplete
		// then run the start script
		if ok {
			m.runStartScript()
		}
	}()
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
		log.Println("unable to serialize hubs", err)
	}

	err = os.WriteFile(path.Join(m.configPath, "hubs.json"), file, 0640)
	if err != nil {
		log.Println("unable to write jubs.json", err)
	}

	m.groups.Save()

}

func (m *Manager) LoadSystem() {
	// var window map[string]map[string]int64

	log.Println("loading system configuration")

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

func (m *Manager) initVMs() {
	m.compiledScripts = js.LoadAllScripts(m.scriptPath)

	for i := 0; i < m.MaxVMCount; i++ {
		// start the preallocated javascrip VMs
		js, err := m.compiledScripts.NewVM()
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
		log.Println("error unable to start enough javascript instances")
	}

}

func (m *Manager) ReloadVMs() {

	log.Println("reloading javascript VMs, please wait")
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

	m.initVMs()

	log.Println("reload complete")

}

func (m *Manager) PushVMID(id int) {
	m.chActiveVM <- id
}

func (m *Manager) GetNextVM() (*js.JavascriptVM, int) {
	tryagain := true

	// TODO: finish ME
	for tryagain {
		select {
		case id := <-m.chActiveVM:
			if len(m.activeVMs) >= id {
				log.Printf("selected javascript VM #%d", id)
				return m.activeVMs[id], id
			}
			tryagain = false
		case <-time.After(time.Second * 5):
			// TODO: this warning needs to be visible in the UI
			log.Println("WARN: not enough javascript VMs avaliable for use")
			tryagain = true
		}
	}

	return nil, 0
}

// Trigger is called one at a time with the deviceid
func (m *Manager) Trigger(deviceid string, timestamp time.Time, props []map[string]interface{}) error {
	log.Println("event triggered")
	//TODO: call client on trigger, need to work out the client script to run

	// if vm := m.actions[deviceid].jsvm; vm == nil {

	//get next avaliable vm
	vm, id := m.GetNextVM()
	// once we have finished we make sure to add the vm id back to the channel list ready for next use
	defer m.PushVMID(id)

	if vm == nil {
		log.Println("invalid javascript vm")
	} else {
		// TODO: somewhere I need to validate the properties so I only save valid states
		log.Println("state:", m.devices)

		// save the current state of all devices
		vm.SaveDeviceState(m.devices)
		// groups are set during vm init

		// register plugins
		for n, v := range m.plugins {
			vm.NewPlugin(n, v)
		}

		// lookup changes, trigger change notifications, what am I supposed
		//  to trigger and how am I supposed to trigger it???

		msg := eventHistory{
			Deviceid:   deviceid,
			Timestamp:  timestamp,
			Properties: props,
		}

		// lookup device, trigger device scripts
		// dev := m.devices[deviceid]

		// process the event
		devList := m.verifyMap2jsDevice(deviceid, timestamp, props)
		vm.Updater = m
		vm.Process(deviceid, timestamp, devList)

		//now we have finished processing save the event to out history list
		m.eventHistory.Add(msg)

		if m.RecordHistory {
			// save history to file, we do this after processing the event so we have a quicker response to the event
			fileData, err := json.Marshal(msg)
			if err != nil {
				log.Println("unable to serialize event", err)
			}
			var f *os.File

			if f, err = os.OpenFile(path.Join(m.configPath, "history.json"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640); err != nil {
				log.Println("unable to open file", err)
			} else {
				defer f.Close()
			}

			_, err = f.Write(append(fileData, []byte("\n")...))
			if err != nil {
				log.Println("unable to write history.json", err)
			}

		}

	}

	log.Println("event finished")
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
			log.Println("recieved property without a name")
			continue
		}
		name := rawName.(string)
		if val, ok := prop["type"]; ok {
			log.Printf("processing %s property: %s", val.(string), name)
			switch val.(string) {
			case "switch":
				swi, err := js.MapToJsSwitch(prop)
				if err != nil {
					log.Println("map error", err)
					continue
				}
				newdev.AddSwitch(name, swi)

			case "dial":
				dial, err := js.MapToJsDial(prop)
				if err != nil {
					log.Println("map error", err)
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

				if dev.DialWindow(name, timestamp) {
					// check we are outside of our repeat window
					newdev.AddDial(name, dial)
					// 	fmt.Println("**>> update allowed")
					// } else {
					// 	fmt.Println("**>> update blocked")
				}

			case "button":
				button, err := js.MapToJsButton(prop)
				if err != nil {
					log.Println("map error", err)
					continue
				}
				newdev.AddButton(name, button)

			case "text":
				text, err := js.MapToJsText(prop)
				if err != nil {
					log.Println("map error", err)
					continue
				}
				newdev.AddText(name, text)

			default:
				log.Println("unknown property type")
			}
		}
	}

	return newdev
}

func (m *Manager) Shutdown() {
	for _, v := range m.actionChannel {
		v.Write(`{"Method": "shutdown"}`)
	}

	time.Sleep(1 * time.Second)
}

func (m *Manager) runStartScript() {
	// called during startup to run the server onstart function

	vm, id := m.GetNextVM()
	defer m.PushVMID(id)

	for n, v := range m.plugins {
		vm.NewPlugin(n, v)
	}

	svr := "homeserver"
	_, err := vm.RunJS(svr, js.StrOnStart, goja.Undefined())
	if err != nil {
		fmt.Println("21>>", err)
	}
	// fmt.Println("!>", v)
}

func (m *Manager) StartPlugins() {
	var pluginList []string
	pluginList = append(pluginList, "telegram")

	connected := make(chan int)
	// TODO: plugins/manager need a rewrite
	go m.startPluginManager(connected)

	<-connected

	for i := 0; i < len(pluginList); i++ {
		go m.startPlugin(pluginList[i])
		<-connected
	}

	fmt.Println("plugins started")
	m.chStartupComplete <- true
}

// func (m *Manager) RunGroupAction(groupId string, fnName string, props []map[string]interface{}) (interface{}, error) {

// 	// if vm := m.actions[groupId].jsvm; vm == nil {
// 	// 	log.Println("js vm not found for group", groupId)
// 	// } else {

// 	// 	// lookup changes, trigger change notifications, what am I supposed
// 	// 	//  to trigger and how am I supposed to trigger it???

// 	// 	// process the event
// 	// 	vm.Updater = m
// 	// 	return vm.RunJSGroupAction(groupId, fnName, props)
// 	// }

// 	// log.Println("event finished")

// 	fmt.Println(">> NOT IMPLEMENTED <<")
// 	return nil, nil

// }
