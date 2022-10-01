package webHandle

import (
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
	default:
		w.Write([]byte("index page"))
	}
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
