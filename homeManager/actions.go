package home

import "log"

type actionsChannel interface {
	IsOpen() bool
	Write(string) (int, error)
}

func (m *Manager) RegisterActionChannel(id string, ch actionsChannel) error {
	log.Println("registerActionChannel id", id)

	for _, v := range m.FindDeviceWithClientID(id) {
		// v.Id
		log.Println("device ID found")
		if len(m.actionChannel) == 0 {
			log.Println("empty action channel")
			m.actionChannel = make(map[string]actionsChannel)
		}

		log.Println("setting channel for device", v.Id)
		m.actionChannel[id] = ch
	}
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
