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
	"os"
	"path"
	"server/homeManager/clientConnector"
	"server/logger"
	"strings"
	"sync"
)

var debugLevel int

type Manager struct {
	lock *sync.RWMutex
	// lock    *lock
	devices map[string]*Device

	configPath string
	// window map[string]*duration

	maxPropertyHistory int
	deviceKeys         []string

	clientConnections *clientConnector.Manager
}

// type duration map[string]int64

type onDiskDevice struct {
	Id          string                    `json:"id"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Help        string                    `json:"help"`
	ClientId    string                    `json:"clientid"`
	Switch      map[string]SwitchProperty `json:"switch"`
	Dial        map[string]DialProperty   `json:"dial"`
	Button      map[string]ButtonProperty `json:"button"`
	Text        map[string]TextProperty   `json:"text"`
}

// type onDiskWindow struct {
// 	Prop map[string]int64
// }

func New(maxPropertyHistory int, configPath string, clientMgr *clientConnector.Manager) *Manager {
	debugLevel = logger.GetDebugLevel()

	return &Manager{
		configPath:         configPath,
		lock:               &sync.RWMutex{},
		devices:            make(map[string]*Device),
		maxPropertyHistory: maxPropertyHistory,
		clientConnections:  clientMgr,
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
	m.lock.RUnlock()

	return deviceList
}

func (m *Manager) GetDeviceMatchClientID(clientId string) []*Device {
	var deviceList []*Device

	m.lock.RLock()
	for _, client := range m.devices {
		if client.ClientId == clientId {
			deviceList = append(deviceList, client)
		}
	}
	m.lock.RUnlock()

	return deviceList
}

func (m *Manager) Save() {
	log := logger.New("deviceManager.Save", &debugLevel)
	log.Debug("saving devices")

	deviceList := make(map[string]onDiskDevice)

	for key, device := range m.devices {
		if strings.HasPrefix(key, "virtual-") {
			continue
		}
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
		log.Error("unable to serialize devices", err)
	}

	err = os.WriteFile(path.Join(m.configPath, "devices.json"), file, 0640)
	if err != nil {
		log.Error("unable to write devices.json", err)
	}

}

func (m *Manager) Load() {
	var deviceList map[string]onDiskDevice
	var virtList map[string]onDiskDevice
	log := logger.New("deviceManager.Load", &debugLevel)

	file, err := os.ReadFile(path.Join(m.configPath, "devices.json"))
	if !errors.Is(err, os.ErrNotExist) {
		if err != nil {
			log.Panic("unable to read devices.json ", err)
		}
		err = json.Unmarshal(file, &deviceList)
		if err != nil {
			log.Panic("unable to read previous system state ", err)
		}
	}

	file, err = os.ReadFile(path.Join(m.configPath, "virtual.json"))
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
				log.Info("non virtual device found in virtual devices")
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
		dev.clientConnection = m.clientConnections

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
