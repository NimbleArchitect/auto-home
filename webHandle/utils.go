package webHandle

import (
	"net/http"
	"server/logger"
)

func cleanPath(rawPath string) string {
	var runePath []rune

	isSep := false
	for _, c := range rawPath {
		if c == []rune("/")[0] {
			if !isSep {
				runePath = append(runePath, c)
			}
			isSep = true
		} else {
			isSep = false
		}

		if !isSep {
			runePath = append(runePath, c)
		}

	}

	return string(runePath)
}

// writeFlush attempts an immediate flush of the buffers after sending text
func writeFlush(w http.ResponseWriter, text string) (int, error) {
	log := logger.New(&debugLevel)

	if len(text) > 0 {
		log.Debug("text", text)
	}

	bytesOut, err := w.Write([]byte(text + "\n"))

	f, canFlush := w.(http.Flusher)
	if canFlush {
		f.Flush()
	} else {
		log.Info("Damn, no flush")
	}
	return bytesOut, err
}
