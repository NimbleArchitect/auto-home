package home

import (
	"log"
	"server/deviceManager"
)

func (m *Manager) AddHub(h Hub) error {
	log.Println("adding Hub:", h.Id)

	if len(m.hubs) == 0 {
		m.hubs = make(map[string]Hub)
	}

	m.hubs[h.ClientId] = h

	return nil
}

// DeviceExistsWithAction finds the device macthing id and compares is action id,
// returns true only if device id exists and action id matches
func (m *Manager) DeviceExistsWithClientId(deviceId string, clientId string) bool {
	var idMatch bool

	if len(clientId) <= 1 {
		return false
	}
	if len(deviceId) <= 1 {
		return false
	}

	idMatch = m.devices.ClientIdMatchesDevice(deviceId, clientId)

	return idMatch
}

func (m *Manager) AddDevice(d *deviceManager.Device, clientId string) error {
	log.Println("adding device:", d.Id)
	// if len(m.devices) == 0 {
	// 	m.devices = make(map[string]Device)
	// }

	//TODO: make sure we dont add properties of different types with the same name

	// TODO: this isnt thread safe
	// if updated, ok := m.devices.Device(d.Id); ok {
	if m.devices.HasDevice(d.Id) {
		// device already exists, so check/update the properties
		log.Println("device exists... updating")

		//TODO: need to find a way to authorize the actionid update as rouge clients could reregister and take over the device from teh servers point of view
		// updated.ClientId = d.ClientId
		// m.setDeviceClient(d.Id, clientId)

		// m.devices[d.Id] = updated
		// } else {
		// 	// device is new, so we create it as is
		// 	m.devices.SetDevice(d.Id) = d
	}

	m.devices.SetDevice(d.Id, d)

	return nil
}

// func (m *Manager) GetDial(id string, name string) (DialProperty, error) {
// 	var err error

// 	if dev, found := m.devices[id]; found {
// 		if prop, ok := dev.PropertyDial[name]; ok {
// 			if prop.Mode == WO {
// 				return DialProperty{}, ErrWriteOnlyProperty
// 			}

// 			return prop, nil
// 		} else {
// 			err = errors.New("missing dial property " + name)
// 		}
// 	} else {
// 		err = errors.New("missing device " + id)
// 	}

// 	return DialProperty{}, err
// }

// func (m *Manager) GetDialValue(id string, name string) (int, bool) {
// 	prop, err := m.GetDial(id, name)
// 	if err != nil {
// 		return 0, false
// 	}

// 	return prop.Value, true
// }

// func (m *Manager) UpdateDial(id string, name string, value int) error {

// 	prop, err := m.GetDial(id, name)

// 	if err == nil {
// 		if prop.Mode == RW || prop.Mode == WO {
// 			jsonOut := m.MakeAction(id, name, DIAL, fmt.Sprint(value))
// 			log.Println("sending action", jsonOut, "to", id)
// 			// TODO: when the client is not online writes are lost, need to replay them somehow?
// 			if _, ok := m.actionChannel[id]; ok {
// 				_, err := m.actionChannel[id].Write(jsonOut)
// 				if err != nil {
// 					log.Println(err)
// 				}

// 				prop := m.devices[id].PropertyDial[name]
// 				prop.Value = value
// 				m.devices[id].PropertyDial[name] = prop
// 			}
// 			return nil
// 		}
// 	}

// 	return err
// }

// func (m *Manager) GetSwitch(id string, name string) (SwitchProperty, error) {
// 	var err error

// 	if dev, found := m.devices[id]; found {
// 		if prop, ok := dev.PropertySwitch[name]; ok {
// 			if prop.Mode == WO {
// 				return SwitchProperty{}, ErrWriteOnlyProperty
// 			}

// 			return prop, nil
// 		} else {
// 			err = errors.New("missing switch property " + name)
// 		}
// 	} else {
// 		err = errors.New("missing device " + id)
// 	}

// 	return SwitchProperty{}, err
// }

// func (m *Manager) GetSwitchValue(id string, name string) (string, bool) {
// 	prop, err := m.GetSwitch(id, name)
// 	if err != nil {
// 		return "", false
// 	}
// 	return prop.Value.String(), true
// }

// func (m *Manager) UpdateSwitch(id string, name string, value string) error {

// 	prop, err := m.GetSwitch(id, name)

// 	if err == nil {
// 		if prop.Mode == RW || prop.Mode == WO {
// 			clientid := m.devices[id].ClientId

// 			jsonOut := m.MakeAction(id, name, SWITCH, value)
// 			log.Println("sending action", jsonOut, "to", clientid)
// 			if _, ok := m.actionChannel[id]; ok {
// 				_, err := m.actionChannel[clientid].Write(jsonOut)
// 				if err != nil {
// 					log.Println(err)
// 				}

// 				prop := m.devices[id].PropertySwitch[name]
// 				prop.Value.Set(value)
// 				m.devices[id].PropertySwitch[name] = prop
// 			}
// 			return nil
// 		}
// 	}

// 	return err
// }

// func (m *Manager) GetButton(id string, name string) (ButtonProperty, error) {
// 	var err error

// 	if dev, found := m.devices[id]; found {
// 		if prop, ok := dev.PropertyButton[name]; ok {
// 			if prop.Mode == WO {
// 				return ButtonProperty{}, ErrWriteOnlyProperty
// 			}

// 			return prop, nil
// 		} else {
// 			err = errors.New("missing button property " + name)
// 		}
// 	} else {
// 		err = errors.New("missing device " + id)
// 	}

// 	return ButtonProperty{}, err
// }

// func (m *Manager) GetButtonValue(id string, name string) (string, bool) {
// 	prop, err := m.GetButton(id, name)
// 	if err != nil {
// 		return "", false
// 	}
// 	return prop.Value.String(), true
// }

// func (m *Manager) UpdateButton(id string, name string, value string) error {
// 	prop, err := m.GetButton(id, name)

// 	if err == nil {
// 		if prop.Mode == RW || prop.Mode == WO {
// 			clientid := m.devices[id].ClientId

// 			jsonOut := m.MakeAction(id, name, BUTTON, value)
// 			log.Println("sending action", jsonOut, "to", clientid)
// 			if _, ok := m.actionChannel[id]; ok {
// 				_, err := m.actionChannel[clientid].Write(jsonOut)
// 				if err != nil {
// 					log.Println(err)
// 				}

// 				// prop := m.devices[id].PropertyButton[name]
// 				// prop.Value.Set(value)
// 				// m.devices[id].PropertyButton[name] = prop
// 			}
// 			return nil
// 		}
// 	}

// 	return err
// }

// func (m *Manager) GetText(id string, name string) (TextProperty, error) {
// 	var err error

// 	if dev, found := m.devices[id]; found {
// 		if prop, ok := dev.PropertyText[name]; ok {
// 			if prop.Mode == WO {
// 				return TextProperty{}, ErrWriteOnlyProperty
// 			}

// 			return prop, nil
// 		} else {
// 			err = errors.New("missing text property " + name)
// 		}
// 	} else {
// 		err = errors.New("missing device " + id)
// 	}

// 	return TextProperty{}, err
// }

// func (m *Manager) GetTextValue(id string, name string) (string, bool) {
// 	prop, err := m.GetText(id, name)
// 	if err != nil {
// 		return "", false
// 	}
// 	return prop.Value, true
// }

// func (m *Manager) UpdateText(id string, name string, value string) error {

// 	prop, err := m.GetText(id, name)

// 	if err == nil {
// 		if prop.Mode == RW || prop.Mode == WO {
// 			clientid := m.devices[id].ClientId

// 			jsonOut := m.MakeAction(id, name, TEXT, value)
// 			log.Println("sending action", jsonOut, "to", clientid)
// 			if _, ok := m.actionChannel[id]; ok {
// 				_, err := m.actionChannel[clientid].Write(jsonOut)
// 				if err != nil {
// 					log.Println(err)
// 				}

// 				prop := m.devices[id].PropertyText[name]
// 				prop.Value = value
// 				m.devices[id].PropertyText[name] = prop
// 			}
// 			return nil
// 		}
// 	}

// 	return err
// }
