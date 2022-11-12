package webHandle

import (
	"log"
	"net/http"
)

// writeFlush returns a writeFlush function if it exists
func writeFlush(w http.ResponseWriter, text string) (int, error) {
	log.Println(">> writeFlush:", text)
	bytesOut, err := w.Write([]byte(text))

	f, canFlush := w.(http.Flusher)
	if canFlush {
		f.Flush()
	} else {
		log.Print("Damn, no flush")
		return svr.Write
	}

}

type serverWriter struct {
	flusher        http.Flusher
	responseWriter http.ResponseWriter
}

func (s *serverWriter) WriteFlush(text string) (int, error) {
	bytesOut, err := s.Write(text)

	log.Println("http response flush")
	s.flusher.Flush()
	return bytesOut, err
}

func (s *serverWriter) Write(text string) (int, error) {
	log.Println("http response Write:", text)
	bytesOut, err := s.responseWriter.Write([]byte(text))

	return bytesOut, err
}
