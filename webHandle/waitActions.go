package webHandle

import (
	"fmt"
	"net/http"
	"sync"
)

type waitActions struct {
	inuse     bool // set this when the struct has been populated with data
	done      chan bool
	write     func(string) (int, error)
	resp      *http.ResponseWriter
	req       *http.Request
	sessionid string
}

func (w *waitActions) IsOpen() bool {
	if w.inuse {
		if w.resp != nil {
			return true
		}
	}
	return false
}

func (w *waitActions) Write(s string) (int, error) {
	fmt.Println("write waitAction:", s)
	if w.inuse {
		if w.resp != nil {
			return w.write(s + "\n")
		}
	}
	return 0, nil
}

// pre sets the action uuid and allocates a space for the actions uri
func (h *Handler) addDeviceActionList(id string, sessionid string) {
	if len(h.deviceActionList) <= 0 {
		h.deviceActionList = make(map[string]waitActions)
		h.lockActionList = sync.RWMutex{}
	}

	// we allocate empty here so we can verify the /events/uuid uri has a match
	h.setActionID(id, &waitActions{
		sessionid: sessionid,
	})
}

// readActionID returns a waitAction that matches the provided id
	h.lockActionList.RLock()
	val, ok := h.deviceActionList[id]
	h.lockActionList.RUnlock()
	return val, ok
}

func (h *Handler) setActionID(id string, waitAction *waitActions) {
	h.lockActionList.Lock()
	fmt.Println(">> setActionID lock:", id)
	h.lockActionList.Unlock()
}

func (h *Handler) deleteActionID(id string) {
	h.lockActionList.Lock()
	delete(h.deviceActionList, id)
	h.lockActionList.Unlock()
}
