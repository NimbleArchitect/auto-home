package webHandle

import (
	"encoding/json"
	"io"
	log "server/logger"
)

func (h *Handler) showPage(req requestInfoBlock) {

	log.Debug(req.Components[0])

	switch req.Components[0] {
	case "public":
		// fmt.Println("fs>>", h.PublicPath)
		// h.fsHandle = http.FileServer(http.Dir(h.PublicPath))
		h.FsHandle.ServeHTTP(req.Response, req.Request)
		// fmt.Println(".>>")
		// h.streamFile(w, r, elements)
	case "reload":
		// TODO: this needs moving/fixing
		h.HomeManager.ReloadVMs()

	case "plugin":
		log.Info("/plugin")
		if req.Request.Method == "POST" {
			var jsonData map[string]interface{}

			if len(req.Body) == 0 {
				log.Error("no data recieved")
				return
			}

			// rawMsg, err := req.JsonMessage.Data.MarshalJSON()
			// if err != nil {
			// 	log.Println("unable to convert json string", err)
			// 	return
			// }

			err := json.Unmarshal(req.Body, &jsonData)

			// err = json.Unmarshal(rawMsg, &jsonData)
			if err != nil && err != io.EOF {
				log.Error("json decode error:", err)
				// TODO: need to return a proper error
				writeFlush(req.Response, "decode error")
				return
			}
			callRet := h.makePluginCall(req.Components[1:], jsonData)
			writeFlush(req.Response, string(callRet))
		}

	default:
		writeFlush(req.Response, "index page")
	}
}

func (h *Handler) makePluginCall(elements []string, postData map[string]interface{}) []byte {
	// TODO: needs safety checks adding

	pluginName := elements[0]
	callName := elements[1]

	out := h.HomeManager.WebCallPlugin(pluginName, callName, postData)

	return out
}

// func (h *Handler) streamFile(w http.ResponseWriter, r *http.Request, elements []string) {
// 	parts := elements[2:len(elements)]
// 	path := strings.Join(parts, "/")
// 	unsafePath := h.PublicPath + "/" + path
// 	info, err := os.Stat(unsafePath)
// 	if err != nil {
// 		fmt.Println("e>>", err)
// 		return
// 	}

// 	if info.IsDir() {
// 		return
// 	}

// 	if info.Mode().IsRegular() {

// 		fmt.Println(">>", info.Name())

// 		// if strings.HasPrefix() {

// 		// }
// 	}
// }
