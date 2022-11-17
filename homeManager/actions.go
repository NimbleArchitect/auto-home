package home

import (
	"net/http"
)

func (m *Manager) SetClient(clientId string, w http.ResponseWriter, r *http.Request) {
	m.clientConnection.SetClient(clientId, w, r)
}

// func (m *Manager) registerActionChannel(id string, client *clientConnector.ClientWriter) error {
// 	log.Println("registerActionChannel id", id)

// 	// TODO: this might not work

// 	// for _, v := range m.FindDeviceWithClientID(id) {
// 	if len(m.devices.FindDeviceWithClientID(id)) > 0 {
// 		// v.Id
// 		log.Println("device ID found")
// 		if len(m.clientConnections) == 0 {
// 			log.Println("empty action channel found, creating new")
// 			m.clientConnections = make(map[string]*clientConnector.ClientWriter)
// 		}

// 		log.Println("setting channel for device:", id)
// 		m.clientConnections[id] = client
// 		m.devices.SetActionWriter(id, client)
// 	}
// 	return nil
// }

// ClientConnection returns the clientWriter that matches the provided clientid
// func (m *Manager) ClientConnection(clientId string) *clientConnector.ClientWriter {
// 	// TODO: add locking
// 	fmt.Println("1.1>>", clientId, "<<")
// 	client, ok := m.clientConnections[clientId]
// 	if ok {
// 		fmt.Println("1.2>>")
// 		fmt.Println("2>> ***", m.clientConnections)
// 		return client
// 	}
// 	fmt.Println("1.3>>")

// 	fmt.Println("creating new uuid")
// 	uid, _ := uuid.NewUUID()
// 	client = &clientConnector.ClientWriter{
// 		Id: uid.String(),
// 	}
// 	fmt.Println("1.4>>")
// 	m.clientConnections[clientId] = client
// 	fmt.Println("1.5>>")
// 	fmt.Println("1>> ***", m.clientConnections)
// 	//
// 	//
// 	// TODO: changed values are not pushed back to the client, this needs fixing asap
// 	//
// 	//
// 	m.devices.SetActionWriter(clientId, client)
// 	return client
// }

// SetClientConnection sets the client to the provided connector
// func (m *Manager) setClientConnection(clientId string, connector *clientConnector.ClientWriter) {
// 	// TODO: add locking
// 	m.clientConnections[clientId] = connector
// 	fmt.Println(">> ******************************************")
// 	fmt.Println(">> ******************************************")
// 	fmt.Println(">> ******************************************")
// 	fmt.Println(">> ***", m.clientConnections)
// 	fmt.Println(">> ******************************************")
// 	fmt.Println(">> ******************************************")
// }
