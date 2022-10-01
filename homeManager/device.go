package home

import (
	"errors"
	"fmt"
	"log"
)

func (m *Manager) AddHub(h Hub) error {
	log.Println("adding Hub:", h.Id)

	if len(m.hubs) == 0 {
		m.hubs = make(map[string]Hub)
	}

	m.hubs[h.ClientId] = h

	return nil
}

// func (h *Handler) AddHubFromJson(json string) (string, error) {
// 	n := home.Hub{}

// 	n.Description = d.Description
// 	n.Id = d.Id
// 	n.ClientId = clientId
// 	n.Name = d.Name
// 	n.Help = d.Help

// 	for _, v := range d.Devices {
// 		if err := h.regDeviceList(v, clientId); err == nil {
// 			n.Devices = append(n.Devices, v.Id)
// 		} else {
// 			return err
// 		}

// 	}

// 	fmt.Println("current-devices:", h.HomeManager.GetDevices())
// 	return h.HomeManager.AddHub(n)
// }

func (m *Manager) GetDevices() map[string]Device {
	return m.devices
}

func (m *Manager) DeviceExists(id string) bool {
	if len(id) <= 1 {
		return false
	}
	_, ok := m.devices[id]
	return ok
}

// DeviceExistsWithAction finds the device macthing id and compares is action id,
// returns true only if device id exists and action id matches
func (m *Manager) DeviceExistsWithClientId(deviceId string, clientId string) bool {
	if len(clientId) <= 1 {
		return false
	}
	if len(deviceId) <= 1 {
		return false
	}

	// fmt.Println("1>>", deviceId)
	// fmt.Println("2>>", m.devices)

	dev, ok := m.devices[deviceId]
	if ok {
		// fmt.Println("3>>", dev.ClientId)
		if dev.ClientId == clientId {
			return true
		}
	}
	return false
}

func (m *Manager) FindDeviceWithClientID(clientId string) []Device {
	var deviceList []Device

	for _, v := range m.devices {
		if v.ClientId == clientId {
			deviceList = append(deviceList, v)
		}
	}

	return deviceList
}

func (m *Manager) AddDevice(d Device) error {
	log.Println("adding device:", d.Id)
	if len(m.devices) == 0 {
		m.devices = make(map[string]Device)
	}

	// TODO: this isnt thread safe
	if updated, ok := m.devices[d.Id]; ok {
		// device already exists, so check/update the properties
		log.Println("device exists... updating")
		updated.Help = d.Help
		updated.Description = d.Description
		updated.PropertyDial = d.PropertyDial
		updated.PropertySwitch = d.PropertySwitch
		updated.Uploads = d.Uploads

		//TODO: need to find a way to authorize the actionid update as rouge clients could reregister and take over the device from teh servers point of view
		// updated.ClientId = d.ClientId

		m.devices[d.Id] = updated
	} else {
		// device is new, so we create it as is
		m.devices[d.Id] = d
	}

	return nil
}

func (m *Manager) GetDial(id string, name string) (DialProperty, error) {
	var err error

	if dev, found := m.devices[id]; found {
		if prop, ok := dev.PropertyDial[name]; ok {
			return prop, nil
		} else {
			err = errors.New("missing dial property " + name)
		}
	} else {
		err = errors.New("missing device " + id)
	}

	return DialProperty{}, err
}

func (m *Manager) GetDialValue(id string, name string) (int, bool) {
	prop, err := m.GetDial(id, name)
	if err != nil {
		return 0, false
	}

	return prop.Value, true
}

func (m *Manager) UpdateDial(id string, name string, value int) error {

	prop, err := m.GetDial(id, name)

	if err == nil {
		if prop.Mode == RW || prop.Mode == WO {
			jsonOut := m.MakeAction(id, name, DIAL, fmt.Sprint(value))
			log.Println("sending action", jsonOut, "to", id)
			// TODO: when the client is not online writes are lost, need to replay them somehow?
			if _, ok := m.actionChannel[id]; ok {
				_, err := m.actionChannel[id].Write(jsonOut)
				if err != nil {
					log.Println(err)
				}

				prop := m.devices[id].PropertyDial[name]
				prop.Value = value
				m.devices[id].PropertyDial[name] = prop
			}
			return nil
		}
	}

	return err
}

func (m *Manager) GetSwitch(id string, name string) (SwitchProperty, error) {
	var err error

	if dev, found := m.devices[id]; found {
		if prop, ok := dev.PropertySwitch[name]; ok {
			return prop, nil
		} else {
			err = errors.New("missing switch property " + name)
		}
	} else {
		err = errors.New("missing device " + id)
	}

	return SwitchProperty{}, err
}

func (m *Manager) GetSwitchValue(id string, name string) (string, bool) {
	prop, err := m.GetSwitch(id, name)
	if err != nil {
		return "", false
	}
	return prop.Value.String(), true
}

func (m *Manager) UpdateSwitch(id string, name string, value string) error {

	prop, err := m.GetSwitch(id, name)

	if err == nil {
		if prop.Mode == RW || prop.Mode == WO {
			clientid := m.devices[id].ClientId

			jsonOut := m.MakeAction(id, name, SWITCH, value)
			log.Println("sending action", jsonOut, "to", clientid)
			if _, ok := m.actionChannel[id]; ok {
				_, err := m.actionChannel[clientid].Write(jsonOut)
				if err != nil {
					log.Println(err)
				}

				prop := m.devices[id].PropertySwitch[name]
				prop.Value.Set(value)
				m.devices[id].PropertySwitch[name] = prop
			}
			return nil
		}
	}

	return err
}
