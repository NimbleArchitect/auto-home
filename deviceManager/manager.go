package deviceManager

// get/set/update device
// get/set/updata properties
// must be thread safe
// Value updates must be quick
// need to be careful with deletes, dont want to lock the system when deleting a device with lots of properties
// needs to support snapshot function to save the device/property state
// needs to record changes so I can call an update func to push the changes back to devices
//   changes can then be processed and devices updates can be sent form the trigger func
// needs to save/load devices list
//

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"
	"sync"
)

type Manager struct {
	lock    sync.RWMutex
	devices map[string]*Device

	// window map[string]*duration

	maxPropertyHistory int
	deviceKeys         []string
}

// type duration map[string]int64

type onDiskDevice struct {
	Id          string
	Name        string
	Description string
	Help        string
	ClientId    string
	Switch      map[string]SwitchProperty
	Dial        map[string]DialProperty
	Button      map[string]ButtonProperty
	Text        map[string]TextProperty
}

// type onDiskWindow struct {
// 	Prop map[string]int64
// }

func New(maxPropertyHistory int) *Manager {
	return &Manager{
		lock:               sync.RWMutex{},
		devices:            make(map[string]*Device),
		maxPropertyHistory: maxPropertyHistory,
	}
}

func (m *Manager) AllDevices() map[string]*Device {
	return m.devices
}

func (m *Manager) Keys() []string {
	arrayCopy := make([]string, len(m.deviceKeys))
	copy(arrayCopy, m.deviceKeys)

	return arrayCopy
}

func (m *Manager) ClientIdMatchesDevice(deviceId string, clientId string) bool {
	var isMatch bool

	dev, ok := m.devices[deviceId]
	if ok {
		if dev.ClientId == clientId {
			isMatch = true
		}
	}

	return isMatch
}

func (m *Manager) FindDeviceWithClientID(clientId string) []string {
	var deviceList []string

	m.lock.RLock()
	for deviceid, client := range m.devices {
		if client.ClientId == clientId {
			deviceList = append(deviceList, deviceid)
		}
	}
	m.lock.RLock()

	return deviceList
}

func (m *Manager) Save() {
	log.Println("saving devices")

	deviceList := make(map[string]onDiskDevice)

	for key, device := range m.devices {
		dial := make(map[string]DialProperty)
		for key, property := range device.PropertyDial {
			property.lock.RLock()
			dial[key] = DialProperty{
				Id:          property.data.Id,
				Name:        property.data.Name,
				Description: property.data.Description,
				Min:         property.data.Min,
				Max:         property.data.Max,
				Value:       property.data.Value,
				Mode:        property.data.Mode,
			}
			property.lock.RUnlock()
		}

		swi := make(map[string]SwitchProperty)
		for key, property := range device.PropertySwitch {
			property.lock.RLock()
			swi[key] = SwitchProperty{
				Id:          property.data.Id,
				Name:        property.data.Name,
				Description: property.data.Description,
				Value:       property.data.Value,
				Mode:        property.data.Mode,
			}
			property.lock.RUnlock()
		}

		button := make(map[string]ButtonProperty)
		for key, property := range device.PropertyButton {
			property.lock.RLock()
			button[key] = ButtonProperty{
				Id:          property.data.Id,
				Name:        property.data.Name,
				Description: property.data.Description,
				Value:       property.data.Value,
				Mode:        property.data.Mode,
			}
			property.lock.RUnlock()
		}

		text := make(map[string]TextProperty)
		for key, property := range device.PropertyText {
			property.lock.RLock()
			text[key] = TextProperty{
				Id:          property.data.Id,
				Name:        property.data.Name,
				Description: property.data.Description,
				Value:       property.data.Value,
				Mode:        property.data.Mode,
			}
			property.lock.RUnlock()
		}

		deviceList[key] = onDiskDevice{
			Id:          device.Id,
			Name:        device.Name,
			Description: device.Description,
			Help:        device.Help,
			ClientId:    device.ClientId,
			Dial:        dial,
			Switch:      swi,
			Button:      button,
			Text:        text,
		}

	}
	file, err := json.Marshal(deviceList)
	if err != nil {
		log.Println("unable to serialize devices", err)
	}
	err = os.WriteFile("devices.json", file, 0640)
	if err != nil {
		log.Println("unable to write devices.json", err)
	}

}

func (m *Manager) Load() {
	var deviceList map[string]onDiskDevice
	var virtList map[string]onDiskDevice

	file, err := os.ReadFile("devices.json")
	if !errors.Is(err, os.ErrNotExist) {
		if err != nil {
			log.Panic("unable to read devices.json ", err)
		}
		err = json.Unmarshal(file, &deviceList)
		if err != nil {
			log.Panic("unable to read previous system state ", err)
		}
	}

	file, err = os.ReadFile("virtual.json")
	if !errors.Is(err, os.ErrNotExist) {
		if err != nil {
			log.Panic("unable to read virtual.json ", err)
		}
		err = json.Unmarshal(file, &virtList)
		if err != nil {
			log.Panic("unable to read previous system state ", err)
		}
		for n, v := range virtList {
			if !strings.HasPrefix(n, "virtual-") {
				log.Println("non virtual device found in virtual devices")
				continue
			}
			if _, ok := deviceList[n]; !ok {
				deviceList[n] = v
			}
		}
	}

	for id, device := range deviceList {
		dev := NewDevice(m.maxPropertyHistory)
		dev.Id = device.Id
		dev.Description = device.Description
		dev.ClientId = device.ClientId
		dev.Name = device.Name
		dev.Help = device.Help

		for name, property := range device.Dial {
			dev.SetDial(name, &property)
		}
		for name, property := range device.Switch {
			dev.SetSwitch(name, &property)
		}
		for name, property := range device.Button {
			dev.SetButton(name, &property)
		}
		for name, property := range device.Text {
			dev.SetText(name, &property)
		}
		m.SetDevice(id, dev)
	}
}
