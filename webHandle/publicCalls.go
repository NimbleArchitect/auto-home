package webHandle

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (h *Handler) showPage(w http.ResponseWriter, r *http.Request, elements []string) {
	switch elements[1] {
	case "public":
		// fmt.Println("fs>>", h.PublicPath)
		// h.fsHandle = http.FileServer(http.Dir(h.PublicPath))
		h.FsHandle.ServeHTTP(w, r)
		// fmt.Println(".>>")
		// h.streamFile(w, r, elements)
	case "reload":
		// TODO: this needs moving/fixing
		h.HomeManager.ReloadVMs()

	case "plugin":
		if r.Method == "POST" {
			var jsonData map[string]interface{}

			err := json.NewDecoder(r.Body).Decode(&jsonData)
			if err != nil && err != io.EOF {
				fmt.Println("json decode error:", err)
				// TODO: need to return a proper error
				w.Write([]byte("decode error"))
				return
			}
			callRet := h.makePluginCall(elements[2:], jsonData)
			w.Write(callRet)
		}

	default:
		w.Write([]byte("index page"))
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
