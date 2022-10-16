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
