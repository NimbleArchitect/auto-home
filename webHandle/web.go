package webHandle

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	log "server/logger"
	"strings"
)

type requestInfoBlock struct {
	Path        string
	Components  []string
	Query       url.Values
	Body        []byte
	Session     string // current session id
	User        string // user id used to retrieve the session
	Request     *http.Request
	Response    http.ResponseWriter
	JsonMessage Generic
}

type clientItem struct {
	ClientId string `json:"clientid"`
	AuthKey  string `json:"authkey"`
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Actual web connection start
	var err error

	req := requestInfoBlock{
		Path:     cleanPath(r.URL.Path),
		Query:    r.URL.Query(),
		Session:  r.Header.Get("session"),
		Request:  r,
		Response: w,
	}
	req.Components = strings.Split(req.Path, "/")[1:]

	if r.Method == "GET" && len(req.Components) == 2 {
		if req.Components[1] == "ping" {
			w.WriteHeader(200)
			writeFlush(req.Response, "")
			return
		}
	}

	req.Body, err = io.ReadAll(r.Body)
	if err != nil {
		log.Error("http body read error:", err)
	}

	json.Unmarshal(req.Body, &req.JsonMessage)

	log.Debug("| path:", req.Path)
	log.Debug("| sessionid:", req.Session)
	log.Debug("| query:", req.Query)
	log.Debug("| body:", string(req.Body))

	if req.Components[0] == "v1" {
		h.callV1api(req)
	} else {
		h.showPage(req)
	}

}
