package webHandle

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type login struct {
	Token string
}

type sessionItem struct {
	actionId string
	ClientId string
}

type clientItem struct {
	ClientId string
	AuthKey  string
}

func (h *Handler) getSessionWithToken(token string) (string, bool) {
	// check if token exists
	log.Println("recieved token:", token)
	if len(h.sessionTable) == 0 {
		h.sessionTable = make(map[string]sessionItem)
		h.lockSessionTable = sync.RWMutex{}
	}

	if clientInfo, ok := h.clientTable[token]; ok {
		session := uuid.New().String()
		h.sessionTable[session] = sessionItem{
			actionId: uuid.New().String(),
			ClientId: clientInfo.ClientId,
		}

		// h.addDeviceActionList(newUuid)
		// log.Println(h.clientTable)
		return session, true
	}
	return "", false
}

func (h *Handler) lookupSessionClient(sessionid string) (sessionItem, bool) {
	// TODO: check if sessionid is empty

	h.lockSessionTable.RLock()
	val, ok := h.sessionTable[sessionid]
	h.lockSessionTable.RUnlock()

	return val, ok
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Println("processing Url:", r.RequestURI)

	if len(r.RequestURI) == 0 {
		h.showPage(w, r, []string{})
	} else {
		elements := strings.Split(r.RequestURI, "/")
		switch elements[1] {
		case "v1":
			h.callV1api(w, r, elements)

		default:
			h.showPage(w, r, elements)
		}
	}
}
