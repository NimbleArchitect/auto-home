package home

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"sync"
	"time"

	js "server/homeManager/js"

	"github.com/dop251/goja"
)

type Manager struct {
	RecordHistory bool
	eventHistory  *historyProcessor //TODO: finish history capture
	devices       map[string]Device
	hubs          map[string]Hub
	// events  event.Manager
	actionChannel map[string]actionsChannel
	groups        map[string]group
	// actions       map[string]Action

	MaxVMCount int
	activeVMs  []*js.JavascriptVM
	chActiveVM chan int

	compiledScripts js.CompiledScripts
	plugins         map[string]rpc.Client
}

type eventHistory struct {
	Deviceid   string
	Timestamp  time.Time
	Properties []map[string]interface{}
}

func NewManager(recordHistory bool, maxHistory int, maxVMs int) *Manager {
	eventProc := historyProcessor{
		lock: sync.RWMutex{},
		max:  maxHistory,
	}

	m := Manager{
		RecordHistory: recordHistory,
		eventHistory:  &eventProc,
		MaxVMCount:    maxVMs,
		chActiveVM:    make(chan int, maxVMs),
	}

	return &m
}

func (m *Manager) SaveSystem() {
	log.Println("saving system configuration")

	file, err := json.Marshal(m.devices)
	if err != nil {
		log.Println("unable to serialize devices", err)
	}
	err = os.WriteFile("devices.json", file, 0640)
	if err != nil {
		log.Println("unable to write devices.json", err)
	}

	file, err = json.Marshal(m.hubs)
	if err != nil {
		log.Println("unable to serialize hubs", err)
	}
	err = os.WriteFile("hubs.json", file, 0640)
	if err != nil {
		log.Println("unable to write jubs.json", err)
	}

	file, err = json.Marshal(m.groups)
	if err != nil {
		log.Println("unable to serialize groups", err)
	}
	err = os.WriteFile("groups.json", file, 0640)
	if err != nil {
		log.Println("unable to write groups.json", err)
	}
}

func (m *Manager) LoadSystem() {
	log.Println("loading system configuration")

	file, err := os.ReadFile("devices.json")
	if !errors.Is(err, os.ErrNotExist) {
		if err != nil {
			log.Panic("unable to read devices.json ", err)
		}
		err = json.Unmarshal(file, &m.devices)
		if err != nil {
			log.Panic("unable to read previous system state ", err)
		}
	}

	file, err = os.ReadFile("hubs.json")
	if !errors.Is(err, os.ErrNotExist) {
		if err != nil {
			log.Panic("unable to read hubs.json ", err)
		}
		err = json.Unmarshal(file, &m.hubs)
		if err != nil {
			log.Panic("unable to read previous system state ", err)
		}
	}

	file, err = os.ReadFile("groups.json")
	if !errors.Is(err, os.ErrNotExist) {
		if err != nil {
			log.Panic("unable to read groups.json ", err)
		}
		err = json.Unmarshal(file, &m.groups)
		if err != nil {
			log.Panic("unable to read previous system state ", err)
		}
	}

	// file, err = os.ReadFile("actions.json")
	// if !errors.Is(err, os.ErrNotExist) {
	// 	if err != nil {
	// 		log.Panic("unable to read actions.json ", err)
	// 	}
	// 	err = json.Unmarshal(file, &m.actions)
	// 	if err != nil {
	// 		log.Panic("unable to read previous system state ", err)
	// 	}
	// }

	m.initVMs()
}

// func (m *Manager) initActions() {
// 	if len(m.actions) == 0 {
// 		m.actions = make(map[string]Action)
// 		return
// 	}

// 	for deviceid, v := range m.actions {
// 		actionFile := v.Location
// 		log.Println("loading script", actionFile, "for device", deviceid)
// 		vm, err := js.NewScript(actionFile)
// 		if err != nil {
// 			log.Println(err)
// 		}

// 		newAction := Action{
// 			jsvm: vm,
// 		}

// 		m.actions[deviceid] = newAction
// 	}
// }

func (m *Manager) initVMs() {
	// fmt.Println("1>>", count)
	m.compiledScripts = js.LoadAllScripts("./scripts/")

	for i := 0; i < m.MaxVMCount; i++ {
		js, err := m.compiledScripts.NewVM()
		if err == nil {
			m.activeVMs = append(m.activeVMs, js)
			m.chActiveVM <- i
		}
	}

	if len(m.activeVMs) != m.MaxVMCount {
		log.Println("error unable to start enough javascript instances")
	}

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
			log.Println("WARN: not enough javascript VMs avaliable for use")
			tryagain = true
		}
	}

	return nil, 0
}

// Trigger is called once at a time, with the deviceid
func (m *Manager) Trigger(deviceid string, timestamp time.Time, props []map[string]interface{}) error {
	log.Println("event triggered")
	//TODO: call client on trigger, need to work out the client script to run

	//get next avaliable vm

	// if vm := m.actions[deviceid].jsvm; vm == nil {

	vm, id := m.GetNextVM()
	defer func() { m.chActiveVM <- id }()

	if vm == nil {
		log.Println("invalid javascript vm")
	} else {
		// TODO: somewhere I need to validate the properties so I only save valid states
		log.Println("state:", m.devices)
		// save the current state of all devices
		err := m.SaveState(vm)
		_ = err
		// lookup changes, trigger change notifications, what am I supposed
		//  to trigger and how am I supposed to trigger it???

		msg := eventHistory{
			Deviceid:   deviceid,
			Timestamp:  timestamp,
			Properties: props,
		}

		// lookup device, trigger device scripts
		// dev := m.devices[deviceid]
		// fmt.Println(">>", deviceid)

		// process the event
		vm.Updater = m
		vm.Process(deviceid, timestamp, props)

		//now we have finished processing save the event to out history list
		m.eventHistory.Add(msg)

		if m.RecordHistory {
			// save history to file, we do this after processing the event so we have a quicker response to the event
			fileData, err := json.Marshal(msg)
			if err != nil {
				log.Println("unable to serialize event", err)
			}
			var f *os.File

			if f, err = os.OpenFile("history.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640); err != nil {
				log.Println("unable to open file", err)
			} else {
				defer f.Close()
			}

			_, err = f.Write(append(fileData, []byte("\n")...))
			if err != nil {
				log.Println("unable to write groups.json", err)
			}

		}

	}

	log.Println("event finished")
	return nil
}

// SaveState copies the current device state into the javascript vm
func (m *Manager) SaveState(js *js.JavascriptVM) error {
	if js == nil {
		return nil
	}

	for k, v := range m.devices {
		dev := js.NewDevice(k, v.Name)

		for dKey, dial := range v.PropertyDial {
			dev.AddDial(dKey, dial.Name, dial.Value, dial.Min, dial.Max)
		}

		for sKey, swi := range v.PropertySwitch {
			dev.AddSwitch(sKey, swi.Name, swi.Value.GetBool(), swi.Value.String())
		}

		for sKey, but := range v.PropertyButton {
			dev.AddButton(sKey, but.Name, but.Value.GetBool(), but.Value.String())
		}

		for sKey, txt := range v.PropertyText {
			dev.AddText(sKey, txt.Name, txt.Value)
		}

		js.SaveDevice(dev)
	}

	return nil
}

func (m *Manager) Shutdown() {
	for _, v := range m.actionChannel {
		v.Write(`{"Method": "shutdown"}`)
	}

	time.Sleep(10 * time.Second)
}

func (m *Manager) runStartScript() {
	// TODO: needs to be called during startup and needs to run the server onsart function

	log.Println("loading script server.js")
	// vm, err := js.NewScript("server.js")
	code := js.LoadAllScripts("scripts/")
	vm, err := code.NewVM()

	if err != nil {
		log.Println(err)
	} else {
		vm.RunJS("server", "server_onStart", goja.Undefined())
	}

}

func (m *Manager) StartPlugins() {
	// go func() {
	// 	m.startPluginManager()
	// 	m.startAllPlugins()
	// }()
}

func (m *Manager) RunGroupAction(groupId string, fnName string, props []map[string]interface{}) (interface{}, error) {

	// if vm := m.actions[groupId].jsvm; vm == nil {
	// 	log.Println("js vm not found for group", groupId)
	// } else {

	// 	// lookup changes, trigger change notifications, what am I supposed
	// 	//  to trigger and how am I supposed to trigger it???

	// 	// process the event
	// 	vm.Updater = m
	// 	return vm.RunJSGroupAction(groupId, fnName, props)
	// }

	// log.Println("event finished")

	fmt.Println(">> NOT IMPLEMENTED <<")
	return nil, nil

}
