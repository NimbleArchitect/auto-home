package home

import (
	"log"
	deviceManager "server/deviceManager"
)

type actionsChannel interface {
	IsOpen() bool
	Write(string) (int, error)
}

func (m *Manager) RegisterActionChannel(id string, ch actionsChannel) error {
	log.Println("registerActionChannel id", id)

	// TODO: this might not work

	// for _, v := range m.FindDeviceWithClientID(id) {
	if len(m.devices.FindDeviceWithClientID(id)) > 0 {
		// v.Id
		log.Println("device ID found")
		if len(m.actionChannel) == 0 {
			log.Println("empty action channel")
			m.actionChannel = make(map[string]actionsChannel)
		}

		log.Println("setting channel for device")
		m.actionChannel[id] = ch
	}
	return nil
}

// MakeAction builds the json string needed to send actions to the device
func (m *Manager) MakeAction(deviceid string, propName string, propType int, value string) string {
	var val string
	var kind string

	switch propType {
	case deviceManager.DIAL:
		kind = "dial"
		val = value
	case deviceManager.SWITCH:
		kind = "switch"
		val = "\"" + value + "\""
		//TODO: add button and text props
	case deviceManager.BUTTON:
		kind = "button"
		val = "\"" + value + "\""
	case deviceManager.TEXT:
		kind = "text"
		val = "\"" + value + "\""

	}
	json := `{"Method": "action","data": {"id": "` + deviceid + `", "properties": [{"name": "` + propName + `","type": "` + kind + `","value": ` + val + `}]}}`
	return json
}
