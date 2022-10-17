package home

import (
	"log"
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
		m.devices.SetActionWriter(id, ch)
	}
	return nil
}
