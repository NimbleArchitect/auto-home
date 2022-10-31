package deviceManager

import (
	"log"
)

type Device struct {
	Id          string
	Name        string
	Description string
	Help        string
	ClientId    string
	// ActionWriter func(s string) (int, error)
	actionWriter ActionWriter
	// Groups      []*group

	PropertySwitch map[string]*Switch
	PropertyDial   map[string]*Dial
	PropertyButton map[string]*Button
	PropertyText   map[string]*Text

	DialNames   []string
	SwitchNames []string
	ButtonNames []string
	TextNames   []string

	repeatWindow       map[string]int64
	maxPropertyHistory int
	// Uploads []*Upload
}

func NewDevice(maxPropertyHistory int) *Device {
	return &Device{
		PropertyDial:   make(map[string]*Dial),
		PropertySwitch: make(map[string]*Switch),
		PropertyButton: make(map[string]*Button),
		PropertyText:   make(map[string]*Text),

		maxPropertyHistory: maxPropertyHistory,
	}
}

func (d *Device) MaxHistory() int {
	return d.maxPropertyHistory
}

func (d *Manager) Device(deviceId string) (*Device, bool) {
	d.lock.RLock()
	out, ok := d.devices[deviceId]
	d.lock.RUnlock()

	if ok {
		return out, true
	}
	return nil, false
}

func (d *Manager) HasDevice(deviceId string) bool {
	d.lock.RLock()
	_, ok := d.devices[deviceId]
	d.lock.RUnlock()

	return ok
}

func (m *Manager) SetDevice(deviceId string, dev *Device) {
	m.lock.Lock()
	if _, ok := m.devices[deviceId]; !ok {
		m.deviceKeys = append(m.deviceKeys, deviceId)
	}

	m.devices[deviceId] = dev
	m.lock.Unlock()
}

// MakeAction builds the json string needed to send actions to the device
func (d Device) MakeAction(deviceid string, propName string, propType int, value string) string {
	var val string
	var kind string

	switch propType {
	case DIAL:
		kind = "dial"
		val = value
	case SWITCH:
		kind = "switch"
		val = "\"" + value + "\""
		//TODO: add button and text props
	case BUTTON:
		kind = "button"
		val = "\"" + value + "\""
	case TEXT:
		kind = "text"
		val = "\"" + value + "\""

	}
	json := `{"Method": "action","data": {"id": "` + deviceid + `", "properties": [{"name": "` + propName + `","type": "` + kind + `","value": ` + val + `}]}}`
	return json
}

func (m *Manager) SetActionWriter(clientId string, writer ActionWriter) {

	devicelist := m.FindDeviceWithClientID(clientId)

	for _, v := range devicelist {
		m.devices[v].SetActionWriter(writer)
	}
}

func (d *Device) SetActionWriter(writer ActionWriter) {
	if writer != nil {
		if d.actionWriter != nil {
			log.Println("ActionWriter has already been set for this device, refusing to set again")
			return
		}
		d.actionWriter = writer
	}
}
