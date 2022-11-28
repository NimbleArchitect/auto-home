package clientConnector

import (
	"errors"

	"net/http"
	"server/logger"
)

var debugLevel int

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
	debugLevel = logger.GetDebugLevel()

	return &Manager{
		writers: make(map[string]*ClientWriter),
	}
}

func (m *Manager) ClientWriter(clientId string) *ClientWriter {
	log := logger.New(&debugLevel)

	if len(m.writers) == 0 {
		log.Debug("writers = nil")
		return nil
	}
	if len(clientId) <= 1 {
		log.Debug("clientid = nil")
		return nil
	}

	log.Debug("m.writers:", m.writers)
	if val, ok := m.writers[clientId]; ok {
		return val
	}
	return nil
}

func (m *Manager) SetClient(clientId string, w http.ResponseWriter, r *http.Request) {
	log := logger.New(&debugLevel)

	log.Debug("clientId", clientId)
	val, ok := m.writers[clientId]
	log.Debug("ok", ok)
	log.Debug("val", val)

	if !ok {
		val := &ClientWriter{
			Id: clientId,
		}

		m.writers[clientId] = val
		log.Debug("val", val)
	}

	log.Debug("m.writers[clientId]", m.writers[clientId])
	val.responseWriter = w
	// time.Sleep(20 * time.Second)

	log.Debug("val.responseWriter", val.responseWriter)
	f, canFlush := w.(http.Flusher)
	log.Debug("f", f)
	val.canFlush = canFlush
	log.Debug("canFlush", canFlush)
	if canFlush {
		val.flusher = f
	} else {
		log.Info("Damn, no flush")
	}
	log.Debug("val", val)
}

// Write writes to the /actions/uuid channel opened from the client
func (c *ClientWriter) Write(text string) (int, error) {
	log := logger.New(&debugLevel)

	log.Debug("http response Write:", text)
	log.Debug("clientWrite:", c.responseWriter)
	if c.responseWriter == nil {
		return 0, errors.New("nil responseWriter")
	}

	bytesOut, err := c.responseWriter.Write([]byte(text + "\n"))

	if c.canFlush {
		log.Debug("http response flush")
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
