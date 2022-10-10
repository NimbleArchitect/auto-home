package webHandle

import (
	"log"
	"net/http"
)

func getWriter(w http.ResponseWriter) func(string) (int, error) {
	svr := serverWriter{
		responseWriter: w,
	}

	f, canFlush := w.(http.Flusher)
	if canFlush {
		svr.flusher = f
		return svr.WriteFlush
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
