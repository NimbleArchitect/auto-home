package home

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	js "server/homeManager/js"
)

type Manager struct {
	devices map[string]Device
	hubs    map[string]Hub
	// events  event.Manager
	actionChannel map[string]actionsChannel
	groups        map[string]group
	actions       map[string]Action
}

func NewManager() *Manager {
	m := Manager{}

	// m.groups = make(map[string]group)
	// m.groups["living-room"] = group{
	// 	Id:     "living-room",
	// 	Name:   "Living Room",
	// 	Device: []string{"123-living-room-321"},
	// }
	// m.groups["downstairs"] = group{
	// 	Id:     "downstairs",
	// 	Name:   "Downstairs",
	// 	Device: []string{"123-living-room-321"},
	// }
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

	file, err = os.ReadFile("actions.json")
	if !errors.Is(err, os.ErrNotExist) {
		if err != nil {
			log.Panic("unable to read actions.json ", err)
		}
		err = json.Unmarshal(file, &m.actions)
		if err != nil {
			log.Panic("unable to read previous system state ", err)
		}
	}

	m.initActions()
}

func (m *Manager) initActions() {
	if len(m.actions) == 0 {
		m.actions = make(map[string]Action)
		return
	}

	for deviceid, v := range m.actions {
		actionFile := v.Location
		log.Println("loading script", actionFile, "for device", deviceid)
		vm, err := js.NewScript(actionFile)
		if err != nil {
			log.Println(err)
		}

		newAction := Action{
			jsvm: vm,
		}

		m.actions[deviceid] = newAction
	}
}

// Trigger is called once at a time, with the deviceid
func (m *Manager) Trigger(deviceid string, timestamp string, props []map[string]interface{}) error {
	// var dev jsDevice

	log.Println("event triggered")

	// if len(dev.propSwitch) == 0 {
	// 	dev.propSwitch = make(map[string]jsSwitch)
	// }

	// if len(dev.propDial) == 0 {
	// 	dev.propDial = make(map[string]jsDial)
	// }
	if vm := m.actions[deviceid].jsvm; vm == nil {
		log.Println("js vm not found for device", deviceid)
	} else {
		log.Println("state:", m.devices)
		// save the current state of all devices
		err := m.SaveState(vm)
		_ = err
		// lookup changes, trigger change notifications, what am I supposed
		//  to trigger and how am I supposed to trigger it???

		// lookup device, trigger device scripts
		// dev := m.devices[deviceid]
		// fmt.Println(">>", deviceid)

		vm.Updater = m
		vm.Process(deviceid, timestamp, props)
	}

	log.Println("event finished")
	return nil
}

// MakeAction builds the json string needed to send actions to the device
func (m *Manager) MakeAction(deviceid string, propName string, propType int, value string) string {
	var val string
	var kind string

	switch propType {
	case DIAL:
		kind = "dial"
		val = value
	case SWITCH:
		kind = "switch"
		val = "\"" + value + "\""
	}
	json := `{"Method": "action","data": {"id": "` + deviceid + `", "properties": [{"name": "` + propName + `","type": "` + kind + `","value": ` + val + `}]}}`
	return json
}

// SaveState copies the current device state into the javascript vm
func (m *Manager) SaveState(js *js.JavascriptVM) error {
	if js == nil {
		return nil
	}

	for k, v := range m.devices {
		dev := js.NewDevice(k, v.Name)

		for dKey, dial := range v.PropertyDial {
			dev.AddDial(dKey, dial.Name, dial.Value)
		}

		for sKey, swi := range v.PropertySwitch {
			dev.AddSwitch(sKey, swi.Name, swi.Value.GetBool(), swi.Value.String())
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
