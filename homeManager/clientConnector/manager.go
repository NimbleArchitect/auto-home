package clientConnector

import (
	"fmt"
	"log"
	"net/http"
)

type Manager struct {
	writers map[string]*ClientWriter
}

type ClientWriter struct {
	responseWriter http.ResponseWriter
	canFlush       bool
	flusher        http.Flusher
	Id             string
}

// NewMagaer returns a new client connection manager object
func NewManager() *Manager {
	return &Manager{
		writers: make(map[string]*ClientWriter),
	}
}

func (m *Manager) ClientWriter(clientId string) *ClientWriter {
	var asd int
	_ = asd
	if len(m.writers) == 0 {
		fmt.Println("writers = nil")
		return nil
	}
	if len(clientId) <= 1 {
		fmt.Println("clientid = nil")
		return nil
	}

	fmt.Println("writers >>", m.writers)
	if val, ok := m.writers[clientId]; ok {
		return val
	}
	return nil
}

func (m *Manager) SetClient(clientId string, w http.ResponseWriter, r *http.Request) {
	fmt.Println("F:SetClient:start")
	fmt.Println("F:SetClient:clientId", clientId)
	val, ok := m.writers[clientId]
	fmt.Println("F:SetClient:ok", ok)
	fmt.Println("F:SetClient:val", val)

	if !ok {
		val := &ClientWriter{
			Id: clientId,
		}

		m.writers[clientId] = val
		fmt.Println("F:SetClient:val", val)
	}

	fmt.Println("F:SetClient:m.writers[clientId]", m.writers[clientId])
	val.responseWriter = w
	// time.Sleep(20 * time.Second)

	fmt.Println("F:SetClient:val.responseWriter", val.responseWriter)
	f, canFlush := w.(http.Flusher)
	fmt.Println("F:SetClient:f", f)
	val.canFlush = canFlush
	fmt.Println("F:SetClient:canFlush", canFlush)
	if canFlush {
		val.flusher = f
	} else {
		log.Print("Damn, no flush")
	}
	fmt.Println("F:SetClient:val", val)
	fmt.Println("F:SetClient:end")

}

// Write writes to the /actions/uuid channel opened from the client
func (c *ClientWriter) Write(text string) (int, error) {
	log.Println("http response Write:", text)
	fmt.Println(">>clientWrite:", c.responseWriter)
	bytesOut, err := c.responseWriter.Write([]byte(text + "\n"))

	if c.canFlush {
		log.Println("http response flush")
		c.flusher.Flush()
	}

	return bytesOut, err
}

func (m *Manager) CloseAll() {

	// for name, v := range m.writers {
	// 	if v == nil {
	// 		continue
	// 	}

	// 	v.Write(`{"Method": "shutdown"}`)
	// 	delete(m.writers, name)
	// }

}

// func (c *ClientWriter) SetReadWrite(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println(">>clientWrite:", w)
// 	c.responseWriter = w
// 	f, canFlush := w.(http.Flusher)
// 	c.canFlush = canFlush
// 	if canFlush {
// 		c.flusher = f
// 	} else {
// 		log.Print("Damn, no flush")
// 	}

// }
