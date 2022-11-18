package webHandle

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type requestInfoBlock struct {
	Path        string
	Components  []string
	Query       url.Values
	Body        []byte
	Session     string
	Request     *http.Request
	Response    http.ResponseWriter
	JsonMessage Generic
}

type clientItem struct {
	ClientId string
	AuthKey  string
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

	if r.Method == "GET" {
		if req.Components[1] == "ping" {
			w.WriteHeader(200)
			writeFlush(req.Response, "")
			return
		}
	}

	req.Body, err = io.ReadAll(r.Body)
	if err != nil {
		log.Println("http body read error:", err)
	}

	json.Unmarshal(req.Body, &req.JsonMessage)

	fmt.Println("| path:", req.Path)
	fmt.Println("| sessionid:", req.Session)
	fmt.Println("| query:", req.Query)
	fmt.Println("| body:", string(req.Body))

	if req.Components[0] == "v1" {
		fmt.Println("is V1")
		h.callV1api(req)
	} else {
		h.showPage(w, r, req.Components)
	}

}
